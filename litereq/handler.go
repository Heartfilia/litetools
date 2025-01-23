package litereq

import (
	"io"
	"net/http"
)

// ResponseHandler is used to validate or handle the response to a request.
type ResponseHandler = func(*http.Response) error

// Or returns the first non-empty argument it receives
// or the zero value for T.
func Or[T comparable](vals ...T) T {
	for _, val := range vals {
		if val != *new(T) {
			return val
		}
	}
	return *new(T)
}

func If[T any](val bool, a, b T) T {
	if val {
		return a
	}
	return b
}

func consumeBody(res *http.Response) (err error) {
	const maxDiscardSize = 640 * 1 << 10
	if _, err = io.CopyN(io.Discard, res.Body, maxDiscardSize); err == io.EOF {
		err = nil
	}
	return err
}
