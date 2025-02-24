package litereq

import (
	"bytes"
	"context"
	"encoding/base64"
	"fmt"
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
	h1         bool // 默认用h2
	ctx        context.Context
	validators []ResponseHandler
	handler    ResponseHandler
	cookieJar  *cookiejar.Jar
	timeout    time.Duration
	verbose    bool
	//tls        TlsGetter
}

// 后面还要增加 tls 指纹的处理

func Build(baseUrl string) *Builder {
	b := &Builder{
		ub: urlBuilder{},
		rb: requestBuilder{
			retry: 1, // setDefault 1
		},
		h1: false,
	}
	b.ub.BaseURL(baseUrl)
	return b
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

func (b *Builder) log(funT string, err error) {
	if b.verbose {
		fmt.Println(litestr.ColorString(funT, "red"), err)
	}
}

// Verbose 用于打印出详细流程日志的
func (b *Builder) Verbose(vb bool) *Builder {
	b.verbose = vb
	return b
}

// H1 默认就是用h2 所以这里是强制改成h1的
func (b *Builder) H1(enable bool) *Builder {
	b.h1 = enable
	return b
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

func (b *Builder) request(ctx context.Context) (req *http.Request, err error) {
	u, err := b.url()
	if err != nil {
		b.log("request", err)
		return nil, err
	}
	req, err = b.rb.Request(ctx, u)
	if err != nil {
		b.log("request", err)
		return nil, joinErrs(ErrRequest, err)
	}
	return req, nil
}

func (b *Builder) do(req *http.Request) (err error) {
	cl := Or(b.hc, &http.Client{
		Transport: createTransport(b.proxy, b.h1),
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
		code, err = do(cl, req, validators, h)
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
	//b.log("do", err)
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

func (b *Builder) Json(v any) *Builder {
	return b.Body(BodyJSON(v)).ContentType("application/json")
}

func (b *Builder) Data(v any) *Builder {
	switch v.(type) {
	case io.Reader:
		b.Body(BodyReader(v.(io.Reader)))
	case func(w io.Writer) error:
		b.Body(BodyWriter(v.(func(w io.Writer) error)))
	case []byte:
		b.Body(BodyBytes(v.([]byte)))
	case url.Values:
		b.Body(BodyForm(v.(url.Values))).ContentType("application/x-www-form-urlencoded")
	case string, map[string]any, map[string]string:
		b.Body(bodyData(v)).ContentType("application/x-www-form-urlencoded")
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

// BodyWriter pipes writes from w to the Builder's request body.
func (b *Builder) bodyWriter(f func(w io.Writer) error) *Builder {
	return b.Body(BodyWriter(f))
}

// 用于请求了之后 清理掉内存里面存的东西 目前单进程没问题 并发还不确定有没有问题 后面看看
func (b *Builder) emptyQueryFields() {
	// 清理配置了的 param
	b.ub.emptyParam()
	// 清理配置了的 body
	b.rb.emptyBody()
	// header和cookie之类的全局共享 不清理
}

func (b *Builder) GetCookies() *Cookies {
	return b.rb.GetCookies()
}

func (b *Builder) URL() (u *url.URL, err error) {
	u, err = b.ub.URL()
	if err != nil {
		return u, joinErrs(ErrURL, err)
	}
	return u, nil
}

// AddValidator adds a response validator to the Builder.
// Adding a validator disables DefaultValidator.
// To disable all validation, just add nil.
func (b *Builder) AddValidator(h ResponseHandler) *Builder {
	b.validators = append(b.validators, h)
	return b
}

// Handle sets the response handler for a Builder.
// To use multiple handlers, use ChainHandlers.
func (b *Builder) Handle(h ResponseHandler) *Builder {
	b.handler = h
	return b
}

// Config allows Builder to be extended by functions that set several options at once.
func (b *Builder) Config(cfgs ...Config) *Builder {
	for _, cfg := range cfgs {
		if cfg != nil {
			cfg(b)
		}
	}
	return b
}

// -----核心入口

func (b *Builder) Fetch(ctx context.Context) (err error) {
	//defer b.emptyQueryFields()
	req, err := b.request(ctx)
	if err != nil {
		return
	}
	return b.do(req)
}

// ---------------------------------------------

func (b *Builder) Head() *Builder {
	b.rb.Method(http.MethodHead)
	return b
}

func (b *Builder) Get() *Builder {
	b.rb.Method(http.MethodGet)
	return b
}

func (b *Builder) Post() *Builder {
	b.rb.Method(http.MethodPost)
	return b
}

func (b *Builder) Put() *Builder {
	b.rb.Method(http.MethodPut)
	return b
}

func (b *Builder) Patch() *Builder {
	b.rb.Method(http.MethodPatch)
	return b
}

func (b *Builder) Options() *Builder {
	b.rb.Method(http.MethodOptions)
	return b
}

func (b *Builder) Delete() *Builder {
	b.rb.Method(http.MethodDelete)
	return b
}

func (b *Builder) Trace() *Builder {
	b.rb.Method(http.MethodTrace)
	return b
}

func (b *Builder) Connect() *Builder {
	b.rb.Method(http.MethodConnect)
	return b
}

// -----------------------

// CheckStatus adds a validator for status code of a response.
func (b *Builder) CheckStatus(acceptStatuses ...int) *Builder {
	return b.AddValidator(CheckStatus(acceptStatuses...))
}

// CheckContentType adds a validator for the content type header of a response.
func (b *Builder) CheckContentType(cts ...string) *Builder {
	return b.AddValidator(CheckContentType(cts...))
}

// CheckPeek adds a validator that peeks at the first n bytes of a response body.
func (b *Builder) CheckPeek(n int, f func([]byte) error) *Builder {
	return b.AddValidator(CheckPeek(n, f))
}

// ToJSON sets the Builder to decode a response as a JSON object
func (b *Builder) ToJSON(v any) *Builder {
	return b.Handle(ToJSON(v))
}

// ToString sets the Builder to write the response body to the provided string pointer.
func (b *Builder) ToString(sp *string) *Builder {
	return b.Handle(ToString(sp))
}

// ToBytesBuffer sets the Builder to write the response body to the provided bytes.Buffer.
func (b *Builder) ToBytesBuffer(buf *bytes.Buffer) *Builder {
	return b.Handle(ToBytesBuffer(buf))
}

// ToWriter sets the Builder to copy the response body into w.
func (b *Builder) ToWriter(w io.Writer) *Builder {
	return b.Handle(ToWriter(w))
}

// ToFile sets the Builder to write the response body to the given file name.
// The file and its parent directories are created automatically.
// For more advanced use cases, use ToWriter.
func (b *Builder) ToFile(name string) *Builder {
	return b.Handle(ToFile(name))
}

// CopyHeaders adds a validator which copies the response headers to h.
// Note that because CopyHeaders adds a validator,
// the DefaultValidator is disabled and must be added back manually
// if status code validation is desired.
func (b *Builder) CopyHeaders(h map[string][]string) *Builder {
	return b.
		AddValidator(CopyHeaders(h))
}

// ToHeaders sets the method to HEAD and adds a handler which copies the response headers to h.
// To just copy headers, see Builder.CopyHeaders.
func (b *Builder) ToHeaders(h map[string][]string) *Builder {
	return b.
		Head().
		Handle(ChainHandlers(CopyHeaders(h), consumeBody))
}

// ErrorJSON adds a validator that applies DefaultValidator
// and decodes the response as a JSON object
// if the DefaultValidator check fails.
func (b *Builder) ErrorJSON(v any) *Builder {
	return b.AddValidator(ErrorJSON(v))
}
