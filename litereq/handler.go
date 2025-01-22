package litereq

import "net/http"

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
