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

	session := litereq.NewSession() // .
	//SetVerbose(true).
	//SetHTTP2(true). // 还没实现
	//SetTimeout(2000).
	//SetRetry(2).
	//SetCookies(map[string]string{"a": "1"}).                     // 全局cookie  后面单独的参数配置的cookie会融合到这里面一起请求
	//SetHeaders(map[string]string{"user-agent": "lite-tools V2", "token": "222222", "xtoken": "11111"}) // . // 兼容map格式和headers对象
	//SetProxy("http://6h65j8:mv2imgwv@61.139.65.104:61063")  // 全局代理 如果option那边传入 按那边为主

	option := reqoptions.NewOption().
		//SetMethod("POST").
		//SetVerify(false).   // 还没实现
		//SetRedirects(true). // 还没实现
		SetHeaders(map[string]string{
			"Referer":         "https://buyin.jinritemai.com/dashboard/servicehall/daren-profile",
			"Accept-Language": "zh-CN,zh;q=0.9,en;q=0.8,en-GB;q=0.7,en-US;q=0.6",
			"Accept":          "application/json, text/plain, */*",
			"User-Agent":      "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/137.0.0.0 Safari/29.36 Edg/27.0.2045.20",
		}).
		SetCookies(map[string]string{
			"buyin_shop_type_v2": "11",
			"buyin_app_id":       "13",
			"buyin_app_id_v2":    "13",
			"buyin_shop_type":    "11",
		}). // cookie兼容 字符串格式和map格式 也兼容cookie对象
		//SetParams([][2]string{
		//	{"fp", "f1118f9f298aa61ff10d3686206c7ccb25700247f7bdffd16c"},
		//	{"msToken", ""},
		//	{"page", "1"},
		//	{"uid", "v2_0a276d5f57b27638eef5cc5ad0d02c51ae9e9f970e1b04b4f14e0d60f9b3bd5fb75c63c31279fe23a61a4b0a3c9a1e037b6fe2584833ab3a330cfda4c0b8b7d40df6aa09ecc3bbaa56a4299993c3500b35de715cfc75b3a35f9f8feb4a22a58759d2207d578f5958ee10c2a2b60d18e5ade4c9012001220103a9160905"},
		//	{"verifyFp", "f1118f9f298aa61ff10d3686206c7ccb25700247f7bdffd16c"},
		//	{"with_product", "false"},
		//	{"a_bogus", "dXm0QDwfdDDkvd8g5vdLfY3qIf13Y/040SVkMDZMBn3A1y39HMOa9exYDHhvwnmjNT/dIeujy4hbYNQprQ/b8ZwfHuix/2xDmESkKl5Q5xSSs1XJtyUgnzUimktUCec2-i-lrOXMw7lHKbmg09oHmhK4bIOwu3GMyD=="},
		//}).
		SetParams("d=4&a=1&b=2&c=3").
		//SetCookieEnable(false). // 设置本次请求不使用cookie
		//SetJson(map[string]any{
		//	"test": map[string]any{
		//		"test1": 123,
		//	}}).
		//SetProxy("http://6h65j8:mv2imgwv@43.248.79.229:64060").
		SetTimeout(3000).
		SetProxy("http://127.0.0.1:9000")
	//UpdateHeaderMap(map[string]string{
	//	"referer": "abc",
	//	"token":   "1111",
	//	//"user-agent": "lite-tools V3",
	//}).ExceptGlobalHeaders([]string{"token"})
	//// 这里优先级高于Fetch里面填写的 如果两边都写了 这里和那边做融合 这里为主

	response := session.Fetch("https://buyin.jinritemai.com/api/authorStatData/authorVideoDetailList", option)

	fmt.Println(response.Text)
	fmt.Println(response.StatusCode)
	fmt.Println(response.Error())
}

func reqTest() {
	s := litereq.NewSession().
		SetHeaders(map[string]string{"user-agent": "lite-tools"}).
		SetCookies("a=1;b=2")
	r := s.Fetch("https://www.baidu.com", nil)
	fmt.Println(1, r.Cookies.String())

	o := reqoptions.NewOption().SetCookies("a=5")
	resp := s.Fetch("http://httpbin.org/get", o)
	fmt.Println(3, resp.Text)
	fmt.Println(4, resp.Cookies.String())
	fmt.Println(5, s.GetCookies().String())

}

func main() {
	//testTime()
	//testJson()
	//testNet()
	//testStr()
	//testTag()
	//testReq()
	reqTest()
}
