package reqoptions

import (
	"fmt"
	netURL "net/url"
	"strings"
)

// SetParams : 设置参数 不同类型 效果不一样
//
// map[string]any | map[string]string | url.Values --> 最后的参数无序
//
// [][2]any | [][2]string | string  --> 最后的参数按照参数顺序拼接
func (o *Option) SetParams(params any) *Option {
	o._tempParams = params
	return o
}
func (o *Option) GetOrderParam() string {
	params := o._tempParams
	var query string
	if params != nil {

		switch params.(type) {
		case string:
			query = params.(string)
		case [][2]any:
			var tempCache []string

			for _, eachParam := range params.([][2]any) {
				tempCache = append(tempCache, fmt.Sprintf("%v=%v", eachParam[0], eachParam[1]))
			}
			if len(tempCache) > 0 {
				query = netURL.QueryEscape(strings.Join(tempCache, "&"))
			}
		case [][2]string:
			var tempCache []string
			for _, eachParam := range params.([][2]string) {
				tempCache = append(tempCache, fmt.Sprintf("%s=%s", eachParam[0], eachParam[1]))
			}
			if len(tempCache) > 0 {
				query = netURL.QueryEscape(strings.Join(tempCache, "&"))
			}
		}
	}
	return query
}

func (o *Option) GetDisorderParams() netURL.Values {
	// 传入
	params := o._tempParams
	if params != nil {
		parse, _ := netURL.Parse(fmt.Sprintf("https://%s?%s", o.domain, o.query))
		query := parse.Query()
		switch params.(type) {
		case map[string]any:
			for k, v := range params.(map[string]any) {
				query.Set(k, fmt.Sprintf("%v", v))
			}
		case map[string]string:
			for k, v := range params.(map[string]string) {
				query.Set(k, v)
			}
		case netURL.Values:
			query = params.(netURL.Values)
		}

		o.params = &query

		return query
	}

	return nil
}
