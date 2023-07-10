package mak

import "database/sql"

type MySQLUMExecutor struct {
	SQLUMExecutor
	MDB *sql.DB
}

func NewMySQLUMExecutor(um *SQLUM, mdb *sql.DB) *MySQLUMExecutor {
	exe := &MySQLUMExecutor{MDB: mdb}
	exe.SQLUMExecutor = *NewSQLUMExecutor(um)
	return exe
}

func (exe *MySQLUMExecutor) Update() sql.Result {
	stmt, err := exe.MDB.Prepare(exe.Mak.Sql())
	if err != nil {
		panic(err)
	}
	defer stmt.Close()

	r, err := stmt.Exec(exe.Mak.Args()...)
	if err != nil {
		panic(err)
	}
	return r
}

func (exe *MySQLUMExecutor) BatchUpdate() []sql.Result {
	stmt, err := exe.MDB.Prepare(exe.Mak.Sql())
	if err != nil {
		panic(err)
	}
	defer stmt.Close()

	rs := []sql.Result{}
	for _, ags := range exe.Mak.BatchArgs() {
		r, err := stmt.Exec(ags...)
		if err != nil {
			panic(err)
		}
		rs = append(rs, r)
	}
	return rs
}
