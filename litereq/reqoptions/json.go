package reqoptions

import (
	"encoding/json"
	"github.com/Heartfilia/litetools/litestr"
	"log"
)

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
