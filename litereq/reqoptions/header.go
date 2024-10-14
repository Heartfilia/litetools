package reqoptions

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

// UpdateHeaderMap  :update Option Header By map
func (o *Option) UpdateHeaderMap(headers map[string]string) *Option {
	if headers != nil {
		rWmu.RLock()
		if o.headers == nil {
			o.headers = NewHeaders()
		}
		for key, value := range headers {
			o.headers.Set(key, value)
		}
		rWmu.RUnlock()
	}
	return o
}

// UpdateHeaderKeyValue  :update Option Header By key-value
func (o *Option) UpdateHeaderKeyValue(key, value string) *Option {
	if key != "" && value != "" {
		rWmu.RLock()
		if o.headers == nil {
			o.headers = NewHeaders()
		}
		o.headers.Set(key, value)
		rWmu.RUnlock()
	}
	return o
}

// ExceptGlobalHeaders  : do not use global header
//
// @Param key: some header you don't want to use in global header area
func (o *Option) ExceptGlobalHeaders(key []string) *Option {
	// 这里主要是用于屏蔽全局的header 在某个条件里面不想使用它的时候 可以用这个把对应的key忽略
	o._exceptHeaders = key
	return o
}

func (o *Option) GetUnHeaderExcept() []string {
	return o._exceptHeaders
}

func (o *Option) GetHeaders() netHTTP.Header {
	if o.headers == nil {
		return nil
	}
	return *o.headers
}
