package god

import (
	"database/sql"

	"github.com/sirupsen/logrus"
)

type SQLExecutor struct {
	Log *logrus.Entry
	God SQLGod
}

type ResultSetExtractor func(rs *sql.Rows) interface{}

func NewSQLExecutor(god SQLGod) *SQLExecutor {
	return &SQLExecutor{
		Log: logrus.WithField("component", "SQLExecutor"),
		God: god,
	}
}
