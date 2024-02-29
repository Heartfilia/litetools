package litedir

import (
	"encoding/json"
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

func FileSaver(str, pathString string) bool {
	if !FileExists(pathString) {
		dstFile, err := os.Create(pathString)
		if err != nil {
			return false
		}
		defer func(file *os.File) {
			err = file.Close()
			if err != nil {
				return
			}
		}(dstFile)
		_, err = dstFile.WriteString(str)
		if err != nil {
			return false
		}
	}
	return true
}

func fileRead(filePath string) []byte {
	if !FileExists(filePath) {
		return nil
	}

	file, err := os.ReadFile(filePath)
	if err != nil {
		return nil
	}
	return file
}

func FileJsonLoader(jsonPath string) map[string][]string {
	if !FileExists(jsonPath) {
		return nil
	}
	var data map[string][]string
	file := fileRead(jsonPath)
	if file == nil {
		return nil
	}
	err := json.Unmarshal(file, &data)
	if err != nil {
		return nil
	}
	return data
}
