package opt

import (
	netHTTP "net/http"
	"strings"
)

// 作为请求参数的配置选项
// 先把基础的一些配置开发了 其它配置后面再优化添加

type Option struct {
	params         string          // 先占位 后续更新
	headers        *netHTTP.Header // 先占位 后续更新
	cookies        *netHTTP.Cookie // 先占位 后续更新
	data           string          // 先占位 后续更新
	json           string          // 先占位 后续更新
	verify         bool            // 默认true
	files          string          // 先占位 后续更新
	proxy          string          // 先占位 后续更新
	method         string          // 默认GET -> 通过 option.SetMethod("POST")调整
	timeout        int             // ms  单位为毫秒
	allowRedirects bool            // 是否允许重定向，默认允许
	stream         string          // 先占位 后续更新
	auth           string          // 先占位 后续更新
	cert           string          // 先占位 后续更新
}

func NewOption() *Option {
	return &Option{
		method:         "GET",
		allowRedirects: true,
		verify:         true,
	}
}

func (o *Option) SetRedirects(allow bool) *Option {
	o.allowRedirects = allow
	return o
}

func (o *Option) SetVerify(enable bool) *Option {
	o.verify = enable
	return o
}

func (o *Option) SetMethod(method string) *Option {
	md := strings.ToUpper(method)
	switch md {
	case OPT:
		o.method = "OPTIONS"
	case GET:
		o.method = "GET"
	case HEAD:
		o.method = "HEAD"
	case POST:
		o.method = "POST"
	case PUT:
		o.method = "PUT"
	case DELETE:
		o.method = "DELETE"
	case TRACE:
		o.method = "TRACE"
	case CONNECT:
		o.method = "CONNECT"
	case PATCH:
		o.method = "PATCH"
	default:
		o.method = "GET"
	}

	return o
}
