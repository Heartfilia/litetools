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

	cursor := parseCursor(t.Cursor)

	if t.Unit == "ms" {
		mis := cursorSecond(cursor, 1000)
		tempTime := time.Now().UnixMicro()
		r.floatMs = float64(tempTime)/1000 + mis
		r.intMs = time.Now().UnixMilli() + int64(mis)
		r.resultString = "ms"
	} else {
		mis := cursorSecond(cursor, 1)
		tempTime := time.Now().UnixMicro()
		r.floatS = float64(tempTime) / 1e6
		r.intS = time.Now().Unix() + int64(mis)
		r.resultString = "s"
	}
}

func (t *Time) number(r *Result) {
	if t.Goal == nil {
		t.unit(r)
	}
}

func stringGoal(goal string, t *Time, r *Result) {
	// 传入了格式化时间 格式化的格式 获得转好的时间戳
	fmtString := getFmt(t.Fmt, false)
	cursor := parseCursor(t.Cursor)

	r.stringFmt = goal
	golangFmt := fmt_time.GetFormat(fmtString)
	location, err := time.LoadLocation(t.Area)
	if err != nil {
		r.err = err
		return
	}
	ts, err := time.ParseInLocation(golangFmt, goal, location)
	if err != nil {
		r.err = err
		return
	}
	if t.Unit == "ms" {
		mis := cursorSecond(cursor, 1000)
		tempTime := ts.UnixMicro()
		r.floatMs = float64(tempTime)/1000 + mis
		r.intMs = ts.UnixMilli() + int64(mis)
		r.resultString = "ms"
	} else {
		mis := cursorSecond(cursor, 1)
		tempTime := ts.UnixMicro()
		r.floatS = float64(tempTime)/1e6 + mis
		r.intS = ts.Unix() + int64(mis)
		r.resultString = "s"
	}
	if cursor != "0h" {
		t, _ := time.ParseDuration(cursor)
		ts = ts.Add(t)
		r.stringFmt = ts.Format(golangFmt)
	}

}

func intGoal(goal int64, t *Time, r *Result) {
	fmtTemp := getFmt(t.Fmt, true)
	cursor := parseCursor(t.Cursor)
	if 1e9 <= goal && goal < 1e10 {
		//
		mis := cursorSecond(cursor, 1)
		r.intS = goal + int64(mis)
		r.intMs = goal*1000 + int64(mis)
	} else if 1e12 <= goal && goal < 1e13 {
		mis := cursorSecond(cursor, 1000)
		r.intMs = goal + int64(mis)
		r.intS = goal/1000 + int64(mis)
	} else {
		panic("只能处理秒或者毫秒级别的数据\npanic: only handle 'm' or 'ms'")
	}

	if fmtTemp == "" {
		// 如果是没有格式化时间的情况
		if t.Unit == "ms" {
			r.stringFmt = fmt.Sprintf("%d", r.intMs)
		} else {
			r.stringFmt = fmt.Sprintf("%d", r.intS)
		}
	} else {
		golangFmt := fmt_time.GetFormat(fmtTemp)
		var ts time.Time
		if t.Unit == "ms" {
			ts = time.Unix(r.intMs/1000, 0)
		} else {
			ts = time.Unix(r.intS, 0)
		}
		r.stringFmt = ts.Format(golangFmt)

	}
}

func cursorSecond(cursorString string, times int64) float64 {
	if cursorString == "0h" {
		return 0
	}
	t, _ := time.ParseDuration(cursorString)

	return t.Seconds() * float64(times)
}

func parseCursor(cursor interface{}) string {
	var resultCursor string
	switch cursor.(type) {
	case int:
		resultCursor = fmt.Sprintf("%dh", cursor.(int)*24)
	case int64:
		resultCursor = fmt.Sprintf("%dh", cursor.(int64)*24)
	case string:
		resultCursor = strings.ToLower(cursor.(string))
		if strings.Contains(resultCursor, "d") {
			// 如果包含了天这个参数 那么需要把天提取出来累加到 h 上面去
			baseH := 0 // 默认的小时
			hasH, _ := regexp.Match("\\d+h", []byte(resultCursor))
			if hasH {
				regH, _ := regexp.Compile("(\\d+)h")
				baseHString := regH.FindStringSubmatch(resultCursor)
				baseH, _ = strconv.Atoi(baseHString[1])
				resultCursor = strings.Replace(resultCursor, baseHString[0], "", -1)
			}
			regD, _ := regexp.Compile("(\\d+)d")
			baseDString := regD.FindStringSubmatch(resultCursor)
			baseD, _ := strconv.Atoi(baseDString[1])
			newHour := fmt.Sprintf("%dh", baseD*24+baseH)
			resultCursor = strings.Replace(resultCursor, baseDString[0], newHour, -1)
		}
	}
	return resultCursor
}

func falseFmt(t *Time, r *Result) {
	var cursor string
	cursor = parseCursor(t.Cursor)
	if t.Unit == "ms" {
		mis := cursorSecond(cursor, 1000)
		tempTime := time.Now().UnixMicro()
		r.floatMs = float64(tempTime)/1000 + mis
		r.intMs = time.Now().UnixMilli() + int64(mis)
		r.stringFmt = fmt.Sprintf("%d", r.intMs)
	} else {
		mis := cursorSecond(cursor, 1)
		tempTime := time.Now().UnixMilli()
		r.floatS = float64(tempTime)/1000 + mis
		r.intS = time.Now().Unix() + int64(mis)
		r.stringFmt = fmt.Sprintf("%d", r.intS)
	}
	r.resultString = t.Unit
}

func noGoal(t *Time, r *Result) {
	var cursor string
	cursor = parseCursor(t.Cursor)

	switch t.Fmt.(type) {
	case bool:
		if t.Fmt.(bool) == true {
			r.stringFmt = fmt_time.FmtType("", cursor)
		} else {
			falseFmt(t, r)
		}
	case string:
		r.stringFmt = fmt_time.FmtType(t.Fmt.(string), cursor)
	}
}

func getFmt(_fmt interface{}, noFmt bool) string {
	var fmtTemp string
	switch _fmt.(type) {
	case bool:
		if _fmt.(bool) == true {
			fmtTemp = defaultTime
		}
	case string:
		fmtTemp = _fmt.(string)
	}
	if (_fmt == nil || fmtTemp == "") && !noFmt {
		fmtTemp = defaultTime
	}
	return fmtTemp
}

func withGoal(t *Time, r *Result) {
	switch t.Goal.(type) {
	case string:
		// 如果goal是字符串相关的
		stringGoal(t.Goal.(string), t, r)
	case float64:
		intGoal(int64(t.Goal.(float64)), t, r)
	case int:
		intGoal(int64(t.Goal.(int)), t, r)
	case int64:
		intGoal(t.Goal.(int64), t, r)
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
	t.Area = "Asia/Shanghai"
}

func (t *Time) GetTime() *Result {
	t.init()
	r := new(Result)
	// 1 直接获取到时间的情况
	if t.Fmt == nil && t.Goal == nil {
		t.number(r)
	} else {
		t.fmtMode(r)
	}

	return r
}
