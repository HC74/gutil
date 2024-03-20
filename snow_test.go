package gutil

import (
	"fmt"
	"testing"
)

func TestNewIdGeneratorOptions(t *testing.T) {
	SnowInit(11)
	id := NextId()
	fmt.Println(id)
}
