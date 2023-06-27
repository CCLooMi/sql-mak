package god

type SQLUMExecutor struct {
	*SQLExecutor
}

func NewSQLUMExecutor(god *SQLGod) *SQLUMExecutor {
	return &SQLUMExecutor{SQLExecutor: NewSQLExecutor(god)}
}

func (exe *SQLUMExecutor) Update() int {
	return 0 // TODO: implement
}

func (exe *SQLUMExecutor) BatchUpdate() []int {
	return nil // TODO: implement
}
