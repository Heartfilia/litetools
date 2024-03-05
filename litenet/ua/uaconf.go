package ua

import (
	"encoding/json"
	"fmt"
	"github.com/Heartfilia/litetools/litenet/request"
	"github.com/Heartfilia/litetools/literand"
	"github.com/Heartfilia/litetools/utils/litedir"
	"github.com/Heartfilia/litetools/utils/types"
	"os"
	"path"
	"strings"
)

const DefaultChoice = "chrome"

var Browser = []string{"chrome", "firefox", "opera", "ie", "edge", "safari"}
var System = []string{"pc", "mobile", "mac", "win", "windows", "linux", "android", "ios", "harmony", "harmonyos"}

var uaTemplateBrowser = map[string]string{
	"chrome":  "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/%s Safari/537.36",
	"firefox": "Mozilla/5.0 (Windows NT 10.0; Win64; x64; rv:%s.0) Gecko/20100101 Firefox/%s.0",
	"opera":   "Opera/9.80 (Macintosh; Intel Mac OS X 10.6.8; U; en) Presto/2.8.131 Version/11.11",
	"ie":      "Mozilla/5.0 (Windows NT 10.0; WOW64; Trident/7.0; rv:11.0) like Gecko",
	"edge":    "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/%s Safari/537.36 Edg/%s",
	"safari":  "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/605.1.15 (KHTML, like Gecko) Version/%s Safari/605.1.15",
}

var defaultSetting = types.ConfigJson{
	Chromium: []string{
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
	Firefox: []string{
		"70", "71", "72", "73", "74", "75", "76", "77", "78", "79",
		"80", "81", "82", "83", "84", "85", "86", "87", "88", "89",
		"90", "91", "92", "93", "94", "95", "96", "97", "98", "99",
		"100", "101", "102", "103", "104", "105", "106", "107", "108", "109",
		"110", "111", "112", "113", "114", "115", "116", "117", "118", "119",
		"120", "121", "122", "123", "124", "125", "126", "127", "128", "129",
	},
	Safari: []string{"10.3.1", "10.6.8", "10.9.2", "10.15.7", "12.1.2", "13.1.1", "14.1.2"},
}
var localCache types.ConfigJson

func sureLocal(configJson string) bool {
	// 确保本地有缓存记录
	if litedir.FileExists(configJson) {
		return true //存在返回存在的状态
	}
	// 不存在这里下载
	// 缓存在这里 后续直接用
	// https://googlechromelabs.github.io/chrome-for-testing/known-good-versions.json
	requestJson := request.DoGet("http://static.litetools.top/source/json/useragent.json")
	if requestJson == "" || !litedir.FileExists(configJson) {
		// 这里下载失败 直接返回false
		return false
	}
	return true
}

func readFromLocal(configJson string) types.ConfigJson {
	file, err := os.ReadFile(configJson)
	if err != nil {
		return types.ConfigJson{}
	}
	var res types.ConfigJson
	err = json.Unmarshal(file, &res)
	if err != nil {
		return types.ConfigJson{}
	}
	return res
}

func configFromCache() types.ConfigJson {
	if !localCache.IsEmpty() {
		// 直接返回缓存好的数据 避免重复加载io
		return localCache
	}
	baseDir := litedir.LiteDir()
	browserDir := path.Join(baseDir, "browser")
	if !litedir.FileExists(browserDir) {
		_ = os.Mkdir(browserDir, 0777)
	}
	configJson := path.Join(browserDir, "config.json")
	local := sureLocal(configJson)
	if local {
		// 本地有 然后也做对应的处理
		res := readFromLocal(configJson)
		if !res.IsEmpty() {
			localCache = res
			return res
		}
	}
	return defaultSetting
}

func CombineString(platform, browser string) string {
	// platform   是哪个终端的系统
	// browser    是哪个浏览器
	version := configFromCache() // 获取浏览器版本合集
	// Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/%s Safari/537.36 Edg/%s
	// headString (osString1; osString2; osString3) kitString1 TailStringWithVersion otherPlug
	headString := "Mozilla/5.0"
	operaOld := false
	if browser == "opera" {
		useOld := literand.RandomChoice([]bool{true, false})
		if useOld {
			headString = "Opera/9.80" // 这个是老引擎的 新版的可以是  Mozilla/5.0
			operaOld = true
		}
	}
	osString1 := "Windows NT " // 后面还需要补version
	osString2 := "Win64"
	osString3 := "; x64"
	if platform == "linux" {
		osString1 = "X11"
		osString2 = "Linux x86_64"
		osString3 = ""
	} else if platform == "mac" {
		osString1 = "Macintosh"
		safariVersion := literand.RandomChoice(version.Safari)
		newSafariVersion := strings.ReplaceAll(safariVersion, ".", "_")
		osString2 = "Intel Mac OS X " + newSafariVersion
		osString3 = literand.RandomChoice([]string{"", "; U; en", "; en-us", "; zh-hans", "; Eu; fr", "; Eu; De"})
	} else if platform == "harmony" {
		// Mozilla/5.0 (Linux; Android 10; K) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/121.0.0.0 Mobile Safari/537.36 EdgA/121.0.0.0
		osString1 = "Linux"
		osString2 = "Android " + literand.RandomChoice([]string{"9", "10", "11", "12"})
		osString3 = literand.RandomChoice([]string{"", "; K", "; L", "; M", "; N", "; O", "; P", "; Q", "; R", "; S", "; T", "; U"})
	} else if platform == "ios" {
		osString1 = "iPhone" + literand.RandomChoice([]string{"", "6s", "8 plus", "x", "11", "12", "13", "14", "15", "15 pro", "14 pro", "13 pro"})
		safariVersion := literand.RandomChoice(version.Safari)
		newSafariVersion := strings.ReplaceAll(safariVersion, ".", "_")
		osString2 = "CPU iPhone OS " + newSafariVersion + " like Mac OS X"
		osString3 = ""
	} else if platform == "android" {
		osString1 = "Linux"
		osString2 = "Android " + literand.RandomChoice([]string{"9.0", "10.0", "11.0", "12.0", "10.1", "10.3", "11.2", "12.1"})
		osString3 = literand.RandomChoice([]string{
			"; zh-cn; BLA-AL00 Build/HUAWEIBLA-AL00",
			"; PAR-AL00 Build/HUAWEIPAR-AL00; wv",
			"; OPPO A57 Build/MMB29M; wv",
			"; EML-AL00 Build/HUAWEIEML-AL00; wv",
			"; DUK-AL20 Build/HUAWEIDUK-AL20; wv",
			"; zh-CN; EML-AL00 Build/HUAWEIEML-AL00",
			"; zh-CN; SM-J3109 Build/LMY47X",
			"; GT1u Build/PI",
			"; ZH960",
			"; SAMSUNG SM-T825Y",
			"; SM-T825Y",
			"; MX10 PRO",
			"; TX3 Build/PPR1.180610.011",
			"; OPPO A73 Build/N6F26Q",
			"; vivo X20Plus A Build/NMF26X; wv",
			"; zh-cn; HUAWEI CAZ-AL10 Build/HUAWEICAZ-AL10",
			"; en-US; SM-G950F Build/PPR1.180610.011",
			"; zh-cn; RVL-AL09 Build/HUAWEIRVL-AL09",
			"; zh-cn; PDBM00 Build/PPR1.180610.011",
			"; INE-AL00; HMSCore 6.1.0.313; GMSCore 19.6.29",
			"; V1913A Build/P00610; wv",
			"; zh-cn; MI 6X Build/PKQ1.180904.001",
			"; COL-AL10; HMSCore 6.1.0.305; GMSCore 17.7.85",
		})
	} else {
		osString1 += literand.RandomChoice([]string{"10.0", "11.0"})
	}

	chromeVersion := literand.RandomChoice(version.Chromium)

	kitString := "AppleWebKit/537.36"
	if browser == "firefox" {
		kitString = "Gecko/20100101"
	} else if platform == "mac" {
		kitString = "AppleWebKit/603.1.30"
	}
	tailStringWithVersion := "(KHTML, like Gecko)"
	if browser == "chrome" || browser == "edge" {
		tailStringWithVersion += " Chrome/" + chromeVersion
	}

	otherPlug := "Safari/537.36"
	if platform == "mac" && browser != "firefox" {
		otherPlug = literand.RandomChoice([]string{"Safari/537.36", "Safari/605.1.15", "Safari/604.3.5"})
	} else if browser == "firefox" {
		otherPlug = "Firefox/" + literand.RandomChoice(version.Firefox)
		tailStringWithVersion = ""
	}

	if platform == "android" || platform == "ios" || platform == "harmony" {
		tailStringWithVersion += " Mobile"
	}

	if !operaOld && browser == "opera" {
		otherPlug += " OPR/" + chromeVersion
	} else if browser == "opera" {
		otherPlug += " Presto/2.1." + literand.RandomChoice([]string{"0", "1", "2", "3", "4", "5", "6", "7", "8", "9", "10", "11", "12", "13", "14", "15", "16", "17", "18", "19", "20"})
	}
	if browser == "edge" {
		otherPlug += " Edg/" + chromeVersion
	}

	return fmt.Sprintf("%s (%s; %s%s) %s %s %s", headString, osString1, osString2, osString3, kitString, tailStringWithVersion, otherPlug)
}
