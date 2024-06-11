package litereq

import "net/http"

type Response struct {
	// StatusCode is the status code of the Response
	StatusCode int
	Body       []byte
	Text       string
	Headers    *http.Header
	Ctx        *Context
}

type History struct {
}
