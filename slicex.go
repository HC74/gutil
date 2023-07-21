package gutil

// Contains 判断切片中是否包含某元素
func Contains(slice []string, element string) bool {
	for _, e := range slice {
		if e == element {
			return true
		}
	}
	return false
}
