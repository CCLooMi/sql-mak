package utils

import (
	"fmt"
	"reflect"
	"sql-mak/mysql/entity"
	"testing"
)

func TestCreateStruct(t *testing.T) {
	u := &entity.User{}
	u = nil
	fmt.Println(reflect.TypeOf(u))
	v := GetValue(u)
	if v == nil {
		fmt.Println("nil")
	}
	tp := GetValueType(u)
	fmt.Println(tp)
	tv := reflect.New(tp)
	fmt.Println(tv.Type().String())
	reflect.ValueOf(u).Set(tv)
	fmt.Println(u)
}
func TestCreateValue(t *testing.T) {
	a := 1
	b := &a
	c := &b
	d := []interface{}{a, b, c}
	dt := GetValueType(&d)
	fmt.Println(dt)
	v := reflect.New(dt)
	fmt.Println(v.Type())
	dt = reflect.SliceOf(dt.Elem())
	fmt.Println(dt)
	dtv := reflect.MakeSlice(dt, 2, 2)
	dtv.Index(1).Set(reflect.ValueOf(9))
	fmt.Println(dtv.Interface())
	fmt.Println(dtv.Type())

	dv := reflect.ValueOf(d)
	fmt.Println(dv.Type())
	dv.Index(0).Set(reflect.ValueOf(9))
	fmt.Println(d[0], GetValue(d[1]), GetValue(d[2]))
}
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
	d := []interface{}{a, b, c}
	fmt.Println(GetValueType(a))
	fmt.Println(GetValueType(b))
	fmt.Println(GetValueType(c))
	fmt.Println(GetValueType(d))
	fmt.Println(GetValueType(&d))
}
