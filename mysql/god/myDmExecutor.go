package god

import "database/sql"

type MySQLDMExecutor struct {
	SQLDMExecutor
	MDB *sql.DB
}

func NewMySQLDMExecutor(um *SQLUM, mdb *sql.DB) *MySQLDMExecutor {
	exe := &MySQLDMExecutor{MDB: mdb}
	exe.SQLUMExecutor = *NewSQLUMExecutor(um)
	return exe
}

func (exe *MySQLDMExecutor) Update() sql.Result {
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

func (exe *MySQLDMExecutor) BatchUpdate() []sql.Result {
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
