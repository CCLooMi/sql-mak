package mak

import (
	"database/sql"

	"github.com/sirupsen/logrus"
)

type SQLExecutor struct {
	Log *logrus.Entry
	Mak SQLMak
}

type ResultSetExtractor func(rs *sql.Rows) interface{}

func NewSQLExecutor(mak SQLMak) *SQLExecutor {
	return &SQLExecutor{
		Log: logrus.WithField("component", "SQLExecutor"),
		Mak: mak,
	}
}
