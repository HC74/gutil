package gutil

import "time"

// Now 获取当前时间 默认为 YYYY-MM-dd HH:ss:mm 格式
func Now() string {
	timeNow := time.Now().Format("2006-01-02 15:04:05")
	return timeNow
}

// ToUnifiedTime 统一返回 YYYY-MM-dd HH:ss:mm
func ToUnifiedTime(v time.Time) string {
	dateTimeStr := v.Format("2006-01-02 15:04:05")
	return dateTimeStr
}
