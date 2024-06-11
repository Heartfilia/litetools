package ua

import (
	"github.com/Heartfilia/litetools/liteslice"
	"strings"
)

type makeChoice struct {
	Browser string
	OsType  string
	Version string
	UA      string
}

func (m *makeChoice) choice() string {
	return CombineString(m.OsType, m.Browser)
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
	var mc makeChoice
	if isBrowser(option) {
		mc.OsType = "win"
		if option == "chrome" || option == "chromium" {
			mc.Browser = "chrome"
		} else if option == "firefox" {
			mc.Browser = "firefox"
		} else if option == "opera" {
			mc.Browser = "opera"
		} else if option == "ie" {
			mc.Browser = "ie"
		} else if option == "edge" {
			mc.Browser = "edge"
		} else if option == "safari" {
			mc.Browser = "safari"
			mc.OsType = "mac"
		}
	} else if isSystem(option) {
		mc.Browser = liteslice.RandomChoice([]string{"chrome", "edge"})
		if option == "pc" {
			mc.OsType = liteslice.RandomChoice([]string{"win", "mac", "linux"})
			if mc.OsType == "mac" {
				mc.Browser = liteslice.RandomChoice([]string{"chrome", "edge", "safari"})
			}
		} else if option == "mobile" {
			mc.OsType = liteslice.RandomChoice([]string{"android", "ios", "harmonyos"})
			if mc.OsType == "harmonyos" {
				mc.OsType = "harmony"
			}
		} else if option == "win" || option == "windows" {
			mc.OsType = "win"
		} else if option == "mac" || option == "macos" {
			mc.OsType = "mac"
		} else if option == "harmony" || option == "harmonyos" {
			mc.OsType = "harmony"
		} else {
			mc.OsType = option
		}
	} else {
		// 否则直接从浏览器里面随机挑返回
		mc.OsType = liteslice.RandomChoice([]string{"win", "mac", "linux"})
		mc.Browser = liteslice.RandomChoice([]string{"chrome", "firefox", "edge"})
	}
	return mc.choice()
}
