package litereq

import (
	"context"
	"encoding/base64"
	"fmt"
	"github.com/Heartfilia/litetools/litereq/utils"
	"github.com/Heartfilia/litetools/litestr"
	"io"
	"log"
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
	ctx        context.Context
	validators []ResponseHandler
	handler    ResponseHandler
	cookieJar  *cookiejar.Jar
	timeout    time.Duration
}

func Build(ctx ...context.Context) *Builder {
	build := &Builder{
		ub: urlBuilder{},
		rb: requestBuilder{
			retry: 1, // setDefault 1
		},
		ctx: If(Or(ctx...) == nil, context.Background(), Or(ctx...)),
	}
	return build
}

func joinErrs(a, b error) error {
	return fmt.Errorf("%w: %w", a, b)
}

// URL builds a *url.URL from the base URL and options set on the Builder.
// If a valid url.URL cannot be built,
// URL() nevertheless returns a new url.URL,
// so it is always safe to call u.String().
func (b *Builder) url() (u *url.URL, err error) {
	u, err = b.ub.URL()
	if err != nil {
		return u, joinErrs(ErrURL, err)
	}
	return u, nil
}

func (b *Builder) Param(key string, values ...string) *Builder {
	if len(values) == 0 {
		b.ub.Param(key, "")
	} else {
		b.ub.Param(key, values...)
	}
	return b
}

func (b *Builder) Params(paramString string) *Builder {
	params := litestr.ParamStringToArray(paramString)
	for _, ps := range params {
		b.ub.Param(ps[0], ps[1])
	}
	return b
}

func (b *Builder) Path(p string) *Builder {
	b.ub.Path(p)
	return b
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
	u, err := b.url()
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
	b.rb.Header(key, values...)
	return b
}

func (b *Builder) Headers(headerMap map[string]string) *Builder {
	for k, v := range headerMap {
		b.rb.Header(k, v)
	}
	return b
}

// Cookie Single cookie value set
func (b *Builder) Cookie(name, value string) *Builder {
	b.rb.Cookie(name, value)
	return b
}

// Cookies Use string cookies
func (b *Builder) Cookies(s string) *Builder {
	b.rb.Header("Cookie", s)
	return b
}

// CookieJar adds a cookieJar to a request.
func (b *Builder) CookieJar(jar *cookiejar.Jar) *Builder {
	b.cookieJar = jar
	return b
}

// BasicAuth sets the Authorization header to a basic auth credential.
func (b *Builder) BasicAuth(username, password string) *Builder {
	auth := username + ":" + password
	v := base64.StdEncoding.EncodeToString([]byte(auth))
	return b.Header("Authorization", "Basic "+v)
}

// Bearer sets the Authorization header to a bearer token.
func (b *Builder) Bearer(token string) *Builder {
	return b.Header("Authorization", "Bearer "+token)
}

// Body 预留给用框架自带的自定义方案格式的
func (b *Builder) Body(src BodyGetter) *Builder {
	b.rb.Body(src)
	return b
}

func (b *Builder) BodyWriter(f func(w io.Writer) error) *Builder {
	b.rb.Body(BodyWriter(f))
	return b
}

func (b *Builder) Json(v any) *Builder {
	return b.Body(BodyJSON(v)).ContentType("application/json")
}

func (b *Builder) Data(v any) *Builder {
	switch v.(type) {
	case io.Reader:
		b.Body(BodyReader(v.(io.Reader)))
	case func(w io.Writer) error:
		b.BodyWriter(v.(func(w io.Writer) error))
	case []byte:
		b.Body(BodyBytes(v.([]byte)))
	case url.Values:
		b.Body(BodyForm(v.(url.Values))).ContentType("application/x-www-form-urlencoded")
	default:
		log.Panicln("wrong body type:", v)
	}
	return b
}

func (b *Builder) File(fp string) *Builder {
	return b.Body(BodyFile(fp))
}

func (b *Builder) ContentType(ct string) *Builder {
	b.rb.Header("Content-Type", ct)
	return b
}

func (b *Builder) UserAgent(ua string) *Builder {
	b.rb.Header("User-Agent", ua)
	return b
}

func (b *Builder) Referer(rf string) *Builder {
	b.rb.Header("Referer", rf)
	return b
}

func (b *Builder) Retry(r int) *Builder {
	b.rb.Retry(r)
	return b
}

func (b *Builder) Timeout(d time.Duration) *Builder {
	b.timeout = d
	return b
}

// -----核心入口

func (b *Builder) fetch(sourceUrl string) *Response {
	u := utils.ParseUrl(sourceUrl)
	b.ub.BaseURL(sourceUrl)
	b.ub.Scheme(u.Scheme)
	b.ub.Host(u.Host)

	resp := &Response{}
	req, err := b.request(b.ctx)
	if err != nil {
		resp.error(err)
		return resp
	}
	err = b.do(req, resp)
	resp.error(err)
	return resp
}

// ---------------------------------------------

func (b *Builder) Head(sourceUrl string) *Response {
	b.rb.Method(http.MethodHead)
	return b.fetch(sourceUrl)
}

func (b *Builder) Get(sourceUrl string) *Response {
	b.rb.Method(http.MethodGet)
	return b.fetch(sourceUrl)
}

func (b *Builder) Post(sourceUrl string) *Response {
	b.rb.Method(http.MethodPost)
	return b.fetch(sourceUrl)
}

func (b *Builder) Put(sourceUrl string) *Response {
	b.rb.Method(http.MethodPut)
	return b.fetch(sourceUrl)
}

func (b *Builder) Patch(sourceUrl string) *Response {
	b.rb.Method(http.MethodPatch)
	return b.fetch(sourceUrl)
}

func (b *Builder) Options(sourceUrl string) *Response {
	b.rb.Method(http.MethodOptions)
	return b.fetch(sourceUrl)
}

func (b *Builder) Delete(sourceUrl string) *Response {
	b.rb.Method(http.MethodDelete)
	return b.fetch(sourceUrl)
}
