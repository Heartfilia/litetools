package litetime

import (
	"github.com/Heartfilia/litetools/litestr"
	"log"
	"time"
)

func Timer(name string) func() {
	startTime := time.Now()
	return func() {
		log.Printf("[ %s ] took --> %v\n", litestr.ColorString(name, "green"), time.Since(startTime))
	}
}
