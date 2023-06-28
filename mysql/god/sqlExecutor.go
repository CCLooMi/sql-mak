package god

import (
	"database/sql"

	"github.com/sirupsen/logrus"
)

type SQLExecutor struct {
	log *logrus.Entry
	god *SQLGod
}

type ResultSetExtractor func(rs *sql.Rows) interface{}

func NewSQLExecutor(god *SQLGod) *SQLExecutor {
	return &SQLExecutor{
		log: logrus.WithField("component", "SQLExecutor"),
		god: god,
	}
}
