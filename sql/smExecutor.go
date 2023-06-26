package sql

import (
	"fmt"
	"reflect"

	"database/sql"
)

type SQLSMExecutor struct {
	SQLExecutor
	sm *SQLSM
}

type PageDataBean struct {
	TotalNumber int64
	Data        interface{}
	Headers     []string
}

type PageCSVDataBean struct {
	TotalNumber int64
	Data        [][]string
}

type ByPageFilter interface {
	DoFilter(rs *sql.Rows)
}

func NewSQLSMExecutor(sm *SQLSM) *SQLSMExecutor {
	return &SQLSMExecutor{SQLExecutor: *NewSQLExecutor(&sm.AbstractSQLGod), sm: sm}
}

func (e *SQLSMExecutor) SaveResultToObject(targetObject interface{}) *SQLSMExecutor {
	e.ExtractorResultSet(func(rs *sql.Rows) interface{} {
		if m, ok := targetObject.(map[string]interface{}); ok {
			RowsToMap(rs, m)
		} else {
			RowsToStruct(rs, targetObject)
		}
		return nil
	})
	return e
}

func (e *SQLSMExecutor) SaveColumnToObject(labelColumnIndex, valueColumnIndex int, targetObject interface{}) *SQLSMExecutor {
	e.ExtractorResultSet(func(rs *sql.Rows) interface{} {
		columns, _ := rs.Columns()
		if len(columns) == 0 {
			return nil
		}
		values := make([]interface{}, len(columns))
		for i := range values {
			var val interface{}
			values[i] = &val
		}
		for rs.Next() {
			rs.Scan(values...)
			if m, ok := targetObject.(map[string]interface{}); ok {
				m[columns[labelColumnIndex]] = values[valueColumnIndex]
			} else {
				e.SetObjectFieldValue(targetObject, columns[labelColumnIndex], values[valueColumnIndex])
			}
		}
		return nil
	})
	return e
}

func (e *SQLSMExecutor) SaveColumnToObjectString(labelColumnName, valueColumnName string, targetObject interface{}) *SQLSMExecutor {
	e.ExtractorResultSet(func(rs *sql.Rows) interface{} {
		columns, _ := rs.Columns()
		if len(columns) == 0 {
			return nil
		}
		var labelIdx = -1
		var valueIdx = -1
		values := make([]interface{}, len(columns))
		for i := range values {
			var val interface{}
			values[i] = &val
			if columns[i] == labelColumnName {
				labelIdx = i
			} else if columns[i] == valueColumnName {
				valueIdx = i
			}
		}
		for rs.Next() {
			rs.Scan(values...)
			if m, ok := targetObject.(map[string]interface{}); ok {
				m[columns[labelIdx]] = values[valueIdx]
			} else {
				e.SetObjectFieldValue(targetObject, columns[labelIdx], values[valueIdx])
			}
		}
		return nil
	})
	return e
}

func (e *SQLSMExecutor) GetResultAsObject(elementType reflect.Type) interface{} {
	panic("not implemented")
}

func (e *SQLSMExecutor) GetResultColumnAsObject(labelColumn, valueColumn int, elementType reflect.Type) interface{} {
	return e.ExtractorResultSet(func(rs *sql.Rows) interface{} {
		for rs.Next() {
			return e.ResultSetColumnToElementType(labelColumn, valueColumn, rs, elementType)
		}
		return nil
	})
}

func (e *SQLSMExecutor) GetResultColumnAsObjectString(labelColumn, valueColumn string, elementType reflect.Type) interface{} {
	return e.ExtractorResultSet(func(rs *sql.Rows) interface{} {
		for rs.Next() {
			return e.ResultSetColumnToElementTypeString(labelColumn, valueColumn, rs, elementType)
		}
		return nil
	})
}

func (e *SQLSMExecutor) GetResultAsMap() map[string]interface{} {
	panic("not implemented")
}

func (e *SQLSMExecutor) GetResultAsList(elementType reflect.Type) []map[string]interface{} {
	panic("not implemented")
}

func (e *SQLSMExecutor) ExtractorResultSet(rse ResultSetExtractor) interface{} {
	rows, err := MYDB.Query("SELECT * FROM user u WHERE u.id = ?", 1)
	if err != nil {
		panic(err)
	}
	rse(rows)
	result, err := MYDB.Exec("INSERT INTO users(id,name) VALUES(?,?)", 1, "test")
	count, err := result.RowsAffected()
	fmt.Printf("rows affected: %d\n", count)
	// panic("not implemented")
}

func (e *SQLSMExecutor) Count() int64 {
	panic("not implemented")
}

func (e *SQLSMExecutor) GetResultAsListByPage(pageNumber, pageSize, totalNumber int, elementType reflect.Type) PageDataBean {
	return e.GetResultAsListByPageWithFilter(pageNumber, pageSize, totalNumber, elementType, nil)
}

func (e *SQLSMExecutor) GetResultAsListByPageWithFilter(pageNumber, pageSize, totalNumber int, elementType reflect.Type, byPageFilter ByPageFilter) PageDataBean {
	pageData := PageDataBean{}
	page := pageNumber - 1
	if page < 0 {
		page = 0
	}
	pageSize = defaultIfNull(pageSize, 16).(int)
	e.sm.LIMIT(page*pageSize, pageSize)

	if page == 0 || totalNumber > -1 {
		pageData.TotalNumber = e.Count()
	} else {
		pageData.TotalNumber = 0
	}

	if byPageFilter != nil {
		pageData.Data = e.ExtractorResultSet(func(rs *sql.Rows) interface{} {
			ls := []interface{}{}
			columns, _ := rs.Columns()
			for rs.Next() {
				byPageFilter.DoFilter(rs)
				ls = append(ls, e.ResultSetToElementType(columns, rs, elementType))
			}
			return ls
		})
	} else {
		pageData.Data = e.GetResultAsList(elementType)
	}

	return pageData
}

func (e *SQLSMExecutor) ResultSetToElementType(columns []string, rs *sql.Rows, elementType reflect.Type) interface{} {
	values := make([]interface{}, len(columns))
	for i := range values {
		var val interface{}
		values[i] = &val
	}
	rs.Scan(values...)
	ee := reflect.New(elementType).Interface()
	if m, ok := ee.(map[string]interface{}); ok {
		for i := range columns {
			m[columns[i]] = values[i]
		}
	} else {
		for i := range columns {
			if values[i] != nil {
				e.SetObjectFieldValue(ee, columns[i], values[i])
			}
		}
	}
	return ee
}

func (e *SQLSMExecutor) ResultSetColumnToElementType(columnLabel, columnValue int, rs *sql.Rows, elementType reflect.Type) interface{} {
	columns, _ := rs.Columns()
	if len(columns) == 0 {
		return nil
	}
	values := make([]interface{}, len(columns))
	for i := range values {
		var val interface{}
		values[i] = &val
	}
	ee := reflect.New(elementType).Interface()
	for rs.Next() {
		err := rs.Scan(values...)
		if err != nil {
			return err
		}
		if m, ok := ee.(map[string]interface{}); ok {
			m[columns[columnLabel]] = values[columnValue]
		} else {
			e.SetObjectFieldValue(ee, columns[columnLabel], values[columnValue])
		}
		//lint:ignore SA4004 just break
		break
	}
	return ee
}

func (e *SQLSMExecutor) ResultSetColumnToElementTypeString(columnLabel, columnValue string, rs *sql.Rows, elementType reflect.Type) interface{} {
	columns, _ := rs.Columns()
	if len(columns) == 0 {
		return nil
	}
	var labelIdx = -1
	var valueIdx = -1
	values := make([]interface{}, len(columns))
	for i := range values {
		var val interface{}
		values[i] = &val
		if columns[i] == columnLabel {
			labelIdx = i
		} else if columns[i] == columnValue {
			valueIdx = i
		}
	}
	ee := reflect.New(elementType).Interface()
	for rs.Next() {
		err := rs.Scan(values...)
		if err != nil {
			return err
		}
		columnLabel = fmt.Sprintf("%s", values[labelIdx])
		if m, ok := ee.(map[string]interface{}); ok {
			m[columnLabel] = values[valueIdx]
		} else {
			e.SetObjectFieldValue(ee, columnLabel, values[valueIdx])
		}
		break //lint:ignore SA4004 just break
	}
	return ee
}

func (e *SQLSMExecutor) SetObjectFieldValue(targetObject interface{}, fieldName string, value interface{}) {
	v := reflect.ValueOf(targetObject).Elem()
	f := e.getClassField(v.Type(), fieldName)
	if f.IsValid() {
		if value != nil {
			// is it need to convert value to f.Type() ?
			f.Set(reflect.ValueOf(value))
		}
	}
}

func (e *SQLSMExecutor) getClassField(c reflect.Type, fieldName string) reflect.Value {
	if c == nil || fieldName == "" {
		return reflect.Value{}
	}
	f, ok := c.FieldByName(fieldName)
	if !ok {
		return e.getClassField(c.Field(0).Type, fieldName)
	}
	return reflect.ValueOf(&f).Elem()
}

func defaultIfNull(value, defaultValue interface{}) interface{} {
	if value == nil {
		return defaultValue
	}
	return value
}

func RowsToMap(rs *sql.Rows, m map[string]interface{}) error {
	columns, _ := rs.Columns()
	if len(columns) == 0 {
		return nil
	}
	values := make([]interface{}, len(columns))
	for i := range values {
		var val interface{}
		values[i] = &val
	}
	for rs.Next() {
		err := rs.Scan(values...)
		if err != nil {
			return err
		}
		for i, col := range columns {
			m[col] = *values[i].(*interface{})
		}
		break //lint:ignore SA4004 just break
	}
	return nil
}

// RowsToMaps converts sql.Rows to a slice of maps
func RowsToMaps(rows *sql.Rows) ([]map[string]interface{}, error) {
	// Get column names
	columns, err := rows.Columns()
	if err != nil {
		return nil, err
	}

	// Create slice of interface{} to hold row values
	values := make([]interface{}, len(columns))
	for i := range values {
		var val interface{}
		values[i] = &val
	}

	// Create slice of maps to hold result
	maps := make([]map[string]interface{}, 0)

	// Iterate over rows
	for rows.Next() {
		// Scan row values into slice
		err := rows.Scan(values...)
		if err != nil {
			return nil, err
		}

		// Create map to hold row data
		rowData := make(map[string]interface{})

		// Iterate over columns and add values to map
		for i, col := range columns {
			rowData[col] = *values[i].(*interface{})
		}

		// Append map to result slice
		maps = append(maps, rowData)
	}

	// Return result
	return maps, nil
}

func RowsToStruct(rs *sql.Rows, v interface{}) error {
	// Get column names
	columns, err := rs.Columns()
	if err != nil {
		return err
	}

	// Create slice of interface{} to hold row values
	values := make([]interface{}, len(columns))
	for i := range values {
		var val interface{}
		values[i] = &val
	}
	vv := reflect.ValueOf(v)
	// Iterate over rows
	for rs.Next() {
		// Scan row values into slice
		err := rs.Scan(values...)
		if err != nil {
			return err
		}

		// Iterate over columns and set struct fields
		for i, col := range columns {
			// Get struct field by column name
			field := vv.FieldByName(col)

			// Set field value
			if field.IsValid() && field.CanSet() {
				val := reflect.ValueOf(values[i])
				field.Set(val)
			}
		}
		break //lint:ignore SA4004 just break
	}
	// Return nil
	return nil
}

// RowsToStructs converts sql.Rows to a slice of structs
func RowsToStructs(rows *sql.Rows, v interface{}) error {
	// Get value of slice
	sliceVal := reflect.ValueOf(v).Elem()

	// Get type of slice element
	elemType := sliceVal.Type().Elem()

	// Get column names
	columns, err := rows.Columns()
	if err != nil {
		return err
	}

	// Create slice of interface{} to hold row values
	values := make([]interface{}, len(columns))
	for i := range values {
		var val interface{}
		values[i] = &val
	}

	// Iterate over rows
	for rows.Next() {
		// Create struct instance
		elem := reflect.New(elemType).Elem()

		// Scan row values into slice
		err := rows.Scan(values...)
		if err != nil {
			return err
		}

		// Iterate over columns and set struct fields
		for i, col := range columns {
			// Get struct field by column name
			field := elem.FieldByName(col)

			// Set field value
			if field.IsValid() && field.CanSet() {
				val := reflect.ValueOf(*values[i].(*interface{}))
				field.Set(val)
			}
		}

		// Append struct to result slice

		sliceVal.Set(reflect.Append(sliceVal, elem))
	}

	// Return nil
	return nil
}
