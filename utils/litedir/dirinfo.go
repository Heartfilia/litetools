package litedir

import (
	"os"
	"path"
)

func LiteDir() string {
	dir, err := os.UserCacheDir()
	if err != nil {
		// 其实不可能报错的
		return ""
	}

	liteTools := path.Join(dir, "lite-tools")
	if !FileExists(liteTools) {
		_ = os.Mkdir(liteTools, 0777)
	}

	return liteTools
}
