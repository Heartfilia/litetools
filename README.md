# lite-tools 同名go包

拉取包
```bash
go get -u github.com/Heartfilia/litetools
```

## 当前功能

### time
```go
package main

import (
    "fmt"
    "github.com/Heartfilia/litetools/litetime"
)

func main(){
	defer litetime.Timer()()   // 用于统计该函数运行耗时
	
    t := litetime.Option{
        //Goal: "2024-01-10 10:43:21", // 如果不传 所有操作基于当前时间 传了字符串 那么基于字符串所示时间处理 不过字符串得对应下面的格式化样式
        // 如果传入了时间戳 基于时间戳处理
        //Unit: "ms",    // 如果是需要获取数字相关的时间 这里配置秒或者毫秒
        //Fmt: true,       // 如果需要格式化时间 这里设置为 true  或者 转入需要的格式化时间样式
        //Cursor: "-1d2h",  // 游标 基于goal或者当前时间进行数据操作 传字符串为精细单位处理
        //Cursor: 2,          // 游标 传数字就为天的单位
    }  // 如果里面什么都不写 默认就是获取当前时间
    
    /*
    可以组合的搭配示例
    */
    t := litetime.Option{}
    t := litetime.Option{
        Goal: "2024-01-10 10:43:21",
    }
    t := litetime.Option{
        Goal: "2024-01-10 10:43:21",
        Uint: "ms",
    }
    t := litetime.Option{
        Uint: "ms",
        Cursor: "-1d3h",
    }
    t := litetime.Option{
        Cursor: 10,
        Fmt: true,
    }
    t := litetime.Option{
        Unit: "ms",
        Fmt: false,
    }
    t := litetime.Option{
        Goal: 1704768201,
    }
    t := litetime.Option{
        Goal: 1704768201123,
        Fmt: "%Y-%m-%d",
    }
    t := litetime.Option{
        Goal: 1704768201123,
        Cursor: -1,
    }
    // 注意得结合预期获得结果 要不然获得的是对应类型的 零值
	defer litetime.Timer()()
	
	fmt.Println(litetime.Time(nil).Int())     // 默认就是获取 13 位的时间戳（毫秒）
	fmt.Println(litetime.Time(nil).String())  // 默认就是 %Y-%m-%d %H:%M:%S 的格式
	fmt.Println(litetime.Time(t).String())    // 传入 litetime.Time{} 结构体  可以自定义输出
	fmt.Println("------------------------------------")
	fmt.Println("错误的情况-->", litetime.Time(123).String())    // 如果传入的不是 nil 或者 litetime.Option{}  会拿不到结果 得到对应的零值
	fmt.Println("------------------------------------")
}
```

### string

```go
package main

import (
	"fmt"
	"github.com/Heartfilia/litetools/litestr"
	"log"
)

func main() {
	// 第二个位置不传入或者传入 "" 均表示 采用标签语法
	fmt.Println(litestr.ColorString("<red>红色</red>还有其它颜色<cyan>其它颜色</cyan>还有一些错误测试<blue>错误的</yellow>"))
	fmt.Println(litestr.ColorString("<red>红色</red>还有其它颜色<cyan>其它颜色</cyan>还有一些错误测试<blue>错误的</yellow>", ""))
	// 如果第一个位置写的标签语法 第二个位置又写上了颜色 以第二个位置为准....
	fmt.Println(litestr.ColorString("整体替换颜色", "blue"))
	fmt.Println(litestr.ColorString("整体替换颜色", "黄"))
	// 如果第二个位置写超过 1个位置的颜色 只采用最先出现的颜色

	// 新增 日志标签 建议配合 log 使用 如下 [D I S W E] -> [debug info success warning error]
	log.Println(litestr.D(), "这里再写自己的日志")
	
	log.Printf("%s %s", litestr.E(), "如果是f的话需要这样子写 要不然有换行异常\n")
	log.Printf(litestr.E() + " %s", "如果是f的话需要这样子写 要不然有换行异常\n")
}
```


### net
```go
package main

import (
	"fmt"
	"github.com/Heartfilia/litetools/litenet"
)

func main() {
	fmt.Println(litenet.GetUA()) // 不输入默认chrome
	fmt.Println(litenet.GetUA("safari")) // 指定浏览器 或者系统
	fmt.Println(litenet.GetUA("safari", "chrome", "linux")) // 从给定的参数里面随机
	
	fmt.Println(litenet.GetWAN())  // 
	fmt.Println(litenet.GetLAN())  // 
}


```

### slice
```go
import (
	"fmt"
	"github.com/Heartfilia/litetools/liteslice"
)

func main(){
	fmt.Println(liteslice.RandomChoice([]string{"a", "b", "c"}))
	fmt.Println(liteslice.RandomChoice([]int64{1, 2, 3}))
	fmt.Println(liteslice.RandomChoice([]float64{1.123, 2.223, 3.333}))
	
	fmt.Println(liteslice.SliceRemove([]string{"a", "b", "a", "c"}, "a"))
	fmt.Println(liteslice.SliceRemove([]int{1, 2, 3, 4, 3}, 3))
}
```

### json
```go
import (
    "fmt"
    "github.com/Heartfilia/litetools/litejson"
)

func main(){
    fmt.Println(litejson.TrgGet(`jsonStringHere`, "ruleHere"))

    baseJson := `{"a-x":{"b_z":[{"c":["x","y","z"]},{"d":[[3,4,5],[6,7,8]]}]}}`
    
    cmd := litejson.TryGet(baseJson, "a-x.b_z[0].c")  // 直接用 . 提取
    fmt.Println(cmd.Value())                                 // 不确定格式的 可以用 Value 取值
    cmd = litejson.TryGet(baseJson, "a-x.b_z[0].e|a.b_z[0].d[-1][-1]")  // 可以用 | 来分割多个rule
    cmd = litejson.TryGet(baseJson, "a-x.b_z[1].d[-1][-1]")             // 可以支持golang不支持的 负数的值的提取
    fmt.Println(cmd.Int())                                 // 确定格式的可以指定某个输出格式
    cmd = litejson.TryGet(baseJson, "a-x.b_z[0].c[-2]")
    fmt.Println(cmd.String())                              // String 是任意类型都可以转成string样式
    cmd := litejson.TryGet(baseJson, "a-x.b_z[6].c[0]")             // 错误的提取可以从两个地方提取
    fmt.Println(cmd.Error())
    fmt.Println(cmd.Err) 
	
    // 支持了第一个位置填入json文件的路径
    cmdOption := litejson.TryGet("your path/package.json", "dependencies.abc")
    // 支持了常用的 一维数组 结果
	for _, v := range cmdOption.StringSlice() {
        fmt.Printf("%T --> %v\n", v, v)
    }
}

```

session
```go
import (
	"fmt"
    "github.com/Heartfilia/litetools/litereq"      // 核心请求包
	"github.com/Heartfilia/litetools/litereq/opt"  // 参数配置包
)

// 下面不同案例的参数配置是一样的
// SetCookies: 支持map[string]string  / string: key=value  // *net.Cookie
// SetHeaders: 支持map[string]string  / net.Header
//
// SetParams : 支持 key=value&key=value / [][2]any{{key, value}} / [][2]string{{key, value}} / netURL.Values / map[string]any / map[string]string
/*
推荐写法
*/
func req(){
    session := litereq.NewSession().
        SetVerbose(true).
        SetHTTP2(true).     // 还没实现
        SetTimeout(2000).   // 全局超时，option有配置以option为主
        SetRetry(2).
        SetCookies(map[string]string{"a": "1"}).                       // 全局cookies  后面单独的参数配置的cookie会融合到这里面一起请求
        SetHeaders(map[string]string{"user-agent": "lite-tools V2"}).  // 全局headers  兼容map格式和headers对象
        SetProxy("http://user:pass@ip:port")         // 全局代理 如果option那边传入 按那边为主 
	
	option := opt.NewOption().
        SetMethod("GET").   // 需要传body之类的方法还没实现 <<<<<<<<
        SetVerify(false).   // 还没实现
        SetRedirects(true). // 还没实现
        //SetHeaders(map[string]string{"user-agent": "from option"}).
        //SetCookies(map[string]string{"b": "", "d": "666"}). // cookie兼容 字符串格式和map格式 也兼容cookie对象
        //SetParams([][2]any{{"k", 1}, {"v", "2"}}).
        SetCookieEnable(false). // 设置本次请求不使用cookie
        //SetProxy("http://ip:port").
        SetTimeout(3000)
	
	response := session.Fetch("http://httpbin.org/get", option)  // 如果需要配置那么就传入option对象
	//response := session.Fetch("http://httpbin.org/get", nil)  // 如果不需要配置 直接传入 nil
	
	// 获取结果
	fmt.Println(response.Error())
	fmt.Println(response.Text)
	fmt.Println(response.StatusCode)
	fmt.Println(response.Headers)
	fmt.Println(response.Body)
	fmt.Println(response.ContentLength)
	
	fmt.Println(response.Cookies)  // 这个cookies是opt对象 所以有下面的二次操作
	
	// 需要 key=value; 格式的cookie
    fmt.Println(response.Cookies.String())
    // 需要 {"key":"value"} 的JSON格式的cookie
    fmt.Println(response.Cookies.Json())
    // 需要 {"key":"value"} 的map[string]string 格式的cookie
    fmt.Println(response.Cookies.Map())
	// 需要 []*net.Cookie 格式的的cookie
    fmt.Println(response.Cookies.Cookies())
}
```


### 项目开发背景
```text
自学了golang一段时间后，苦于没有实践，然后发现golang里面很多操作重复写的东西太多

然后也看了一下目前已有的工具，都太老套了，实现一个功能写N行才能拿到想要的结果
然后我就打算 开发一些适合自己使用习惯的功能

项目经验不够丰富，所以代码没有那么好看，等以后经验多了，会持续修复优化代码

>>> 如果有bug或者建议 欢迎 issue 或者 pr
```
