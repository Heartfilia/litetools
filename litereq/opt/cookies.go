package opt

import netHTTP "net/http"

func NewCookies() *netHTTP.Cookie {
	return &netHTTP.Cookie{}
}
