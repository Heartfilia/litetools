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

func normalizePlatform(option string) string {
	switch option {
	case "windows":
		return "win"
	case "macos":
		return "mac"
	case "harmonyos":
		return "harmony"
	default:
		return option
	}
}

func randomDesktopBrowser(platform string) string {
	switch platform {
	case "mac":
		return liteslice.RandomChoice([]string{"chrome", "firefox", "edge", "safari"})
	case "linux":
		return liteslice.RandomChoice([]string{"chrome", "firefox", "edge"})
	default:
		return liteslice.RandomChoice([]string{"chrome", "firefox", "edge"})
	}
}

func browserForPlatform(platform string, preferred string) string {
	switch platform {
	case "ios":
		if preferred == "" {
			return "safari"
		}
		if preferred == "safari" {
			return preferred
		}
		return "safari"
	case "android", "harmony":
		if preferred == "chrome" || preferred == "firefox" || preferred == "edge" {
			return preferred
		}
		return "chrome"
	case "mac":
		if preferred == "" {
			return randomDesktopBrowser(platform)
		}
		if preferred == "ie" {
			return "safari"
		}
		return preferred
	case "linux":
		if preferred == "" {
			return randomDesktopBrowser(platform)
		}
		if preferred == "safari" || preferred == "ie" {
			return "chrome"
		}
		return preferred
	case "win":
		if preferred == "" {
			return randomDesktopBrowser(platform)
		}
		if preferred == "safari" {
			return "chrome"
		}
		return preferred
	default:
		if preferred == "" {
			return randomDesktopBrowser("win")
		}
		return preferred
	}
}

func Options(option string) string {
	option = normalizePlatform(strings.ToLower(option)) // 先全部弄成小写
	var mc makeChoice
	if isBrowser(option) {
		switch option {
		case "chrome", "chromium":
			mc.Browser = "chrome"
			mc.OsType = "win"
		case "firefox":
			mc.Browser = "firefox"
			mc.OsType = liteslice.RandomChoice([]string{"win", "mac", "linux"})
		case "opera":
			mc.Browser = "opera"
			mc.OsType = liteslice.RandomChoice([]string{"win", "mac", "linux"})
		case "ie":
			mc.Browser = "ie"
			mc.OsType = "win"
		case "edge":
			mc.Browser = "edge"
			mc.OsType = liteslice.RandomChoice([]string{"win", "mac", "linux"})
		case "safari":
			mc.Browser = "safari"
			mc.OsType = liteslice.RandomChoice([]string{"mac", "ios"})
		}
	} else if isSystem(option) {
		if option == "pc" {
			mc.OsType = liteslice.RandomChoice([]string{"win", "mac", "linux"})
			mc.Browser = browserForPlatform(mc.OsType, "")
		} else if option == "mobile" {
			mc.OsType = liteslice.RandomChoice([]string{"android", "ios", "harmony"})
			mc.Browser = browserForPlatform(mc.OsType, "")
		} else {
			mc.OsType = option
			mc.Browser = browserForPlatform(mc.OsType, "")
		}
	} else {
		// 否则直接从浏览器里面随机挑返回
		mc.OsType = liteslice.RandomChoice([]string{"win", "mac", "linux"})
		mc.Browser = browserForPlatform(mc.OsType, "")
	}
	mc.OsType = normalizePlatform(mc.OsType)
	mc.Browser = browserForPlatform(mc.OsType, mc.Browser)
	return mc.choice()
}
