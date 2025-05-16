package audiofileutil

import (
	"os"
	"path/filepath"

	"github.com/duke-git/lancet/strutil"
)

func CheckFileIsExist(filename string) bool {
	var exist = true
	if _, err := os.Stat(filename); os.IsNotExist(err) {
		exist = false
	}
	return exist
}

// 将基础数据路径转换为绝对路径
// @param path 路径
func AbsBasedir(path string) string {
	if strutil.IsBlank(path) {
		panic("Basedir is blank")
	}

	info, err := os.Stat(path)

	if err != nil && !os.IsExist(err) {
		panic("Basedir not exists: " + path)
	}

	if !info.IsDir() {
		panic("Basedir is not dir: " + path)
	}

	pathAbs, err := filepath.Abs(path)
	if err != nil {
		panic(err)
	}

	return pathAbs
}
