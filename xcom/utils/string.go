package utils

import "unicode"

// 判断字符串是否包含小写字母
func StrContainsLower(str string) bool {
	for _, char := range str {
		if unicode.IsLower(char) {
			return true
		}
	}
	return false
}

// 判断字符串是否包含大写字母
func StrContainsUpper(str string) bool {
	for _, char := range str {
		if unicode.IsUpper(char) {
			return true
		}
	}
	return false
}

// 判断字符串是否包含数字
func StrContainsDigit(str string) bool {
	for _, char := range str {
		if unicode.IsDigit(char) {
			return true
		}
	}
	return false
}
