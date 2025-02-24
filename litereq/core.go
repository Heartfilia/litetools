package litereq

import (
	"compress/flate"
	"compress/gzip"
	"context"
	"github.com/andybalholm/brotli"
	"io"
	"net/http"
	"net/url"
)

// 借鉴 https://github.com/earthboundkid/requests

type nopCloser struct {
	io.Reader
}

func rc(r io.Reader) nopCloser {
	return nopCloser{r}
}

func (nopCloser) Close() error { return nil }

type multimap struct {
	key    string
	values []string
}

type kvPair struct {
	key, value string
}

type requestBuilder struct {
	headers []multimap
	cookies []kvPair
	tls     any // 先放着 后面需要搞指纹
	getBody BodyGetter
	retry   int
	method  string
}

func (rb *requestBuilder) Header(key string, values ...string) {
	rb.headers = append(rb.headers, multimap{key, values})
}

func (rb *requestBuilder) Cookie(name, value string) {
	rb.cookies = append(rb.cookies, kvPair{name, value})
}

func (rb *requestBuilder) GetCookies() *Cookies {
	var cookies []*http.Cookie
	for _, kv := range rb.cookies {
		cookies = append(cookies, &http.Cookie{
			Name:  kv.key,
			Value: kv.value,
		})
	}
	return &Cookies{jar: cookies}
}

func (rb *requestBuilder) Method(method string) {
	rb.method = method
}

func (rb *requestBuilder) Body(src BodyGetter) {
	rb.getBody = src
}

// Clone creates a new Builder suitable for independent mutation.
func (rb *requestBuilder) Clone() *requestBuilder {
	rb2 := *rb
	Clip(&rb2.headers)
	Clip(&rb2.cookies)
	return &rb2
}

func (rb *requestBuilder) Retry(r int) {
	if r > 1 {
		rb.retry = r
	}
}

func do(cl *http.Client, req *http.Request, validators []ResponseHandler, h ResponseHandler) (doResponse, error) {
	res, err := cl.Do(req)
	if err != nil {
		return doConnect, err
	}
	defer func(Body io.ReadCloser) {
		err = Body.Close()
		if err != nil {
		}
	}(res.Body)

	for _, v := range validators {
		if v == nil {
			continue
		}
		if err = v(res); err != nil {
			return doValidate, err
		}
	}

	err = switchContentEncoding(res)

	if err = h(res); err != nil {
		return doHandle, err
	}

	return doOK, nil
}

func Clip[T any](sp *[]T) {
	s := *sp
	*sp = s[:len(s):len(s)]
}

func switchContentEncoding(res *http.Response) (err error) {
	var bodyReader io.ReadCloser
	switch res.Header.Get("Content-Encoding") {
	case "gzip":
		bodyReader, err = gzip.NewReader(res.Body)
	case "deflate":
		bodyReader = flate.NewReader(res.Body)
	case "br":
		bodyReader = io.NopCloser(brotli.NewReader(res.Body))
	default:
		bodyReader = res.Body
	}
	res.Body = bodyReader
	return
}

// ------------------ rb ------------------------

// Request builds a new http.Request with its context set.
func (rb *requestBuilder) Request(ctx context.Context, u *url.URL) (req *http.Request, err error) {
	var body io.Reader
	if rb.getBody != nil {
		if body, err = rb.getBody(); err != nil {
			return nil, err
		}
		if nopPer, ok := body.(nopCloser); ok {
			body = nopPer.Reader
		}
	}
	method := Or(rb.method,
		If(rb.getBody == nil, "GET", "POST"))

	req, err = http.NewRequestWithContext(ctx, method, u.String(), body)
	if err != nil {
		return nil, err
	}
	req.GetBody = rb.getBody

	for _, kv := range rb.headers {
		req.Header[http.CanonicalHeaderKey(kv.key)] = kv.values
	}
	for _, kv := range rb.cookies {
		req.AddCookie(&http.Cookie{
			Name:  kv.key,
			Value: kv.value,
		})
	}
	return req, nil
}

func (rb *requestBuilder) emptyBody() {
	rb.getBody = nil
}
