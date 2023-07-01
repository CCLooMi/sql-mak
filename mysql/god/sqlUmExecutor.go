package god

import "database/sql"

type SQLUMExecutor struct {
	SQLExecutor
	SQLUMExecutorChild
}
type SQLUMExecutorChild interface {
	Update() sql.Result
	BatchUpdate() []sql.Result
}

func NewSQLUMExecutor(um *SQLUM) *SQLUMExecutor {
	exe := &SQLUMExecutor{}
	god := um.toSQLGod()
	exe.SQLExecutor = *NewSQLExecutor(god)
	return exe
}
