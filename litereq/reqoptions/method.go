package reqoptions

import "strings"

// 下面优先适配基本的一些操作

const (
	options = "OPTIONS" // 允许客户端查看服务器的性能
	get     = "GET"     // 请求指定的页面信息，并返回实体主体
	head    = "HEAD"    // 类似于GET请求，响应中没有具体的内容，用于获取报头
	post    = "POST"    // 向指定资源提交数据并进行处理请求。数据被包含在请求体中，POST请求可能会导致新的资源的建立或已有资源的修改
	put     = "PUT"     // 从客户端向服务器传送新的数据到指定的位置中
	delete  = "DELETE"  // 请求服务器删除指定的页面
	trace   = "TRACE"   // 回显服务器收到的请求，主要用于测试或诊断
	connect = "CONNECT" // 可以开启一个客户端与所请求资源之间的双向沟通的通道。它可以用来创建隧道（tunnel）
	patch   = "PATCH"   // 是对请求方式中的PUT补充，用来对已知资源进行局部更新。
)

func (o *Option) SetMethod(method string) *Option {
	md := strings.ToUpper(method)
	switch md {
	case options:
		o.method = "OPTIONS"
	case get:
		o.method = "GET"
	case head:
		o.method = "HEAD"
	case post:
		o.method = "POST"
	case put:
		o.method = "PUT"
	case delete:
		o.method = "DELETE"
	case trace:
		o.method = "TRACE"
	case connect:
		o.method = "CONNECT"
	case patch:
		o.method = "PATCH"
	default:
		o.method = "GET"
	}

	return o
}

func (o *Option) GetMethod() string {
	return o.method
}
