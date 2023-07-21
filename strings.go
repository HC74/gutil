package gutil

import (
	"crypto/sha256"
	"encoding/hex"
	"math/rand"
	"strings"
)

// StringIsEmpty 判断字符串是否为空
func StringIsEmpty(v string) bool {
	newName := strings.TrimSpace(v)
	if len(newName) == 0 {
		return true
	}
	return false
}

// StringGetHash 计算字符串的哈希值
func StringGetHash(s string) string {
	h := sha256.New()
	h.Write([]byte(s))
	return hex.EncodeToString(h.Sum(nil))
}

// StringRandom 生成随机字符串
func StringRandom(n int) string {
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	b := make([]byte, n)
	if _, err := rand.Read(b); err != nil {
		panic(err)
	}
	for i := range b {
		b[i] = charset[b[i]%byte(len(charset))]
	}
	return string(b)
}
