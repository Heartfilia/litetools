package litedir

import (
	"encoding/json"
	"github.com/Heartfilia/litetools/utils/types"
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

func FileJsonLoader(jsonPath string) types.ConfigJson {
	if !FileExists(jsonPath) {
		return types.ConfigJson{}
	}
	var data types.ConfigJson
	file := fileRead(jsonPath)
	if file == nil {
		return types.ConfigJson{}
	}
	err := json.Unmarshal(file, &data)
	if err != nil {
		return types.ConfigJson{}
	}
	return data
}
