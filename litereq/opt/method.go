package opt

import "strings"

// 下面优先适配基本的一些操作

const (
	OPT     = "OPTIONS" // 允许客户端查看服务器的性能
	GET     = "GET"     // 请求指定的页面信息，并返回实体主体
	HEAD    = "HEAD"    // 类似于GET请求，响应中没有具体的内容，用于获取报头
	POST    = "POST"    // 向指定资源提交数据并进行处理请求。数据被包含在请求体中，POST请求可能会导致新的资源的建立或已有资源的修改
	PUT     = "PUT"     // 从客户端向服务器传送新的数据到指定的位置中
	DELETE  = "DELETE"  // 请求服务器删除指定的页面
	TRACE   = "TRACE"   // 回显服务器收到的请求，主要用于测试或诊断
	CONNECT = "CONNECT" // 可以开启一个客户端与所请求资源之间的双向沟通的通道。它可以用来创建隧道（tunnel）
	PATCH   = "PATCH"   // 是对请求方式中的PUT补充，用来对已知资源进行局部更新。
)

func (o *Option) SetMethod(method string) *Option {
	md := strings.ToUpper(method)
	switch md {
	case OPT:
		o.method = "OPTIONS"
	case GET:
		o.method = "GET"
	case HEAD:
		o.method = "HEAD"
	case POST:
		o.method = "POST"
	case PUT:
		o.method = "PUT"
	case DELETE:
		o.method = "DELETE"
	case TRACE:
		o.method = "TRACE"
	case CONNECT:
		o.method = "CONNECT"
	case PATCH:
		o.method = "PATCH"
	default:
		o.method = "GET"
	}

	return o
}

func (o *Option) GetMethod() string {
	return o.method
}
