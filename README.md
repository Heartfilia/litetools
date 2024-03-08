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
	defer litetime.Timer("main")()   // 用于统计该函数运行耗时
    t := litetime.Time{
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
    t := litetime.Time{}
    t := litetime.Time{
        Goal: "2024-01-10 10:43:21",
    }
    t := litetime.Time{
        Goal: "2024-01-10 10:43:21",
        Uint: "ms",
    }
    t := litetime.Time{
        Uint: "ms",
        Cursor: "-1d3h",
    }
    t := litetime.Time{
        Cursor: 10,
        Fmt: true,
    }
    t := litetime.Time{
        Unit: "ms",
        Fmt: false,
    }
    t := litetime.Time{
        Goal: 1704768201,
    }
    t := litetime.Time{
        Goal: 1704768201123,
        Fmt: "%Y-%m-%d",
    }
    t := litetime.Time{
        Goal: 1704768201123,
        Cursor: -1,
    }
    // 注意得结合预期获得结果 要不然获得的是对应类型的 零值
    fmt.Println(t.GetTime().Int())
    fmt.Println(t.GetTime().Float())
    fmt.Println(t.GetTime().String())
}
```

### string
```go
package main

import (
    "fmt"
    "github.com/Heartfilia/litetools/litestring"
)

func main(){
    fmt.Println(litestring.ColorString("<red>红色</red>还有其它颜色<cyan>其它颜色</cyan>还有一些错误测试<blue>错误的</yellow>", ""))
    fmt.Println(litestring.ColorString("整体替换颜色", "blue"))
    fmt.Println(litestring.ColorString("整体替换颜色", "黄"))
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
	
	fmt.Println(litenet.GetWAN())  // 还没弄 后续会调整用法
	fmt.Println(litenet.GetLAN())  // 还没弄 后续会调整用法
}


```

### rand
```go
import (
	"fmt"
	"github.com/Heartfilia/litetools/literand"
)

func main(){
	fmt.Println(literand.RandomChoice([]string{"a", "b", "c"}))
	fmt.Println(literand.RandomChoice([]int64{1, 2, 3}))
	fmt.Println(literand.RandomChoice([]float64{1.123, 2.223, 3.333}))
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
}

```