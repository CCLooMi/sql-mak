package main

import (
	"fmt"
	"reflect"
	"sql-mak/mysql"
	"sql-mak/mysql/base"
	"sql-mak/mysql/god"
)

func main() {
	god.RegisterExecutorProvider("sm", reflect.ValueOf(base.NewMySQLSMExecutor))

	r := mysql.
		SELECT("*").
		FROM("users", "u").
		Execute().
		GetResultAsMapList()
	fmt.Printf("%s\n", r)
}
func sqlTest() {
	sm := god.NewSQLSM()
	sm.SELECT("*").
		FROM("users", "u").
		WHERE("u.id=?", 110).
		LOGSQL(true).
		Sql()

	sm = god.NewSQLSM()
	sm.SELECT("*").
		FROM("users", "u").
		LEFT_JOIN("info", "i", "u.id=i.user_id").
		WHERE("u.id=?", 110).
		AND("i.email=?", "cx@wios.com").
		GROUP_BY("u.id", "i.id").
		LOGSQL(true).
		Sql()

	um := god.NewSQLUM()
	um.UPDATE("users", "u").
		SET("u.name=?").
		WHERE("u.id=?", 110).
		AND("u.email=?", "cx@wios.com").
		LOGSQL(true).
		Sql()

	mysql.SELECT("u.id", "u.name").FROM("users", "u").LOGSQL(true).Sql()
	mysql.
		INSERT_INTO_TABLE("users").
		INTO_COLUMNS("u.id", "u.name").
		VALUES_SM(sm).
		LOGSQL(true).
		Sql()
}
