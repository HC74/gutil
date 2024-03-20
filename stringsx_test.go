package gutil

import (
	"fmt"
	"testing"
)

func TestStringIsEmpty(t *testing.T) {

}

func TestThreadPool(t *testing.T) {

}

func TestStringGetHash(t *testing.T) {
	hash := StringGetHash("hello")
	fmt.Println(hash)
}

func TestStringRandom(t *testing.T) {
	hash := StringRandom(10)
	fmt.Println(hash)
}
