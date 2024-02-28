package litenet

import (
	"github.com/Heartfilia/litetools/litenet/request"
	"log"
	"net"
	"regexp"
)

var ipV4Rule *regexp.Regexp
var internalAddresses []string // 缓存

func GetLAN() []string {
	if ipV4Rule == nil {
		rule, err := regexp.Compile("^\\d{1,3}\\.\\d{1,3}\\.\\d{1,3}\\.\\d{1,3}$")
		if err != nil {
			log.Fatal("创建正则表达柿失败（我故意写错的，因为这里不可能错")
		}
		ipV4Rule = rule
	}
	if internalAddresses == nil {
		adders, err := net.InterfaceAddrs()
		if err != nil {
			panic(err)
		}

		for _, addr := range adders {
			ipNet, ok := addr.(*net.IPNet)
			if !ok || ipNet == nil {
				continue
			}
			lan := ipNet.IP.String()
			if ipV4Rule.Match([]byte(lan)) && lan != "127.0.0.1" {
				// 我只弄IPV4
				internalAddresses = append(internalAddresses, ipNet.IP.String())
			}
		}
	}
	return internalAddresses
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
