package main

import (
	"fmt"
	"github.com/Heartfilia/litetools/litetime"
)

func main() {
	t := litetime.Time{
		//Unit: "ms",
		Fmt:    true,
		Cursor: "-24h",
	}
	//fmt.Println(t.GetTime().Int())
	//fmt.Println(t.GetTime().Float())
	fmt.Println(t.GetTime().String())

	//stringTime = stringTime.Add(86400)
	//newString := stringTime.Format("2006-01-02 15:04:05.000")
	//fmt.Println(newString)

	//fmt.Println(litestring.ColorString("<red>红色</red>还有其它颜色<cyan>其它颜色</cyan>还有一些错误测试<blue>错误的</yellow>", ""))
	//fmt.Println(litestring.ColorString("整体替换颜色", "blue"))
	//fmt.Println(litestring.ColorString("整体替换颜色", "黄"))

}
