package utils

import (
	"fmt"
	"testing"
)

func TestSetValue(t *testing.T) {
	a := 1
	b := &a
	c := &b
	SetValue(a, 2)
	fmt.Println(GetValue(a), GetValue(b), GetValue(c))
	SetValue(&a, 2)
	fmt.Println(GetValue(a), GetValue(b), GetValue(c))
	SetValue(b, 3)
	fmt.Println(GetValue(a), GetValue(b), GetValue(c))
	SetValue(c, 4)
	fmt.Println(GetValue(a), GetValue(b), GetValue(c))
}
func TestGetType(t *testing.T) {
	a := 1
	b := &a
	c := &b
	fmt.Println(GetValueType(a))
	fmt.Println(GetValueType(b))
	fmt.Println(GetValueType(c))
}
