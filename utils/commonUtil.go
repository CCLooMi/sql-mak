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

// 递归获取指针的值类型
func GetType(t reflect.Type) reflect.Type {
	for t.Kind() == reflect.Ptr {
		t = t.Elem()
	}
	return t
}

// 递归设置指针的值
func SetIValue(o interface{}, v interface{}) {
	ov := reflect.ValueOf(o)
	for ov.Type().Kind() == reflect.Ptr {
		ov = ov.Elem()
	}
	if ov.CanAddr() {
		ov.Set(reflect.ValueOf(v))
	}
}

// 递归设置指针的值
func SetValue(ov reflect.Value, v interface{}) {
	for ov.Type().Kind() == reflect.Ptr {
		ov = ov.Elem()
	}
	if ov.CanAddr() {
		ov.Set(reflect.ValueOf(v))
	}
}

// 设置指针的值的属性值
func SetFValues(ov reflect.Value, fields *[]string, values *[]interface{}) {
	for ov.Type().Kind() == reflect.Ptr {
		ov = ov.Elem()
	}
	if ov.CanAddr() {
		vL := len(*values)
		var field reflect.Value
		var targetType reflect.Type
		var value reflect.Value
		for i, fi := range *fields {
			if i < vL {
				field = ov.FieldByName(fi)
				targetType = field.Type()
				value = reflect.ValueOf((*values)[i])
				if value.Type().ConvertibleTo(targetType) {
					field.Set(value.Convert(targetType))
				}
			}
		}
	}
}

// 递归获取指针的值
func GetReflectValue(o interface{}) reflect.Value {
	ov := reflect.ValueOf(o)
	for ov.Type().Kind() == reflect.Ptr {
		ov = ov.Elem()
	}
	return ov
}

// 递归获取指针的值
func GetIValue(o interface{}) interface{} {
	ov := reflect.ValueOf(o)
	for ov.Type().Kind() == reflect.Ptr {
		ov = ov.Elem()
	}
	return ov.Interface()
}

// 递归获取指针的值
func GetValue(ov reflect.Value) reflect.Value {
	for ov.Type().Kind() == reflect.Ptr {
		ov = ov.Elem()
	}
	return ov
}
