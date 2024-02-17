package utils

import (
	"os"

	"github.com/beego/beego/logs"
)

// 读取文件全部文本
func ReadAllText(path string) string {
	bytes, err := os.ReadFile(path)
	if err != nil {
		logs.Error(err)
		return ""
	}
	return string(bytes)
}

// 写入文件全部文本
func WriteText(path string, text string) {
	err := os.WriteFile(path, []byte(text), 0644)
	if err != nil {
		logs.Error(err)
	}
}
