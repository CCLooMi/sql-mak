package mak

import (
	"database/sql"
	"go.uber.org/zap"
)

type SQLExecutor struct {
	Log *zap.Logger
	Mak SQLMak
}

type ResultSetExtractor func(rs *sql.Rows) interface{}

func NewSQLExecutor(mak SQLMak) *SQLExecutor {
	log, _ := zap.NewProduction()
	return &SQLExecutor{
		Log: log,
		Mak: mak,
	}
}
