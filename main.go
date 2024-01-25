package main

import (
	"fmt"
	"github.com/Heartfilia/litetools/litetime"
)

func main() {
	t := litetime.Time{
		//Unit: "ms",
		Fmt: true,
	}
	fmt.Println(t.GetTime().Int())
	fmt.Println(t.GetTime().Float())
	fmt.Println(t.GetTime().String())

	//stringTime := time.Now().String()
	//newString := strings.Split(stringTime, ".")
	//fmt.Println(newString[0])

	//fmt.Println(litestring.ColorString("<red>红色</red>还有其它颜色<cyan>其它颜色</cyan>还有一些错误测试<blue>错误的</yellow>", ""))
	//fmt.Println(litestring.ColorString("整体替换颜色", "blue"))
	//fmt.Println(litestring.ColorString("整体替换颜色", "黄"))

}
