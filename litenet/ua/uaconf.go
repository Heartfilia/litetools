package ua

import (
	"encoding/json"
	"github.com/Heartfilia/litetools/litenet/request"
	"github.com/Heartfilia/litetools/utils/litedir"
	"os"
	"path"
)

const DefaultChoice = "chrome"

var Browser = []string{"chrome", "firefox", "opera", "ie", "edge", "safari"}
var System = []string{"pc", "mobile", "mac", "win", "windows", "linux", "android", "ios", "harmony", "harmonyos"}

var UATemplateBrowser = map[string]string{
	"chrome":  "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/%s Safari/537.36",
	"firefox": "Mozilla/5.0 (Windows NT 10.0; Win64; x64; rv:%s.0) Gecko/20100101 Firefox/%s.0",
	"opera":   "Opera/9.80 (Macintosh; Intel Mac OS X 10.6.8; U; en) Presto/2.8.131 Version/11.11",
	"ie":      "Mozilla/5.0 (Windows NT 10.0; WOW64; Trident/7.0; rv:11.0) like Gecko",
	"edge":    "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/%s Safari/537.36 Edg/%s",
	"safari":  "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/605.1.15 (KHTML, like Gecko) Version/%s Safari/605.1.15",
}

var DefaultSetting = map[string][]string{
	"chromium": {
		"70.0.3538.16", "70.0.3538.67", "70.0.3538.97",
		"71.0.3578.137", "71.0.3578.30", "71.0.3578.33", "71.0.3578.80",
		"72.0.3626.69", "72.0.3626.7",
		"73.0.3683.20", "73.0.3683.68",
		"74.0.3729.6",
		"75.0.3770.140", "75.0.3770.8", "75.0.3770.90",
		"76.0.3809.12", "76.0.3809.126", "76.0.3809.25", "76.0.3809.68",
		"77.0.3865.10", "77.0.3865.40",
		"78.0.3904.105", "78.0.3904.11", "78.0.3904.70",
		"79.0.3945.16", "79.0.3945.36",
		"80.0.3987.106", "80.0.3987.16",
		"81.0.4044.138", "81.0.4044.20", "81.0.4044.69",
		"83.0.4103.14", "83.0.4103.39",
		"84.0.4147.30",
		"85.0.4183.38", "85.0.4183.83", "85.0.4183.87",
		"86.0.4240.22",
		"87.0.4280.20", "87.0.4280.87", "87.0.4280.88",
		"88.0.4324.27", "88.0.4324.96",
		"89.0.4389.23",
		"90.0.4430.24",
		"91.0.4472.101", "91.0.4472.19",
		"92.0.4515.107", "92.0.4515.43",
		"93.0.4577.15", "93.0.4577.63",
		"94.0.4606.113", "94.0.4606.41", "94.0.4606.61",
		"95.0.4638.10", "95.0.4638.17", "95.0.4638.54", "95.0.4638.69",
		"96.0.4664.18", "96.0.4664.35", "96.0.4664.45", "96.0.1054.53",
		"97.0.4692.20", "97.0.4692.36", "97.0.4692.71",
		"98.0.4758.48", "98.0.4758.80", "98.0.4758.102",
		"99.0.4844.17", "99.0.4844.35", "99.0.4844.51",
		"100.0.4896.20", "100.0.4896.60",
		"101.0.4951.15", "101.0.4951.41",
		"102.0.5005.27", "102.0.5005.61",
		"103.0.5060.24", "103.0.5060.53", "103.0.5060.134",
		"104.0.5112.20", "104.0.5112.29", "104.0.5112.79", "104.0.5112.81",
		"105.0.5195.19", "105.0.5195.52",
		"106.0.5249.21", "106.0.5249.61",
		"107.0.5304.18", "107.0.5304.62",
		"108.0.5359.22", "108.0.5359.71",
		"109.0.5414.25", "109.0.5414.74",
		"110.0.5481.30", "110.0.5481.77",
		"111.0.5563.19", "111.0.5563.41", "111.0.5563.64",
		"112.0.5615.28", "112.0.5615.49",
		"113.0.5672.24", "113.0.5672.63",
		"114.0.5735.16",
	},
	"firefox": {
		"70", "71", "72", "73", "74", "75", "76", "77", "78", "79",
		"80", "81", "82", "83", "84", "85", "86", "87", "88", "89",
		"90", "91", "92", "93", "94", "95", "96", "97", "98", "99",
		"100", "101", "102", "103", "104", "105", "106", "107", "108", "109",
		"110", "111", "112", "113", "114", "115", "116", "117", "118", "119",
		"120", "121", "122", "123", "124", "125", "126", "127", "128", "129",
	},
}
var requestJson = ""

func ConfigFromCache() map[string][]string {
	baseDir := litedir.LiteDir()
	browserDir := path.Join(baseDir, "browser")
	if !litedir.FileExists(browserDir) {
		_ = os.Mkdir(browserDir, 0777)
	}
	configJson := path.Join(browserDir, "config.json")
	var result map[string][]string
	if !litedir.FileExists(configJson) {
		// 如果不存在 那么就联网下载到本地
		if requestJson == "" {
			// 缓存在这里 后续直接用
			// https://googlechromelabs.github.io/chrome-for-testing/known-good-versions.json
			requestJson = request.DoGet("http://static.litetools.top/source/json/useragent.json")
		}

		if requestJson != "" {
			// 下载成功 先缓存到本地
			litedir.FileSaver(requestJson, configJson)
			err := json.Unmarshal([]byte(requestJson), &result)
			if err != nil {
				result = DefaultSetting
			}
		}

	} else {
		// 如果有 那么就读取本地记录
		tempResult := litedir.FileJsonLoader(configJson)
		if tempResult != nil {
			result = tempResult
		}
	}

	if result != nil {
		return result
	}

	return DefaultSetting
}
