package litereq

import "net/http"

type Wrap[T any] struct {
	Req     *http.Request
	Resp    *http.Response
	Curl    string
	RespStr string
	Data    T
	Err     error
}
