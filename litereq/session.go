package litereq

import (
	"errors"
	"github.com/Heartfilia/litetools/litereq/opt"
	"github.com/Heartfilia/litetools/litestr"
	"io"
	"log"
	netHTTP "net/http"
	"strings"
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
	client       *netHTTP.Client
	headers      *netHTTP.Header   // 全局headers
	cookies      []*netHTTP.Cookie // 全局的cookies
	_tempCookies any               // 用于临时记录
	_tempProxy   string            // 临时记录proxy
}

func NewSession() *Session {
	return &Session{
		maxRetry: 1,
		verbose:  false,
		http2:    false,
		client:   &netHTTP.Client{},
	}
}

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
				return nil
			}
			response.Body = respBytes
			response.Text = string(respBytes)
			response.Headers = resp.Header
			response.StatusCode = resp.StatusCode
			response.Proto = resp.Proto
			response.Status = resp.Status
			response.ContentLength = int(resp.ContentLength)
			suc = true
		}
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

func (s *Session) http1Request(url string, o *opt.Option) (*netHTTP.Response, []byte, error) {
	var req *netHTTP.Request
	var err error
	if o.GetMethod() == "POST" {
		body := strings.NewReader("测试")
		req, err = netHTTP.NewRequest(o.GetMethod(), url, body)
		if err != nil {
			return nil, nil, err
		}
	} else {
		req, err = netHTTP.NewRequest("GET", url, nil)
		if err != nil {
			return nil, nil, err
		}
	}
	if o.GetHeaders() != nil {
		req.Header = o.GetHeaders()
	} else if s.GetHeaders() != nil {
		req.Header = s.GetHeaders()
	}
	if o.GetParams() != nil {
		req.URL.RawQuery = o.GetParams().Encode()
	}

	if o.GetCookieEnable() {
		if o.GetCookies() != nil {
			for _, ck := range o.GetCookies() {
				req.AddCookie(ck)
			}
		} else if s.GetCookies() != nil {
			for _, ck := range s.GetCookies() {
				req.AddCookie(ck)
			}
		}
	}

	resp, err := s.client.Do(req)
	if err != nil {
		return nil, nil, err
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {

		}
	}(resp.Body)
	respByte, _ := io.ReadAll(resp.Body)

	return resp, respByte, nil
}

func (s *Session) handle3XXResponse() {
	// 处理 30X 的响应
}

func (s *Session) SetProxy(proxy string) *Session {
	// 这里 是全局代理 优先级低于独立配置的代理位置： 这里更加适合放隧道代理或者长效代理
	s._tempProxy = proxy
	return s
}

func (s *Session) setProxy() {

}

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

func (s *Session) GetHeaders() netHTTP.Header {
	if s.headers == nil {
		return nil
	}
	return *s.headers
}

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

func (s *Session) GetCookies() []*netHTTP.Cookie {
	return s.cookies
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
