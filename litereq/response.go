package litereq

type Response struct {
	Status        int
	Header        map[string]string // 先占位
	Cookie        map[string]string // 先占位
	Body          []byte
	Text          string
	ContentLength int
	Json          string
	err           error
}

func (r *Response) error(err error) {
	r.err = err
}

func (r *Response) Error() error {
	return r.err
}
