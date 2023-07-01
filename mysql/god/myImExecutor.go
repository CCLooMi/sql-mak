package god

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

func (exe *MySQLIMExecutor) BatchUpdate() []sql.Result {
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
