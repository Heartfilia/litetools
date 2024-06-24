package litereq

import (
	"github.com/Heartfilia/litetools/litestr"
	"log"
	netURL "net/url"
)

func parseDomain(rawURL string) string {
	parsedURL, err := netURL.Parse(rawURL)
	if err != nil {
		return ""
	}
	return parsedURL.Host // 只管当前域名
}

func mustParseURL(rawURL string) *netURL.URL {
	u, err := netURL.Parse(rawURL)
	if err != nil {
		log.Println(litestr.E(), ":", err)
		return nil
	}
	return u
}
