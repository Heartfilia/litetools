package opt

import "encoding/json"

func gen(x any) {
	json.Marshal(x)
}
