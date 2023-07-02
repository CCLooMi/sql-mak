package main

import (
	"fmt"
	"sql-mak/mysql"
)

func main() {
	r := mysql.
		SELECT("*").
		FROM("users", "u").
		Execute(mysql.MYDB).
		GetResultAsMapList()
	fmt.Printf("%s\n", r)
}
