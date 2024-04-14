package litereq

import (
	"github.com/Heartfilia/litetools/litereq/opt"
	"net/url"
)

type Session struct {
	MaxRetry int // max retry, default 1
	request  *Request
	response *Response
}

type Request struct {
	URL     *url.URL
	Ctx     *Context
	Options *opt.Option
}

func NewSession() *Session {
	return &Session{
		MaxRetry: 1,
		request:  &Request{},
		response: &Response{},
	}
}
