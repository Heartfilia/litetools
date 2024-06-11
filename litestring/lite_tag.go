package litestring

import (
	"fmt"
	"runtime"
	"strings"
)

func D() string {
	name, line := getFunName(2)
	return ColorString(fmt.Sprintf("[ DEBUG  - %10s<%d>]", name, line), "blue")
}

func I() string {
	name, line := getFunName(2)
	return ColorString(fmt.Sprintf("[  INFO  - %10s<%d>]", name, line), "cyan")
}

func S() string {
	name, line := getFunName(2)
	return ColorString(fmt.Sprintf("[SUCCESS - %10s<%d>]", name, line), "green")
}

func W() string {
	name, line := getFunName(2)
	return ColorString(fmt.Sprintf("[  WARN  - %10s<%d>]", name, line), "yellow")
}

func E() string {
	name, line := getFunName(2)
	return ColorString(fmt.Sprintf("[ ERROR  - %10s<%d>]", name, line), "red")
}

func getFunName(l int) (string, int) {
	pc, _, line, _ := runtime.Caller(l)
	name := runtime.FuncForPC(pc).Name()
	split := strings.Split(name, ".")
	//fmt.Printf("第%d层函数,函数名称是:%s\n", l, name)
	//return split[len(split)-1]
	return split[len(split)-1], line
}
