package main

import (
	"context"
	"fmt"
	"github.com/Heartfilia/litetools/litejs"
	"github.com/Heartfilia/litetools/litejson"
	"github.com/Heartfilia/litetools/litenet"
	"github.com/Heartfilia/litetools/litereq"
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

	//res := litestr.CookieStringToMap("x-auth-token=9e712b07dc404fe7b384e7f3dce7bbba; x-auth-app=Demo; x-auth-brand=; client_version=5.2.2.123; client_build_version=95228; client_flags=tabs; ")
	//fmt.Println(res)
	//newRes := litestr.CookieMapToString(res)
	//fmt.Println(newRes)

	//res := litestr.ParamStringToMap("a=1&b=sdfsdfsd&c=&d=jjjj")
	res := litestr.ParamStringToArray("a=1&b=sdfsdfsd&c=&d=jjjj&e=")
	fmt.Println(res)
}

func testTag() {
	log.Println(litestr.D(), "测试D的状态")
	testStr()
	log.Println(litestr.E(), "测试E的状态")
}

func FetchMsg() {

}

func testReq() {
	// 采用自定义接收方案 any位置可以放你那边任意格式的struct
	jar := litereq.NewCookieJar()
	//jar, _ := cookiejar.New(nil)  // 或者这里自定义
	var resp litereq.Wrap[any]
	_ = litereq.Build("https://www.coingecko.com/accounts/csrf_meta.json").
		UserAgent("Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/133.0.0.0 Safari/537.36 Edg/133.0.0.0").
		Referer("https://www.coingecko.com/").
		Handle(litereq.ToWrap(&resp)).
		Proxy("http://127.0.0.1:7890").
		CookieJar(jar).
		Fetch(context.Background())
	fmt.Println(resp.Data)
	fmt.Println("-->", jar)
}

func reqTest() {
	//s := litereq.NewSession().
	//	SetHeaders(map[string]string{"user-agent": "lite-tools"})
	////r := s.Fetch("https://www.baidu.com", nil)
	////fmt.Println(1, r.Cookies.String())
	//
	//o := reqoptions.NewOption().SetCookies("a=5")
	//resp := s.Fetch("http://httpbin.org/get", o)
	//fmt.Println(3, resp.Text)
	//fmt.Println(4, resp.Cookies.String())
	//fmt.Println(5, s.GetCookies().String())

}

func testJS() {
	cmd := litejs.CmdNode{
		JsPath:  "test\\rbzid.js", // 最好写绝对路径 相对路径我还没测试 晚点...
		Verbose: true,             // 如果结果不对的时候打印提示
	}
	res := cmd.Call("DrmK8giIdcVRRZq8XRNp5aRujILYxhBQEuOQc\\/q1b0nKD2IvT0P5u3TqHeTtOSocz0p2pLFb+0A\\/eMMeQX6ImMapoAurPqEPC7uAvM104ZYAfP9k1fSjJz+d0\\/EPLwepf9M5CDPG6sDXkA03wOnGs4H\\/9FtxvE5DUCxRoCaV0JWIU5L8M+ywiTilxLJnFKbFVeY+46g\\/OhdWp8WI6+\\/ynI1+1QWQwV8bFgM1sNNGUG2Jql7JpofNkQn1LQb1SHP5") // 有参数就传 没有就不管  因为是node那边调用 所以不需要传入函数名

	// 特别注意：如果js返回的是一个 对象 最好是  JSON.stringify() 包着处理一次 要不然 golang的json包会报错
	fmt.Println(string(res))
}

type qrcode struct {
	Data struct {
		Redirect string `json:"redirect_url"`
		Status   string `json:"status"`
	} `json:"data"`
}

func checkQrcode(session *litereq.Builder, token string) {
	//resp := session.
	//	Param("client_key", "aw7tduvjdk1a0x3r").
	//	Param("scope", "mobile,user_info,video.create,video.data").
	//	Param("next", "https://buyin.jinritemai.com/mpa/account/login?log_out=1").
	//	Param("state", "douyin").
	//	Param("token", token).
	//	Param("jump_type", "h5").
	//	Param("customize_params", "{\"comment_id\":\"\",\"source\":\"qr_only\",\"not_skip_confirm\":\"true\",\"open_extra\":\"\"}").
	//	Get("https://open.douyin.com/oauth/check_qrcode/")
	//if resp.Error() != nil {
	//	log.Println(resp.Error())
	//} else {
	//	fmt.Println(resp.Cookie().String())
	//	fmt.Println(resp.Text)
	//	var qr qrcode
	//	err := resp.Json(&qr)
	//	if err != nil {
	//		fmt.Println(err)
	//		return
	//	}
	//	fmt.Println("--->", qr.Data.Status)
	//	if qr.Data.Status == "confirmed" {
	//		rUrl := qr.Data.Redirect
	//		resp2 := session.Get(rUrl)
	//		fmt.Println("---------------------------------------------------")
	//		fmt.Println(resp2.Cookie().String())
	//
	//		fmt.Println("---------------------------------------------------")
	//		resp3 := session.
	//			Referer("https://buyin.jinritemai.com/mpa/account/login?log_out=1&type=24").
	//			Get("https://buyin.jinritemai.com/ecomauth/loginv1/callback?login_source=buyin&subject_aid=6077&bus_child_type=1128&entry_source=0&ecom_login_extra=&_lid=705167126569")
	//		fmt.Println(resp3.Cookie().String())
	//		fmt.Println(resp3.Text)
	//	}
	//}
}

func red(session *litereq.Builder, rUrl string) {

}

func testDy() {
	//token := "1eec6310182ed34a8c06980ca0187226_hl"
	//se := litereq.Build().Headers(
	//	map[string]string{
	//		"referer":        "https://open.douyin.com/qrconnect?container=[object%20HTMLDivElement]&client_key=aw7tduvjdk1a0x3r&response_type=code&scope=mobile,user_info,video.create,video.data&redirect_uri=https://buyin.jinritemai.com/mpa/account/login?log_out=1&type=24&jump_type=h5&state=douyin&iframeAttrs=[object%20Object]&style=//lf3-eduinfra-tos.elabstatic.com/obj/elab-statics/ecdn/file/955d04c6b61cbb98e3efeeaa41fc5db7.css",
	//		"user-agent":     "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/132.0.0.0 Safari/537.36 Edg/132.0.0.0",
	//		"sec-fetch-site": "same-origin",
	//		"sec-fetch-mode": "cors",
	//		"sec-fetch-dest": "empty",
	//	}).Cookies("passport_csrf_token=f0bc1b484f10962eec2487a3c4e3563e; passport_csrf_token_default=f0bc1b484f10962eec2487a3c4e3563e")
	//checkQrcode(se, token)
	//fmt.Println(se.GetCookies().String())
	//red(se, "https://buyin.jinritemai.com/mpa/account/login?code=2b0fbd0dfbbd5efdsS0AXkGy1ie0cLIOUwlG_hl\\u0026log_out=1\\u0026state=douyin")
}

func main() {
	//testTime()
	//testJson()
	//testNet()
	//testStr()
	//testTag()
	testReq()
	//reqTest()
	//testJS()
	//testDy()
}
