package litedir

import (
	"os"
)

func FileExists(pathName string) bool {
	// 判断文件 或者文件夹 在不在
	if _, err := os.Stat(pathName); os.IsNotExist(err) {
		return false
	} else {
		return true
	}
}
