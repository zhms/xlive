package xutils

import (
	"crypto/md5"
	"encoding/hex"
)

// 获取字符串md5值 eg:test -> 098f6bcd4621d373cade4e832627b4f6
func Md5(str string) string {
	h := md5.New()
	h.Write([]byte(str))
	return hex.EncodeToString(h.Sum(nil))
}
