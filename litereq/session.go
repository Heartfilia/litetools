package litereq

import (
	"errors"
	"github.com/Heartfilia/litetools/litereq/opt"
	"log"
	netHTTP "net/http"
)

/*
这个项目的宗旨：不是为了创造一个框架，而是创建一个便携请求
核心部分均采用系统的 net  之类的这种包 减少不必要的兼容麻烦
（更加是为了自己以后接单啥的，能写请求更快，目前市面上的感觉都用不习惯
*/

type Session struct {
	MaxRetry int  // max retry, default 1
	HTTP2    bool // default false   先不忙支持 后面我会弄的
	client   *netHTTP.Client
	headers  *netHTTP.Header // 全局headers
	//globalCookie  // 需要记录下来全局的cookie信息
}

func NewSession() *Session {
	return &Session{
		MaxRetry: 1,
		client:   &netHTTP.Client{},
	}
}

func (s *Session) Do(url string, option *opt.Option) *Response {
	if option == nil {
		option = opt.NewOption()
	}

	return &Response{}
}

func (s *Session) sendRequest(url string, option *opt.Option) *Response {
	response := NewResponse()
	suc := false
	for r := 0; r < s.MaxRetry; r++ {

		if suc == true {
			break
		}
	}
	if suc == false {
		response.err = errors.New("bad requests")
	}
	return response
}

func (s *Session) SetHeaders(header any) {
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
		log.Panicln("Headers only support <>")
	}
}

func (s *Session) SetRetry(retry int) {
	s.MaxRetry = retry
}

func (s *Session) SetHTTP2(h2 bool) {
	s.HTTP2 = h2
}

func (s *Session) storeCookie() {
	// 用于每次请求后保存当前的请求的cookie的

}
