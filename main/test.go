package main

import (
	"fmt"
	"sql-mak/mysql"
	"sql-mak/mysql/god"
)

func main() {
	r := mysql.
		SELECT("*").
		FROM("users", "u").
		Execute(mysql.MYDB).
		GetResultAsMapList()
	fmt.Printf("%s\n", r)

	mysql.INSERT_INTO("users").
		INTO_COLUMNS("u.id", "u.name").
		BatchArgs([]interface{}{1, 2}).
		Execute(mysql.MYDB)
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
		INSERT_INTO("users").
		INTO_COLUMNS("u.id", "u.name").
		VALUES_SM(sm).
		LOGSQL(true).
		Sql()
}
