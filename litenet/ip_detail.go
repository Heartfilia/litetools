package litenet

import (
	"github.com/Heartfilia/litetools/litenet/request"
)

func GetLAN() string {

	return ""
}

var ip = "" // 缓存结果到全局变量后面重复使用

func GetWAN() string {
	if ip == "" {
		result := request.DoGet("http://httpbin.org/ip")
		if result == "" {
			return ""
		}
		ip = result
	}
	return ip
}
