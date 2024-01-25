package main

import (
	"fmt"
	"github.com/Heartfilia/litetools/litetime"
)

func main() {
	t := litetime.Time{
		//Unit: "ms",
		Fmt: true,
	}
	//fmt.Println(t.GetTime().Int())
	//fmt.Println(t.GetTime().Float())
	fmt.Println(t.GetTime().String())

	//stringTime := time.Now().String()
	//newString := strings.Split(stringTime, ".")
	//fmt.Println(newString[0])

}
