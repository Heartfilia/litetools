package litereq

import (
	"bufio"
	"bytes"
	"encoding/json"
	"golang.org/x/net/html"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

// ResponseHandler is used to validate or handle the response to a request.
type ResponseHandler = func(*http.Response) error

// Or returns the first non-empty argument it receives
// or the zero value for T.
func Or[T comparable](vals ...T) T {
	for _, val := range vals {
		if val != *new(T) {
			return val
		}
	}
	return *new(T)
}

func If[T any](val bool, a, b T) T {
	if val {
		return a
	}
	return b
}

func consumeBody(res *http.Response) (err error) {
	const maxDiscardSize = 640 * 1 << 10
	if _, err = io.CopyN(io.Discard, res.Body, maxDiscardSize); err == io.EOF {
		err = nil
	}
	return err
}

// ToWrap 下面俩是自定义的返回格式 实际用不到 忽略
func ToWrap[T any](w *Wrap[T]) ResponseHandler {
	return func(res *http.Response) error {
		w.Response = res
		w.Request = res.Request
		data, err := io.ReadAll(res.Body)
		if err != nil {
			w.Err = err
			return err
		}
		w.RespStr = string(data)
		w.Curl = toCurl(res.Request)
		if w.RespStr == "" {
			return nil
		}
		if err = json.Unmarshal(data, &w.Data); err != nil {
			w.Err = err
			return err
		}
		return nil
	}
}

func toCurl(req *http.Request) string {
	if req == nil {
		return ""
	}
	// Start building the cURL command
	curlCmd := "curl"

	// Set the request method
	curlCmd += " -X " + req.Method

	// Add headers
	for key, values := range req.Header {
		for _, value := range values {
			curlCmd += " -H \"" + key + ": " + value + "\""
		}
	}

	// Add cookies
	if len(req.Cookies()) > 0 {
		var cookieHeader string
		for _, cookie := range req.Cookies() {
			cookieHeader += cookie.Name + "=" + cookie.Value + "; "
		}
		curlCmd += " -H \"Cookie: " + strings.TrimRight(cookieHeader, "; ") + "\""
	}

	// Set the URL
	curlCmd += " \"" + req.URL.String() + "\""

	return curlCmd
}

// ChainHandlers allows for the composing of validators or response handlers.
func ChainHandlers(handlers ...ResponseHandler) ResponseHandler {
	return func(r *http.Response) error {
		for _, h := range handlers {
			if h == nil {
				continue
			}
			if err := h(r); err != nil {
				return err
			}
		}
		return nil
	}
}

// ToJSON decodes a response as a JSON object.
func ToJSON(v any) ResponseHandler {
	return func(res *http.Response) error {
		data, err := io.ReadAll(res.Body)
		if err != nil {
			return err
		}
		if err = json.Unmarshal(data, v); err != nil {
			return err
		}
		return nil
	}
}

// ToString writes the response body to the provided string pointer.
func ToString(sp *string) ResponseHandler {
	return func(res *http.Response) error {
		var buf strings.Builder
		_, err := io.Copy(&buf, res.Body)
		if err == nil {
			*sp = buf.String()
		}
		return err
	}
}

// ToBytesBuffer writes the response body to the provided bytes.Buffer.
func ToBytesBuffer(buf *bytes.Buffer) ResponseHandler {
	return func(res *http.Response) error {
		_, err := io.Copy(buf, res.Body)
		return err
	}
}

// ToBufioReader takes a callback which wraps the response body in a bufio.Reader.
func ToBufioReader(f func(r *bufio.Reader) error) ResponseHandler {
	return func(res *http.Response) error {
		return f(bufio.NewReader(res.Body))
	}
}

// ToBufioScanner takes a callback which wraps the response body in a bufio.Scanner.
func ToBufioScanner(f func(r *bufio.Scanner) error) ResponseHandler {
	return func(res *http.Response) error {
		return f(bufio.NewScanner(res.Body))
	}
}

// ToHTML parses the page with x/net/html.Parse.
func ToHTML(n *html.Node) ResponseHandler {
	return ToBufioReader(func(r *bufio.Reader) error {
		n2, err := html.Parse(r)
		if err != nil {
			return err
		}
		*n = *n2
		return nil
	})
}

// ToWriter copies the response body to w.
func ToWriter(w io.Writer) ResponseHandler {
	return ToBufioReader(func(r *bufio.Reader) error {
		_, err := io.Copy(w, r)

		return err
	})
}

// ToFile writes the response body at the provided file path.
// The file and its parent directories are created automatically.
func ToFile(name string) ResponseHandler {
	return func(res *http.Response) error {
		_ = os.MkdirAll(filepath.Dir(name), 0777)

		f, err := os.Create(name)
		if err != nil {
			return err
		}
		defer func(f *os.File) {
			err := f.Close()
			if err != nil {

			}
		}(f)

		_, err = io.Copy(f, res.Body)
		return err
	}
}

// ToHeaders is an alias for backwards compatibility.
var ToHeaders = CopyHeaders
