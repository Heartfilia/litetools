package opt

import netHTTP "net/http"

func NewHeaders() *netHTTP.Header {
	return &netHTTP.Header{}
}
