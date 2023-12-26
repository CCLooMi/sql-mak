package utils

import (
	"reflect"
	"strings"
)

type EntityInfo struct {
	TableName  string
	PrimaryKey string
	Fields     []string
	Columns    []string
	Tags       []reflect.StructTag
	FCMap      map[string]string
	CFMap      map[string]string
	IExpMap    map[string]string
}

var infoCache = make(map[string]*EntityInfo)
var nameCache = make(map[string]string)

// 传入结构体指针或结构体，返回实体信息
func GetEntityInfo(table interface{}) EntityInfo {
	var tableName string
	switch t := table.(type) {
	case string:
		tableName = t
	default:
		tableName = TableName(table)
	}
	ei := infoCache[tableName]
	if ei != nil {
		return *ei
	}
	primaryKey := "id"
	fields := GetFields(table)
	fL := len(fields)
	fnames := make([]string, fL)
	fcmap := make(map[string]string)
	cfmap := make(map[string]string)
	iExpMap := make(map[string]string)
	cols := make([]string, fL)
	tags := make([]reflect.StructTag, fL)
	for i := 0; i < fL; i++ {
		fi := fields[i]
		fn := fi.Name
		cn, orm := getTableColumnName(fi)
		//判断orm中是否有primaryKey
		if strings.Contains(orm, "primaryKey") {
			primaryKey = cn
		}
		fcmap[fn] = cn
		cfmap[cn] = fn
		insertExp := fi.Tag.Get("insertExp")
		if insertExp != "" {
			iExpMap[fn] = insertExp
			iExpMap[cn] = insertExp
		}
		fnames[i] = fn
		cols[i] = cn
		tags[i] = fi.Tag
	}
	ei = &EntityInfo{
		TableName:  tableName,
		PrimaryKey: primaryKey,
		Fields:     fnames,
		Columns:    cols,
		Tags:       tags,
		FCMap:      fcmap,
		CFMap:      cfmap,
		IExpMap:    iExpMap,
	}
	infoCache[tableName] = ei
	return *ei
}
func getTableColumnName(field reflect.StructField) (string, string) {
	orm := field.Tag.Get("orm")
	name := field.Tag.Get("column")
	jname := field.Tag.Get("json")
	if name != "" {
		return name, orm
	}
	if jname != "" {
		return jname, orm
	}
	return field.Name, orm
}
func TableName(table interface{}) string {
	t := reflect.TypeOf(table)
	name := t.Name()
	if t.Kind() == reflect.Ptr {
		nm := nameCache[t.String()]
		if nm != "" {
			return nm
		}
		method, ok := t.MethodByName("TableName")
		if ok {
			if method.Type.In(0) == t {
				name = method.Func.Call([]reflect.Value{reflect.ValueOf(table)})[0].String()
			}
		}
		nameCache[t.String()] = name
		return name
	}
	t = reflect.PtrTo(t)
	nm := nameCache[t.String()]
	if nm != "" {
		return nm
	}
	method, ok := t.MethodByName("TableName")
	if ok {
		if method.Type.In(0) == t {
			/*
				Create an instance of a pointer type
				and assign the original non-pointer type parameter
				to the newly created instance of the pointer type.
			*/
			ptr := reflect.New(t.Elem())
			ptr.Elem().Set(reflect.ValueOf(table))
			name = method.Func.Call([]reflect.Value{ptr})[0].String()
		}
	}
	nameCache[t.String()] = name
	return name
}
func GetFieldNames(o interface{}) []string {
	t := reflect.TypeOf(o)
	fds := make([]string, 0)
	for i := 0; i < t.NumField(); i++ {
		fd := t.Field(i)
		if fd.Type.Kind() == reflect.Struct {
			// 如果字段是结构体类型，需要判断是否是指针类型
			if fd.Type.Kind() == reflect.Ptr {
				// 如果是指针类型，需要解引用后再递归调用
				fds = append(fds, GetFieldNames(reflect.ValueOf(o).Field(i).Elem().Interface())...)
			} else if fd.Anonymous {
				// 如果是匿名结构体类型，直接递归调用
				fds = append(fds, GetFieldNames(reflect.ValueOf(o).Field(i).Interface())...)
			} else {
				fds = append(fds, fd.Name)
			}
			continue
		}
		if fd.Anonymous && fd.Type.Kind() == reflect.Interface {
			continue
		}
		fds = append(fds, fd.Name)
	}
	return fds
}
func GetFields(o interface{}) []reflect.StructField {
	t := reflect.TypeOf(o)
	if t.Kind() == reflect.Ptr {
		t = t.Elem()
		o = reflect.ValueOf(o).Elem().Interface()
	}
	fds := make([]reflect.StructField, 0)
	for i := 0; i < t.NumField(); i++ {
		fd := t.Field(i)
		if fd.Type.Kind() == reflect.Struct {
			// 如果是匿名结构体类型，直接递归调用
			if fd.Anonymous {
				fds = append(fds, GetFields(reflect.ValueOf(o).Field(i).Interface())...)
			} else {
				fds = append(fds, fd)
			}
			continue
		}
		if fd.Anonymous && fd.Type.Kind() == reflect.Interface {
			continue
		}
		fds = append(fds, fd)
	}
	return fds
}

func GetFieldValue(o interface{}, field string) interface{} {
	value := reflect.ValueOf(o)
	if value.Kind() == reflect.Ptr {
		value = value.Elem()
	}
	fieldValue := value.FieldByName(field)
	return fieldValue.Interface()
}
func GetFieldValueByIndex(o interface{}, index int) interface{} {
	value := reflect.ValueOf(o)
	if value.Kind() == reflect.Ptr {
		value = value.Elem()
	}
	fieldValue := value.Field(index)
	return fieldValue.Interface()
}
func GetFieldValues(o interface{}, fields []string) []interface{} {
	value := reflect.ValueOf(o)
	if value.Kind() == reflect.Ptr {
		value = value.Elem()
	}
	fieldValues := make([]interface{}, len(fields))
	for i := 0; i < len(fields); i++ {
		fieldValues[i] = value.FieldByName(fields[i]).Interface()
	}
	return fieldValues
}
