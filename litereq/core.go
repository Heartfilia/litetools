package litereq

import "io"

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

type urlBuilder struct {
	baseurl                       string
	scheme, host                  string
	paths                         []string
	params, footParams, godParams []multimap
}

type requestBuilder struct {
	headers []multimap
	cookies []kvPair
	getBody BodyGetter
	method  string
}

func (rb *requestBuilder) Header(key string, values ...string) {
	rb.headers = append(rb.headers, multimap{key, values})
}

func (rb *requestBuilder) Cookie(name, value string) {
	rb.cookies = append(rb.cookies, kvPair{name, value})
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

func Clip[T any](sp *[]T) {
	s := *sp
	*sp = s[:len(s):len(s)]
}
