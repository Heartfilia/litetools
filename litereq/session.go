package litereq

import (
	"errors"
	"github.com/Heartfilia/litetools/litereq/opt"
	"log"
	netHTTP "net/http"
)

/*
这个项目的宗旨：不是为了创造一个框架，而是创建一个便携整体流程的通用请求
核心部分均采用系统的 net  之类的这种包 减少不必要的兼容麻烦
（更加是为了自己以后接单啥的，能写请求更快，目前市面上的感觉都用不习惯
*/

type Session struct {
	maxRetry     int  // max retry, default 1
	http2        bool // default false   先不忙支持 后面我会弄的
	verbose      bool // default false 就是用于打印详细日志的
	option       *opt.Option
	client       *netHTTP.Client
	headers      *netHTTP.Header   // 全局headers
	cookies      []*netHTTP.Cookie // 全局的cookies
	_tempCookies any               // 用于临时记录
	//globalCookie  // 需要记录下来全局的cookie信息
}

func NewSession() *Session {
	return &Session{
		maxRetry: 1,
		verbose:  false,
		http2:    false,
		client:   &netHTTP.Client{},
	}
}

func (s *Session) Do(url string) *Response {
	// main : 这里可以处理一些额外的操作 但是目前我这里先省略
	s.setCookies(url)
	return s.sendRequest(url)
}

func (s *Session) SetOption(o *opt.Option) *Session {
	if o == nil {
		s.option = opt.NewOption()
	} else {
		s.option = o
	}
	return s
}

func (s *Session) sendRequest(url string) *Response {
	response := NewResponse()
	suc := false
	for r := 0; r < s.maxRetry; r++ {

		if s.verbose {
			// 这里是在过程中遇到的报错打印出来
		}
		if suc == true {
			break
		}
	}
	if suc == false {
		response.err = errors.New("bad requests")
	}
	return response
}

func (s *Session) handle3XXResponse() {
	// 处理 30X 的响应
}

func (s *Session) SetHeaders(header any) *Session {
	// 这个方法是直接操作类似 option里面的操作了
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
	return s
}

func (s *Session) SetCookies(cookie any) *Session {
	s._tempCookies = cookie
	return s
}

func (s *Session) setCookies(rawUrl string) *Session {
	// 这个地方才是主要的操作 option里面的操作了 --> 这里其实属于慢操作，核心的
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
			// if string -->  k=v; k=v  --> map[string]string --> save

		default:
			log.Panicln("Cookies only support <[]*http.Cookie || *http.Cookie || map[string]string || string>")
		}
	}
	return s
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
