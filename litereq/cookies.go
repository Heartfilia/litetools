package litereq

import (
	"context"
	"fmt"
	"net/http"
	"net/http/cookiejar"
	"strings"

	"golang.org/x/net/publicsuffix"
)

type Cookies struct {
	jar []*http.Cookie
	ctx context.Context
}

// NewCookieJar returns a cookie jar using the standard public suffix list.
func NewCookieJar() http.CookieJar {
	jar, err := cookiejar.New(&cookiejar.Options{PublicSuffixList: publicsuffix.List})
	// As of Go 1.16, cookiejar.New err is hardcoded nil
	if err != nil {
		panic(err)
	}
	return jar
}

func (c *Cookies) String() string {
	rawCk := make([]string, 0)
	for _, ck := range c.jar {
		rawCk = append(rawCk, fmt.Sprintf("%s=%s", ck.Name, ck.Value))
	}
	return strings.Join(rawCk, ";")
}

func (c *Cookies) Jar() []*http.Cookie {
	return c.jar
}
