package god

type SQLDMExecutor struct {
	SQLExecutor
	SQLUMExecutor
}

func NewSQLDMExecutor(god SQLGod, child *SQLUMExecutorChild) *SQLDMExecutor {
	exe := &SQLDMExecutor{}
	exe.SQLExecutor = *NewSQLExecutor(god)
	exe.SQLUMExecutor = *NewSQLUMExecutor(god, child)
	return exe
}
