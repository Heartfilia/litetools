package litereq

import (
	"context"
	"io"
	"net/http"
	"time"
)

// 借鉴 https://github.com/earthboundkid/requests

type Builder struct {
	hc      *http.Client
	rb      requestBuilder
	timeout time.Duration
}

func Build(urlPath string) *Builder {
	build := &Builder{}
	return build
}

func (b *Builder) request(ctx context.Context) (req *http.Request, err error) {

	return nil, nil
}

func (b *Builder) do(req *http.Request) (err error) {

	return nil
}

func (b *Builder) Method(method string) *Builder {

	return b
}

func (b *Builder) Header(key string, values ...string) *Builder {

	return b
}

func (b *Builder) Body(src BodyGetter) *Builder {
	b.rb.Body(src)
	return b
}

func (b *Builder) BodyWriter(f func(w io.Writer) error) *Builder {
	b.Body(BodyWriter(f))
	return b
}

func (b *Builder) ContentType(ct string) *Builder {
	return b.Header("Content-Type", ct)
}

func (b *Builder) Fetch(ctx context.Context) *Response {
	resp := &Response{}
	req, err := b.request(ctx)
	if err != nil {
		resp.error(err)
		return resp
	}
	err = b.do(req)
	resp.error(err)
	return resp
}
