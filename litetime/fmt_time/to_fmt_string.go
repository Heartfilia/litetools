package fmt_time

import (
	"strings"
	"time"
)

func NowFmt(cursor int) string {
	// 第一种是直接获取当前时间的 格式化时间的
	if cursor == 0 {
		stringTime := time.Now().String()
		return strings.Split(stringTime, ".")[0]
	}

	return ""

}

func FmtType(fStr string, cursor int) string {
	// 这里是获取格式化时间的地方

	return ""
}
