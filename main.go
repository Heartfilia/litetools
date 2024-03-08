package main

import (
	"fmt"
	"github.com/Heartfilia/litetools/litejson"
)

func main() {
	value := litejson.TryGet(`{"a":{"b":[{"c":[0,1,2],"d":["x","y","z"]}]}}`, "a.b[1].d[-1]")
	fmt.Printf("%T, %v\n", value.Value, value.Value)
	fmt.Println(value.Error)
	value = litejson.TryGet(`{"a":{"b":[{"c":[0,1,2]},{"d":["x","y","z"]}]}}`, "a.b[0].d[-1]")
	fmt.Printf("%T, %v\n", value.Value, value.Int())
	fmt.Println(value.Error)
	value = litejson.TryGet(`"a":[]`, "a.b[0].d[-1]")
	fmt.Printf("%T, %v\n", value.Value, value.String())
	fmt.Println(value.Error)
	value = litejson.TryGet(`{"a":{"b":[{"c":[4,5,6],"d":["x","y","z"]}]}}`, "a.b[0].c[-1]")
	fmt.Printf("%T, %v\n", value.Value, value.Int())
	fmt.Println(value.Error)
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
	//	Unit:   "ms",
	//	Fmt:    true,
	//	Cursor: -10,
	//}
	//fmt.Println(t.GetTime().Int()) // 1704854601
	//fmt.Println(t.GetTime().Float())
	//fmt.Println(t.GetTime().String())

	//fmt.Println(litestring.ColorString("<red>红色</red>还有其它颜色<cyan>其它颜色</cyan>还有一些错误测试<blue>错误的</yellow>", ""))
	//fmt.Println(litestring.ColorString("整体替换颜色", "blue"))
	//fmt.Println(litestring.ColorString("整体替换颜色", "黄"))

}
