package mysql

import (
	"database/sql"
	"sql-mak/mysql/god"

	_ "github.com/go-sql-driver/mysql"
)

var MYDB *sql.DB

func init() {
	_db, err := sql.Open("mysql", "root:apple@tcp(127.0.0.1:3308)/wios?charset=utf8mb4&parseTime=True&loc=Local")
	if err != nil {
		panic(err)
	}
	if err = _db.Ping(); err != nil {
		_db, _ := sql.Open("mysql", "root:123456@tcp(127.0.0.1:3306)/test?charset=utf8mb4&parseTime=True&loc=Local")
		if err = _db.Ping(); err != nil {
			panic(err)
		}
		MYDB = _db
	} else {
		MYDB = _db
	}

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

func INSERT_INTO(table interface{}, columns ...string) *god.SQLIM {
	return god.NewSQLIM().INSERT_INTO(table, columns...)
}

func UPDATE(table interface{}, alias string) *god.SQLUM {
	return god.NewSQLUM().UPDATE(table, alias)
}

func DELETE() *god.SQLDM {
	return god.NewSQLDM()
}
