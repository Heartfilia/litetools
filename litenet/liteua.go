package litenet

import (
	"github.com/Heartfilia/litetools/litenet/ua"
	"math/rand"
	"time"
)

func GetUA(options ...string) string {
	var platform string
	if options == nil {
		platform = "chrome"
	} else {
		if len(options) == 1 {
			platform = options[0]
		} else {
			r := rand.New(rand.NewSource(time.Now().Unix()))
			newInd := r.Intn(len(options))
			platform = options[newInd]
		}
	}

	return ua.Options(platform)
}
