package gutil

import (
	"fmt"
	"testing"
)

func TestNewRing(t *testing.T) {
	ring := NewRing[int](2)
	ring.Do(func(t int) {
		fmt.Println(t)
	})
}

func TestSlice(t *testing.T) {
	arr := make([]int, 0, 3)
	prevCap := cap(arr)
	for i := 0; i < 3000; i++ {
		arr = append(arr, 1)
		nextCap := cap(arr)
		if nextCap != prevCap {
			fmt.Println(fmt.Sprintf("%v -> %v", prevCap, nextCap))
			prevCap = nextCap
		}
	}
	fmt.Println(prevCap)
}
