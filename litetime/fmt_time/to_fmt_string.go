package fmt_time

import (
	"fmt"
	"strings"
	"time"
)

const (
	Y string = "2006"
	y string = "06"
	m string = "01"
	d string = "02"
	H string = "15"
	M string = "04"
	S string = "05"
	a string = "Mon"
	A string = "Monday"
	b string = "Jan"
	B string = "January"
)

func getFormat(format *string) {
	for _, rule := range []string{"%Y", "%m", "%d", "%H", "%M", "%S"} {
		compare := strings.Index(*format, rule)
		if compare == -1 {
			continue
		}
		var realTime string
		switch rule {
		case "%Y":
			realTime = Y
		case "%y":
			realTime = y
		case "%m":
			realTime = m
		case "%d":
			realTime = d
		case "%H":
			realTime = H
		case "%M":
			realTime = M
		case "%S":
			realTime = S
		}
		if realTime != "" {
			*format = strings.Replace(*format, rule, realTime, -1)
		}
	}
}

func NowFmt(cursor int) string {
	// 第一种是直接获取当前时间的 格式化时间的
	if cursor == 0 {
		stringTime := time.Now().String()
		return strings.Split(stringTime, ".")[0]
	}
	x := time.Unix(86400, 0)
	fmt.Println(x)

	return ""

}

func FmtType(fStr string, cursor int) string {
	// 这里是获取格式化时间的地方

	return ""
}
