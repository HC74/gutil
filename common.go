package gutil

import (
	"crypto/md5"
	"encoding/hex"
)

// DataSlice 数据分片
func DataSlice[T any](v []T, size int) [][]T {
	count := len(v) / size
	var results [][]T
	var temp = 0
	var tempSize = size
	for i := 0; i < count; i++ {
		results = append(results, v[temp:tempSize])
		temp += size
		tempSize += size
	}
	if len(v)%size != 0 {
		results = append(results, v[temp:])
	}
	return results
}

// Md5 MD5
func Md5(str string) string {
	hash := md5.Sum([]byte(str))
	md5str := hex.EncodeToString(hash[:])
	return md5str
}
