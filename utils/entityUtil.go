package utils

import "reflect"

type EntityInfo struct {
	TableName string
	Fields    []string
	Columns   []string
	Tags      []reflect.StructTag
	FCMap     map[string]string
	CFMap     map[string]string
}

var infoCache = make(map[string]*EntityInfo)
var nameCache = make(map[string]string)

// 传入结构体指针或结构体，返回实体信息
func GetEntityInfo(table interface{}) EntityInfo {
	tableName := TableName(table)
	ei := infoCache[tableName]
	if ei != nil {
		return *ei
	}
	fields := GetFields(table)
	fL := len(fields)
	fnames := make([]string, fL)
	fcmap := make(map[string]string)
	cfmap := make(map[string]string)
	cols := make([]string, fL)
	tags := make([]reflect.StructTag, fL)
	for i := 0; i < fL; i++ {
		fi := fields[i]
		fn := fi.Name
		cn := getTableColumnName(fi)
		fcmap[fn] = cn
		cfmap[cn] = fn
		fnames[i] = fn
		cols[i] = cn
		tags[i] = fi.Tag
	}
	ei = &EntityInfo{
		TableName: tableName,
		Fields:    fnames,
		Columns:   cols,
		Tags:      tags,
		FCMap:     fcmap,
		CFMap:     cfmap,
	}
	infoCache[tableName] = ei
	return *ei
}
func getTableColumnName(field reflect.StructField) string {
	name := field.Tag.Get("column")
	if name != "" {
		return name
	}
	return field.Name
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
			name = method.Func.Call([]reflect.Value{reflect.ValueOf(table)})[0].String()
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
		name = method.Func.Call([]reflect.Value{reflect.New(t).Elem()})[0].String()
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
