package mysql

import (
	"database/sql"
	"reflect"
	"sql-mak/mysql/god"

	_ "github.com/go-sql-driver/mysql"
)

var MYDB *sql.DB

func init() {
	_db, err := sql.Open("mysql", "root:apple@tcp(127.0.0.1:3308)/wios")
	if err != nil {
		panic(err)
	}
	MYDB = _db

	MYDB.SetMaxOpenConns(10)
	MYDB.SetMaxIdleConns(5)
}

func SELECT(columns ...interface{}) *god.SQLSM {
	return god.NewSQLSM().SELECT(columns...)
}

func SELECT_EXP(exp *god.EXP) *god.SQLSM {
	return god.NewSQLSM().SELECT(exp)
}

func SELECT_AS(column, alias string) *god.SQLSM {
	return god.NewSQLSM().SELECT_AS(column, alias)
}

func SELECT_SM_AS(column *god.SQLSM, alias string) *god.SQLSM {
	return god.NewSQLSM().SELECT_SM(column, alias)
}

func SELECT_EXP_AS(exp *god.EXP, alias string) *god.SQLSM {
	return god.NewSQLSM().SELECT_EXP(exp, alias)
}

func INSERT_INTO(c reflect.Type, columns ...string) *god.SQLIM {
	return god.NewSQLIM().INSERT_INTO(c, columns...)
}

func INSERT_INTO_TABLE(table string, columns ...string) *god.SQLIM {
	return god.NewSQLIM().INSERT_INTO(table, columns...)
}

func UPDATE(c reflect.Type, alias string) *god.SQLUM {
	return god.NewSQLUM().UPDATE(c, alias)
}

func UPDATE_TABLE(table, alias string) *god.SQLUM {
	return god.NewSQLUM().UPDATE(table, alias)
}

func DELETE() *god.SQLDM {
	return god.NewSQLDM()
}
