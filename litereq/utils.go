package litereq

import (
	netURL "net/url"
)

func parseDomain(rawURL string) string {
	parsedURL, err := netURL.Parse(rawURL)
	if err != nil {
		return ""
	}
	return parsedURL.Host // 只管当前域名
}
