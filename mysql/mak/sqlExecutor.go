package mak

import (
	"database/sql"

	"github.com/sirupsen/logrus"
)

type SQLExecutor struct {
	Log *logrus.Logger
	Mak SQLMak
}

type ResultSetExtractor func(rs *sql.Rows) interface{}

func NewSQLExecutor(mak SQLMak) *SQLExecutor {
	return &SQLExecutor{
		Log: logrus.New(),
		Mak: mak,
	}
}
