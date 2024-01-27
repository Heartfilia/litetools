package litetime

import (
	"fmt"
	"github.com/Heartfilia/litetools/litetime/fmt_time"
	"regexp"
	"strconv"
	"strings"
	"time"
)

type Time struct {
	Goal   interface{} // 基础数据类型 不传入 默认进行的是时间戳获取
	Fmt    interface{} // 格式化样式 不传入 默认不操作
	Unit   string      // 时间样式 s为秒 ms为毫秒
	Cursor interface{} // 游标 默认为0 传入数字为天 支持的单位 d天 H小时 M分钟 S秒 可以组合传递 如： -1d12H 一天12小时前
	Area   string      // 时区配置 不传入 默认配置到 Asia/Shanghai
}

type Result struct {
	result       interface{} // 统一的结果
	stringFmt    string      // 格式化时间最后的结果
	resultString string      // 用于判断情况的中间参数
	intMs        int64       // 整形的毫秒
	intS         int64       // 整形的秒
	floatMs      float64     // 浮点型的毫秒
	floatS       float64     // 浮点型的秒
	err          error       // 携带的错误 一般出现在格式化时间转时间戳的情况
}

const defaultTime = "%Y-%m-%d %H:%M:%S"

func (t *Time) init() {
	if t.Unit == "" {
		t.Unit = "s"
	}
	if t.Area == "" {
		//_, err := time.LoadLocation("Asia/Shanghai")
		//if err != nil {
		//	return
		//}
		t.Area = "Asia/Shanghai"
	}
	if t.Cursor == nil {
		t.Cursor = 0
	}
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

func stringGoal(goal, fmtString, unit, area string, r *Result) {
	// 传入了格式化时间 格式化的格式 获得转好的时间戳
	r.stringFmt = goal
	golangFmt := fmt_time.GetFormat(fmtString)
	location, err := time.LoadLocation(area)
	if err != nil {
		r.err = err
		return
	}
	ts, err := time.ParseInLocation(golangFmt, goal, location)
	if err != nil {
		r.err = err
		return
	}
	if unit == "ms" {
		tempTime := ts.UnixMicro()
		r.floatMs = float64(tempTime) / 1000
		r.intMs = ts.UnixMilli()
		r.resultString = "ms"
	} else {
		tempTime := ts.UnixMicro()
		r.floatS = float64(tempTime) / 1e6
		r.intS = ts.Unix()
		r.resultString = "s"
	}
}

func intGoal() {

}

func noGoal(t *Time, r *Result) {
	var cursor string
	switch t.Cursor.(type) {
	case int:
		cursor = fmt.Sprintf("%dh", t.Cursor.(int)*24)
	case int64:
		cursor = fmt.Sprintf("%dh", t.Cursor.(int64)*24)
	case string:
		cursor = strings.ToLower(t.Cursor.(string))
		if strings.Contains(cursor, "d") {
			// 如果包含了天这个参数 那么需要把天提取出来累加到 h 上面去
			baseH := 0 // 默认的小时
			hasH, _ := regexp.Match("\\d+h", []byte(cursor))
			if hasH {
				regH, _ := regexp.Compile("(\\d+)h")
				baseHString := regH.FindStringSubmatch(cursor)
				baseH, _ = strconv.Atoi(baseHString[1])
				cursor = strings.Replace(cursor, baseHString[0], "", -1)
			}
			regD, _ := regexp.Compile("(\\d+)d")
			baseDString := regD.FindStringSubmatch(cursor)
			baseD, _ := strconv.Atoi(baseDString[1])
			newHour := fmt.Sprintf("%dh", baseD*24+baseH)
			cursor = strings.Replace(cursor, baseDString[0], newHour, -1)
		}
	}

	switch t.Fmt.(type) {
	case bool:
		if t.Fmt.(bool) == true {
			r.stringFmt = fmt_time.FmtType("", cursor)
		}
	case string:
		r.stringFmt = fmt_time.FmtType(t.Fmt.(string), cursor)
	}
}

func withGoal(t *Time, r *Result) {
	switch t.Goal.(type) {
	case string:
		// 如果goal是字符串相关的
		var fmtTemp string
		switch t.Fmt.(type) {
		case bool:
			fmtTemp = defaultTime
		case string:
			fmtTemp = t.Fmt.(string)
		}
		if t.Fmt == nil || fmtTemp == "" {
			fmtTemp = defaultTime
		}
		stringGoal(t.Goal.(string), fmtTemp, t.Unit, t.Area, r)
	case int:
		intGoal()
	case int64:
		intGoal()
	}
	// 如果goal是数字相关的

}

func (t *Time) fmtMode(r *Result) {
	if t.Goal == nil {
		// 如果没有传入的话 那么就是获得 当前时间的格式化时间 或者 添加了游标之后的格式化时间
		// 1. 什么都没传  fmt样式也没有传的
		noGoal(t, r)
	} else {
		// 否则就是 传入了数字的时间戳 然后我这里需要对它进行格式化
		// 需要 秒 或者 毫秒 单位的时间  +  fmt的格式(默认 %Y-%m-%d %H:%M:%S)
		withGoal(t, r)
	}
}

//------------- 主入口 -----------------

func (t *Time) Default() {
	/*
		快速恢复成默认状态
	*/
	t.Goal = nil
	t.Fmt = nil
	t.Unit = "s"
	t.Cursor = 0
	t.Area = ""
}

func (t *Time) GetTime() *Result {
	t.init()
	r := Result{}
	// 1 直接获取到时间的情况
	if t.Fmt == nil && t.Goal == nil {
		t.number(&r)
	} else {
		t.fmtMode(&r)
	}

	return &r
}
