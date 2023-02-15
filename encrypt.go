package gutil

import (
	"golang.org/x/crypto/bcrypt"
)

// BcryptEncrypt Bcrypt加密 @params password : 被加密对象
func BcryptEncrypt(password string) string {
	bytes, _ := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(bytes)
}

// BcryptCheck Bcrypt解密 @params password : 被解密对象
func BcryptCheck(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}
