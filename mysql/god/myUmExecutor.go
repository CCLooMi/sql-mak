package god

import "database/sql"

type MyUMExecutor struct {
	SQLUMExecutor
	MDB *sql.DB
}

func NewMyUMExecutor(um *SQLUM, mdb *sql.DB) *MyUMExecutor {
	exe := &MyUMExecutor{MDB: mdb}
	exe.SQLUMExecutor = *NewSQLUMExecutor(um)
	return exe
}

func (exe *MyUMExecutor) Update() sql.Result {
	stmt, err := exe.MDB.Prepare(exe.God.Sql())
	if err != nil {
		panic(err)
	}
	r, err := stmt.Exec(exe.God.Args()...)
	if err != nil {
		panic(err)
	}
	return r
}
func (exe *MyUMExecutor) BatchUpdate() []sql.Result {
	stmt, err := exe.MDB.Prepare(exe.God.Sql())
	if err != nil {
		panic(err)
	}
	rs := []sql.Result{}
	for _, ags := range exe.God.BatchArgs() {
		r, _ := stmt.Exec(ags...)
		rs = append(rs, r)
	}
	return rs
}
