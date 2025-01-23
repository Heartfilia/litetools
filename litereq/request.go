package litereq

import (
	"context"
	"fmt"
	"github.com/Heartfilia/litetools/litereq/utils"
	"io"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"time"
)

// 借鉴 https://github.com/earthboundkid/requests

type Builder struct {
	hc         *http.Client
	rb         requestBuilder
	ub         urlBuilder
	proxy      ProxyGetter
	validators []ResponseHandler
	handler    ResponseHandler
	cookieJar  *cookiejar.Jar
	timeout    time.Duration
}

func Build(urlPath string) *Builder {
	u := utils.ParseUrl(urlPath)
	build := &Builder{
		ub: urlBuilder{},
		rb: requestBuilder{
			retry: 1, // setDefault 1
		},
	}
	build.ub.BaseURL(urlPath)
	build.ub.Scheme(u.Scheme)
	build.ub.Host(u.Host)
	build.ub.Path(u.Path)
	return build
}

func joinErrs(a, b error) error {
	return fmt.Errorf("%w: %w", a, b)
}

// URL builds a *url.URL from the base URL and options set on the Builder.
// If a valid url.URL cannot be built,
// URL() nevertheless returns a new url.URL,
// so it is always safe to call u.String().
func (b *Builder) URL() (u *url.URL, err error) {
	u, err = b.ub.URL()
	if err != nil {
		return u, joinErrs(ErrURL, err)
	}
	return u, nil
}

func (b *Builder) BodyMd5() []byte {
	if b.rb.getBody == nil {
		return nil
	}
	bd, err := b.rb.getBody()
	if err != nil {
		return nil
	}
	bt, err := io.ReadAll(bd)
	defer func(bd io.ReadCloser) {
		err := bd.Close()
		if err != nil {

		}
	}(bd)
	if err != nil {
		return nil
	}
	return utils.GetMd5(bt)
}

func (b *Builder) BodyString() string {
	if b.rb.getBody == nil {
		return ""
	}
	bd, err := b.rb.getBody()
	if err != nil {
		return ""
	}
	bt, err := io.ReadAll(bd)
	defer func(bd io.ReadCloser) {
		err := bd.Close()
		if err != nil {

		}
	}(bd)
	return string(bt)
}

func (b *Builder) request(ctx context.Context) (req *http.Request, err error) {
	u, err := b.URL()
	if err != nil {
		return nil, err
	}
	req, err = b.rb.Request(ctx, u)
	if err != nil {
		return nil, joinErrs(ErrRequest, err)
	}
	return req, nil
}

func (b *Builder) do(req *http.Request, resp *Response) (err error) {
	cl := Or(b.hc, &http.Client{
		Transport: createTransport(b.proxy),
		Timeout:   Or(b.timeout, DefaultTimeout),
	})
	if b.cookieJar != nil {
		cl.Jar = b.cookieJar
	}
	validators := b.validators
	if len(validators) == 0 {
		validators = []ResponseHandler{DefaultValidator}
	}
	h := If(b.handler != nil,
		b.handler,
		consumeBody)
	var code doResponse

	for i := 0; i < b.rb.retry; i++ {
		code, err = do(cl, req, validators, h, resp)
		if code == doOK {
			break
		}
	}

	switch code {
	case doOK:
		return nil
	case doConnect:
		err = joinErrs(ErrTransport, err)
	case doValidate:
		err = joinErrs(ErrValidator, err)
	case doHandle:
		err = joinErrs(ErrHandler, err)
	}
	return err
}

func (b *Builder) Proxy(p string) *Builder {
	b.proxy = func() Proxy {
		return &ProxyInfo{
			ProxyIp: p,
		}
	}
	return b
}

func (b *Builder) ProxyFunc(src ProxyGetter) *Builder {
	b.proxy = src
	return b
}

func (b *Builder) Header(key string, values ...string) *Builder {

	return b
}

func (b *Builder) Body(src BodyGetter) *Builder {
	b.rb.Body(src)
	return b
}

func (b *Builder) BodyWriter(f func(w io.Writer) error) *Builder {
	b.Body(BodyWriter(f))
	return b
}

func (b *Builder) ContentType(ct string) *Builder {
	return b.Header("Content-Type", ct)
}

func (b *Builder) Retry(r int) *Builder {
	b.rb.Retry(r)
	return b
}

// ---------------------------------------------

func (b *Builder) Fetch(ctx ...context.Context) *Response {
	resp := &Response{}
	req, err := b.request(If(Or(ctx...) == nil, context.Background(), Or(ctx...)))
	if err != nil {
		resp.error(err)
		return resp
	}
	err = b.do(req, resp)
	resp.error(err)
	return resp
}
