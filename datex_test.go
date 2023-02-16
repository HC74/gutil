package gutil

import (
	"fmt"
	"testing"
	"time"
)

// 将时间戳转换为指定格式的字符串
func FormatTimestamp(timestamp int64, format string) string {
	t := time.Unix(timestamp, 0)
	return t.Format(format)
}

// 将字符串形式的日期时间转换为指定格式的字符串
func FormatDateString(dateStr string, format string) (string, error) {
	t, err := time.Parse("2006-01-02 15:04:05", dateStr)
	if err != nil {
		return "", err
	}
	return t.Format(format), nil
}

func TestD(t *testing.T) {
	// 将时间戳转换为指定格式的字符串
	timestamp := int64(1645077993)
	dateTimeStr := FormatTimestamp(timestamp, "2006-01-02 15:04:05")
	fmt.Println(dateTimeStr) // 输出：2022-02-17 10:59:53

	// 将字符串形式的日期时间转换为指定格式的字符串
	dateStr := "2022-02-17 10:59:53"
	dateTimeStr, err := FormatDateString(dateStr, "2006-01-02 15:04:05")
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(dateTimeStr) // 输出：2022-02-17 10:59:53
}

func TestNow(t *testing.T) {

}
func TestToUnifiedTime(t *testing.T) {

}
