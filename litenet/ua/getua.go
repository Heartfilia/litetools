package ua

import "strings"

type makeChoice struct {
	OsType  string
	Version string
	UA      string
}

func (m *makeChoice) choice() {

}

func isBrowser(option string) bool {
	// 判断是不是浏览器
	for _, browser := range Browser {
		if browser == option {
			return true
		}
	}
	return false
}

func isSystem(option string) bool {
	// 判断是不是系统
	for _, system := range System {
		if system == option {
			return true
		}
	}
	return false
}

func Options(option string) string {
	option = strings.ToLower(option) // 先全部弄成小写
	if isBrowser(option) {

	} else if isSystem(option) {

	} else {
		// 否则直接从浏览器里面随机挑返回
	}
	return option
}
