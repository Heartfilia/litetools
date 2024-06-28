package opt

import (
	"log"
	netHTTP "net/http"
)

func NewHeaders() *netHTTP.Header {
	return &netHTTP.Header{}
}

func (o *Option) SetHeaders(headers any) *Option {
	if headers != nil {
		rWmu.RLock()
		switch headers.(type) {
		case map[string]string:
			baseHeaders := NewHeaders()
			for key, value := range headers.(map[string]string) {
				baseHeaders.Set(key, value)
			}
			o.headers = baseHeaders
		case *netHTTP.Header:
			o.headers = headers.(*netHTTP.Header)
		default:
			log.Panicln("Headers only support <*http.Header || map[string]string>")
		}
		rWmu.RUnlock()
	}
	return o
}

func (o *Option) GetHeaders() netHTTP.Header {
	if o.headers == nil {
		return nil
	}
	return *o.headers
}
