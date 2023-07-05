package utils

import (
	"reflect"
)

// 递归获取指针的值类型
func GetValueType(o interface{}) reflect.Type {
	t := reflect.TypeOf(o)
	for t.Kind() == reflect.Ptr {
		t = t.Elem()
	}
	return t
}

// 递归设置指针的值
func SetValue(o interface{}, v interface{}) {
	ov := reflect.ValueOf(o)
	for ov.Type().Kind() == reflect.Ptr {
		ov = ov.Elem()
	}
	if ov.CanAddr() {
		ov.Set(reflect.ValueOf(v))
	}
}

// 递归获取指针的值
func GetValue(o interface{}) interface{} {
	ov := reflect.ValueOf(o)
	for ov.Type().Kind() == reflect.Ptr {
		ov = ov.Elem()
	}
	return ov.Interface()
}
