package jsonPath

import (
	"encoding/json"
	"fmt"
	"log"
	"regexp"
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

type Result struct {
	resultString string
	Int64        int64
	Int          int
	Float64      float64
	Float32      float32
	Bool         bool
	Object       interface{}
}

func (r *Result) String() string {
	return r.resultString
}

func (r *Result) SetString(str string) {
	r.resultString = str
}

type sReplace string

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

type resultCache struct {
	Base    string         // 缓存当前处理的段的数据
	BaseAny any            // 把用json提取的Base放到这里
	Array   []any          // 如果下一段是 任意类型数组
	Object  map[string]any // 如果下一段是任意类型的map对象
}

type regexRule struct {
	OnlyArray *regexp.Regexp // [0]
	OnlyKey   *regexp.Regexp // a
	KeyArray  *regexp.Regexp // a[0]
	// 不支持的格式 [0]b  --> not support
}

func (r *regexRule) isEmpty() bool {
	// 判断是否为未初始化的数据结构
	if r.OnlyArray == nil || r.OnlyKey == nil || r.KeyArray == nil {
		return true
	}
	return false
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

var regRuleCache regexRule

func (r *resultCache) parse(rule string) {
	if regRuleCache.isEmpty() {
		// 如果是没有初始化的情况下 需要初始化内部的参数
		regRuleCache.OnlyArray = regexp.MustCompile("^\\[\\d+]$")
		regRuleCache.OnlyKey = regexp.MustCompile("^[^[]+")
		regRuleCache.KeyArray = regexp.MustCompile(".+\\[\\d+]$")
	}

	if regRuleCache.onlyArray(rule) {
		// 如果当前目标key是 [0]  那么预期当前节点能获取到的key样式应该是 []any
		fmt.Println("当前格式是 onlyArray   :", rule)
	} else if regRuleCache.keyWithArray(rule) {
		// 如果当前目标key是 a[0] 那么预期当前节点能获取到的样式应该是 map[string]any  需要拆分然后继续处理一次
		fmt.Println("当前格式是 keyWithArray:", rule)

	} else if regRuleCache.onlyKey(rule) {
		// 如果当前目标key是 a    那么预期当前节点能获取到的样式应该是 map[string]any
		fmt.Println("当前格式是 onlyKey     :", rule)
	} else {
		log.Fatal("m没有匹配到的格式是:", rule)
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

func parseRule(jsonString, rule string) {
	res, err := transToObj(jsonString)
	if err != nil {
		log.Fatal(err)
		return // 这里东西要返回 后面处理
	}

	var js = resultCache{
		Base:    jsonString,
		BaseAny: res,
	}

	rule = replaceTo(rule, 0) // 恢复成正常的rule格式
	eachBlock := strings.Split(rule, ".")
	for _, each := range eachBlock {
		js.parse(each)
	}
}

func JudgeAndExtractEachRule(jsonString string, rules []string) {
	for _, rule := range rules {
		parseRule(jsonString, rule)
	}
}
