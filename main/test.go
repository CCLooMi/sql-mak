package main

import (
	"sql-mak/sql/god"
)

func main() {
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

}
