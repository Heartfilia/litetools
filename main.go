package main

import (
	"fmt"
	"github.com/Heartfilia/litetools/litejson"
	"github.com/Heartfilia/litetools/litenet"
	"github.com/Heartfilia/litetools/litestring"
	"github.com/Heartfilia/litetools/litetime"
)

func testTime() {
	defer litetime.Timer("main")()
	t := litetime.Time{
		Unit:   "ms",
		Fmt:    true,
		Cursor: -10,
	}
	fmt.Println(t.GetTime().Float())
	fmt.Println(t.GetTime().String())
}

func testJson() {
	baseJson := `{"a":{"b":[{"c":["x","y","z"]},{"d":[[3,4,5],[6,7,8]]}]}}`

	value, _ := litejson.TryGet(baseJson, "a.b[0].c")
	fmt.Println(value.Value)
	value, _ = litejson.TryGet(baseJson, "a.b[0].e|a.b[0].d[-1][-1]")
	value, _ = litejson.TryGet(baseJson, "a.b[1].d[-1][-1]")
	fmt.Println(value.Int())
	value, _ = litejson.TryGet(baseJson, "a.b[0].c[-2]")
	fmt.Println(value.String())
	value, _ = litejson.TryGet(baseJson, "a.b[6].c[-5]")
	fmt.Println(value.Error)
}

func testNet() {
	fmt.Println(litenet.GetLAN())
	fmt.Println(litenet.GetWAN())
	fmt.Println(litenet.GetUA())
	fmt.Println(litenet.GetUA("pc"))
	fmt.Println(litenet.GetUA("mobile"))
	fmt.Println(litenet.GetUA("mac"))
	fmt.Println(litenet.GetUA("ios", "android"))
	fmt.Println(litenet.GetUA("chrome", "edge", "opera"))
	fmt.Println(litenet.GetUA("chrome", "edge"))
}

func testString() {
	fmt.Println(litestring.ColorString("<red>红色</red>还有其它颜色<cyan>其它颜色</cyan>还有一些错误测试<blue>错误的</yellow>", ""))
	fmt.Println(litestring.ColorString("整体替换颜色", "blue"))
	fmt.Println(litestring.ColorString("整体替换颜色", "黄"))
}

func main() {
	//testTime()
	testJson()
	//testNet()
	//testString()
}
