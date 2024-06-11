package litepage

import (
	"fmt"
	"github.com/Heartfilia/litetools/utils/litedir"
	"log"
	"strings"
)

// 先借鉴dp的api风格: https://gitee.com/g1879/DrissionPage/blob/dev/DrissionPage/_configs/chromium_options.py

type ChromiumOptions struct {
	yamlConf      string // 如果传入这个 则表示从这个yaml文件里面读取配置文件 不过我先不实现这个操作
	userDataPath  string
	user          string
	headless      bool
	downloadPath  string
	tmpPath       string
	arguments     []string
	browserPath   string
	extensions    []string
	prefs         map[string]string
	flags         map[string]string
	address       string
	loadMode      string
	proxy         string
	port          int
	retry         int
	retryInterval float64
}

func (c *ChromiumOptions) DownloadPath() string {
	// 默认下载路径文件路径
	return c.downloadPath
}

func (c *ChromiumOptions) BrowserPath() string {
	// 浏览器启动文件路径
	return c.browserPath
}

func (c *ChromiumOptions) UserDataPath() string {
	// 返回用户数据文件夹路径
	return c.userDataPath
}

func (c *ChromiumOptions) TmpPath() string {
	// 返回临时文件夹路径
	return c.tmpPath
}

func (c *ChromiumOptions) User() string {
	// 返回用户配置文件夹名称
	return c.user
}

func (c *ChromiumOptions) LoadMode() string {
	// 返回页面加载策略，'normal', 'eager', 'none'
	if c.loadMode == "" {
		return "normal"
	}
	return c.loadMode
}

func (c *ChromiumOptions) Proxy() string {
	// 返回代理设置
	return c.proxy
}

func (c *ChromiumOptions) Address() string {
	// 返回浏览器地址，ip:port
	return c.address
}

func (c *ChromiumOptions) Arguments() []string {
	// 返回浏览器命令行设置列表
	return c.arguments
}

func (c *ChromiumOptions) Extensions() []string {
	// 以list形式返回要加载的插件路径
	return c.extensions
}

func (c *ChromiumOptions) Preferences() map[string]string {
	// 返回用户首选项配置
	return c.prefs
}

func (c *ChromiumOptions) Flags() map[string]string {
	// 返回实验项配置
	return c.flags
}

func (c *ChromiumOptions) Retry() int {
	// 返回连接失败时的重试次数
	return c.retry
}

func (c *ChromiumOptions) RetryInterval() float64 {
	// 返回连接失败时的重试间隔（秒）
	return c.retryInterval
}

func (c *ChromiumOptions) SetRetry(ts int, interval float64) *ChromiumOptions {
	// 设置连接失败时的重试操作
	c.retry = ts
	c.retryInterval = interval
	return c
}

func (c *ChromiumOptions) RemoveArgument(value string) *ChromiumOptions {
	// 移除一个argument项
	if c.arguments == nil {
		c.arguments = make([]string, 0)
	}
	if len(c.arguments) == 0 {
		return c
	}

	delList := make([]string, 0)

	for _, argument := range c.arguments {
		if argument == value || strings.Index(argument, fmt.Sprintf("%v=", value)) == 0 {
			delList = append(delList, argument)
		}
	}

	//for _, del := range delList {
	//	c.arguments = liteslice.SliceRemove(c.arguments, del)
	//}
	return c
}

func (c *ChromiumOptions) SetArgument(arg string, value any) *ChromiumOptions {
	// 设置浏览器配置的argument属性
	// arg  : 属性名
	// value: 属性值，如果有值的传入值，没有值的传入 "" ,需要删除这个选项的传入 nil
	if c.arguments == nil {
		c.arguments = make([]string, 0)
	}
	c.RemoveArgument(arg)

	switch value.(type) {
	case string:
		if arg == "--headless" && value.(string) == "" {
			c.arguments = append(c.arguments, "--headless=new")
		} else if arg != "" && value.(string) == "" {
			c.arguments = append(c.arguments, arg)
		} else if arg != "" && value.(string) != "" {
			c.arguments = append(c.arguments, fmt.Sprintf("%s=%s", arg, value.(string)))
		}
	}

	return c
}

func (c *ChromiumOptions) AddExtension(path string) *ChromiumOptions {
	if !litedir.FileExists(path) {
		log.Panicln("插件路径不存在")
	} else {
		c.extensions = append(c.extensions, path)
	}
	return c
}

func (c *ChromiumOptions) RemoveExtension() *ChromiumOptions {
	c.extensions = []string{}
	return c
}

func (c *ChromiumOptions) SetConfig(configPath string) {
	// 用于配置配置文件的路径，目前只是配置了基础操作 但是没有实现后面的提取操作
	if litedir.FileExists(configPath) {
		c.yamlConf = configPath // 如果路径存在那么久设置到配置文件里面
	} else {
		log.Panicln("file not found:", configPath)
	}
}
