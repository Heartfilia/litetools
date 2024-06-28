package opt

import (
	"fmt"
	"log"
	netURL "net/url"
)

func (o *Option) SetParams(params any) *Option {
	o._tempParams = params
	return o
}

func (o *Option) GetParams() netURL.Values {
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
		case [][2]any:
			for _, eachParam := range params.([][2]any) {
				query.Set(fmt.Sprintf("%v", eachParam[0]), fmt.Sprintf("%v", eachParam[1]))
			}
		case [][2]string:
			for _, eachParam := range params.([][2]string) {
				query.Set(eachParam[0], eachParam[1])
			}
		case netURL.Values:
			query = params.(netURL.Values)
		case string:
			items := parseStringParams(params.(string))
			if items != nil {
				for k, v := range items {
					query.Set(k, v)
				}
			}
		default:
			log.Panicln("Params only support <url.Values || map[string]string || map[string]any || string || [][2]any || [][2]string>")
		}

		o.params = &query

		return query
	}

	return nil
}
