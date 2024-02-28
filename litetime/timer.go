package litetime

import (
	"log"
	"time"
)

func Timer(name string) func() {
	startTime := time.Now()
	return func() {
		log.Printf("[ %s ] took --> %v\n", name, time.Since(startTime))
	}
}
