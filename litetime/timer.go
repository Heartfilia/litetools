package litetime

import (
	"github.com/Heartfilia/litetools/litestring"
	"log"
	"time"
)

func Timer(name string) func() {
	startTime := time.Now()
	return func() {
		log.Printf("[ %s ] took --> %v\n", litestring.ColorString(name, "green"), time.Since(startTime))
	}
}
