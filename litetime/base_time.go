package litetime

import (
	"fmt"
	"github.com/Heartfilia/litetools/litetime/fmt_time"
	"time"
)

type Time struct {
	Goal       interface{} // 基础数据类型 不传入 默认进行的是时间戳获取
	Fmt        interface{} // 格式化样式 不传入 默认不操作
	Unit       string      // 时间样式 s为秒 ms为毫秒
	Cursor     int         // 游标 默认为0  当前版本只兼容
	CursorUnit string      // 游标单位 默认 天d 还可以配置比天更小的单位 时H 分M 秒S
	Area       string      // 时区配置 不传入 默认配置到 Asia/Shanghai
}

type Result struct {
	result       interface{} // 统一的结果
	stringFmt    string      // 格式化时间最后的结果
	resultString string      // 用于判断情况的中间参数
	intMs        int64       // 整形的毫秒
	intS         int64       // 整形的秒
	floatMs      float64     // 浮点型的毫秒
	floatS       float64     // 浮点型的秒
}

func (t *Time) init() {
	if t.Unit == "" {
		t.Unit = "s"
	}
	if t.Area == "" {
		_, err := time.LoadLocation("Asia/Shanghai")
		if err != nil {
			return
		}
	}
	if t.CursorUnit == "" {
		t.CursorUnit = "d"
	}
	// 下面是interface版本 后续做更加精细兼容的时候再打开
	//if t.Cursor == nil {
	//	t.Cursor = 0
	//}
}

// -------------------------------------

func (r *Result) String() string {
	// 对于结果是字符串类型的统一返回
	if r.stringFmt == "" {
		// 如果格式化时间这里没有结果 那么再判断 结果那里面
		switch r.result.(type) {
		case string:
			return r.result.(string)
		case int:
		case int32:
		case int64:
			return fmt.Sprintf("%d", r.result)
		case float32:
		case float64:
			return fmt.Sprintf("%f", r.result)
		default:
			return ""
		}
	} else {
		return r.stringFmt
	}
	return ""
}

func (r *Result) Int64() int64 {
	if r.resultString == "ms" {
		return r.intMs
	} else if r.resultString == "s" {
		return r.intS
	}
	return 0
}

func (r *Result) Int() int {
	// 对于结果是数值类型的统一返回
	//res := r.result
	//switch res.(type) {
	//case int64:
	//	return int(res.(int64))
	//case int:
	//	return res.(int)
	//}
	//return 0
	if r.resultString == "ms" {
		return int(r.intMs)
	} else if r.resultString == "s" {
		return int(r.intS)
	}
	return 0
}

func (r *Result) Float() float64 {
	// 对于结果是浮点数的类型返回
	if r.resultString == "ms" {
		return r.floatMs
	} else if r.resultString == "s" {
		return r.floatS
	}
	return 0.
}

// -------------------------------------
func (t *Time) unit(r *Result) {
	//fmt.Println("1 s :", time.Now().Unix())
	//fmt.Println("2 ms:", time.Now().UnixMilli())
	//fmt.Println("3 ns:", time.Now().UnixMicro())
	//fmt.Println("4 ps:", time.Now().UnixNano())
	if t.Unit == "ms" {
		tempTime := time.Now().UnixMicro()
		r.floatMs = float64(tempTime) / 1000
		r.intMs = time.Now().UnixMilli()
		r.resultString = "ms"
	} else {
		tempTime := time.Now().UnixMicro()
		r.floatS = float64(tempTime) / 1e6
		r.intS = time.Now().Unix()
		r.resultString = "s"
	}
}

func (t *Time) number(r *Result) {
	if t.Goal == nil {
		t.unit(r)
	}
}

func (t *Time) fmtMode(r *Result) {
	if t.Goal == nil {
		// 如果没有传入的话 那么就是获得 当前时间的格式化时间 或者 添加了游标之后的格式化时间
		// 1. 什么都没传  fmt样式也没有传的
		var cursor int64
		if t.CursorUnit == "S" {
			cursor = 1
		} else if t.CursorUnit == "M" {
			cursor = 60
		} else if t.CursorUnit == "H" {
			cursor = 3600
		} else {
			cursor = 86400
		}
		switch t.Fmt.(type) {
		case bool:
			if t.Fmt.(bool) == true {
				r.stringFmt = fmt_time.NowFmt(t.Cursor, cursor)
			}
		case string:
			r.stringFmt = fmt_time.FmtType(t.Fmt.(string), t.Cursor, cursor)
		}

	} else {
		// 否则就是 传入了数字的时间戳 然后我这里需要对它进行格式化
		// 需要 秒 或者 毫秒 单位的时间  +  fmt的格式(默认 %Y-%m-%d %H:%M:%S)
	}
}

//------------- 主入口 -----------------

func (t *Time) GetTime() *Result {
	t.init()
	r := Result{}
	// 1 直接获取到时间的情况
	if t.Fmt == nil {
		t.number(&r)
	} else {
		t.fmtMode(&r)
	}

	return &r
}
