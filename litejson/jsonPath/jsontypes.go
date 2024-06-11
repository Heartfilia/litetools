package jsonPath

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/Heartfilia/litetools/litestring"
	"github.com/Heartfilia/litetools/utils/litedir"
	"log"
	"reflect"
	"regexp"
	"strconv"
	"strings"
)

const (
	_     = iota + '﮽'
	dot   // .
	left  // [
	right // ]
	line  // |
	blank //
)

type resultCache struct {
	//Base    string         // 缓存当前处理的段的数据   // 没啥用 先处理了
	BaseAny any            // 把用json提取的Base放到这里
	Array   []any          // 如果下一段是 任意类型数组
	Result  any            // 结果
	Object  map[string]any // 如果下一段是任意类型的map对象
	OK      bool           // 是否是最终结果
	Error   error          // 程序过程中的错误
}

type Result struct {
	Error  error
	value  any
	isLast bool
}

var regRuleCache regexRule

type regexRule struct {
	OnlyArray     *regexp.Regexp // [0]
	OnlyKey       *regexp.Regexp // a
	KeyArray      *regexp.Regexp // a[0]
	SplitKeyArray *regexp.Regexp // a [0] [1] 拆分这个
	// 不支持的格式 [0]b  --> not support
}

func (r *Result) setLast(value bool) {
	r.isLast = value
}

func (r *Result) getLast() bool {
	return r.isLast
}

func (r *Result) setValue(value any) {
	r.value = value
}

func (r *Result) Value() any {
	return r.value
}

func (r *Result) StringSlice() []string {
	switch r.Value().(type) {
	case []interface{}:
		result := make([]string, len(r.Value().([]any)))
		for i, v := range r.Value().([]any) {
			result[i] = fmt.Sprintf("%v", v)
		}
		return result
	}
	return nil
}

func (r *Result) String() string {
	if r.Value() != nil {
		return fmt.Sprintf("%v", r.Value())
	}
	return ""
}

func (r *Result) Int() int {
	switch r.Value().(type) {
	case int:
		return r.Value().(int)
	case int64:
		return int(r.Value().(int64))
	case int32:
		return int(r.Value().(int32))
	case float64:
		return int(r.Value().(float64))
	case float32:
		return int(r.Value().(float32))
	}
	r.Error = errors.New(colorPanic("type conversion failed"))
	return 0
}

func (r *Result) IntSlice() []int {
	switch r.Value().(type) {
	case []any:
		result := make([]int, len(r.Value().([]any)))
		for i, v := range r.Value().([]any) {
			realValue, err := strconv.Atoi(fmt.Sprintf("%v", v))
			if err != nil {
				r.Error = errors.New(colorPanic("bad slice types"))
				return nil
			}
			result[i] = realValue
		}
		return result
	}
	return nil
}

func (r *Result) Int64() int64 {
	switch r.Value().(type) {
	case int:
		return int64(r.Value().(int))
	case int64:
		return r.Value().(int64)
	case int32:
		return int64(r.Value().(int32))
	case float64:
		return int64(r.Value().(float64))
	case float32:
		return int64(r.Value().(float32))
	}
	r.Error = errors.New(colorPanic("type conversion failed"))
	return 0
}

func (r *Result) Int64Slice() []int64 {
	switch r.Value().(type) {
	case []any:
		result := make([]int64, len(r.Value().([]any)))
		for i, v := range r.Value().([]any) {
			realValue, err := strconv.Atoi(fmt.Sprintf("%v", v))
			if err != nil {
				r.Error = errors.New(colorPanic("bad slice types"))
				return nil
			}
			result[i] = int64(realValue)
		}
		return result
	}
	return nil
}

func (r *Result) Int32() int32 {
	switch r.Value().(type) {
	case int:
		return int32(r.Value().(int))
	case int64:
		return int32(r.Value().(int64))
	case int32:
		return r.Value().(int32)
	case float64:
		return int32(r.Value().(float64))
	case float32:
		return int32(r.Value().(float32))

	}
	r.Error = errors.New(colorPanic("type conversion failed"))
	return 0
}

func (r *Result) Int32Slice() []int32 {
	switch r.Value().(type) {
	case []any:
		result := make([]int32, len(r.Value().([]any)))
		for i, v := range r.Value().([]any) {
			realValue, err := strconv.Atoi(fmt.Sprintf("%v", v))
			if err != nil {
				r.Error = errors.New(colorPanic("bad slice types"))
				return nil
			}
			result[i] = int32(realValue)
		}
		return result
	}
	return nil
}

func (r *Result) Float() float64 {
	switch r.Value().(type) {
	case float64:
		return r.Value().(float64)
	case float32:
		return float64(r.Value().(float32))
	case int:
		return float64(r.Value().(int))
	case int64:
		return float64(r.Value().(int64))
	case int32:
		return float64(r.Value().(int32))
	}
	r.Error = errors.New(colorPanic("type conversion failed"))
	return 0.0
}

func (r *Result) FloatSlice() []float64 {
	switch r.Value().(type) {
	case []any:
		result := make([]float64, len(r.Value().([]any)))
		for i, v := range r.Value().([]any) {
			switch v.(type) {
			case float64:
				result[i] = v.(float64)
			case float32:
				result[i] = float64(v.(float32))
			default:
				r.Error = errors.New(colorPanic("bad slice types"))
				return nil
			}

		}
		return result
	}
	return nil
}

func (r *Result) Float32() float32 {
	switch r.Value().(type) {
	case float64:
		return float32(r.Value().(float64))
	case float32:
		return r.Value().(float32)
	case int:
		return float32(r.Value().(int))
	case int64:
		return float32(r.Value().(int64))
	case int32:
		return float32(r.Value().(int32))
	}
	r.Error = errors.New(colorPanic("type conversion failed"))
	return 0.0
}

func (r *Result) Float32Slice() []float32 {
	switch r.Value().(type) {
	case []any:
		result := make([]float32, len(r.Value().([]any)))
		for i, v := range r.Value().([]any) {
			switch v.(type) {
			case float32:
				result[i] = v.(float32)
			case float64:
				result[i] = float32(v.(float64))
			default:
				r.Error = errors.New(colorPanic("bad slice types"))
				return nil
			}

		}
		return result
	}
	return nil
}

func (r *Result) Bool() bool {
	otherBool := false
	switch r.Value().(type) {
	case bool:
		return r.Value().(bool)
	case int64:
		if r.Value().(int64) != 0 {
			otherBool = true
		}
	case int32:
		if r.Value().(int32) != 0 {
			otherBool = true
		}
	case int:
		if r.Value().(int) != 0 {
			otherBool = true
		}
	case float64:
		if r.Value().(float64) != 0 {
			otherBool = true
		}
	case float32:
		if r.Value().(float32) != 0 {
			otherBool = true
		}
	case string:
		if r.Value().(string) != "" {
			otherBool = true
		}
	default:
		if r.Value() != nil {
			otherBool = true
		}
	}
	return otherBool
}

func (r *Result) BoolSlice() []bool {
	switch r.Value().(type) {
	case []bool:
		return r.Value().([]bool)
	}
	return nil
}

func colorPanic(msg string) string {
	return litestring.ColorString("panic: ", "red") + msg
}

func replaceTo(str string, mode int) string {
	if mode == 1 {
		// 把不需要转义的字符串替换了 如 \.
		str = strings.ReplaceAll(str, "\\.", string(dot))
		str = strings.ReplaceAll(str, "\\[", string(left))
		str = strings.ReplaceAll(str, "\\]", string(right))
		str = strings.ReplaceAll(str, "\\|", string(line))
		str = strings.ReplaceAll(str, "\\ ", string(blank))
	} else {
		str = strings.ReplaceAll(str, string(dot), "\\.")
		str = strings.ReplaceAll(str, string(left), "\\[")
		str = strings.ReplaceAll(str, string(right), "\\]")
		str = strings.ReplaceAll(str, string(line), "\\|")
		str = strings.ReplaceAll(str, string(blank), "\\ ")
	}
	return str
}

func SplitRule(rule string) []string {
	// 替换特殊符号 如 \.   \[   \|
	rule = replaceTo(rule, 1)
	return strings.Split(rule, "|")
}

func (r *regexRule) isEmpty() bool {
	// 判断是否为未初始化的数据结构
	if r.OnlyArray == nil || r.OnlyKey == nil || r.KeyArray == nil || r.SplitKeyArray == nil {
		return true
	}
	return false
}

func (r *regexRule) initRegex() {
	r.OnlyArray = regexp.MustCompile("^\\[-?\\d+]$")
	r.OnlyKey = regexp.MustCompile("^[^[]+")
	r.KeyArray = regexp.MustCompile(".+\\[-?\\d+]$")
	r.SplitKeyArray = regexp.MustCompile("\\[-?\\d+]")
}

func (r *regexRule) onlyArray(rule string) bool {
	return r.OnlyArray.MatchString(rule)
}

func (r *regexRule) onlyKey(rule string) bool {
	return r.OnlyKey.MatchString(rule)
}

func (r *regexRule) keyWithArray(rule string) bool {
	return r.KeyArray.MatchString(rule)
}

func clearArray(rule *string, array []string) {
	// 避免出现 a[bc.c[1]
	for _, item := range array {
		*rule = strings.ReplaceAll(*rule, item, "")
	}
}

func (r *resultCache) parse(rule string, lastKey bool) {
	nowObj := r.BaseAny
	if nowObj == nil && lastKey != true {
		r.Result = nil
		r.OK = true
		r.Error = errors.New(colorPanic("there is a next node, but no object is available"))
		return
	}
	if regRuleCache.isEmpty() {
		// 如果是没有初始化的情况下 需要初始化内部的参数
		regRuleCache.initRegex()
	}

	if regRuleCache.onlyArray(rule) {
		// 如果当前目标key是 [0]  那么预期当前节点能获取到的key样式应该是 []any
		//fmt.Println("当前格式是 onlyArray   :", rule, r.BaseAny)
		numberRule := strings.ReplaceAll(rule, "[", "")
		numberRule = strings.ReplaceAll(numberRule, "]", "") // 去除前后括号
		numberInRule, err := strconv.Atoi(numberRule)
		if err != nil {
			r.Result = nil
			r.OK = true
			r.Error = errors.New(colorPanic(" wrong extraction sequence number"))
			return
		}
		if nowObj == nil {
			r.Result = nil
			r.OK = true
			r.Error = errors.New(colorPanic("runtime error: invalid memory address or nil pointer dereference"))
			return
		}
		typeAny := reflect.TypeOf(nowObj)
		if typeAny != nil {
			objKind := typeAny.Kind()
			switch objKind {
			case reflect.Slice:
				array := nowObj.([]any)
				if numberInRule < 0 {
					// 兼容 [-1]  [-2]
					numberInRule += len(array)
					if numberInRule < 0 {
						r.Result = nil
						r.OK = true
						r.Error = errors.New(colorPanic(
							fmt.Sprintf("runtime error: index out of range [%d] with length %d",
								numberInRule-len(array), len(array)),
						))
						return
					}
				} else if numberInRule > len(array)-1 {
					r.Result = nil
					r.OK = true
					r.Error = errors.New(colorPanic(
						fmt.Sprintf("runtime error: index out of range [%d] with length %d",
							numberInRule, len(array)),
					))
					return
				}

				if array == nil || len(array)-1 < numberInRule {
					r.Result = nil
					r.OK = true
					r.Error = errors.New(colorPanic("no slice or wrong extraction sequence number"))
					return
				}
				r.BaseAny = array[numberInRule]

				if lastKey {
					r.Result = array[numberInRule]
				}
			}
		} else {
			r.OK = true
		}
	} else if regRuleCache.keyWithArray(rule) {
		// 如果当前目标key是 a[0]   a[1][2] 那么预期当前节点能获取到的样式应该是 map[string]any  需要拆分然后继续处理一次
		//fmt.Println("当前格式是 keyWithArray:", rule)
		items := regRuleCache.SplitKeyArray.FindAllString(rule, -1)
		clearArray(&rule, items) // 清理掉每个数组
		r.parse(rule, false)
		lastKeyInd := false
		for ind, item := range items {
			if lastKey && ind == len(items)-1 {
				lastKeyInd = true
			}
			lastKeyInd = lastKeyInd && lastKey // 要是最后一个 并且是最后一个格子
			r.parse(item, lastKeyInd)
		}
	} else if regRuleCache.onlyKey(rule) {
		// 如果当前目标key是 a    那么预期当前节点能获取到的样式应该是 map[string]any
		//fmt.Println("当前格式是 onlyKey     :", rule, r.BaseAny)
		typeAny := reflect.TypeOf(nowObj)
		if typeAny != nil {
			objKind := typeAny.Kind()
			switch objKind {
			case reflect.Map:
				value := nowObj.(map[string]any)[rule]
				r.BaseAny = value
				if lastKey {
					r.Result = value
				}
			}
		} else {
			r.OK = true
		}
	} else {
		log.Fatal("没有匹配到的格式是:", rule, lastKey)
	}

}

func transToObj(jsonString string) (any, error) {
	var s any
	err := json.Unmarshal([]byte(jsonString), &s)
	if err != nil {
		return nil, err
	}
	return s, nil
}

func parseRule(jsonString, rule string, resultObj *Result) {
	res, err := transToObj(jsonString)
	if err != nil {
		resultObj.Error = err
		return // 这里东西要返回 后面处理
	}

	var js = resultCache{
		//Base:    jsonString,  // 这一行缓存没啥用 不占用资源了
		BaseAny: res,
	}

	rule = replaceTo(rule, 0) // 恢复成正常的rule格式
	eachBlock := strings.Split(rule, ".")
	lastKey := false
	for ind, each := range eachBlock {
		if ind == len(eachBlock)-1 {
			lastKey = true // 如果是最后一个key 那么筛选情况可能不一样
		}
		js.parse(each, lastKey)
		if js.OK {
			break
		}
	}
	resultObj.setValue(js.Result)
	resultObj.Error = js.Error
	if js.Error != nil {
		resultObj.setLast(false)
	} else if js.Result != nil {
		resultObj.setLast(true)
	}

}

func JudgeAndExtractEachRule(jsonString string, rules []string, resultObj *Result) {
	if strings.Index(jsonString, ".json") != -1 && litedir.FileExists(jsonString) {
		// 只是猜测是不是json的文件 如果是的话 直接读取处理
		tempString := litedir.FileReader(jsonString)
		if tempString != "" {
			jsonString = tempString
		}
	}
	for _, rule := range rules {
		parseRule(jsonString, rule, resultObj)
		if resultObj.getLast() {
			break
		}
	}
}
