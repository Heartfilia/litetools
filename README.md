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
    t := litetime.Time{
        //Unit: "ms",    // 如果是需要获取数字相关的时间 这里配置秒或者毫秒
        Fmt: true,       // 如果需要格式化时间 这里设置为 true  或者 转入需要的格式化时间样式
    }
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
