package base

import (
	"reflect"
	"sql-mak/mysql"
	"sql-mak/mysql/god"
	"unsafe"
)

type MySQLSMExecutor struct {
	god.SQLSMExecutor
	god.SQLSMExecutorChild
}

func NewMySQLSMExecutor(sm *god.SQLSM) *MySQLSMExecutor {
	exe := &MySQLSMExecutor{}
	exe.SQLSMExecutor = *god.NewSQLSMExecutor(sm, exe.toSQLSMExecutorChild())
	return exe
}

func (exe *MySQLSMExecutor) toSQLSMExecutorChild() *god.SQLSMExecutorChild {
	return (*god.SQLSMExecutorChild)(unsafe.Pointer(exe))
}
func (exe *MySQLSMExecutor) GetResultAsStruct(c reflect.Type) interface{} {
	sm := exe.SM
	rows, err := mysql.MYDB.Query(sm.Sql(), sm.Args()...)
	if err != nil {
		panic(err)
	}
	instance := reflect.New(c).Elem()
	exe.RowsToStruct(rows, instance)
	return instance
}
func (exe *MySQLSMExecutor) GetResultAsMap() map[string]interface{} {
	sm := exe.SM
	rows, err := mysql.MYDB.Query(sm.Sql(), sm.Args()...)
	if err != nil {
		panic(err)
	}
	m := map[string]interface{}{}
	exe.RowsToMap(rows, m)
	return m
}
func (exe *MySQLSMExecutor) GetResultAsMapList(c reflect.Type) []map[string]interface{} {
	sm := exe.SM
	rows, err := mysql.MYDB.Query(sm.Sql(), sm.Args()...)
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
	rows, err := mysql.MYDB.Query(sm.Sql(), sm.Args()...)
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

func (exe *MySQLSMExecutor) Count() int64 {
	sm := exe.SM
	rows, err := mysql.MYDB.Query(sm.CountSql(), sm.Args()...)
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
