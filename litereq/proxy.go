package litereq

import (
	"github.com/Heartfilia/litetools/litereq/utils"
	"net/url"
	"strings"
	"time"
)

type Proxy interface {
	URL() *url.URL
	Expired() time.Duration
	String() string
}

type ProxyStatus int
type ProxyChange int
type ProxyTimelyType int

type ProxyGetter = func() Proxy

type ProxyInfo struct {
	ProxyIp     string
	JustIp      string
	Area        string
	ExpireTime  time.Time
	Status      ProxyStatus
	ProxyTimely ProxyTimelyType
	Change      ProxyChange
}

func (p *ProxyInfo) URL() *url.URL {
	if p == nil {
		return nil
	}
	return parseURL(p.ProxyIp)
}

func parseURL(proxyIp string) *url.URL {
	if proxyIp == "" {
		return nil
	}
	prefix := "http" + "://" // 避免pycharm的http提示而已
	if strings.HasPrefix(proxyIp, "socks5://") {
		prefix = "socks5://"
	}
	proxy := strings.ReplaceAll(proxyIp, prefix, "")
	s := strings.Split(proxy, ":")
	pUrl, err := url.Parse(prefix + s[0] + ":" + s[1])
	if err != nil {
		panic("invalid ProxyUrl")
	}
	if len(s) == 4 {
		pUrl.User = url.UserPassword(s[2], s[3])
	}
	return pUrl
}

func (p *ProxyInfo) Expired() time.Duration {
	if p != nil {
		return utils.Max(p.ExpireTime.Sub(time.Now()), minProxyExpired)
	}
	return maxProxyExpired
}

func (p *ProxyInfo) String() string {
	if p == nil {
		return ""
	}
	return p.ProxyIp
}
