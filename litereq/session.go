package litereq

import (
	"errors"
	"github.com/Heartfilia/litetools/litereq/opt"
	"github.com/Heartfilia/litetools/litestr"
	"io"
	"log"
	netHTTP "net/http"
	"strings"
	"time"
)

/*
这个项目的宗旨：不是为了创造一个框架，而是创建一个便携整体流程的通用请求
核心部分均采用系统的 net  之类的这种包 减少不必要的兼容麻烦
（更加是为了自己以后能写请求更快，目前市面上的感觉都用不习惯

TODO 下面设置session的全局参数的时候 需要枷锁 比如 cookie  headers 后面弄
*/

type Session struct {
	maxRetry     int  // max retry, default 1
	http2        bool // default false --> 先不忙支持 后面我会弄的
	verbose      bool // default false 就是用于打印详细日志的
	timeout      int  // 毫秒的单位 不传不管 这里是全局参数 在单独请求哪里也有这个控制
	client       *netHTTP.Client
	headers      *netHTTP.Header   // 全局headers
	cookies      []*netHTTP.Cookie // 全局的cookies
	_tempCookies any               // 用于临时记录
	_tempProxy   string            // 临时记录proxy
}

// NewSession : create base session object that can be chained
func NewSession() *Session {
	return &Session{
		maxRetry: 1,
		verbose:  false,
		http2:    false,
		client:   &netHTTP.Client{},
	}
}

// Fetch     : do the last request
//
// @Param url: The target page you are requesting
//
// @Param o  : Single request parameter option <or> nil
func (s *Session) Fetch(url string, o *opt.Option) *Response {
	// main : 这里可以处理一些额外的操作 但是目前我这里先省略
	s.setCookies(url)
	if o == nil {
		o = opt.NewOption()
	}
	o.SetURLDetail(url)
	return s.sendRequest(url, o)
}

func (s *Session) sendRequest(url string, o *opt.Option) *Response {
	response := NewResponse()
	suc := false
	for r := 0; r < s.maxRetry; r++ {
		if s.http2 {
			// 如果是http2模式下走这个地方 现在先不兼容
		} else {
			resp, respBytes, err := s.http1Request(url, o)
			if err != nil {
				if s.verbose {
					log.Println(litestr.E(), err)
				}
				response.err = err
				continue
			}
			response.Body = respBytes
			response.Text = string(respBytes)
			response.Headers = resp.Header
			response.StatusCode = resp.StatusCode
			response.Proto = resp.Proto
			response.Status = resp.Status
			response.ContentLength = int(resp.ContentLength)
			response.err = nil
			suc = true
		}
		if s.verbose && response.Error() != nil {
			// 这里是在过程中遇到的报错打印出来
			log.Println(litestr.E(), "error:", response.Error())
		}
		if suc == true {
			break
		}
	}
	if suc == false && response.Error() == nil {
		// 如果失败的时候 并且没有失败的日志记录 那么补充一个错误提示
		response.err = errors.New("bad requests with this packages: help me fix it with debug")
	}
	return response
}

func (s *Session) http1Request(url string, o *opt.Option) (*netHTTP.Response, []byte, error) {
	var req *netHTTP.Request
	var err error
	switch o.GetMethod() {
	case "POST", "PUT", "DELETE", "PATCH":
		body := strings.NewReader("测试")
		req, err = netHTTP.NewRequest(o.GetMethod(), url, body)
	case "GET", "HEAD":
		req, err = netHTTP.NewRequest(o.GetMethod(), url, nil)
	default:
		log.Panicf("not support your method: %s", o.GetMethod())
	}
	if err != nil {
		return nil, nil, err
	}

	if o.GetParams() != nil {
		req.URL.RawQuery = o.GetParams().Encode()
	}
	s.setReqHeaders(req, o.GetHeaders())
	if o.GetCookieEnable() {
		s.setReqCookies(req, o.GetCookies())
	}
	s.setTimeout(o.GetTimeout())
	s.setProxy(o.GetProxy())

	resp, err := s.client.Do(req)
	if err != nil {
		return nil, nil, err
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {

		}
	}(resp.Body)
	respByte, err := io.ReadAll(resp.Body)

	return resp, respByte, err
}

func (s *Session) handle3XXResponse() {
	// 处理 30X 的响应
}

// SetProxy : Set global proxy: example > http://name:pass@ip:port <or> http://ip:port
func (s *Session) SetProxy(proxy string) *Session {
	// 这里 是全局代理 优先级低于独立配置的代理位置： 这里更加适合放隧道代理或者长效代理
	s._tempProxy = proxy
	return s
}

func (s *Session) setProxy(optionProxy string) {
	if optionProxy != "" { // 优先使用option里面的代理
		transport := &netHTTP.Transport{
			Proxy: netHTTP.ProxyURL(mustParseURL(optionProxy)),
		}
		s.client.Transport = transport
	} else if s._tempProxy != "" { // 其次使用全局的代理
		transport := &netHTTP.Transport{
			Proxy: netHTTP.ProxyURL(mustParseURL(s._tempProxy)),
		}
		s.client.Transport = transport
	}
}

func (s *Session) SetTimeout(timeout int) *Session {
	s.timeout = timeout
	return s
}

func (s *Session) setTimeout(optionTimeout int) {
	if optionTimeout > 0 { // 优先使用option里面的超时
		s.client.Timeout = time.Duration(optionTimeout) * time.Millisecond
	} else if s.timeout > 0 { // 其次使用全局的超时
		s.client.Timeout = time.Duration(s.timeout) * time.Millisecond
	}
}

// SetHeaders : Set global headers: support
//
// >>> map[string]string | http.header
func (s *Session) SetHeaders(header any) *Session {
	// 这个方法是直接操作类似 option里面的操作了
	if header != nil {
		switch header.(type) {
		case map[string]string:
			baseHeaders := opt.NewHeaders()
			for key, value := range header.(map[string]string) {
				baseHeaders.Set(key, value)
			}
			s.headers = baseHeaders
		case *netHTTP.Header:
			s.headers = header.(*netHTTP.Header)
		default:
			log.Panicln("Headers only support <*http.Header || map[string]string>")
		}
	}
	return s
}

func (s *Session) setReqHeaders(req *netHTTP.Request, headers netHTTP.Header) {
	if headers != nil {
		req.Header = headers
	} else if s.headers != nil && *s.headers != nil {
		req.Header = *s.headers
	}
}

func (s *Session) setReqCookies(req *netHTTP.Request, cookies []*netHTTP.Cookie) {
	if cookies != nil {
		for _, ck := range cookies {
			req.AddCookie(ck)
		}
	} else if s.cookies != nil {
		for _, ck := range s.cookies {
			req.AddCookie(ck)
		}
	}
}

// SetCookies : Set global Cookies: support
//
// >>> string<a=1;b=2> | map[string]string<map[string]string{"a":"1","b":"2"}> | *http.Cookie
func (s *Session) SetCookies(cookie any) *Session {
	s._tempCookies = cookie
	return s
}

func (s *Session) setCookies(rawUrl string) *Session {
	// 这个地方才是主要的操作 option里面的操作了 --> 这里其实属于慢操作，核心的
	// 这里不需要判断是否在cookie里面已经存在的值了，因为初始化这里的时候才会添加cookie  但是option那边不是 会出现同样的cookie值 避免猛增
	if s.cookies == nil {
		s.cookies = make([]*netHTTP.Cookie, 0)
	}
	cookie := s._tempCookies
	if cookie != nil {
		domain := parseDomain(rawUrl)
		switch cookie.(type) {
		case map[string]string:
			for key, value := range cookie.(map[string]string) {
				baseCookies := opt.NewCookies()
				baseCookies.Name = key
				baseCookies.Value = value
				baseCookies.Path = "/"
				baseCookies.Domain = domain
				s.cookies = append(s.cookies, baseCookies)
			}
		case []*netHTTP.Cookie:
			s.cookies = cookie.([]*netHTTP.Cookie)
		case *netHTTP.Cookie:
			s.cookies = append(s.cookies, cookie.(*netHTTP.Cookie))
		case string:
			mapCookie := litestr.CookieStringToMap(cookie.(string))
			for key, value := range mapCookie {
				baseCookies := opt.NewCookies()
				baseCookies.Name = key
				baseCookies.Value = value
				baseCookies.Path = "/"
				baseCookies.Domain = domain
				s.cookies = append(s.cookies, baseCookies)
			}
		default:
			log.Panicln("Cookies only support <[]*http.Cookie || *http.Cookie || map[string]string || string>")
		}
	}
	return s
}

// GetCookies : return global store cookie >>> all saved cookie
func (s *Session) GetCookies() *opt.Cookie {
	ck := &opt.Cookie{}
	ck.StoreCookies(s.cookies)
	return ck
}

func (s *Session) SetRetry(retry int) *Session {
	if retry < 1 {
		retry = 1 // 如果设置是0或者负数 那么就改为1
	}
	s.maxRetry = retry
	return s
}

func (s *Session) SetHTTP2(h2 bool) *Session {
	s.http2 = h2
	return s
}

func (s *Session) SetVerbose(verbose bool) *Session {
	s.verbose = verbose
	return s
}

func (s *Session) storeCookie() {
	// 用于每次请求后保存当前的请求的cookie的

}
