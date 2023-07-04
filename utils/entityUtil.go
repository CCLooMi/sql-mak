package utils

import "reflect"

type EntityInfo struct {
	TableName string
	Fields    []string
	Columns   []string
	FCMap     map[string]string
	CFMap     map[string]string
}

var infoCache = make(map[string]*EntityInfo)
var nameCache = make(map[string]string)

// 传入结构体指针或结构体，返回实体信息
func GetEntityInfo(table interface{}) *EntityInfo {
	tableName := TableName(table)
	ei := infoCache[tableName]
	if ei != nil {
		return ei
	}

	fields := GetFields(table)
	fnames := make([]string, len(fields))
	fcmap := make(map[string]string)
	cfmap := make(map[string]string)
	cols := make([]string, len(fields))
	for i := 0; i < len(fields); i++ {
		fn := fields[i].Name
		cn := getTableColumnName(fields[i])
		fcmap[fn] = cn
		cfmap[cn] = fn
		fnames[i] = fn
		cols[i] = cn
	}
	ei = &EntityInfo{
		TableName: tableName,
		Fields:    fnames,
		Columns:   cols,
		FCMap:     fcmap,
		CFMap:     cfmap,
	}
	infoCache[tableName] = ei
	return ei
}
func getTableColumnName(field reflect.StructField) string {
	name := field.Tag.Get("column")
	if name == "" {
		return name
	}
	return field.Name
}

func TableName(table interface{}) string {
	t := reflect.TypeOf(table)
	name := t.Name()
	if t.Kind() == reflect.Ptr { //指针
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
	//获取值地址类型、
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
		//解引用
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

func GetFieldValue(o interface{}, field reflect.StructField) interface{} {
	// 将interface{}类型转换为具体的结构体类型
	value := reflect.ValueOf(o)
	// 如果o是指针类型，则获取指针所指向的元素
	if value.Kind() == reflect.Ptr {
		value = value.Elem()
	}
	// 获取指定字段的值
	fieldValue := value.FieldByName(field.Name)
	// 返回字段的值
	return fieldValue.Interface()
}
func GetFieldValueByIndex(o interface{}, index int) interface{} {
	// 将interface{}类型转换为具体的结构体类型
	value := reflect.ValueOf(o)
	// 如果o是指针类型，则获取指针所指向的元素
	if value.Kind() == reflect.Ptr {
		value = value.Elem()
	}
	// 获取指定字段的值
	fieldValue := value.Field(index)
	// 返回字段的值
	return fieldValue.Interface()
}
func GetFieldValues(o interface{}, fields []reflect.StructField) []interface{} {
	// 将interface{}类型转换为具体的结构体类型
	value := reflect.ValueOf(o)
	// 如果o是指针类型选拿指针所指向的元素
	if value.Kind() == reflect.Ptr {
		value = value.Elem()
	}
	// 获取指定字段的值
	fieldValues := make([]interface{}, len(fields))
	for i := 0; i < len(fields); i++ {
		fieldValues[i] = value.FieldByName(fields[i].Name).Interface()
	}
	// 返回字段的值
	return fieldValues
}

func GetFieldValuesByIndexs(o interface{}, indexs []int) []interface{} {
	// 将interface{}类型转换为具体的结构体类型
	value := reflect.ValueOf(o)
	// 如果o是指针类型选拿指针所指向的元素
	if value.Kind() == reflect.Ptr {
		value = value.Elem()
	}
	// 获取指定字段的值
	fieldValues := make([]interface{}, len(indexs))
	for i := 0; i < len(indexs); i++ {
		fieldValues[i] = value.Field(indexs[i]).Interface()
	}
	// 返回字段的值
	return fieldValues
}

func GetNFieldValues(o interface{}, n int) []interface{} {
	// 将interface{}类型转换为具体的结构体类型
	value := reflect.ValueOf(o)
	// 如果o是指针类型选拿指针所指向的元素
	if value.Kind() == reflect.Ptr {
		value = value.Elem()
	}
	// 获取n个字段的值
	fieldValues := make([]interface{}, n)
	for i := 0; i < n; i++ {
		fieldValues[i] = value.Field(i).Interface()
	}
	// 返回字段的值
	return fieldValues
}
