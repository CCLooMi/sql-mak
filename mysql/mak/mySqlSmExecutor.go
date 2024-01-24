package mak

import (
	"database/sql"
	"reflect"

	"github.com/CCLooMi/sql-mak/utils"
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

func (exe *MySQLSMExecutor) GetResultAsMap() map[string]interface{} {
	stmp, err := exe.MDB.Prepare(exe.Mak.Sql())
	if err != nil {
		panic(err)
	}
	rows, err := stmp.Query(exe.Mak.Args()...)
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
	stmp, err := exe.MDB.Prepare(exe.Mak.Sql())
	if err != nil {
		panic(err)
	}
	rows, err := stmp.Query(exe.Mak.Args()...)
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
	stmp, err := exe.MDB.Prepare(exe.Mak.Sql())
	if err != nil {
		panic(err)
	}
	rows, err := stmp.Query(exe.Mak.Args()...)
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

	outType := utils.GetType(out.Type())
	out = utils.GetValue(out)
	if outType.Kind() == reflect.Slice {
		eleType := utils.GetType(outType.Elem())
		var ei *utils.EntityInfo
		var fnames []string
		for rs.Next() {
			ele := reflect.New(eleType)
			if reflect.ValueOf(ei).IsZero() {
				ei = utils.GetEntityInfo(ele.Elem().Interface())
				if ei == nil {
					//rs.scan result to ele
					rs.Scan(ele.Interface())
					if outType.Elem().Kind() == reflect.Ptr {
						out.Set(reflect.Append(out, ele))
					} else {
						out.Set(reflect.Append(out, ele.Elem()))
					}
					continue
				}
				fnames = make([]string, cL)
				for i, col := range columns {
					fnames[i] = ei.CFMap[col]
				}
			}
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
	if ei == nil {
		//rs.scan result to out
		rs.Scan(out.Interface())
		return
	}
	fnames := make([]string, cL)
	for i, col := range columns {
		fnames[i] = ei.CFMap[col]
	}
	for rs.Next() {
		utils.SetValuesWithRows(out, &fnames, rs)
	}
}

func (exe *MySQLSMExecutor) Count() int64 {
	stmp, err := exe.MDB.Prepare(exe.Mak.CountSql())
	if err != nil {
		panic(err)
	}
	rows, err := stmp.Query(exe.Mak.Args()...)
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
