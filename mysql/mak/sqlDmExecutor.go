package mak

type SQLDMExecutor struct {
	SQLExecutor
	SQLUMExecutor
}

func NewSQLDMExecutor(dm *SQLDM) *SQLDMExecutor {
	exe := &SQLDMExecutor{}
	god := dm.toSQLMak()
	exe.SQLExecutor = *NewSQLExecutor(god)
	// exe.SQLUMExecutor = *NewSQLUMExecutor(dm)
	return exe
}
