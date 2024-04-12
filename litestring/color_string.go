package litestring

import (
	myColor "github.com/Heartfilia/litetools/litestring/color"
	"regexp"
	"strings"
)

// 颜色表

var reflectColor = map[string]string{
	"红":      myColor.Red,
	"红色":     myColor.Red,
	"RED":    myColor.Red,
	"绿":      myColor.Green,
	"绿色":     myColor.Green,
	"GREEN":  myColor.Green,
	"黄":      myColor.Yellow,
	"黄色":     myColor.Yellow,
	"YELLOW": myColor.Yellow,
	"蓝":      myColor.Blue,
	"蓝色":     myColor.Blue,
	"BLUE":   myColor.Blue,
	"紫":      myColor.Purple,
	"紫色":     myColor.Purple,
	"PURPLE": myColor.Purple,
	"青":      myColor.Cyan,
	"青色":     myColor.Cyan,
	"靛":      myColor.Cyan,
	"靛色":     myColor.Cyan,
	"CYAN":   myColor.Cyan,
	"白":      myColor.White,
	"白色":     myColor.White,
	"WHITE":  myColor.White,
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

func isOriginal(color string) (ok bool) {
	switch color {
	case myColor.Red:
		ok = true
	case myColor.Green:
		ok = true
	case myColor.Yellow:
		ok = true
	case myColor.Blue:
		ok = true
	case myColor.Purple:
		ok = true
	case myColor.Cyan:
		ok = true
	case myColor.White:
		ok = true
	case myColor.Base:
		ok = true
	}
	return
}

func totalReplace(goal, color string) string {
	var newColor string
	if !isOriginal(color) {
		color = strings.ToUpper(color)
		colorString, ok := isValidColor(color)
		if !ok {
			return goal
		}
		newColor = colorString
	} else {
		newColor = color
	}

	obj := strings.Builder{}
	obj.WriteString(newColor)
	obj.WriteString(goal)
	obj.WriteString(myColor.Base)
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

func ColorString(goal, color string) string {
	// 仅支持 颜色 设置  其它样式不兼容
	// 如果color传入为 空字符串 那么则表示 需要解析 goal里面是否有字符串颜色标签语法
	if color == "" {
		// 标签语法
		return parseStringWithTag(goal)
	} else {
		// 直接整体换色
		return totalReplace(goal, color)
	}
}
