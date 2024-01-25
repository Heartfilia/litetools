package fmt_time

import (
	"strings"
	"time"
)

var fmtMap = map[string]string{
	// 暂时兼容这些
	"%a":  "Mon",
	"%A":  "Monday",
	"%b":  "Jan",
	"%B":  "January",
	"%Y":  "2006",
	"%y":  "06",
	"%m":  "01",
	"%d":  "02",
	"%-d": "2",
	"%H":  "15",
	"%-H": "3",
	"%M":  "04",
	"%-M": "4",
	"%S":  "05",
	"%-S": "5",
}

func getFormat(format string) string {
	for rule, value := range fmtMap {
		compare := strings.Index(format, rule)
		if compare == -1 {
			continue
		}
		format = strings.Replace(format, rule, value, -1)
	}
	return format
}

func NowFmt(cursor int, cursorUnit int64) string {
	// 第一种是直接获取当前时间的 格式化时间的
	if cursor == 0 {
		stringTime := time.Now().String()
		return strings.Split(stringTime, ".")[0]
	}
	// 这里需要对 时间 进行二次转换 然后再变成格式化时间
	nowS := time.Now().Unix() + int64(cursor)*cursorUnit

	return time.Unix(nowS, 0).Format(getFormat("%Y-%m-%d %H:%M:%S"))

}

func FmtType(fStr string, cursor int, cursorUnit int64) string {
	// 这里是获取格式化时间的地方

	return ""
}
