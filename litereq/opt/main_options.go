package opt

import (
	"encoding/json"
	"fmt"
	"github.com/Heartfilia/litetools/litestr"
	"log"
	netHTTP "net/http"
	netURL "net/url"
	"strings"
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
	enableCookie   bool   // 默认使用 用于某些情况下是否使用cookie的情况
	data           string // 先占位 后续更新
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

func (o *Option) SetParams(params any) *Option {
	o._tempParams = params
	return o
}

func (o *Option) GetParams() netURL.Values {
	// 传入
	params := o._tempParams
	if params != nil {
		parse, _ := netURL.Parse(fmt.Sprintf("https://%s?%s", o.domain, o.query))
		query := parse.Query()
		switch params.(type) {
		case map[string]any:
			for k, v := range params.(map[string]any) {
				query.Set(k, fmt.Sprintf("%v", v))
			}
		case map[string]string:
			for k, v := range params.(map[string]string) {
				query.Set(k, v)
			}
		case [][2]any:
			for _, eachParam := range params.([][2]any) {
				query.Set(fmt.Sprintf("%v", eachParam[0]), fmt.Sprintf("%v", eachParam[1]))
			}
		case [][2]string:
			for _, eachParam := range params.([][2]string) {
				query.Set(eachParam[0], eachParam[1])
			}
		case netURL.Values:
			query = params.(netURL.Values)
		case string:
			items := parseStringParams(params.(string))
			if items != nil {
				for k, v := range items {
					query.Set(k, v)
				}
			}
		default:
			log.Panicln("Params only support <url.Values || map[string]string || map[string]any || string || [][2]any || [][2]string>")
		}

		o.params = &query

		return query
	}

	return nil
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

func (o *Option) GetMethod() string {
	return o.method
}

func (o *Option) SetProxy(proxy string) *Option {
	o.proxy = proxy
	return o
}

func (o *Option) GetProxy() string {
	return o.proxy
}

func (o *Option) SetHeaders(headers any) *Option {
	if headers != nil {
		rWmu.RLock()
		switch headers.(type) {
		case map[string]string:
			baseHeaders := NewHeaders()
			for key, value := range headers.(map[string]string) {
				baseHeaders.Set(key, value)
			}
			o.headers = baseHeaders
		case *netHTTP.Header:
			o.headers = headers.(*netHTTP.Header)
		default:
			log.Panicln("Headers only support <*http.Header || map[string]string>")
		}
		rWmu.RUnlock()
	}
	return o
}

func (o *Option) GetHeaders() netHTTP.Header {
	if o.headers == nil {
		return nil
	}
	return *o.headers
}

func (o *Option) SetCookies(cookie any) *Option {
	o._tempCookies = cookie
	o.setCookies()
	return o
}

func (o *Option) SetCookieEnable(enable bool) *Option {
	o.enableCookie = enable
	return o
}

func (o *Option) GetCookieEnable() bool {
	return o.enableCookie
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

func (o *Option) setCookies() *Option {
	// 这个地方才是主要的操作 option里面的操作了 --> 这里其实属于慢操作，核心的
	if o.cookies == nil {
		o.cookies = make([]*netHTTP.Cookie, 0)
	}
	cookie := o._tempCookies
	if cookie != nil {
		switch cookie.(type) {
		case map[string]string:
			for key, value := range cookie.(map[string]string) {
				baseCookies := NewCookies()
				baseCookies.Name = key
				baseCookies.Value = value
				baseCookies.Path = "/"
				baseCookies.Domain = o.domain

				exists := false
				if len(o.cookies) > 0 {
					for ind, ck := range o.cookies {
						if ck.Name == key && ck.Domain == o.domain { // 如果是存在的cookie那么就要替换
							o.cookies[ind] = baseCookies
							exists = true
						}
					}
				}
				if !exists {
					o.cookies = append(o.cookies, baseCookies)
				}
			}
		case []*netHTTP.Cookie:
			o.cookies = cookie.([]*netHTTP.Cookie)
		case *netHTTP.Cookie:
			exists := false
			for ind, ck := range o.cookies {
				if ck.Name == cookie.(*netHTTP.Cookie).Name && ck.Domain == cookie.(*netHTTP.Cookie).Domain {
					o.cookies[ind] = ck
					exists = true
				}
			}
			if !exists {
				o.cookies = append(o.cookies, cookie.(*netHTTP.Cookie))
			}
		case string:
			mapCookie := litestr.CookieStringToMap(cookie.(string))
			for key, value := range mapCookie {
				baseCookies := NewCookies()
				baseCookies.Name = key
				baseCookies.Value = value
				baseCookies.Path = "/"
				baseCookies.Domain = o.domain

				exists := false
				if len(o.cookies) > 0 {
					for ind, ck := range o.cookies {
						if ck.Name == key && ck.Domain == o.domain { // 如果是存在的cookie那么就要替换
							o.cookies[ind] = baseCookies
							exists = true
						}
					}
				}
				if !exists {
					o.cookies = append(o.cookies, baseCookies)
				}
			}
		default:
			log.Panicln("Cookies only support <[]*http.Cookie || *http.Cookie || map[string]string || string>")
		}
	}
	return o
}

func (o *Option) GetCookies() []*netHTTP.Cookie {
	return o.cookies
}

func (o *Option) SetJson(object any) *Option {
	if object != nil {
		marshal, err := json.Marshal(object)
		if err != nil {
			log.Panicln(litestr.E(), "error json object:", err)
			return o
		}
		o.json = marshal
	}
	return o
}

func (o *Option) GetJson() []byte {
	return o.json
}
