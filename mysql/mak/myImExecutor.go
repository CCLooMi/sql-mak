package mak

import "database/sql"

type MySQLIMExecutor struct {
	SQLIMExecutor
	MDB *sql.DB
}

func NewMySQLIMExecutor(im *SQLIM, mdb *sql.DB) *MySQLIMExecutor {
	exe := &MySQLIMExecutor{MDB: mdb}
	exe.SQLIMExecutor = *NewSQLIMExecutor(im)
	return exe
}

func (exe *MySQLIMExecutor) Update() sql.Result {
	stmt, err := exe.MDB.Prepare(exe.Mak.Sql())
	if err != nil {
		panic(err)
	}
	r, err := stmt.Exec(exe.Mak.Args()...)
	if err != nil {
		panic(err)
	}
	return r
}

func (exe *MySQLIMExecutor) BatchUpdate() []sql.Result {
	stmt, err := exe.MDB.Prepare(exe.Mak.Sql())
	if err != nil {
		panic(err)
	}
	rs := []sql.Result{}
	for _, ags := range exe.Mak.BatchArgs() {
		r, _ := stmt.Exec(ags...)
		rs = append(rs, r)
	}
	return rs
}
