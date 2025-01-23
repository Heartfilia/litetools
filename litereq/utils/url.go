package utils

import "net/url"

func ParseUrl(u string) *url.URL {
	parse, err := url.Parse(u)
	if err != nil {
		return nil
	}
	return parse
}
