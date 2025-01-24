package litereq

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
)

type Response struct {
	Status        int
	Proto         string
	Header        http.Header
	cb            *Cookies
	Body          io.ReadCloser
	Content       []byte
	Text          string
	ContentLength int
	json          string
	err           error
}

func (r *Response) error(err error) {
	r.err = err
}

func (r *Response) Error() error {
	return r.err
}

func (r *Response) cookie(c []*http.Cookie) {
	r.cb = &Cookies{jar: c}
}

func (r *Response) Cookie() *Cookies {
	return r.cb
}

func (r *Response) detail(rc io.ReadCloser) {
	bodyBytes, err := io.ReadAll(rc)
	if err != nil {
		log.Println(err)
		return
	}
	r.Body = rc
	r.Content = bodyBytes
	r.ContentLength = len(bodyBytes)
	r.Text = string(bodyBytes)
}

func (r *Response) Json(v any) error {
	err := json.Unmarshal(r.Content, &v)
	if err != nil {
		return err
	}
	return nil
}
