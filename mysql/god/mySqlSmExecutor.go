package god

import (
	"database/sql"
	"reflect"
	"sql-mak/utils"
)

type MySQLSMExecutor struct {
	SQLSMExecutor
	MDB *sql.DB
}

func NewMySQLSMExecutor(sm *SQLSM, mdb *sql.DB) *MySQLSMExecutor {
	exe := &MySQLSMExecutor{MDB: mdb}
	exe.SQLSMExecutor = *NewSQLSMExecutor(sm)
	return exe
}

func (exe *MySQLSMExecutor) INSERT_INTO_TABLExtractorResultSet(rse ResultSetExtractor) interface{} {
	stmp, err := exe.MDB.Prepare(exe.God.Sql())
	if err != nil {
		panic(err)
	}
	rows, err := stmp.Query(exe.God.Args()...)
	if err != nil {
		panic(err)
	}
	return rse(rows)
}
func (exe *MySQLSMExecutor) GetResultAsMap() map[string]interface{} {
	stmp, err := exe.MDB.Prepare(exe.God.Sql())
	if err != nil {
		panic(err)
	}
	rows, err := stmp.Query(exe.God.Args()...)
	if err != nil {
		panic(err)
	}
	m, err := exe.RowsToMap(rows)
	if err != nil {
		panic(err)
	}
	return m
}
func (exe *MySQLSMExecutor) GetResultAsMapList() []map[string]interface{} {
	stmp, err := exe.MDB.Prepare(exe.God.Sql())
	if err != nil {
		panic(err)
	}
	rows, err := stmp.Query(exe.God.Args()...)
	if err != nil {
		panic(err)
	}
	list, err := exe.RowsToMaps(rows)
	if err != nil {
		panic(err)
	}
	return list
}
func (exe *MySQLSMExecutor) ExtractorResultSet(rse ResultSetExtractor) interface{} {
	stmp, err := exe.MDB.Prepare(exe.God.Sql())
	if err != nil {
		panic(err)
	}
	rows, err := stmp.Query(exe.God.Args()...)
	if err != nil {
		panic(err)
	}
	defer rows.Close() // finally close rows
	return rse(rows)
}

func (e *MySQLSMExecutor) ExtractorResultTo(out interface{}) *MySQLSMExecutor {
	e.ExtractorResultSet(func(rs *sql.Rows) interface{} {
		RowsToOut(rs, reflect.ValueOf(out))
		return nil
	})
	return e
}
func RowsToOut(rs *sql.Rows, out reflect.Value) {
	columns, _ := rs.Columns()
	cL := len(columns)

	// values := make([]interface{}, cL)
	// for i := range values {
	// 	var val interface{}
	// 	values[i] = &val
	// }

	outType := utils.GetType(out.Type())
	out = utils.GetValue(out)
	if outType.Kind() == reflect.Slice {
		eleType := utils.GetType(outType.Elem())
		var ei utils.EntityInfo
		var fnames []string
		for rs.Next() {
			ele := reflect.New(eleType)
			if reflect.ValueOf(ei).IsZero() {
				ei = utils.GetEntityInfo(ele.Elem().Interface())
				fnames = make([]string, cL)
				for i, col := range columns {
					fnames[i] = ei.CFMap[col]
				}
			}
			// rs.Scan(values...)
			// utils.SetFValues(ele, &fnames, &values)
			utils.SetValuesWithRows(ele, &fnames, rs)
			if outType.Elem().Kind() == reflect.Ptr {
				out.Set(reflect.Append(out, ele))
			} else {
				out.Set(reflect.Append(out, ele.Elem()))
			}
		}
		return
	}
	ei := utils.GetEntityInfo(out.Interface())
	fnames := make([]string, cL)
	for i, col := range columns {
		fnames[i] = ei.CFMap[col]
	}
	for rs.Next() {
		// rs.Scan(values...)
		// utils.SetFValues(out, &fnames, &values)
		utils.SetValuesWithRows(out, &fnames, rs)
	}
}

func (exe *MySQLSMExecutor) Count() int64 {
	stmp, err := exe.MDB.Prepare(exe.God.Sql())
	if err != nil {
		panic(err)
	}
	rows, err := stmp.Query(exe.God.Args()...)
	if err != nil {
		panic(err)
	}
	defer rows.Close() // finally close rows
	var count int64
	for rows.Next() {
		err := rows.Scan(&count)
		if err != nil {
			panic(err)
		}
		return count //lint:ignore SA4004 just break
	}
	return 0
}
