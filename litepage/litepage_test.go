package litepage

import (
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"
)

func TestNormalizeDebugAddress(t *testing.T) {
	tests := map[string]string{
		"":                   "127.0.0.1:9222",
		"9223":               "127.0.0.1:9223",
		"127.0.0.1:9555":     "127.0.0.1:9555",
		"http://1.2.3.4:80":  "1.2.3.4:80",
		"ws://1.2.3.4:81/ws": "1.2.3.4:81",
	}

	for input, want := range tests {
		if got := normalizeDebugAddress(input); got != want {
			t.Fatalf("normalizeDebugAddress(%q) = %q, want %q", input, got, want)
		}
	}
}

func TestBuildLaunchArgs(t *testing.T) {
	opts := NewChromiumOptions().
		SetLocalPort(9333).
		SetProxy("http://127.0.0.1:8080").
		SetArgument("--window-size", "1280,720").
		SetHeadless(true)

	args := buildLaunchArgs(opts, "/tmp/profile")
	joined := strings.Join(args, " ")

	for _, want := range []string{
		"--remote-debugging-port=9333",
		"--user-data-dir=/tmp/profile",
		"--proxy-server=http://127.0.0.1:8080",
		"--window-size=1280,720",
		"--headless=new",
	} {
		if !strings.Contains(joined, want) {
			t.Fatalf("expected launch args to contain %q, got %q", want, joined)
		}
	}
}

func TestChromiumOptionsSetRetryDuration(t *testing.T) {
	opts := NewChromiumOptions().SetRetryDuration(3, 500*time.Millisecond)
	if opts.Retry() != 3 {
		t.Fatalf("unexpected retry: %d", opts.Retry())
	}
	if opts.RetryInterval() != 500*time.Millisecond {
		t.Fatalf("unexpected retry interval: %s", opts.RetryInterval())
	}
}

func TestElementSelector(t *testing.T) {
	page := &Page{}
	el := page.Find("#app")
	if el == nil {
		t.Fatal("expected element")
	}
	if el.Selector() != "#app" {
		t.Fatalf("unexpected selector: %q", el.Selector())
	}
}

func TestFindAllBuildsIndexedElements(t *testing.T) {
	page := &Page{}
	items := make([]*Element, 0, 3)
	for i := 0; i < 3; i++ {
		items = append(items, &Element{page: page, selector: ".item", index: i, mode: "css", rootExpr: "document"})
	}

	if items[0].Index() != 0 || items[2].Index() != 2 {
		t.Fatalf("unexpected indexes: %d %d", items[0].Index(), items[2].Index())
	}
	if items[1].jsExpr() != `document.querySelectorAll(".item")[1]` {
		t.Fatalf("unexpected js expr: %s", items[1].jsExpr())
	}
}

func TestXPathElementExpr(t *testing.T) {
	el := &Element{
		page:     &Page{},
		selector: `//div[@id="app"]`,
		mode:     "xpath",
		rootExpr: "document",
	}
	expr := el.jsExpr()
	if !strings.Contains(expr, "document.evaluate") {
		t.Fatalf("expected xpath expr, got %s", expr)
	}
}

func TestKeyCode(t *testing.T) {
	if keyCode("a") != 65 {
		t.Fatalf("unexpected keycode for a: %d", keyCode("a"))
	}
	if keyCode("Enter") != 13 {
		t.Fatalf("unexpected keycode for enter: %d", keyCode("Enter"))
	}
}

func TestCookieToParams(t *testing.T) {
	params := cookieToParams(Cookie{
		Name:   "sid",
		Value:  "abc",
		Domain: "example.com",
		Path:   "/",
		Secure: true,
	})
	if params["name"] != "sid" || params["value"] != "abc" {
		t.Fatalf("unexpected cookie params: %#v", params)
	}
	if params["domain"] != "example.com" {
		t.Fatalf("unexpected cookie domain: %#v", params)
	}
}

func TestElementUploadFilesMissing(t *testing.T) {
	el := &Element{page: &Page{}}
	err := el.UploadFiles("/not-found-file")
	if err == nil {
		t.Fatal("expected missing file error")
	}
}

func TestSaveScreenshotWritesFile(t *testing.T) {
	dir := t.TempDir()
	fp := filepath.Join(dir, "a.png")
	err := os.WriteFile(fp, []byte("x"), 0o644)
	if err != nil {
		t.Fatalf("seed file: %v", err)
	}
	data, err := os.ReadFile(fp)
	if err != nil || string(data) != "x" {
		t.Fatalf("unexpected file contents: %q err=%v", string(data), err)
	}
}

func TestCDPConnSubscribeDispatch(t *testing.T) {
	conn := &cdpConn{
		waiters: make(map[int64]chan cdpResponse),
		subs:    make(map[string][]chan json.RawMessage),
	}
	ch, cancel := conn.subscribe("Network.requestWillBeSent", 1)
	defer cancel()

	payload := json.RawMessage(`{"requestId":"1"}`)
	conn.dispatchEvent("Network.requestWillBeSent", payload)

	got, ok := <-ch
	if !ok {
		t.Fatal("expected event payload")
	}
	if string(got) != string(payload) {
		t.Fatalf("unexpected payload: %s", string(got))
	}
}

func TestFindXPathMode(t *testing.T) {
	page := &Page{}
	el := page.FindXPath("//div")
	if el.mode != "xpath" {
		t.Fatalf("unexpected mode: %s", el.mode)
	}
}

func TestInterceptActionHeaders(t *testing.T) {
	action := InterceptAction{
		URL:      "https://example.com/api",
		Method:   "POST",
		PostData: `{"a":1}`,
		Headers: map[string]string{
			"X-Test": "1",
		},
	}
	if action.URL == "" || action.Method == "" || action.Headers["X-Test"] != "1" {
		t.Fatalf("unexpected action: %#v", action)
	}
}
