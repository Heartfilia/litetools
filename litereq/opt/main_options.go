package opt

import (
	netHTTP "net/http"
	netURL "net/url"
	"sync"
)

// 作为请求参数的配置选项
// 先把基础的一些配置开发了 其它配置后面再优化添加
var rWmu sync.RWMutex

type Option struct {
	domain         string
	path           string
	query          string
	params         *netURL.Values // 先占位 后续更新
	_tempParams    any
	headers        *netHTTP.Header
	_tempCookies   any
	cookies        []*netHTTP.Cookie
	enableCookie   bool // 默认使用 用于某些情况下是否使用cookie的情况
	data           *netURL.Values
	_tempData      any
	json           []byte // 这里传入任何可以转成json的对象 然后我会记录在这里
	verify         bool   // 默认true
	files          string // 先占位 后续更新
	proxy          string
	method         string // 默认GET -> 通过 option.SetMethod("POST")调整
	timeout        int    // ms  单位为毫秒
	allowRedirects bool   // 是否允许重定向，默认允许
	stream         string // 先占位 后续更新
	auth           string // 先占位 后续更新
	cert           string // 先占位 后续更新
}

func NewOption() *Option {
	//这个是控制单次请求的一些东西，部分参数会和全局重复 优先级为这里优先
	return &Option{
		method:         "GET",
		allowRedirects: true,
		verify:         true,
		enableCookie:   true,
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

func (o *Option) SetProxy(proxy string) *Option {
	o.proxy = proxy
	return o
}

func (o *Option) GetProxy() string {
	return o.proxy
}

func (o *Option) SetTimeout(timeout int) *Option {
	o.timeout = timeout
	return o
}

func (o *Option) GetTimeout() int {
	return o.timeout
}

func (o *Option) SetURLDetail(rawUrl string) *Option {
	host, path, query := parseDomain(rawUrl)
	o.domain = host
	o.path = path
	o.query = query
	return o
}
