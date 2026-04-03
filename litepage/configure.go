package litepage

import (
	"errors"
	"fmt"
	"github.com/Heartfilia/litetools/utils/litedir"
	"net"
	"net/url"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"
)

const (
	defaultHost          = "127.0.0.1"
	defaultPort          = 9222
	defaultRetryTimes    = 20
	defaultRetryInterval = 250 * time.Millisecond
)

// ChromiumOptions 借鉴 DrissionPage 的 ChromiumOptions 概念，
// 但这里按 Go 的方式保留浏览器接管和启动配置。
type ChromiumOptions struct {
	yamlConf      string
	userDataPath  string
	user          string
	headless      bool
	downloadPath  string
	tmpPath       string
	arguments     []string
	browserPath   string
	extensions    []string
	prefs         map[string]string
	flags         map[string]string
	address       string
	loadMode      string
	proxy         string
	port          int
	retry         int
	retryInterval time.Duration
}

func NewChromiumOptions() *ChromiumOptions {
	return &ChromiumOptions{
		address:       net.JoinHostPort(defaultHost, strconv.Itoa(defaultPort)),
		port:          defaultPort,
		loadMode:      "normal",
		retry:         defaultRetryTimes,
		retryInterval: defaultRetryInterval,
		prefs:         make(map[string]string),
		flags:         make(map[string]string),
		arguments:     make([]string, 0),
		extensions:    make([]string, 0),
	}
}

func (c *ChromiumOptions) clone() *ChromiumOptions {
	if c == nil {
		return NewChromiumOptions()
	}
	cp := *c
	cp.arguments = append([]string(nil), c.arguments...)
	cp.extensions = append([]string(nil), c.extensions...)
	cp.prefs = copyStringMap(c.prefs)
	cp.flags = copyStringMap(c.flags)
	return &cp
}

func copyStringMap(src map[string]string) map[string]string {
	if len(src) == 0 {
		return make(map[string]string)
	}
	dst := make(map[string]string, len(src))
	for k, v := range src {
		dst[k] = v
	}
	return dst
}

func (c *ChromiumOptions) DownloadPath() string           { return c.downloadPath }
func (c *ChromiumOptions) BrowserPath() string            { return c.browserPath }
func (c *ChromiumOptions) UserDataPath() string           { return c.userDataPath }
func (c *ChromiumOptions) TmpPath() string                { return c.tmpPath }
func (c *ChromiumOptions) User() string                   { return c.user }
func (c *ChromiumOptions) Proxy() string                  { return c.proxy }
func (c *ChromiumOptions) Address() string                { return c.address }
func (c *ChromiumOptions) Arguments() []string            { return append([]string(nil), c.arguments...) }
func (c *ChromiumOptions) Extensions() []string           { return append([]string(nil), c.extensions...) }
func (c *ChromiumOptions) Preferences() map[string]string { return copyStringMap(c.prefs) }
func (c *ChromiumOptions) Flags() map[string]string       { return copyStringMap(c.flags) }
func (c *ChromiumOptions) Retry() int                     { return c.retry }
func (c *ChromiumOptions) RetryInterval() time.Duration   { return c.retryInterval }

func (c *ChromiumOptions) LoadMode() string {
	if c.loadMode == "" {
		return "normal"
	}
	return c.loadMode
}

func (c *ChromiumOptions) SetRetry(ts int, interval float64) *ChromiumOptions {
	c.retry = ts
	if interval <= 0 {
		c.retryInterval = defaultRetryInterval
		return c
	}
	c.retryInterval = time.Duration(interval * float64(time.Second))
	return c
}

func (c *ChromiumOptions) SetRetryDuration(ts int, interval time.Duration) *ChromiumOptions {
	c.retry = ts
	if interval <= 0 {
		interval = defaultRetryInterval
	}
	c.retryInterval = interval
	return c
}

func (c *ChromiumOptions) SetAddress(address string) *ChromiumOptions {
	c.address = normalizeDebugAddress(address)
	host, port, err := net.SplitHostPort(c.address)
	if err == nil {
		_ = host
		if p, convErr := strconv.Atoi(port); convErr == nil {
			c.port = p
		}
	}
	return c
}

func (c *ChromiumOptions) SetLocalPort(port int) *ChromiumOptions {
	if port <= 0 {
		port = defaultPort
	}
	c.port = port
	c.address = net.JoinHostPort(defaultHost, strconv.Itoa(port))
	return c
}

func (c *ChromiumOptions) SetBrowserPath(browserPath string) *ChromiumOptions {
	c.browserPath = browserPath
	return c
}

func (c *ChromiumOptions) SetUserDataPath(userDataPath string) *ChromiumOptions {
	c.userDataPath = userDataPath
	return c
}

func (c *ChromiumOptions) SetDownloadPath(downloadPath string) *ChromiumOptions {
	c.downloadPath = downloadPath
	return c
}

func (c *ChromiumOptions) SetTmpPath(tmpPath string) *ChromiumOptions {
	c.tmpPath = tmpPath
	return c
}

func (c *ChromiumOptions) SetUser(user string) *ChromiumOptions {
	c.user = user
	return c
}

func (c *ChromiumOptions) SetLoadMode(mode string) *ChromiumOptions {
	switch strings.ToLower(mode) {
	case "", "normal", "eager", "none":
		if mode == "" {
			c.loadMode = "normal"
		} else {
			c.loadMode = strings.ToLower(mode)
		}
	}
	return c
}

func (c *ChromiumOptions) SetHeadless(headless bool) *ChromiumOptions {
	c.headless = headless
	if headless {
		return c.SetArgument("--headless", "")
	}
	return c.RemoveArgument("--headless")
}

func (c *ChromiumOptions) SetProxy(proxy string) *ChromiumOptions {
	// 不支持账号密码的配置 并且每次都需要从头启动浏览器才有效果哦
	c.proxy = proxy
	if proxy == "" {
		return c.RemoveArgument("--proxy-server")
	}
	return c.SetArgument("--proxy-server", proxy)
}

func (c *ChromiumOptions) RemoveArgument(value string) *ChromiumOptions {
	if len(c.arguments) == 0 {
		return c
	}
	filtered := c.arguments[:0]
	for _, argument := range c.arguments {
		if argument == value || strings.HasPrefix(argument, value+"=") {
			continue
		}
		filtered = append(filtered, argument)
	}
	c.arguments = filtered
	return c
}

func (c *ChromiumOptions) SetArgument(arg string, value any) *ChromiumOptions {
	if arg == "" {
		return c
	}
	c.RemoveArgument(arg)
	if value == nil {
		return c
	}

	switch v := value.(type) {
	case string:
		if v == "" {
			if arg == "--headless" {
				c.arguments = append(c.arguments, "--headless=new")
			} else {
				c.arguments = append(c.arguments, arg)
			}
			return c
		}
		c.arguments = append(c.arguments, fmt.Sprintf("%s=%s", arg, v))
	case bool:
		if v {
			c.arguments = append(c.arguments, arg)
		}
	case fmt.Stringer:
		c.arguments = append(c.arguments, fmt.Sprintf("%s=%s", arg, v.String()))
	default:
		c.arguments = append(c.arguments, fmt.Sprintf("%s=%v", arg, v))
	}
	return c
}

func (c *ChromiumOptions) AddExtension(path string) *ChromiumOptions {
	if path == "" {
		return c
	}
	if litedir.FileExists(path) {
		c.extensions = append(c.extensions, path)
	}
	return c
}

func (c *ChromiumOptions) RemoveExtension() *ChromiumOptions {
	c.extensions = []string{}
	return c
}

func (c *ChromiumOptions) SetPreference(key, value string) *ChromiumOptions {
	if c.prefs == nil {
		c.prefs = make(map[string]string)
	}
	c.prefs[key] = value
	return c
}

func (c *ChromiumOptions) SetFlag(key, value string) *ChromiumOptions {
	if c.flags == nil {
		c.flags = make(map[string]string)
	}
	c.flags[key] = value
	return c
}

func (c *ChromiumOptions) SetConfig(configPath string) error {
	if !litedir.FileExists(configPath) {
		return fmt.Errorf("file not found: %s", configPath)
	}
	c.yamlConf = configPath
	return nil
}

func (c *ChromiumOptions) ensureUserDataPath() (string, error) {
	if c.userDataPath != "" {
		if err := os.MkdirAll(c.userDataPath, 0o755); err != nil {
			return "", err
		}
		return c.userDataPath, nil
	}
	base := filepath.Join(os.TempDir(), "litepage")
	if err := os.MkdirAll(base, 0o755); err != nil {
		return "", err
	}
	dir, err := os.MkdirTemp(base, "profile-*")
	if err != nil {
		return "", err
	}
	c.userDataPath = dir
	return dir, nil
}

func (c *ChromiumOptions) validate() error {
	if c.retry < 0 {
		return errors.New("retry must be >= 0")
	}
	if c.retryInterval < 0 {
		return errors.New("retry interval must be >= 0")
	}
	return nil
}

func (c *ChromiumOptions) hasLaunchSpecificSettings() bool {
	if c == nil {
		return false
	}
	return c.proxy != "" ||
		c.headless ||
		c.browserPath != "" ||
		c.userDataPath != "" ||
		c.downloadPath != "" ||
		c.tmpPath != "" ||
		len(c.extensions) > 0 ||
		len(c.arguments) > 0
}

func normalizeDebugAddress(address string) string {
	address = strings.TrimSpace(address)
	if address == "" {
		return net.JoinHostPort(defaultHost, strconv.Itoa(defaultPort))
	}
	if strings.HasPrefix(address, "http://") || strings.HasPrefix(address, "https://") {
		u, err := url.Parse(address)
		if err == nil && u.Host != "" {
			return u.Host
		}
	}
	if strings.HasPrefix(address, "ws://") || strings.HasPrefix(address, "wss://") {
		u, err := url.Parse(address)
		if err == nil && u.Host != "" {
			return u.Host
		}
	}
	if strings.Contains(address, ":") {
		if _, _, err := net.SplitHostPort(address); err == nil {
			return address
		}
	}
	if _, err := strconv.Atoi(address); err == nil {
		return net.JoinHostPort(defaultHost, address)
	}
	return net.JoinHostPort(address, strconv.Itoa(defaultPort))
}
