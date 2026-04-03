package litepage

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os/exec"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"
	"sync"
	"time"
)

type Browser struct {
	Addr       string
	options    *ChromiumOptions
	httpClient *http.Client
	process    *exec.Cmd
	browserWS  string
	conn       *cdpConn
	owned      bool
	mu         sync.Mutex
	closed     bool
}

type browserVersion struct {
	Browser              string `json:"Browser"`
	ProtocolVersion      string `json:"Protocol-Version"`
	UserAgent            string `json:"User-Agent"`
	WebSocketDebuggerURL string `json:"webSocketDebuggerUrl"`
}

type targetInfo struct {
	Description          string `json:"description"`
	DevtoolsFrontendURL  string `json:"devtoolsFrontendUrl"`
	ID                   string `json:"id"`
	Title                string `json:"title"`
	Type                 string `json:"type"`
	URL                  string `json:"url"`
	WebSocketDebuggerURL string `json:"webSocketDebuggerUrl"`
}

type Cookie struct {
	Name     string  `json:"name"`
	Value    string  `json:"value"`
	Domain   string  `json:"domain,omitempty"`
	Path     string  `json:"path,omitempty"`
	Expires  float64 `json:"expires,omitempty"`
	HTTPOnly bool    `json:"httpOnly,omitempty"`
	Secure   bool    `json:"secure,omitempty"`
	SameSite string  `json:"sameSite,omitempty"`
	URL      string  `json:"url,omitempty"`
}

type createTargetResult struct {
	TargetID string `json:"targetId"`
}

func NewBrowser(opts *ChromiumOptions) (*Browser, error) {
	options := opts.clone()
	if err := options.validate(); err != nil {
		return nil, err
	}

	b := &Browser{
		Addr:       options.Address(),
		options:    options,
		httpClient: &http.Client{Timeout: 5 * time.Second},
	}

	if err := b.ensureRunning(context.Background()); err != nil {
		_ = b.Quit()
		return nil, err
	}
	return b, nil
}

func CreateBrowser(browser *Browser) *Browser {
	if browser == nil {
		return nil
	}
	return browser
}

func (b *Browser) ensureRunning(ctx context.Context) error {
	if b.isReachable(ctx) {
		return b.waitUntilAttachable(ctx)
	}
	if err := b.launch(ctx); err != nil {
		return err
	}
	return b.waitUntilAttachable(ctx)
}

func (b *Browser) isReachable(ctx context.Context) bool {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, b.httpURL("/json/version"), nil)
	if err != nil {
		return false
	}
	resp, err := b.httpClient.Do(req)
	if err != nil {
		return false
	}
	defer func() { _ = resp.Body.Close() }()
	return resp.StatusCode == http.StatusOK
}

func (b *Browser) launch(ctx context.Context) error {
	browserPath, err := resolveBrowserPath(b.options.BrowserPath())
	if err != nil {
		return err
	}
	userDataDir, err := b.options.ensureUserDataPath()
	if err != nil {
		return err
	}

	args := buildLaunchArgs(b.options, userDataDir)
	cmd := exec.CommandContext(ctx, browserPath, args...)
	if err := cmd.Start(); err != nil {
		return err
	}
	b.process = cmd
	b.owned = true

	retries := b.options.Retry()
	if retries <= 0 {
		retries = defaultRetryTimes
	}
	interval := b.options.RetryInterval()
	if interval <= 0 {
		interval = defaultRetryInterval
	}
	for i := 0; i < retries; i++ {
		if b.isReachable(ctx) {
			return nil
		}
		time.Sleep(interval)
	}
	return fmt.Errorf("chromium did not expose devtools endpoint at %s", b.Addr)
}

func (b *Browser) waitUntilAttachable(ctx context.Context) error {
	retries := b.options.Retry()
	if retries <= 0 {
		retries = defaultRetryTimes
	}
	interval := b.options.RetryInterval()
	if interval <= 0 {
		interval = defaultRetryInterval
	}

	var lastErr error
	for i := 0; i < retries; i++ {
		if ctx.Err() != nil {
			return ctx.Err()
		}
		if err := b.attach(); err == nil {
			return nil
		} else {
			lastErr = err
		}
		time.Sleep(interval)
	}
	if lastErr != nil {
		return fmt.Errorf("devtools endpoint is reachable but browser is not attachable yet: %w", lastErr)
	}
	return errors.New("devtools endpoint is reachable but browser is not attachable")
}

func buildLaunchArgs(opts *ChromiumOptions, userDataDir string) []string {
	args := []string{
		"--remote-debugging-address=" + defaultHost,
		"--remote-debugging-port=" + strconv.Itoa(opts.port),
		"--remote-allow-origins=*",
		"--user-data-dir=" + userDataDir,
		"--disable-background-networking",
		"--disable-renderer-backgrounding",
		"--disable-background-timer-throttling",
		"--disable-popup-blocking",
		"--no-default-browser-check",
		"--no-first-run",
	}
	if opts.downloadPath != "" {
		args = append(args, "--download-default-directory="+opts.downloadPath)
	}
	if opts.tmpPath != "" {
		args = append(args, "--disk-cache-dir="+opts.tmpPath)
	}
	if opts.proxy != "" {
		args = append(args, "--proxy-server="+opts.proxy)
	}
	if len(opts.extensions) > 0 {
		joined := strings.Join(opts.extensions, ",")
		args = append(args, "--disable-extensions-except="+joined, "--load-extension="+joined)
	}
	args = append(args, opts.arguments...)
	return dedupeArgs(args)
}

func dedupeArgs(args []string) []string {
	seen := make(map[string]struct{}, len(args))
	out := make([]string, 0, len(args))
	for _, arg := range args {
		if _, ok := seen[arg]; ok {
			continue
		}
		seen[arg] = struct{}{}
		out = append(out, arg)
	}
	return out
}

func resolveBrowserPath(configured string) (string, error) {
	if configured != "" {
		return configured, nil
	}
	candidates := []string{"google-chrome", "chromium", "chromium-browser", "chrome", "msedge", "microsoft-edge"}
	if runtime.GOOS == "darwin" {
		candidates = append([]string{
			"/Applications/Google Chrome.app/Contents/MacOS/Google Chrome",
			"/Applications/Microsoft Edge.app/Contents/MacOS/Microsoft Edge",
			"/Applications/Chromium.app/Contents/MacOS/Chromium",
		}, candidates...)
	}
	for _, candidate := range candidates {
		if strings.Contains(candidate, string(filepath.Separator)) {
			if _, err := exec.LookPath(candidate); err == nil {
				return candidate, nil
			}
			continue
		}
		if path, err := exec.LookPath(candidate); err == nil {
			return path, nil
		}
	}
	return "", errors.New("no chromium-compatible browser binary found")
}

func (b *Browser) attach() error {
	version, err := b.Version()
	if err != nil {
		return err
	}
	if version.WebSocketDebuggerURL == "" {
		return errors.New("webSocketDebuggerUrl is empty")
	}
	b.browserWS = version.WebSocketDebuggerURL
	if b.conn != nil {
		_ = b.conn.close()
		b.conn = nil
	}
	conn, err := newCDPConn(version.WebSocketDebuggerURL)
	if err != nil {
		return err
	}
	b.conn = conn
	return nil
}

func (b *Browser) httpURL(path string) string {
	return "http://" + normalizeDebugAddress(b.Addr) + path
}

func (b *Browser) Version() (*browserVersion, error) {
	var version browserVersion
	if err := b.getJSON("/json/version", &version); err != nil {
		return nil, err
	}
	return &version, nil
}

func (b *Browser) Targets() ([]targetInfo, error) {
	var targets []targetInfo
	if err := b.getJSON("/json/list", &targets); err != nil {
		return nil, err
	}
	return targets, nil
}

func (b *Browser) NewPage(rawURL string) (*Page, error) {
	if rawURL == "" {
		rawURL = "about:blank"
	}
	var result createTargetResult
	if err := b.RunCDP("Target.createTarget", map[string]any{"url": rawURL}, &result); err != nil {
		return nil, err
	}
	return b.Page(result.TargetID)
}

func (b *Browser) Pages() ([]*Page, error) {
	targets, err := b.Targets()
	if err != nil {
		return nil, err
	}
	pages := make([]*Page, 0, len(targets))
	for _, target := range targets {
		if target.Type != "page" {
			continue
		}
		page, pageErr := newPage(b, target)
		if pageErr != nil {
			return nil, pageErr
		}
		pages = append(pages, page)
	}
	return pages, nil
}

func (b *Browser) Page(targetID string) (*Page, error) {
	targets, err := b.Targets()
	if err != nil {
		return nil, err
	}
	for _, target := range targets {
		if target.ID == targetID {
			return newPage(b, target)
		}
	}
	return nil, fmt.Errorf("page %s not found", targetID)
}

func (b *Browser) LatestPage() (*Page, error) {
	targets, err := b.Targets()
	if err != nil {
		return nil, err
	}
	for i := len(targets) - 1; i >= 0; i-- {
		if targets[i].Type == "page" {
			return newPage(b, targets[i])
		}
	}
	return nil, errors.New("no page target found")
}

func (b *Browser) RunCDP(cmd string, cmdArgs map[string]any, out ...any) error {
	if b.conn == nil {
		return errors.New("browser cdp connection is not ready")
	}
	var result any
	if len(out) > 0 {
		result = out[0]
	}
	return b.conn.call(cmd, cmdArgs, result)
}

func (b *Browser) ClosePage(targetID string) error {
	req, err := http.NewRequest(http.MethodGet, b.httpURL("/json/close/"+targetID), nil)
	if err != nil {
		return err
	}
	resp, err := b.httpClient.Do(req)
	if err != nil {
		return err
	}
	defer func() { _ = resp.Body.Close() }()
	if resp.StatusCode >= 400 {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("close page failed: %s", bytes.TrimSpace(body))
	}
	return nil
}

func (b *Browser) ActivatePage(targetID string) error {
	return b.RunCDP("Target.activateTarget", map[string]any{"targetId": targetID})
}

func (b *Browser) Cookies(urls ...string) ([]Cookie, error) {
	params := map[string]any{}
	if len(urls) > 0 {
		params["urls"] = urls
	}
	var result struct {
		Cookies []Cookie `json:"cookies"`
	}
	if err := b.RunCDP("Network.getCookies", params, &result); err != nil {
		return nil, err
	}
	return result.Cookies, nil
}

func (b *Browser) SetCookies(cookies ...Cookie) error {
	if len(cookies) == 0 {
		return nil
	}
	var result struct {
		Success bool `json:"success"`
	}
	for _, cookie := range cookies {
		if err := b.RunCDP("Network.setCookie", cookieToParams(cookie), &result); err != nil {
			return err
		}
		if !result.Success {
			return fmt.Errorf("set cookie %q failed", cookie.Name)
		}
	}
	return nil
}

func (b *Browser) ClearCookies() error {
	return b.RunCDP("Network.clearBrowserCookies", nil)
}

func (b *Browser) Quit() error {
	b.mu.Lock()
	if b.closed {
		b.mu.Unlock()
		return nil
	}
	b.closed = true
	conn := b.conn
	process := b.process
	owned := b.owned
	b.mu.Unlock()

	if conn != nil {
		_ = conn.close()
	}
	if owned && process != nil && process.Process != nil {
		return process.Process.Kill()
	}
	return nil
}

func cookieToParams(cookie Cookie) map[string]any {
	params := map[string]any{
		"name":  cookie.Name,
		"value": cookie.Value,
	}
	if cookie.URL != "" {
		params["url"] = cookie.URL
	}
	if cookie.Domain != "" {
		params["domain"] = cookie.Domain
	}
	if cookie.Path != "" {
		params["path"] = cookie.Path
	}
	if cookie.Expires != 0 {
		params["expires"] = cookie.Expires
	}
	if cookie.Secure {
		params["secure"] = cookie.Secure
	}
	if cookie.HTTPOnly {
		params["httpOnly"] = cookie.HTTPOnly
	}
	if cookie.SameSite != "" {
		params["sameSite"] = cookie.SameSite
	}
	return params
}

func (b *Browser) getJSON(path string, out any) error {
	resp, err := b.httpClient.Get(b.httpURL(path))
	if err != nil {
		return err
	}
	defer func() { _ = resp.Body.Close() }()
	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("devtools http %s returned %d: %s", path, resp.StatusCode, bytes.TrimSpace(body))
	}
	return json.NewDecoder(resp.Body).Decode(out)
}

func originFromWSURL(raw string) string {
	u, err := url.Parse(raw)
	if err != nil {
		return "http://127.0.0.1/"
	}
	scheme := "http"
	if u.Scheme == "wss" {
		scheme = "https"
	}
	return scheme + "://" + u.Host + "/"
}
