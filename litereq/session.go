package litereq

import (
	"github.com/Heartfilia/litetools/litereq/opt"
	netHTTP "net/http"
)

type Session struct {
	MaxRetry int // max retry, default 1
	client   *netHTTP.Client
	headers  map[string]string // 全局headers
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
	return response
}

func (s *Session) storeCookie() {
	// 用于每次请求后保存当前的请求的cookie的

}
