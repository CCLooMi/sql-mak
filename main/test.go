package main

import (
	"fmt"
	"sql-mak/sql"
)

func main() {
	sm := sql.NewSQLSM()

	fmt.Println(sm.SELECT("*").FROM("users", "u").WHERE("u.id=?", 110).Sql())

}
