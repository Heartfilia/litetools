package litestr

import (
	"fmt"
	"strings"
)

func CookieStringToMap(ck string) map[string]string {
	rawList := strings.SplitN(ck, ";", -1)
	baseCookie := make(map[string]string, 0)
	for _, raw := range rawList {
		raw = strings.TrimSpace(raw)
		keyValue := strings.SplitN(raw, "=", 2)
		if len(keyValue) == 2 {
			baseCookie[keyValue[0]] = keyValue[1]
		}
	}
	return baseCookie
}

func CookieMapToString(ck map[string]string) string {
	baseCookieList := make([]string, 0)
	for key, value := range ck {
		baseCookieList = append(baseCookieList, fmt.Sprintf("%s=%s", key, value))
	}
	res := strings.Join(baseCookieList, "; ")
	return res
}
