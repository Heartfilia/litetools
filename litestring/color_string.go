package litestring

import (
	"regexp"
	"strings"
)

// 颜色表
const (
	base   string = "\033[0m"
	red    string = "\033[31m"
	green  string = "\033[32m"
	yellow string = "\033[33m"
	blue   string = "\033[34m"
	purple string = "\033[35m"
	cyan   string = "\033[36m"
	white  string = "\033[37m"
)

var reflectColor = map[string]string{
	"红":      red,
	"红色":     red,
	"RED":    red,
	"绿":      green,
	"绿色":     green,
	"GREEN":  green,
	"黄":      yellow,
	"黄色":     yellow,
	"YELLOW": yellow,
	"蓝":      blue,
	"蓝色":     blue,
	"BLUE":   blue,
	"紫":      purple,
	"紫色":     purple,
	"PURPLE": purple,
	"青":      cyan,
	"青色":     cyan,
	"靛":      cyan,
	"靛色":     cyan,
	"CYAN":   cyan,
	"白":      white,
	"白色":     white,
	"WHITE":  white,
}

func isValidColor(color string) (string, bool) {
	newColor, ok := reflectColor[color]
	return newColor, ok
}

type pattern struct {
	ColorPattern     *regexp.Regexp
	BlockPatternHead *regexp.Regexp
	BlockPatternTail *regexp.Regexp
}

func (p *pattern) initPattern() {
	fullPattern, _ := regexp.Compile("<\\w+>.*?</\\w+>")
	p.ColorPattern = fullPattern

	headPattern, _ := regexp.Compile("<(\\w+)>")
	p.BlockPatternHead = headPattern

	tailPattern, _ := regexp.Compile("</(\\w+)>")
	p.BlockPatternTail = tailPattern
}

func totalReplace(goal, color string) string {
	color = strings.ToUpper(color)
	newColor, ok := isValidColor(color)
	if !ok {
		return goal
	}
	obj := strings.Builder{}
	obj.WriteString(newColor)
	obj.WriteString(goal)
	obj.WriteString(base)
	return obj.String()
}

func parseStringWithTag(goal string) string {
	/*
		这里只能装饰指定的颜色 只能修改字体 颜色装饰方案如下
		   <red>xxx</red>          -- 红色
		   <yellow>xxx</yellow>    -- 黄色
		   <blue>xxx</blue>        -- 蓝色
		   <green>xxx</green>      -- 绿色
		   <cyan>xxx</cyan>        -- 青色/靛色
		   <purple>xxx</purple>    -- 紫色
		   <pink>xxx</pink>        -- 粉丝
		   <black>xxx</black>      -- 黑色/不用配置这个没啥用
		   <white>xxx</white>      -- 白色
	*/
	nowPattern := pattern{}
	nowPattern.initPattern()

	allPattern := nowPattern.ColorPattern.FindAll([]byte(goal), -1)
	for _, rule := range allPattern {
		tempStringHead := nowPattern.BlockPatternHead.FindStringSubmatch(string(rule))
		tempStringTail := nowPattern.BlockPatternTail.FindStringSubmatch(string(rule))
		if tempStringHead[1] != tempStringTail[1] {
			continue
		}
		tempRule := strings.Replace(string(rule), tempStringHead[0], "", -1)
		tempRule = strings.Replace(tempRule, tempStringTail[0], "", -1)
		newString := totalReplace(tempRule, tempStringHead[1])
		goal = strings.Replace(goal, string(rule), newString, -1)
	}
	return goal
}

func ColorString(goal string, color ...string) string {
	// 仅支持 颜色 设置  其它样式不兼容
	// 如果color传入为 空字符串 那么则表示 需要解析 goal里面是否有字符串颜色标签语法
	if color == nil || (len(color) <= 1 && color[0] == "") {
		// 标签语法
		return parseStringWithTag(goal)
	} else {
		// 直接整体换色
		realColor := ""
		if len(color) >= 1 {
			realColor = color[0]
			return totalReplace(goal, realColor)
		}
		return parseStringWithTag(goal) // 一般不会走到这里来的
	}
}
