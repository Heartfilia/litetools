package litenet

import (
	"github.com/Heartfilia/litetools/litenet/ua"
	"github.com/Heartfilia/litetools/literand"
)

func GetUA(options ...string) string {
	var platform string
	if options == nil {
		platform = ua.DefaultChoice
	} else {
		if len(options) == 1 {
			platform = options[0]
		} else {
			platform = literand.RandomChoice(options)
		}
	}

	return ua.Options(platform)
}
