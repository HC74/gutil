package gutil

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"crypto/des"
	"crypto/rand"
	"golang.org/x/crypto/bcrypt"
	"io"
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

// SecureAESEncrypt ASE加密
func SecureAESEncrypt(key, text []byte) ([]byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}
	ciphertext := make([]byte, aes.BlockSize+len(text))
	iv := ciphertext[:aes.BlockSize]
	if _, err := io.ReadFull(rand.Reader, iv); err != nil {
		return nil, err
	}
	cfb := cipher.NewCFBEncrypter(block, iv)
	cfb.XORKeyStream(ciphertext[aes.BlockSize:], text)
	return ciphertext, nil
}

// SecureAESDecrypt ASE解密
func SecureAESDecrypt(key, ciphertext []byte) ([]byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}
	iv := ciphertext[:aes.BlockSize]
	ciphertext = ciphertext[aes.BlockSize:]
	cfb := cipher.NewCFBDecrypter(block, iv)
	cfb.XORKeyStream(ciphertext, ciphertext)
	return ciphertext, nil
}

// SecureDESEncrypt DES加密
func SecureDESEncrypt(plaintext, key []byte) ([]byte, error) {
	// 创建DES加密器
	block, err := des.NewCipher(key)
	if err != nil {
		return nil, err
	}
	// 使用CBC模式加密，IV长度必须等于块大小
	iv := []byte("12345678")
	mode := cipher.NewCBCEncrypter(block, iv)
	// 对数据进行填充
	plaintext = PKCS5Padding(plaintext, block.BlockSize())
	// 加密数据
	ciphertext := make([]byte, len(plaintext))
	mode.CryptBlocks(ciphertext, plaintext)
	return ciphertext, nil
}

// SecureDESDecrypt DES解密
func SecureDESDecrypt(ciphertext, key []byte) ([]byte, error) {
	// 创建DES解密器
	block, err := des.NewCipher(key)
	if err != nil {
		return nil, err
	}
	// 使用CBC模式解密，IV长度必须等于块大小
	iv := []byte("12345678")
	mode := cipher.NewCBCDecrypter(block, iv)
	// 解密数据
	plaintext := make([]byte, len(ciphertext))
	mode.CryptBlocks(plaintext, ciphertext)
	// 去除填充数据
	plaintext = PKCS5UnPadding(plaintext)
	return plaintext, nil
}

// PKCS5Padding 对数据进行填充
func PKCS5Padding(data []byte, blockSize int) []byte {
	padding := blockSize - len(data)%blockSize
	padtext := bytes.Repeat([]byte{byte(padding)}, padding)
	return append(data, padtext...)
}

// PKCS5UnPadding 去除填充数据
func PKCS5UnPadding(data []byte) []byte {
	length := len(data)
	unpadding := int(data[length-1])
	return data[:(length - unpadding)]
}
