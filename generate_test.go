package gutil

import (
	"fmt"
	"testing"
	"time"
)

func TestA(t *testing.T) {
	var ss []string
	for i := 0; i < 100000; i++ {
		id := GuidID()
		for _, s := range ss {
			if id == s {
				panic(1)
			}
		}
		ss = append(ss, id)
	}
}

func TestB(t *testing.T) {
	fmt.Println(time.Now().Location().String())
}
