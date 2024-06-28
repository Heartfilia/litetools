package opt

import (
	"encoding/json"
	"fmt"
	"github.com/Heartfilia/litetools/litestr"
	"log"
	netHTTP "net/http"
	"strings"
)

type Cookie struct {
	originCookie []*netHTTP.Cookie
}

func NewCookies() *netHTTP.Cookie {
	return &netHTTP.Cookie{}
}

// -------------------------------------------------------------------------

func (o *Option) setCookies() *Option {
	// 这个地方才是主要的操作 option里面的操作了 --> 这里其实属于慢操作，核心的
	if o.cookies == nil {
		o.cookies = make([]*netHTTP.Cookie, 0)
	}
	cookie := o._tempCookies
	if cookie != nil {
		switch cookie.(type) {
		case map[string]string:
			for key, value := range cookie.(map[string]string) {
				baseCookies := NewCookies()
				baseCookies.Name = key
				baseCookies.Value = value
				baseCookies.Path = "/"
				baseCookies.Domain = o.domain

				exists := false
				if len(o.cookies) > 0 {
					for ind, ck := range o.cookies {
						if ck.Name == key && ck.Domain == o.domain { // 如果是存在的cookie那么就要替换
							o.cookies[ind] = baseCookies
							exists = true
						}
					}
				}
				if !exists {
					o.cookies = append(o.cookies, baseCookies)
				}
			}
		case []*netHTTP.Cookie:
			o.cookies = cookie.([]*netHTTP.Cookie)
		case *netHTTP.Cookie:
			exists := false
			for ind, ck := range o.cookies {
				if ck.Name == cookie.(*netHTTP.Cookie).Name && ck.Domain == cookie.(*netHTTP.Cookie).Domain {
					o.cookies[ind] = ck
					exists = true
				}
			}
			if !exists {
				o.cookies = append(o.cookies, cookie.(*netHTTP.Cookie))
			}
		case string:
			mapCookie := litestr.CookieStringToMap(cookie.(string))
			for key, value := range mapCookie {
				baseCookies := NewCookies()
				baseCookies.Name = key
				baseCookies.Value = value
				baseCookies.Path = "/"
				baseCookies.Domain = o.domain

				exists := false
				if len(o.cookies) > 0 {
					for ind, ck := range o.cookies {
						if ck.Name == key && ck.Domain == o.domain { // 如果是存在的cookie那么就要替换
							o.cookies[ind] = baseCookies
							exists = true
						}
					}
				}
				if !exists {
					o.cookies = append(o.cookies, baseCookies)
				}
			}
		default:
			log.Panicln("Cookies only support <[]*http.Cookie || *http.Cookie || map[string]string || string>")
		}
	}
	return o
}

func (o *Option) GetCookies() []*netHTTP.Cookie {
	return o.cookies
}

func (o *Option) SetCookies(cookie any) *Option {
	o._tempCookies = cookie
	o.setCookies()
	return o
}

func (o *Option) SetCookieEnable(enable bool) *Option {
	o.enableCookie = enable
	return o
}

func (o *Option) GetCookieEnable() bool {
	return o.enableCookie
}

// -------------------------------------------------------------------------

func (c *Cookie) Cookies() []*netHTTP.Cookie {
	return c.originCookie
}

// StoreCookies : do not use it
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

func (c *Cookie) Map() map[string]string {
	if c.Cookies() == nil {
		return map[string]string{}
	}
	baseCK := make(map[string]string)
	for _, ck := range c.Cookies() {
		baseCK[ck.Name] = ck.Value
	}
	if len(baseCK) == 0 {
		return map[string]string{}
	}
	return baseCK
}

func (c *Cookie) Json() string {
	baseCK := c.Map()
	marshal, err := json.Marshal(baseCK)
	if err != nil {
		return "{}"
	}
	return string(marshal)
}
