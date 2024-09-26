package litetime

import (
	"github.com/Heartfilia/litetools/litestr"
	"log"
	"runtime"
	"strings"
	"time"
)

func Timer() func() {
	name, _ := getFunName(2)
	startTime := time.Now()
	return func() {
		log.Printf("[ %s ] took --> %v\n", litestr.ColorString(name, "green"), time.Since(startTime))
	}
}

func getFunName(l int) (string, int) {
	pc, _, line, _ := runtime.Caller(l)
	name := runtime.FuncForPC(pc).Name()
	split := strings.Split(name, ".")
	//fmt.Printf("第%d层函数,函数名称是:%s\n", l, name)
	if len(split) > 0 {
		return split[len(split)-1], line
	} else {
		return "", 0
	}
}
