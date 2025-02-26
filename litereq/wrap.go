package litereq

import "net/http"

type Wrap[T any] struct {
	Request  *http.Request  // 原始请求信息
	Response *http.Response // 响应信息港i
	Curl     string         // 忽略 拼接的curl
	RespStr  string         // 结果的json类型数据
	Data     T
	Err      error
}
