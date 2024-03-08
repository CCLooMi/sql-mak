package mak

import (
	"database/sql"
)

type SQLSMExecutor struct {
	SQLExecutor
	SQLSMExecutorChild
	SM *SQLSM
}

type SQLSMExecutorChild interface {
	GetResultAsMap() map[string]interface{}
	GetResultAsMapList() []map[string]interface{}
	GetResultAsList() []interface{}
	ExtractorResultSet(rse ResultSetExtractor) interface{}
	ExtractorResultTo(out interface{})
	Count() int64
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
	sme := &SQLSMExecutor{SM: sm}
	mak := sm.toSQLMak()
	sme.SQLExecutor = *NewSQLExecutor(mak)
	return sme
}

func (e *SQLSMExecutor) GetResultAsListByPage(pageNumber, pageSize, totalNumber int) PageDataBean {
	return e.GetResultAsListByPageWithFilter(pageNumber, pageSize, totalNumber, nil)
}

func (e *SQLSMExecutor) GetResultAsListByPageWithFilter(pageNumber, pageSize, totalNumber int, byPageFilter ByPageFilter) PageDataBean {
	pageData := PageDataBean{}
	page := pageNumber - 1
	if page < 0 {
		page = 0
	}
	pageSize = e.DefaultIfNull(pageSize, 16).(int)
	e.SM.LIMIT(page*pageSize, pageSize)

	if page == 0 || totalNumber > -1 {
		pageData.TotalNumber = e.Count()
	} else {
		pageData.TotalNumber = 0
	}

	if byPageFilter != nil {
		pageData.Data = e.ExtractorResultSet(func(rs *sql.Rows) interface{} {
			columns, _ := rs.Columns()
			values := make([]interface{}, len(columns))
			for i := range values {
				var val interface{}
				values[i] = &val
			}
			// Create slice of maps to hold result
			maps := make([]map[string]interface{}, 0)

			defer rs.Close() // finally close rows
			for rs.Next() {
				byPageFilter.DoFilter(rs)
				// Create map to hold row data
				rowData := make(map[string]interface{})

				// Iterate over columns and add values to map
				for i, col := range columns {
					rowData[col] = *values[i].(*interface{})
				}
				// Append map to result slice
				maps = append(maps, rowData)
			}
			return maps
		})
	} else {
		pageData.Data = e.ExtractorResultSet(func(rs *sql.Rows) interface{} {
			data, _ := e.RowsToMaps(rs)
			return data
		})
	}
	return pageData
}

func (e *SQLSMExecutor) DefaultIfNull(value, defaultValue interface{}) interface{} {
	if value == nil {
		return defaultValue
	}
	return value
}

func (e *SQLSMExecutor) RowsToMap(rows *sql.Rows) (map[string]interface{}, error) {
	columns, err := rows.Columns()
	if err != nil {
		return nil, err
	}
	values := make([]interface{}, len(columns))
	for i := range values {
		var val interface{}
		values[i] = &val
	}
	defer rows.Close() // finally close rows
	for rows.Next() {
		err := rows.Scan(values...)
		if err != nil {
			return nil, err
		}
		// Create map to hold row data
		rowData := make(map[string]interface{})
		for i, col := range columns {
			rowData[col] = *values[i].(*interface{})
		}
		return rowData, nil //lint:ignore SA4004 just break
	}
	return make(map[string]interface{}), nil
}
func getGoType(dbType string) interface{} {
	switch dbType {
	case "VARCHAR", "TEXT", "NVARCHAR", "VARCHAR2", "CHAR", "XML", "JSON":
		return new(*string)
	case "INT":
		return new(*int)
	case "BIGINT":
		return new(*int64)
	case "DECIMAL":
		return new(*float64)
	case "BOOL":
		return new(*bool)
	case "JSONB":
		return map[string]interface{}{}
	default:
		return new(*interface{})
	}
}

// RowsToMaps converts sql.Rows to a slice of maps
func (e *SQLSMExecutor) RowsToMaps(rows *sql.Rows) ([]map[string]interface{}, error) {
	// Get column names
	columns, err := rows.Columns()
	if err != nil {
		return nil, err
	}
	colTypes, err := rows.ColumnTypes()
	if err != nil {
		return nil, err
	}

	// Create slice of interface{} to hold row values
	values := make([]interface{}, len(columns))
	for i := range values {
		values[i] = getGoType(colTypes[i].DatabaseTypeName())
	}

	// Create slice of maps to hold result
	maps := make([]map[string]interface{}, 0)

	defer rows.Close() // finally close rows
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
			rowData[col] = values[i]
			values[i] = getGoType(colTypes[i].DatabaseTypeName())
		}

		// Append map to result slice
		maps = append(maps, rowData)
	}

	// Return result
	return maps, nil
}

func (e *SQLSMExecutor) ColumnToList(rows *sql.Rows, col int) ([]interface{}, error) {
	colTypes, err := rows.ColumnTypes()
	if err != nil {
		return nil, err
	}
	result := make([]interface{}, 0)
	i := 0
	for rows.Next() {
		result = append(result, getGoType(colTypes[col].DatabaseTypeName()))
		// Scan row values into slice
		err := rows.Scan(result[i])
		i++
		if err != nil {
			return nil, err
		}
	}
	return result, nil
}
