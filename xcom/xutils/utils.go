package xutils

import (
	"github.com/beego/beego/logs"
	"github.com/spf13/viper"
)

func BaseDecrypt(data string) string {
	xbase := viper.GetString("server.xbase")
	bytes, err := HttpPostEx(xbase+"/decrypt", []byte(data), map[string]string{})
	if err != nil {
		logs.Error("BaseDecrypt error:", err)
		return ""
	}
	return string(bytes)
}
