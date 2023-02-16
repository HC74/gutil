package gutil

import (
	"fmt"
	"testing"
)

func TestRIsEmpty(t *testing.T) {
	empty := RIsEmpty("1")
	fmt.Println(empty)
}
func TestRIsPointer(t *testing.T) {
	type A struct {
	}
	a := &A{}
	pointer := RIsPointer(a)
	fmt.Println(pointer)
}
func TestRGetTypeName(t *testing.T) {
	type A struct {
	}
	a := int64(1)
	name := RGetTypeName(a)
	fmt.Println(name)
}
func TestRGetFieldsValues(t *testing.T) {
	type A struct {
		Name string
		Age  int
	}
	a := &A{
		Name: "zs",
		Age:  12,
	}
	values := RGetFieldsValues(a)
	fmt.Println(values)
}
