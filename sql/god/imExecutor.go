package god

import (
	"database/sql"
	"unsafe"
)

type SQLIMExecutor struct {
	SQLUMExecutor
	SQLUMExecutorChild
	child SQLIMExecutorChild
}
type SQLIMExecutorChild interface {
	SQLUMExecutorChild
	UpdateAndGetGeneratedKey() *sql.Result
}

func NewSQLIMExecutor(god *SQLGod, child *SQLIMExecutorChild) *SQLIMExecutor {
	exe := &SQLIMExecutor{}
	exe.SQLUMExecutor = *NewSQLUMExecutor(god, exe.toSQLUMExecutorChild())
	exe.child = *child
	return exe
}
func (im *SQLIMExecutor) toSQLUMExecutorChild() *SQLUMExecutorChild {
	return (*SQLUMExecutorChild)(unsafe.Pointer(im))
}
func (im *SQLIMExecutor) UpdateAndGetGeneratedKey() *sql.Result {
	return im.child.UpdateAndGetGeneratedKey()
}
