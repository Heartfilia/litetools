package opt

import (
	"fmt"
	netHTTP "net/http"
	"strings"
)

type Cookie struct {
	originCookie []*netHTTP.Cookie
}

func NewCookies() *netHTTP.Cookie {
	return &netHTTP.Cookie{}
}

func (c *Cookie) Cookies() []*netHTTP.Cookie {
	return c.originCookie
}

func (c *Cookie) StoreCookies(ck []*netHTTP.Cookie) {
	c.originCookie = ck
}

func (c *Cookie) String() string {
	if c.Cookies() == nil {
		return ""
	}
	baseCK := make([]string, 0)
	for _, ck := range c.Cookies() {
		baseCK = append(baseCK, fmt.Sprintf("%s=%s", ck.Name, ck.Value))
	}
	if len(baseCK) == 0 {
		return ""
	}
	return strings.Join(baseCK, "; ")
}
