package opt

// 作为请求参数的配置选项
// 先把基础的一些配置开发了 其它配置后面再优化添加

type Option struct {
	Params         string // 先占位 后续更新
	Headers        string // 先占位 后续更新
	Cookies        string // 先占位 后续更新
	Data           string // 先占位 后续更新
	Json           string // 先占位 后续更新
	Verify         string // 先占位 后续更新
	Files          string // 先占位 后续更新
	Proxy          string // 先占位 后续更新
	Method         string // 默认GET
	Timeout        int    // ms  单位为毫秒
	AllowRedirects bool
	Stream         string // 先占位 后续更新
	Auth           string // 先占位 后续更新
	Cert           string // 先占位 后续更新
}

func NewOption() *Option {
	return &Option{
		Method: "GET",
	}
}

func (o *Option) Reset() {

}
