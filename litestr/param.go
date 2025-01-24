package litestr

import (
	"fmt"
	"strings"
)

func ParamStringToMap(ps string) map[string]string {
	params := make(map[string]string)
	for _, p := range strings.Split(ps, "&") {
		kv := strings.SplitN(p, "=", 2)
		params[kv[0]] = kv[1]
	}
	return params
}

func ParamStringToArray(ps string) [][2]string {
	slice := strings.Split(ps, "&")
	fmt.Println(len(slice))
	params := make([][2]string, 0, len(slice))
	for _, p := range slice {
		kv := strings.SplitN(p, "=", 2)
		params = append(params, [2]string{kv[0], kv[1]})
	}
	return params
}
