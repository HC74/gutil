package gutil

import "strings"

// StringIsEmpty 判断字符串是否为空
func StringIsEmpty(v string) bool {
	newName := strings.TrimSpace(v)
	if len(newName) == 0 {
		return true
	}
	return false
}
