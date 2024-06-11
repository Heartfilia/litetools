package opt

// 作为请求参数的配置选项

type Option struct {
	Params         string // 先占位 后续更新
	Headers        string // 先占位 后续更新
	Cookies        string // 先占位 后续更新
	Data           string // 先占位 后续更新
	Json           string // 先占位 后续更新
	Verify         string // 先占位 后续更新
	Files          string // 先占位 后续更新
	Proxy          string // 先占位 后续更新
	Method         string
	Timeout        string // 先占位 后续更新
	AllowRedirects bool
	Stream         string // 先占位 后续更新
	Auth           string // 先占位 后续更新
	Cert           string // 先占位 后续更新
}
