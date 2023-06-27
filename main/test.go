package main

import (
	"sql-mak/sql/god"
)

func main() {
	sm := god.NewSQLSM()
	sm.SELECT("*").FROM("users", "u").WHERE("u.id=?", 110).LOGSQL(true).Sql()
}
