package litereq

import netHTTP "net/http"

type Response struct {
	// StatusCode is the status code of the Response
	StatusCode int
	Body       []byte
	Text       string
	Headers    netHTTP.Header
	Ctx        *Context
	err        error // 记录错误详情
}

type History struct {
}

func NewResponse() *Response {
	return &Response{
		StatusCode: netHTTP.StatusOK, // 先默认一下状态值
	}
}

func (r *Response) Error() error {
	// 返回响应记录的错误
	return r.err
}

func (r *Response) setError(err error) {
	r.err = err
}
