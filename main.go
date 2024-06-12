package main

import (
	"fmt"
	"github.com/Heartfilia/litetools/litejson"
	"github.com/Heartfilia/litetools/litenet"
	"github.com/Heartfilia/litetools/litestr"
	"github.com/Heartfilia/litetools/litetime"
	"log"
)

func testTime() {
	defer litetime.Timer("main")()
	//t := litetime.Option{
	//	Unit:   "ms",
	//	Fmt:    true,
	//	Cursor: -10,
	//}
	//fmt.Println(litetime.Time(nil).Int())
	//fmt.Println(litetime.Time(nil).String())
	//fmt.Println(litetime.Time(t).String())
	//fmt.Println("------------------------------------")
	//fmt.Println("错误的情况-->", litetime.Time(123).String())
	//fmt.Println("------------------------------------")
	fmt.Println(litetime.Time(litetime.Option{
		Goal: "2024-03-14 10:36:00",
		Unit: "ms",
	}).Int())
}

func testJson() {
	//baseJson := `{"a-x":{"b_z":[{"c":["x","y","z"]},{"d":[[3,4,5],[6,7,8]]}]}}`

	//value, _ := litejson.TryGet(baseJson, "a-x.b_z[0].c")
	//fmt.Println(value.Value())
	//value, _ = litejson.TryGet(baseJson, "a-x.b_z[0].e|a.b_z[0].d[-1][-1]")
	//value, _ = litejson.TryGet(baseJson, "a-x.b_z[1].d[-1][-1]")
	//fmt.Println(value.Int())
	//value, _ = litejson.TryGet(baseJson, "a-x.b_z[0].c[-2]")
	//fmt.Println(value.String())
	//value, err := litejson.TryGet(baseJson, "a-x.b_z[6].c[0]")
	//fmt.Println(value.Error)
	//fmt.Println(err)

	value, _ := litejson.TryGet("your path/package.json", "dependencies.abc")
	for _, v := range value.StringSlice() {
		fmt.Printf("%T --> %v\n", v, v)
	}
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
	fmt.Println(litestr.ColorString("<red>红色</red>还有其它颜色<cyan>其它颜色</cyan>还有一些错误测试<blue>错误的</yellow>", ""))
	fmt.Println(litestr.ColorString("整体替换颜色", "blue"))
	fmt.Println(litestr.ColorString("整体替换颜色", "黄"))
	log.Println(litestr.S(), "测试S的状态")
}

func testTag() {

	log.Println(litestr.D(), "测试D的状态")
	testString()
	log.Println(litestr.E(), "测试E的状态")
}

func main() {
	//testTime()
	//testJson()
	//testNet()
	//testString()
	testTag()
}
