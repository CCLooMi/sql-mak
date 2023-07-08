package mak

type SQLIMExecutor struct {
	SQLExecutor
	SQLUMExecutorChild
}

func NewSQLIMExecutor(im *SQLIM) *SQLIMExecutor {
	exe := &SQLIMExecutor{}
	mak := im.toSQLMak()
	exe.SQLExecutor = *NewSQLExecutor(mak)
	return exe
}
