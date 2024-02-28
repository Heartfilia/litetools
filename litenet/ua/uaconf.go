package ua

const DefaultChoice = "chrome"

var Browser = []string{"chrome", "firefox", "opera", "ie", "edge", "safari"}
var System = []string{"pc", "mobile", "mac", "windows", "linux", "android", "ios", "harmony"}

var baseString = "%s (%s; %s) %s"
var UATemplateBrowser = map[string]string{
	"chrome":  "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/%s Safari/537.36",
	"firefox": "Mozilla/5.0 (Windows NT 10.0; Win64; x64; rv:%s.0) Gecko/20100101 Firefox/%s.0",
	"opera":   "Opera/9.80 (Macintosh; Intel Mac OS X 10.6.8; U; en) Presto/2.8.131 Version/11.11",
	"ie":      "Mozilla/5.0 (Windows NT 10.0; WOW64; Trident/7.0; rv:11.0) like Gecko",
	"edge":    "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/%s Safari/537.36 Edg/%s",
	"safari":  "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/605.1.15 (KHTML, like Gecko) Version/%s Safari/605.1.15",
}

var DefaultSetting = map[string]map[string]interface{}{}
