package controllers

import (
	"fmt"
	"reflect"
	"testing"
	"unsafe"
)

func TestFuncValuePtr(t *testing.T) {
	t1 := NewSocketController()
	t2 := NewSocketController()
	v1 := reflect.ValueOf(t1.Connect)
	fmt.Println(v1.Interface())
	v2 := reflect.ValueOf(t2.Connect)
	fmt.Println(v2.Interface())

	fmt.Println(unsafe.Sizeof(t2.Connect))
	fmt.Println(unsafe.Sizeof(t1.Connect))
}
