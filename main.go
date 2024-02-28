package main

import (
	"fmt"
	"github.com/Heartfilia/litetools/litenet"
)

func main() {
	fmt.Println(litenet.GetLAN())
	//fmt.Println(litenet.GetWAN())
	//fmt.Println(litenet.GetUA())
	//fmt.Println(litenet.GetUA("mac", "ios", "pc", "linux", "windows", "chrome"))
	//t := litetime.Time{
	//Unit: "ms",
	//Fmt: true,
	//Cursor: -10,
	//}
	//fmt.Println(t.GetTime().Int()) // 1704854601
	//fmt.Println(t.GetTime().Float())
	//fmt.Println(t.GetTime().String())

	//fmt.Println(litestring.ColorString("<red>红色</red>还有其它颜色<cyan>其它颜色</cyan>还有一些错误测试<blue>错误的</yellow>", ""))
	//fmt.Println(litestring.ColorString("整体替换颜色", "blue"))
	//fmt.Println(litestring.ColorString("整体替换颜色", "黄"))

}
