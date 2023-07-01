package god

type SQLIMExecutor struct {
	SQLExecutor
	SQLUMExecutorChild
}

func NewSQLIMExecutor(im *SQLIM) *SQLIMExecutor {
	exe := &SQLIMExecutor{}
	god := im.toSQLGod()
	exe.SQLExecutor = *NewSQLExecutor(god)
	return exe
}
