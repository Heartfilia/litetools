package opt

import (
	"github.com/Heartfilia/litetools/litestr"
	"log"
	netURL "net/url"
	"strings"
)

func parseDomain(rawURL string) (string, string, string) {
	parsedURL, err := netURL.Parse(rawURL)
	if err != nil {
		return "", "", ""
	}
	return parsedURL.Host, parsedURL.Path, parsedURL.RawQuery
}

func parseStringParams(raw string) map[string]string {
	query, err := netURL.QueryUnescape(raw)
	if err != nil {
		log.Println(litestr.E(), "parse params error:", err)
		return nil
	}
	rawList := strings.SplitN(query, "&", -1)
	res := make(map[string]string)
	for _, rawItem := range rawList {
		items := strings.SplitN(rawItem, "=", 2)
		if len(items) == 2 {
			res[items[0]] = items[1]
		}
	}
	return res
}
