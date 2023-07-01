package god

import (
	"database/sql"
	"reflect"
)

type MySQLSMExecutor struct {
	SQLSMExecutor
	MDB *sql.DB
	// SQLSMExecutorChild
}

func NewMySQLSMExecutor(sm *SQLSM, mdb *sql.DB) *MySQLSMExecutor {
	exe := &MySQLSMExecutor{MDB: mdb}
	exe.SQLSMExecutor = *NewSQLSMExecutor(sm, exe)
	return exe
}

func (exe *MySQLSMExecutor) INSERT_INTO_TABLExtractorResultSet(rse ResultSetExtractor) interface{} {
	sm := exe.SM
	rows, err := exe.MDB.Query(sm.Sql(), sm.Args()...)
	if err != nil {
		panic(err)
	}
	return rse(rows)
}
func (exe *MySQLSMExecutor) GetResultAsStruct(c reflect.Type) interface{} {
	sm := exe.SM
	rows, err := exe.MDB.Query(sm.Sql(), sm.Args()...)
	if err != nil {
		panic(err)
	}
	instance := reflect.New(c).Elem()
	exe.RowsToStruct(rows, instance)
	return instance
}
func (exe *MySQLSMExecutor) GetResultAsMap() map[string]interface{} {
	sm := exe.SM
	rows, err := exe.MDB.Query(sm.Sql(), sm.Args()...)
	if err != nil {
		panic(err)
	}
	m := map[string]interface{}{}
	exe.RowsToMap(rows, m)
	return m
}
func (exe *MySQLSMExecutor) GetResultAsMapList() []map[string]interface{} {
	sm := exe.SM
	rows, err := exe.MDB.Query(sm.Sql(), sm.Args()...)
	if err != nil {
		panic(err)
	}
	list, err := exe.RowsToMaps(rows)
	if err != nil {
		panic(err)
	}
	return list
}
func (exe *MySQLSMExecutor) GetResultAsStructList(c reflect.Type) []interface{} {
	sm := exe.SM
	rows, err := exe.MDB.Query(sm.Sql(), sm.Args()...)
	if err != nil {
		panic(err)
	}
	sliceType := reflect.SliceOf(c)
	list := reflect.MakeSlice(sliceType, 0, 0)

	if exe.RowsToStructs(rows, &list) != nil {
		return nil
	}
	return list.Interface().([]interface{})
}
func (exe *MySQLSMExecutor) ExtractorResultSet(rse ResultSetExtractor) interface{} {
	sm := exe.SM
	rows, err := exe.MDB.Query(sm.Sql(), sm.Args()...)
	if err != nil {
		panic(err)
	}
	return rse(rows)
}
func (exe *MySQLSMExecutor) Count() int64 {
	sm := exe.SM
	rows, err := exe.MDB.Query(sm.CountSql(), sm.Args()...)
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
		return count
	}
	return 0
}
