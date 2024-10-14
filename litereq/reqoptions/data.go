package reqoptions

import (
	"bytes"
	"fmt"
	"log"
	netURL "net/url"
)

// 这个地方 需要兼容的模板比较多

func (o *Option) SetData(data any) *Option {
	o._tempData = data
	return o
}

func (o *Option) GetData() (any, string) {
	// 传入
	data := o._tempData
	if data != nil {
		payload := netURL.Values{}
		switch data.(type) {
		case map[string]any:
			for k, v := range data.(map[string]any) {
				payload.Set(k, fmt.Sprintf("%v", v))
			}
		case map[string]string:
			for k, v := range data.(map[string]string) {
				payload.Set(k, v)
			}
		case [][2]any:
			for _, eachData := range data.([][2]any) {
				payload.Set(fmt.Sprintf("%v", eachData[0]), fmt.Sprintf("%v", eachData[1]))
			}
		case [][2]string:
			for _, eachData := range data.([][2]string) {
				payload.Set(eachData[0], eachData[1])
			}
		case netURL.Values:
			payload = data.(netURL.Values)
		case string:
			items := parseStringParams(data.(string))
			if items != nil {
				for k, v := range items {
					payload.Set(k, v)
				}
			}
		case []byte:
			// 这里需要测试
			var requestBody = new(bytes.Buffer)
			_, _ = requestBody.Write(data.([]byte))
			return requestBody, "bytes"
		default:
			log.Panicln("Params only support <url.Values || map[string]string || map[string]any || string || [][2]any || [][2]string>")
		}

		o.data = &payload

		return payload, "form"
	}

	return nil, ""
}
