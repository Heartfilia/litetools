package litepage

import (
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"
)

type Page struct {
	browser *Browser
	id      string
	title   string
	url     string
	wsURL   string
	conn    *cdpConn
}

type Element struct {
	page     *Page
	selector string
	index    int
	mode     string
	rootExpr string
}

type Frame struct {
	page     *Page
	rootExpr string
	selector string
}

type RequestEvent struct {
	RequestID   string `json:"requestId"`
	LoaderID    string `json:"loaderId"`
	DocumentURL string `json:"documentURL"`
	Request     struct {
		URL      string            `json:"url"`
		Method   string            `json:"method"`
		Headers  map[string]string `json:"headers"`
		PostData string            `json:"postData,omitempty"`
	} `json:"request"`
	Type string `json:"type"`
}

type ResponseEvent struct {
	RequestID string `json:"requestId"`
	LoaderID  string `json:"loaderId"`
	Type      string `json:"type"`
	Response  struct {
		URL        string            `json:"url"`
		Status     int               `json:"status"`
		StatusText string            `json:"statusText"`
		Headers    map[string]string `json:"headers"`
		MimeType   string            `json:"mimeType"`
	} `json:"response"`
}

type RequestPattern struct {
	URLPattern   string
	ResourceType string
	RequestStage string
}

type RequestPausedEvent struct {
	RequestID string `json:"requestId"`
	Request   struct {
		URL      string            `json:"url"`
		Method   string            `json:"method"`
		Headers  map[string]string `json:"headers"`
		PostData string            `json:"postData,omitempty"`
	} `json:"request"`
	ResourceType string `json:"resourceType,omitempty"`
}

type InterceptAction struct {
	FailReason string
	URL        string
	Method     string
	PostData   string
	Headers    map[string]string
}

type runtimeEvaluateResult struct {
	Result struct {
		Type        string `json:"type"`
		Subtype     string `json:"subtype,omitempty"`
		Value       any    `json:"value"`
		Description string `json:"description,omitempty"`
	} `json:"result"`
	ExceptionDetails any `json:"exceptionDetails"`
}

type pageNavigateResult struct {
	FrameID string `json:"frameId"`
}

type captureScreenshotResult struct {
	Data string `json:"data"`
}

func newPage(browser *Browser, target targetInfo) (*Page, error) {
	if target.WebSocketDebuggerURL == "" {
		return nil, fmt.Errorf("target %s has no websocket debugger url", target.ID)
	}
	conn, err := newCDPConn(target.WebSocketDebuggerURL)
	if err != nil {
		return nil, err
	}
	p := &Page{
		browser: browser,
		id:      target.ID,
		title:   target.Title,
		url:     target.URL,
		wsURL:   target.WebSocketDebuggerURL,
		conn:    conn,
	}
	_ = p.conn.call("Page.enable", nil, nil)
	_ = p.conn.call("Runtime.enable", nil, nil)
	_ = p.conn.call("DOM.enable", nil, nil)
	return p, nil
}

func (p *Page) ID() string        { return p.id }
func (p *Page) URL() string       { return p.url }
func (p *Page) TitleHint() string { return p.title }

func (p *Page) Navigate(rawURL string) error {
	if rawURL == "" {
		return errors.New("url is empty")
	}
	if !strings.Contains(rawURL, "://") && rawURL != "about:blank" {
		rawURL = "http://" + rawURL
	}
	if err := p.conn.call("Page.navigate", map[string]any{"url": rawURL}, &pageNavigateResult{}); err != nil {
		return err
	}
	p.url = rawURL
	return p.WaitReady(10 * time.Second)
}

func (p *Page) WaitReady(timeout time.Duration) error {
	deadline := time.Now().Add(timeout)
	for {
		state, err := p.EvalString("document.readyState")
		if err == nil {
			switch p.browser.options.LoadMode() {
			case "none":
				return nil
			case "eager":
				if state == "interactive" || state == "complete" {
					return nil
				}
			default:
				if state == "complete" {
					return nil
				}
			}
		}
		if time.Now().After(deadline) {
			return fmt.Errorf("wait page ready timeout after %s", timeout)
		}
		time.Sleep(100 * time.Millisecond)
	}
}

func (p *Page) Eval(js string) (any, error) {
	var result runtimeEvaluateResult
	if err := p.conn.call("Runtime.evaluate", map[string]any{
		"expression":    js,
		"returnByValue": true,
		"awaitPromise":  true,
	}, &result); err != nil {
		return nil, err
	}
	if result.ExceptionDetails != nil {
		return nil, fmt.Errorf("javascript exception: %v", result.ExceptionDetails)
	}
	return result.Result.Value, nil
}

func (p *Page) evalBool(js string) (bool, error) {
	value, err := p.Eval(js)
	if err != nil {
		return false, err
	}
	v, ok := value.(bool)
	if !ok {
		return false, fmt.Errorf("javascript result is %T, not bool", value)
	}
	return v, nil
}

func (p *Page) EvalString(js string) (string, error) {
	value, err := p.Eval(js)
	if err != nil {
		return "", err
	}
	switch v := value.(type) {
	case string:
		return v, nil
	default:
		return fmt.Sprintf("%v", v), nil
	}
}

func (p *Page) Find(selector string) *Element {
	return &Element{page: p, selector: selector, mode: "css", rootExpr: "document"}
}

func (p *Page) MustFind(selector string) *Element {
	return p.Find(selector)
}

func (p *Page) FindXPath(xpath string) *Element {
	return &Element{page: p, selector: xpath, mode: "xpath", rootExpr: "document"}
}

func (p *Page) FindAll(selector string) ([]*Element, error) {
	count, err := p.Count(selector)
	if err != nil {
		return nil, err
	}
	items := make([]*Element, 0, count)
	for i := 0; i < count; i++ {
		items = append(items, &Element{
			page:     p,
			selector: selector,
			index:    i,
			mode:     "css",
			rootExpr: "document",
		})
	}
	return items, nil
}

func (p *Page) FindAllXPath(xpath string) ([]*Element, error) {
	count, err := p.CountXPath(xpath)
	if err != nil {
		return nil, err
	}
	items := make([]*Element, 0, count)
	for i := 0; i < count; i++ {
		items = append(items, &Element{
			page:     p,
			selector: xpath,
			index:    i,
			mode:     "xpath",
			rootExpr: "document",
		})
	}
	return items, nil
}

func (p *Page) Exists(selector string) (bool, error) {
	return p.Find(selector).Exists()
}

func (p *Page) Count(selector string) (int, error) {
	value, err := p.Eval(fmt.Sprintf("document.querySelectorAll(%s).length", strconv.Quote(selector)))
	if err != nil {
		return 0, err
	}
	switch v := value.(type) {
	case float64:
		return int(v), nil
	case int:
		return v, nil
	default:
		return 0, fmt.Errorf("javascript result is %T, not count", value)
	}
}

func (p *Page) CountXPath(xpath string) (int, error) {
	value, err := p.Eval(fmt.Sprintf(`(() => {
		const result = document.evaluate(%s, document, null, XPathResult.ORDERED_NODE_SNAPSHOT_TYPE, null);
		return result.snapshotLength;
	})()`, strconv.Quote(xpath)))
	if err != nil {
		return 0, err
	}
	switch v := value.(type) {
	case float64:
		return int(v), nil
	case int:
		return v, nil
	default:
		return 0, fmt.Errorf("javascript result is %T, not count", value)
	}
}

func (p *Page) Frame(selector string) *Frame {
	return &Frame{
		page:     p,
		selector: selector,
		rootExpr: fmt.Sprintf(`(() => {
			const frame = document.querySelector(%s);
			if (!frame || !frame.contentDocument) return null;
			return frame.contentDocument;
		})()`, strconv.Quote(selector)),
	}
}

func (p *Page) HTML() (string, error) {
	return p.EvalString("document.documentElement.outerHTML")
}

func (p *Page) Title() (string, error) {
	return p.EvalString("document.title")
}

func (p *Page) Reload() error {
	if err := p.conn.call("Page.reload", map[string]any{"ignoreCache": false}, nil); err != nil {
		return err
	}
	return p.WaitReady(10 * time.Second)
}

func (p *Page) Refresh() error {
	return p.Reload()
}

func (p *Page) Back() error {
	_, err := p.Eval(`history.back(); true`)
	if err != nil {
		return err
	}
	time.Sleep(200 * time.Millisecond)
	return p.WaitReady(10 * time.Second)
}

func (p *Page) Forward() error {
	_, err := p.Eval(`history.forward(); true`)
	if err != nil {
		return err
	}
	time.Sleep(200 * time.Millisecond)
	return p.WaitReady(10 * time.Second)
}

func (p *Page) Screenshot() ([]byte, error) {
	var result captureScreenshotResult
	if err := p.conn.call("Page.captureScreenshot", map[string]any{"format": "png"}, &result); err != nil {
		return nil, err
	}
	return base64.StdEncoding.DecodeString(result.Data)
}

func (p *Page) SaveScreenshot(path string) error {
	data, err := p.Screenshot()
	if err != nil {
		return err
	}
	return os.WriteFile(path, data, 0o644)
}

func (p *Page) RunCDP(method string, params map[string]any, out ...any) error {
	var result any
	if len(out) > 0 {
		result = out[0]
	}
	return p.conn.call(method, params, result)
}

func (p *Page) Close() error {
	if p.conn != nil {
		_ = p.conn.close()
	}
	return p.browser.ClosePage(p.id)
}

func (p *Page) MouseClick(x, y float64) error {
	if err := p.conn.call("Input.dispatchMouseEvent", map[string]any{
		"type":       "mousePressed",
		"x":          x,
		"y":          y,
		"button":     "left",
		"clickCount": 1,
	}, nil); err != nil {
		return err
	}
	return p.conn.call("Input.dispatchMouseEvent", map[string]any{
		"type":       "mouseReleased",
		"x":          x,
		"y":          y,
		"button":     "left",
		"clickCount": 1,
	}, nil)
}

func (p *Page) MouseMove(x, y float64) error {
	return p.conn.call("Input.dispatchMouseEvent", map[string]any{
		"type": "mouseMoved",
		"x":    x,
		"y":    y,
	}, nil)
}

func (p *Page) KeyDown(key string) error {
	return p.conn.call("Input.dispatchKeyEvent", map[string]any{
		"type":                  "keyDown",
		"key":                   key,
		"text":                  key,
		"unmodifiedText":        key,
		"windowsVirtualKeyCode": keyCode(key),
		"nativeVirtualKeyCode":  keyCode(key),
	}, nil)
}

func (p *Page) KeyUp(key string) error {
	return p.conn.call("Input.dispatchKeyEvent", map[string]any{
		"type":                  "keyUp",
		"key":                   key,
		"windowsVirtualKeyCode": keyCode(key),
		"nativeVirtualKeyCode":  keyCode(key),
	}, nil)
}

func (p *Page) Press(key string) error {
	if err := p.KeyDown(key); err != nil {
		return err
	}
	return p.KeyUp(key)
}

func (p *Page) EnableNetwork() error {
	return p.conn.call("Network.enable", nil, nil)
}

func (p *Page) SetExtraHTTPHeaders(headers map[string]string) error {
	if err := p.EnableNetwork(); err != nil {
		return err
	}
	return p.conn.call("Network.setExtraHTTPHeaders", map[string]any{
		"headers": headers,
	}, nil)
}

func (p *Page) SetBlockedURLs(patterns ...string) error {
	if err := p.EnableNetwork(); err != nil {
		return err
	}
	return p.conn.call("Network.setBlockedURLs", map[string]any{
		"urls": patterns,
	}, nil)
}

func (p *Page) DisableCache(disabled bool) error {
	if err := p.EnableNetwork(); err != nil {
		return err
	}
	return p.conn.call("Network.setCacheDisabled", map[string]any{
		"cacheDisabled": disabled,
	}, nil)
}

func (p *Page) ListenRequests(buffer int) (<-chan RequestEvent, func(), error) {
	if err := p.EnableNetwork(); err != nil {
		return nil, nil, err
	}
	rawCh, cancelRaw := p.conn.subscribe("Network.requestWillBeSent", buffer)
	out := make(chan RequestEvent, buffer)
	done := make(chan struct{})
	go func() {
		defer close(out)
		for {
			select {
			case <-done:
				return
			case payload, ok := <-rawCh:
				if !ok {
					return
				}
				var evt RequestEvent
				if err := json.Unmarshal(payload, &evt); err != nil {
					continue
				}
				out <- evt
			}
		}
	}()
	cancel := func() {
		close(done)
		cancelRaw()
	}
	return out, cancel, nil
}

func (p *Page) ListenResponses(buffer int) (<-chan ResponseEvent, func(), error) {
	if err := p.EnableNetwork(); err != nil {
		return nil, nil, err
	}
	rawCh, cancelRaw := p.conn.subscribe("Network.responseReceived", buffer)
	out := make(chan ResponseEvent, buffer)
	done := make(chan struct{})
	go func() {
		defer close(out)
		for {
			select {
			case <-done:
				return
			case payload, ok := <-rawCh:
				if !ok {
					return
				}
				var evt ResponseEvent
				if err := json.Unmarshal(payload, &evt); err != nil {
					continue
				}
				out <- evt
			}
		}
	}()
	cancel := func() {
		close(done)
		cancelRaw()
	}
	return out, cancel, nil
}

func (p *Page) EnableRequestInterception(patterns ...RequestPattern) error {
	fetchPatterns := make([]map[string]any, 0, len(patterns))
	for _, pattern := range patterns {
		item := map[string]any{}
		if pattern.URLPattern != "" {
			item["urlPattern"] = pattern.URLPattern
		}
		if pattern.ResourceType != "" {
			item["resourceType"] = pattern.ResourceType
		}
		if pattern.RequestStage != "" {
			item["requestStage"] = pattern.RequestStage
		}
		fetchPatterns = append(fetchPatterns, item)
	}
	params := map[string]any{}
	if len(fetchPatterns) > 0 {
		params["patterns"] = fetchPatterns
	}
	return p.conn.call("Fetch.enable", params, nil)
}

func (p *Page) DisableRequestInterception() error {
	return p.conn.call("Fetch.disable", nil, nil)
}

func (p *Page) InterceptRequests(buffer int, handler func(event RequestPausedEvent) InterceptAction) (func(), error) {
	if handler == nil {
		return nil, errors.New("intercept handler is nil")
	}
	if err := p.EnableRequestInterception(); err != nil {
		return nil, err
	}
	rawCh, cancelRaw := p.conn.subscribe("Fetch.requestPaused", buffer)
	done := make(chan struct{})
	go func() {
		for {
			select {
			case <-done:
				return
			case payload, ok := <-rawCh:
				if !ok {
					return
				}
				var evt RequestPausedEvent
				if err := json.Unmarshal(payload, &evt); err != nil {
					continue
				}
				action := handler(evt)
				if action.FailReason != "" {
					_ = p.conn.call("Fetch.failRequest", map[string]any{
						"requestId":   evt.RequestID,
						"errorReason": action.FailReason,
					}, nil)
					continue
				}

				params := map[string]any{
					"requestId": evt.RequestID,
				}
				if action.URL != "" {
					params["url"] = action.URL
				}
				if action.Method != "" {
					params["method"] = action.Method
				}
				if action.PostData != "" {
					params["postData"] = action.PostData
				}
				if len(action.Headers) > 0 {
					headers := make([]map[string]string, 0, len(action.Headers))
					for k, v := range action.Headers {
						headers = append(headers, map[string]string{
							"name":  k,
							"value": v,
						})
					}
					params["headers"] = headers
				}
				_ = p.conn.call("Fetch.continueRequest", params, nil)
			}
		}
	}()
	cancel := func() {
		close(done)
		cancelRaw()
		_ = p.DisableRequestInterception()
	}
	return cancel, nil
}

func (p *Page) Cookies() ([]Cookie, error) {
	if err := p.EnableNetwork(); err != nil {
		return nil, err
	}
	params := map[string]any{}
	if p.url != "" {
		params["urls"] = []string{p.url}
	}
	var result struct {
		Cookies []Cookie `json:"cookies"`
	}
	if err := p.conn.call("Network.getCookies", params, &result); err != nil {
		return nil, err
	}
	return result.Cookies, nil
}

func (p *Page) SetCookies(cookies ...Cookie) error {
	if err := p.EnableNetwork(); err != nil {
		return err
	}
	for i := range cookies {
		if cookies[i].URL == "" {
			cookies[i].URL = p.url
		}
		var result struct {
			Success bool `json:"success"`
		}
		if err := p.conn.call("Network.setCookie", cookieToParams(cookies[i]), &result); err != nil {
			return err
		}
		if !result.Success {
			return fmt.Errorf("set cookie %q failed", cookies[i].Name)
		}
	}
	return nil
}

func (p *Page) ClearCookies() error {
	if err := p.EnableNetwork(); err != nil {
		return err
	}
	return p.conn.call("Network.clearBrowserCookies", nil, nil)
}

func (p *Page) LocalStorageItem(key string) (string, error) {
	return p.EvalString(fmt.Sprintf(`window.localStorage.getItem(%s) ?? ""`, strconv.Quote(key)))
}

func (p *Page) SetLocalStorageItem(key, value string) error {
	_, err := p.Eval(fmt.Sprintf(`window.localStorage.setItem(%s, %s); true`, strconv.Quote(key), strconv.Quote(value)))
	return err
}

func (p *Page) RemoveLocalStorageItem(key string) error {
	_, err := p.Eval(fmt.Sprintf(`window.localStorage.removeItem(%s); true`, strconv.Quote(key)))
	return err
}

func (p *Page) SessionStorageItem(key string) (string, error) {
	return p.EvalString(fmt.Sprintf(`window.sessionStorage.getItem(%s) ?? ""`, strconv.Quote(key)))
}

func (p *Page) SetSessionStorageItem(key, value string) error {
	_, err := p.Eval(fmt.Sprintf(`window.sessionStorage.setItem(%s, %s); true`, strconv.Quote(key), strconv.Quote(value)))
	return err
}

func (p *Page) RemoveSessionStorageItem(key string) error {
	_, err := p.Eval(fmt.Sprintf(`window.sessionStorage.removeItem(%s); true`, strconv.Quote(key)))
	return err
}

func (p *Page) WaitURLContains(substr string, timeout time.Duration) error {
	deadline := time.Now().Add(timeout)
	for {
		current, err := p.EvalString("window.location.href")
		if err == nil {
			p.url = current
			if strings.Contains(current, substr) {
				return nil
			}
		}
		if time.Now().After(deadline) {
			return fmt.Errorf("wait url containing %q timeout after %s", substr, timeout)
		}
		time.Sleep(100 * time.Millisecond)
	}
}

func (p *Page) WaitTitleContains(substr string, timeout time.Duration) error {
	deadline := time.Now().Add(timeout)
	for {
		title, err := p.Title()
		if err == nil && strings.Contains(title, substr) {
			return nil
		}
		if time.Now().After(deadline) {
			return fmt.Errorf("wait title containing %q timeout after %s", substr, timeout)
		}
		time.Sleep(100 * time.Millisecond)
	}
}

func (e *Element) Selector() string {
	return e.selector
}

func (e *Element) Index() int {
	return e.index
}

func (e *Element) jsExpr() string {
	root := e.rootExpr
	if root == "" {
		root = "document"
	}
	if e.mode == "xpath" {
		return fmt.Sprintf(`(() => {
			const root = %s;
			if (!root) return null;
			const result = document.evaluate(%s, root, null, XPathResult.ORDERED_NODE_SNAPSHOT_TYPE, null);
			return result.snapshotItem(%d);
		})()`, root, strconv.Quote(e.selector), e.index)
	}
	query := fmt.Sprintf("%s.querySelectorAll(%s)", root, strconv.Quote(e.selector))
	if e.index <= 0 {
		return query + "[0]"
	}
	return fmt.Sprintf("%s[%d]", query, e.index)
}

func (e *Element) Exists() (bool, error) {
	return e.page.evalBool(fmt.Sprintf("%s !== undefined && %s !== null", e.jsExpr(), e.jsExpr()))
}

func (e *Element) Wait(timeout time.Duration) error {
	deadline := time.Now().Add(timeout)
	for {
		ok, err := e.Exists()
		if err == nil && ok {
			return nil
		}
		if time.Now().After(deadline) {
			return fmt.Errorf("wait element %q timeout after %s", e.selector, timeout)
		}
		time.Sleep(100 * time.Millisecond)
	}
}

func (e *Element) IsVisible() (bool, error) {
	return e.page.evalBool(fmt.Sprintf(`(() => {
		const el = %s;
		if (!el) return false;
		const style = window.getComputedStyle(el);
		const rect = el.getBoundingClientRect();
		return style.display !== "none" &&
			style.visibility !== "hidden" &&
			style.opacity !== "0" &&
			rect.width > 0 &&
			rect.height > 0;
	})()`, e.jsExpr()))
}

func (e *Element) WaitVisible(timeout time.Duration) error {
	deadline := time.Now().Add(timeout)
	for {
		ok, err := e.IsVisible()
		if err == nil && ok {
			return nil
		}
		if time.Now().After(deadline) {
			return fmt.Errorf("wait element %q visible timeout after %s", e.selector, timeout)
		}
		time.Sleep(100 * time.Millisecond)
	}
}

func (e *Element) IsClickable() (bool, error) {
	return e.page.evalBool(fmt.Sprintf(`(() => {
		const el = %s;
		if (!el) return false;
		const style = window.getComputedStyle(el);
		const rect = el.getBoundingClientRect();
		const disabled = "disabled" in el && !!el.disabled;
		return !disabled &&
			style.display !== "none" &&
			style.visibility !== "hidden" &&
			style.pointerEvents !== "none" &&
			rect.width > 0 &&
			rect.height > 0;
	})()`, e.jsExpr()))
}

func (e *Element) WaitClickable(timeout time.Duration) error {
	deadline := time.Now().Add(timeout)
	for {
		ok, err := e.IsClickable()
		if err == nil && ok {
			return nil
		}
		if time.Now().After(deadline) {
			return fmt.Errorf("wait element %q clickable timeout after %s", e.selector, timeout)
		}
		time.Sleep(100 * time.Millisecond)
	}
}

func (e *Element) Text() (string, error) {
	return e.page.EvalString(fmt.Sprintf(`(() => {
		const el = %s;
		if (!el) throw new Error("element not found");
		return el.innerText ?? el.textContent ?? "";
	})()`, e.jsExpr()))
}

func (e *Element) HTML() (string, error) {
	return e.page.EvalString(fmt.Sprintf(`(() => {
		const el = %s;
		if (!el) throw new Error("element not found");
		return el.outerHTML;
	})()`, e.jsExpr()))
}

func (e *Element) Attribute(name string) (string, error) {
	return e.page.EvalString(fmt.Sprintf(`(() => {
		const el = %s;
		if (!el) throw new Error("element not found");
		const value = el.getAttribute(%s);
		return value === null ? "" : value;
	})()`, e.jsExpr(), strconv.Quote(name)))
}

func (e *Element) Rect() (map[string]float64, error) {
	value, err := e.page.Eval(fmt.Sprintf(`(() => {
		const el = %s;
		if (!el) throw new Error("element not found");
		const r = el.getBoundingClientRect();
		return {x: r.x, y: r.y, width: r.width, height: r.height};
	})()`, e.jsExpr()))
	if err != nil {
		return nil, err
	}
	obj, ok := value.(map[string]any)
	if !ok {
		return nil, fmt.Errorf("javascript result is %T, not rect", value)
	}
	rect := make(map[string]float64, 4)
	for k, v := range obj {
		if f, ok := v.(float64); ok {
			rect[k] = f
		}
	}
	return rect, nil
}

func (e *Element) ScrollIntoView() error {
	_, err := e.page.Eval(fmt.Sprintf(`(() => {
		const el = %s;
		if (!el) throw new Error("element not found");
		el.scrollIntoView({behavior: "instant", block: "center", inline: "center"});
		return true;
	})()`, e.jsExpr()))
	return err
}

func (e *Element) Hover() error {
	if err := e.ScrollIntoView(); err != nil {
		return err
	}
	rect, err := e.Rect()
	if err != nil {
		return err
	}
	return e.page.conn.call("Input.dispatchMouseEvent", map[string]any{
		"type": "mouseMoved",
		"x":    rect["x"] + rect["width"]/2,
		"y":    rect["y"] + rect["height"]/2,
	}, nil)
}

func (e *Element) Click() error {
	if err := e.ScrollIntoView(); err != nil {
		return err
	}
	rect, err := e.Rect()
	if err != nil {
		return err
	}
	return e.page.MouseClick(rect["x"]+rect["width"]/2, rect["y"]+rect["height"]/2)
}

func (e *Element) Input(value string) error {
	_, err := e.page.Eval(fmt.Sprintf(`(() => {
		const el = %s;
		if (!el) throw new Error("element not found");
		el.scrollIntoView({behavior: "instant", block: "center", inline: "center"});
		el.focus();
		el.value = %s;
		el.dispatchEvent(new Event("input", {bubbles: true}));
		el.dispatchEvent(new Event("change", {bubbles: true}));
		return true;
	})()`, e.jsExpr(), strconv.Quote(value)))
	return err
}

func (e *Element) Clear() error {
	return e.Input("")
}

func (e *Element) Type(value string) error {
	if _, err := e.page.Eval(fmt.Sprintf(`(() => {
		const el = %s;
		if (!el) throw new Error("element not found");
		el.scrollIntoView({behavior: "instant", block: "center", inline: "center"});
		el.focus();
		return true;
	})()`, e.jsExpr())); err != nil {
		return err
	}
	return e.page.conn.call("Input.insertText", map[string]any{"text": value}, nil)
}

func (e *Element) UploadFiles(paths ...string) error {
	if len(paths) == 0 {
		return nil
	}
	for _, path := range paths {
		if _, err := os.Stat(path); err != nil {
			return err
		}
	}

	var result struct {
		Result struct {
			ObjectID string `json:"objectId"`
		} `json:"result"`
	}
	if err := e.page.conn.call("Runtime.evaluate", map[string]any{
		"expression": fmt.Sprintf(`(() => {
			const el = %s;
			if (!el) throw new Error("element not found");
			return el;
		})()`, e.jsExpr()),
	}, &result); err != nil {
		return err
	}
	if result.Result.ObjectID == "" {
		return errors.New("failed to resolve element object id for upload")
	}

	var domNode struct {
		NodeID int `json:"nodeId"`
	}
	if err := e.page.conn.call("DOM.requestNode", map[string]any{
		"objectId": result.Result.ObjectID,
	}, &domNode); err != nil {
		return err
	}
	if domNode.NodeID == 0 {
		return errors.New("failed to resolve node id for upload")
	}
	return e.page.conn.call("DOM.setFileInputFiles", map[string]any{
		"nodeId": domNode.NodeID,
		"files":  paths,
	}, nil)
}

func (e *Element) ContentFrame() *Frame {
	return &Frame{
		page:     e.page,
		selector: e.selector,
		rootExpr: fmt.Sprintf(`(() => {
			const frame = %s;
			if (!frame || !frame.contentDocument) return null;
			return frame.contentDocument;
		})()`, e.jsExpr()),
	}
}

func (f *Frame) Exists() (bool, error) {
	return f.page.evalBool(fmt.Sprintf("%s !== null", f.rootExpr))
}

func (f *Frame) Find(selector string) *Element {
	return &Element{page: f.page, selector: selector, mode: "css", rootExpr: f.rootExpr}
}

func (f *Frame) FindXPath(xpath string) *Element {
	return &Element{page: f.page, selector: xpath, mode: "xpath", rootExpr: f.rootExpr}
}

func keyCode(key string) int {
	if len(key) == 1 {
		return int(strings.ToUpper(key)[0])
	}
	switch strings.ToLower(key) {
	case "enter":
		return 13
	case "tab":
		return 9
	case "escape":
		return 27
	case "backspace":
		return 8
	case "space":
		return 32
	case "arrowleft":
		return 37
	case "arrowup":
		return 38
	case "arrowright":
		return 39
	case "arrowdown":
		return 40
	default:
		return 0
	}
}
