package god

import (
	"database/sql"
	"reflect"
)

// import "reflect"

type SQL struct {
}

var MYDB *sql.DB

func init() {
	// _db, err := sql.Open("mysql", "root:root@tcp(127.0.0.1:3306)/test")
	// if err != nil {
	// 	panic(err)
	// }
	// MYDB = _db
}

func SELECT(columns ...interface{}) *SQLSM {
	return NewSQLSM().SELECT(columns...)
}

func SELECT_EXP(exp *EXP) *SQLSM {
	return NewSQLSM().SELECT(exp)
}

func SELECT_AS(column, alias string) *SQLSM {
	return NewSQLSM().SELECT_AS(column, alias)
}

func SELECT_SM_AS(column *SQLSM, alias string) *SQLSM {
	return NewSQLSM().SELECT_SM(column, alias)
}

func SELECT_EXP_AS(exp *EXP, alias string) *SQLSM {
	return NewSQLSM().SELECT_EXP(exp, alias)
}

func INSERT_INTO(c reflect.Type, columns ...string) *SQLIM {
	return NewSQLIM().INSERT_INTO(c, columns...)
}

func INSERT_INTO_TABLE(table string, columns ...string) *SQLIM {
	return NewSQLIM().INSERT_INTO(table, columns...)
}

func UPDATE(c reflect.Type, alias string) *SQLUM {
	return NewSQLUM().UPDATE(c, alias)
}

func UPDATE_TABLE(table, alias string) *SQLUM {
	return NewSQLUM().UPDATE(table, alias)
}

func DELETE() *SQLDM {
	return NewSQLDM()
}
