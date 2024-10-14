package main

import (
	"fmt"
	"github.com/Heartfilia/litetools/litejson"
	"github.com/Heartfilia/litetools/litenet"
	"github.com/Heartfilia/litetools/litereq"
	"github.com/Heartfilia/litetools/litereq/reqoptions"
	"github.com/Heartfilia/litetools/litestr"
	"github.com/Heartfilia/litetools/litetime"
	"log"
)

func testTime() {
	defer litetime.Timer()()
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

	cmd := litejson.TryGet("your path/package.json", "dependencies.abc")
	for _, v := range cmd.StringSlice() {
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

func testStr() {
	//fmt.Println(litestr.ColorString("<red>红色</red>还有其它颜色<cyan>其它颜色</cyan>还有一些错误测试<blue>错误的</yellow>", ""))
	//fmt.Println(litestr.ColorString("整体替换颜色", "blue"))
	//fmt.Println(litestr.ColorString("整体替换颜色", "黄"))
	//log.Println(litestr.S(), "测试S的状态")

	res := litestr.CookieStringToMap("x-auth-token=9e712b07dc404fe7b384e7f3dce7bbba; x-auth-app=Demo; x-auth-brand=; client_version=5.2.2.123; client_build_version=95228; client_flags=tabs")
	fmt.Println(res)
	newRes := litestr.CookieMapToString(res)
	fmt.Println(newRes)
}

func testTag() {
	log.Println(litestr.D(), "测试D的状态")
	testStr()
	log.Println(litestr.E(), "测试E的状态")
}

func testReq() {
	// 数据安全 加锁 啥的 后面整体流程实现了之后再去处理
	defer litetime.Timer()()

	session := litereq.NewSession().
		SetVerbose(true).
		//SetHTTP2(true). // 还没实现
		//SetTimeout(2000).
		//SetRetry(2).
		//SetCookies(map[string]string{"a": "1"}).                     // 全局cookie  后面单独的参数配置的cookie会融合到这里面一起请求
		SetHeaders(map[string]string{"user-agent": "lite-tools V2", "token": "222222", "xtoken": "11111"}) // . // 兼容map格式和headers对象
	//SetProxy("http://6h65j8:mv2imgwv@61.139.65.104:61063")  // 全局代理 如果option那边传入 按那边为主

	option := reqoptions.NewOption().
		//SetMethod("POST").
		SetVerify(false).   // 还没实现
		SetRedirects(true). // 还没实现
		//SetHeaders(map[string]string{"user-agent": "from option"}).
		//SetCookies(map[string]string{"b": "", "d": "666"}). // cookie兼容 字符串格式和map格式 也兼容cookie对象
		//SetParams([][2]any{{"k", 1}, {"v", "2"}}).
		SetCookieEnable(false). // 设置本次请求不使用cookie
		//SetJson(map[string]any{
		//	"test": map[string]any{
		//		"test1": 123,
		//	}}).
		//SetProxy("http://6h65j8:mv2imgwv@43.248.79.229:64060").
		SetTimeout(3000).UpdateHeaderMap(map[string]string{
		"referer": "abc",
		"token":   "1111",
		//"user-agent": "lite-tools V3",
	}).ExceptGlobalHeaders([]string{"token"})
	//// 这里优先级高于Fetch里面填写的 如果两边都写了 这里和那边做融合 这里为主

	//response := session.Fetch("http://httpbin.org/post", option)
	response := session.Fetch("http://httpbin.org/redirect/1", option)

	fmt.Println(response.Text)
	fmt.Println(response.StatusCode)
	fmt.Println(response.Error())
}

func main() {
	//testTime()
	//testJson()
	//testNet()
	//testStr()
	//testTag()
	testReq()
}
