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
	"%-m": "1",
	"%d":  "02",
	"%-d": "2",
	"%H":  "15",
	"%h":  "3", // 这个不太确定
	"%M":  "04",
	"%-M": "4",
	"%S":  "05",
	"%-S": "5",
	"%f":  "000",
	//"%F":   "999",
	"%.3f": "000",
	"%.6f": "000000",
	"%.9f": "000000000",
	//"%Z":   "MST",
	"%z":    "Z0700",  // +0800
	"%-z":   "-07",    // +08
	"%z:00": "Z07:00", // +08:00
}

func GetFormat(format string) string {
	for rule, value := range fmtMap {
		compare := strings.Index(format, rule)
		if compare == -1 {
			continue
		}
		format = strings.Replace(format, rule, value, -1)
	}
	return format
}

func nowFmt(cursor, formatSample string) string {
	nowS := time.Now()
	if cursor != "0h" {
		t, _ := time.ParseDuration(cursor)
		nowS = nowS.Add(t)
	}
	//return time.Unix(nowS.Unix(), 0).Format(GetFormat(formatSample))
	return nowS.Format(GetFormat(formatSample))

}

func FmtType(fStr string, cursor string) string {
	// 这里是获取格式化时间的地方
	if fStr == "" {
		fStr = "%Y-%m-%d %H:%M:%S"
	}
	return nowFmt(cursor, fStr)
}
