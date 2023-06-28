package god

type SQLUMExecutor struct {
	SQLExecutor
	child SQLUMExecutorChild
}
type SQLUMExecutorChild interface {
	Update() int
	BatchUpdate() []int
}

func NewSQLUMExecutor(god *SQLGod, child *SQLUMExecutorChild) *SQLUMExecutor {
	exe := &SQLUMExecutor{}
	exe.SQLExecutor = *NewSQLExecutor(god)
	exe.child = *child
	return exe
}

func (exe *SQLUMExecutor) Update() int {
	return exe.child.Update()
}

func (exe *SQLUMExecutor) BatchUpdate() []int {
	return exe.child.BatchUpdate()
}
