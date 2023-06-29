package base

import (
	"reflect"
	"sql-mak/mysql"
	"sql-mak/mysql/god"
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

func (exe *MySQLSMExecutor) toSQLSMExecutorChild() god.SQLSMExecutorChild {
	return exe
}

func (exe *MySQLSMExecutor) INSERT_INTO_TABLExtractorResultSet(rse *god.ResultSetExtractor) interface{} {
	sm := exe.SM
	rows, err := mysql.MYDB.Query(sm.Sql(), sm.Args()...)
	if err != nil {
		panic(err)
	}
	return (*rse)(rows)
}
func (exe *MySQLSMExecutor) _getResultAsStruct(c reflect.Type) interface{} {
	sm := exe.SM
	rows, err := mysql.MYDB.Query(sm.Sql(), sm.Args()...)
	if err != nil {
		panic(err)
	}
	instance := reflect.New(c).Elem()
	exe.RowsToStruct(rows, instance)
	return instance
}
func (exe *MySQLSMExecutor) _getResultAsMap() map[string]interface{} {
	sm := exe.SM
	rows, err := mysql.MYDB.Query(sm.Sql(), sm.Args()...)
	if err != nil {
		panic(err)
	}
	m := map[string]interface{}{}
	exe.RowsToMap(rows, m)
	return m
}
func (exe *MySQLSMExecutor) _getResultAsMapList(c reflect.Type) []map[string]interface{} {
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
func (exe *MySQLSMExecutor) _getResultAsStructList(c reflect.Type) []interface{} {
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

func (exe *MySQLSMExecutor) _count() int64 {
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
