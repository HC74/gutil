package gutil

import (
	"bytes"
	"encoding/binary"
)

// Int64ToBytes int64转[]byte
func Int64ToBytes(n int64) (b []byte) {
	bytesBuffer := bytes.NewBuffer([]byte{})
	_ = binary.Write(bytesBuffer, binary.LittleEndian, n)
	return bytesBuffer.Bytes()
}

// BytesToInt64 []byte转换int64
func BytesToInt64(b []byte) (tmp int64) {
	bytesBuffer := bytes.NewBuffer(b)
	_ = binary.Read(bytesBuffer, binary.LittleEndian, &tmp)
	return
}

// Uint64ToBytes uint64转[]byte
func Uint64ToBytes(n uint64) (b []byte) {
	bytesBuffer := bytes.NewBuffer([]byte{})
	_ = binary.Write(bytesBuffer, binary.LittleEndian, n)
	return bytesBuffer.Bytes()
}

// BytesTouInt64 []byte转换uint64
func BytesTouInt64(b []byte) (tmp uint64) {
	bytesBuffer := bytes.NewBuffer(b)
	_ = binary.Read(bytesBuffer, binary.LittleEndian, &tmp)
	return
}

// BytesTouInt8 []byte转换uint8
func BytesTouInt8(b []byte) (tmp uint8) {
	bytesBuffer := bytes.NewBuffer(b)
	_ = binary.Read(bytesBuffer, binary.LittleEndian, &tmp)
	return
}
