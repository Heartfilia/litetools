package main

import "github.com/Heartfilia/litetools/liteparser"

func main() {
	_, err := liteparser.TryGet(`{"a":{"b":[1,2,3]}}`, "a.b[0][1].c[2]")
	//_, err = liteparser.TryGet(`{"a":{"b":[1,2,3]}}`, "a.b[0].c")
	if err != nil {
		return
	}
	//fmt.Println(litedir.LiteDir())
	//fmt.Println(litedir.FileJsonLoader("/Users/lodge/Library/Caches/lite-tools/browser/config.json"))
	//fmt.Println(litenet.GetLAN())
	//fmt.Println(litenet.GetWAN())
	//fmt.Println(litenet.GetUA())
	//fmt.Println(litenet.GetUA("pc"))
	//fmt.Println(litenet.GetUA("mobile"))
	//fmt.Println(litenet.GetUA("mac"))
	//fmt.Println(litenet.GetUA("ios", "android"))
	//fmt.Println(litenet.GetUA("chrome", "edge", "opera"))
	//fmt.Println(litenet.GetUA("chrome", "edge"))
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
