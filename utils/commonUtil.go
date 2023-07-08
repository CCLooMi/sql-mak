package utils

import (
	"database/sql"
	"log"
	"reflect"
)

// 递归获取指针的值类型
func GetType(t reflect.Type) reflect.Type {
	for t.Kind() == reflect.Ptr {
		t = t.Elem()
	}
	return t
}

// 递归设置指针的值
func SetIValue(o interface{}, v interface{}) {
	SetValue(reflect.ValueOf(o), v)
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
func SetValuesWithRows(ov reflect.Value, fields *[]string, rs *sql.Rows) {
	for ov.Type().Kind() == reflect.Ptr {
		ov = ov.Elem()
	}
	if ov.CanAddr() {
		values := make([]interface{}, len(*fields))
		for i, fi := range *fields {
			values[i] = ov.FieldByName(fi).Addr().Interface()
		}
		rs.Scan(values...)
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
				value = GetValue(reflect.ValueOf((*values)[i]))
				if value.Type().ConvertibleTo(targetType) {
					field.Set(value.Convert(targetType))
				} else {
					log.Printf("Can't convert %s to %s", value.Type(), targetType)
				}
			}
		}
	}
}
func SetValues(ov reflect.Value, fvs ...interface{}) {
	for ov.Type().Kind() == reflect.Ptr {
		ov = ov.Elem()
	}
	if ov.CanAddr() {
		vl := len(fvs)
		var field reflect.Value
		var targetType reflect.Type
		var value reflect.Value
		for i := 0; i < vl; i += 2 {
			field = ov.FieldByName(fvs[i].(string))
			targetType = field.Type()
			value = GetValue(reflect.ValueOf(fvs[i+1]))
			if value.Type().ConvertibleTo(targetType) {
				field.Set(value.Convert(targetType))
			} else {
				log.Printf("Can't convert %s to %s", value.Type(), targetType)
			}
		}
	}
}

// 递归获取指针的值
func GetReflectValue(o interface{}) reflect.Value {
	return GetValue(reflect.ValueOf(o))
}

// 递归获取指针的值
func GetIValue(o interface{}) interface{} {
	return GetValue(reflect.ValueOf(o)).Interface()
}

// 递归获取指针的值
func GetValue(ov reflect.Value) reflect.Value {
	for ov.Type().Kind() == reflect.Ptr {
		ov = ov.Elem()
	}
	for ov.Type().Kind() == reflect.Interface {
		ov = reflect.ValueOf(ov.Interface())
	}
	return ov
}
