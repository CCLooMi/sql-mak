package mak

import (
	"database/sql"
)

type SQLExecutor struct {
	Mak SQLMak
}

type ResultSetExtractor func(rs *sql.Rows) interface{}

func NewSQLExecutor(mak SQLMak) *SQLExecutor {
	return &SQLExecutor{
		Mak: mak,
	}
}
