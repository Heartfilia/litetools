package opt

import (
	netHTTP "net/http"
	"strings"
)

// 作为请求参数的配置选项
// 先把基础的一些配置开发了 其它配置后面再优化添加

type Option struct {
	Params         string          // 先占位 后续更新
	Headers        *netHTTP.Header // 先占位 后续更新
	Cookies        *netHTTP.Cookie // 先占位 后续更新
	Data           string          // 先占位 后续更新
	Json           string          // 先占位 后续更新
	Verify         bool            // 默认true
	Files          string          // 先占位 后续更新
	Proxy          string          // 先占位 后续更新
	Method         string          // 默认GET -> 通过 option.SetMethod("POST")调整 或者 option.Method = "POST"
	Timeout        int             // ms  单位为毫秒
	AllowRedirects bool            // 是否允许重定向，默认允许
	Stream         string          // 先占位 后续更新
	Auth           string          // 先占位 后续更新
	Cert           string          // 先占位 后续更新
}

func NewOption() *Option {
	return &Option{
		Method:         "GET",
		AllowRedirects: true,
		Verify:         true,
	}
}

func (o *Option) SetRedirects(allow bool) {
	o.AllowRedirects = allow
}

func (o *Option) SetVerify(enable bool) {
	o.Verify = enable
}

func (o *Option) SetMethod(method string) {
	md := strings.ToUpper(method)
	if md == OPT {
		o.Method = "OPTIONS"
	} else if md == GET {
		o.Method = "GET"
	} else if md == HEAD {
		o.Method = "HEAD"
	} else if md == POST {
		o.Method = "POST"
	} else if md == PUT {
		o.Method = "PUT"
	} else if md == DELETE {
		o.Method = "DELETE"
	} else if md == TRACE {
		o.Method = "TRACE"
	} else if md == CONNECT {
		o.Method = "CONNECT"
	} else if md == PATCH {
		o.Method = "PATCH"
	} else {
		o.Method = "GET"
	}
}
