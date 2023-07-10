package mak

import "database/sql"

type MySQLDMExecutor struct {
	SQLDMExecutor
	MDB *sql.DB
}

func NewMySQLDMExecutor(dm *SQLDM, mdb *sql.DB) *MySQLDMExecutor {
	exe := &MySQLDMExecutor{MDB: mdb}
	exe.SQLDMExecutor = *NewSQLDMExecutor(dm)
	return exe
}

func (exe *MySQLDMExecutor) Update() sql.Result {
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

func (exe *MySQLDMExecutor) BatchUpdate() []sql.Result {
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
