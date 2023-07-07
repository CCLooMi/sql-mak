package utils

import (
	"encoding/json"
	"fmt"
	"reflect"
	"sql-mak/mysql/entity"
	"testing"
)

func TestCreateSlice(t *testing.T) {
	var uList = &[]*entity.User{}
	lsType := GetValueType(uList)
	lsValue := GetReflectValue(uList)
	fmt.Println(lsType, lsValue.Type())
	//判断lsType是否是数组
	if lsType.Kind() == reflect.Slice {
		//获取数组中元素类型
		eleType := GetType(lsType.Elem())
		fmt.Println(eleType)
		//创建eleType实例
		ele := reflect.New(eleType)
		SetFValues(ele, &[]string{"Username", "Password"}, &[]interface{}{"Seemie", "123456"})
		fmt.Println(ele.Type())
		//在lsValue中添加元素ele
		if lsType.Elem().Kind() == reflect.Ptr {
			lsValue.Set(reflect.Append(lsValue, ele))
		} else {
			lsValue.Set(reflect.Append(lsValue, ele.Elem()))
		}
		//print uList为json string
		fmt.Println(toJSONString(uList))
	}
}
func toJSONString(v interface{}) string {
	b, _ := json.Marshal(v)
	return string(b)
}
func TestCreateStruct(t *testing.T) {
	var u = &entity.User{}
	SetFValues(reflect.ValueOf(u), &[]string{"Username", "Password"}, &[]interface{}{"Seemie", "123456"})
	fmt.Println(toJSONString(u))
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
	fmt.Println(d[0], GetIValue(d[1]), GetIValue(d[2]))
}
func TestSetValue(t *testing.T) {
	a := 1
	b := &a
	c := &b
	SetIValue(a, 2)
	fmt.Println(GetIValue(a), GetIValue(b), GetIValue(c))
	SetIValue(&a, 2)
	fmt.Println(GetIValue(a), GetIValue(b), GetIValue(c))
	SetIValue(b, 3)
	fmt.Println(GetIValue(a), GetIValue(b), GetIValue(c))
	SetIValue(c, 4)
	fmt.Println(GetIValue(a), GetIValue(b), GetIValue(c))
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
