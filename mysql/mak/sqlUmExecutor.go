package mak

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
	mak := um.toSQLMak()
	exe.SQLExecutor = *NewSQLExecutor(mak)
	return exe
}
