package reqoptions

import (
	"encoding/json"
	"github.com/Heartfilia/litetools/litestr"
	"log"
)

// SetJson : recommend struct
//
// 使用 map[string]any 或者 struct 结构的数据 推荐后面的
func (o *Option) SetJson(object any) *Option {
	if object != nil {
		marshal, err := json.Marshal(object)
		if err != nil {
			log.Panicln(litestr.E(), "error json object:", err)
			return o
		}
		o.json = marshal
	}
	return o
}

func (o *Option) GetJson() []byte {
	return o.json
}
