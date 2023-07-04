package god

import (
	"database/sql"
)

type MySQLSMExecutor struct {
	SQLSMExecutor
	MDB *sql.DB
}

func NewMySQLSMExecutor(sm *SQLSM, mdb *sql.DB) *MySQLSMExecutor {
	exe := &MySQLSMExecutor{MDB: mdb}
	exe.SQLSMExecutor = *NewSQLSMExecutor(sm)
	return exe
}

func (exe *MySQLSMExecutor) INSERT_INTO_TABLExtractorResultSet(rse ResultSetExtractor) interface{} {
	stmp, err := exe.MDB.Prepare(exe.God.Sql())
	if err != nil {
		panic(err)
	}
	rows, err := stmp.Query(exe.God.Args()...)
	if err != nil {
		panic(err)
	}
	return rse(rows)
}
func (exe *MySQLSMExecutor) GetResultAsMap() map[string]interface{} {
	stmp, err := exe.MDB.Prepare(exe.God.Sql())
	if err != nil {
		panic(err)
	}
	rows, err := stmp.Query(exe.God.Args()...)
	if err != nil {
		panic(err)
	}
	m, err := exe.RowsToMap(rows)
	if err != nil {
		panic(err)
	}
	return m
}
func (exe *MySQLSMExecutor) GetResultAsMapList() []map[string]interface{} {
	stmp, err := exe.MDB.Prepare(exe.God.Sql())
	if err != nil {
		panic(err)
	}
	rows, err := stmp.Query(exe.God.Args()...)
	if err != nil {
		panic(err)
	}
	list, err := exe.RowsToMaps(rows)
	if err != nil {
		panic(err)
	}
	return list
}
func (exe *MySQLSMExecutor) ExtractorResultSet(rse ResultSetExtractor) interface{} {
	stmp, err := exe.MDB.Prepare(exe.God.Sql())
	if err != nil {
		panic(err)
	}
	rows, err := stmp.Query(exe.God.Args()...)
	if err != nil {
		panic(err)
	}
	return rse(rows)
}

func (exe *MySQLSMExecutor) ExtractorRows(rse RowsExtractor) error {
	stmp, err := exe.MDB.Prepare(exe.God.Sql())
	if err != nil {
		return err
	}
	rows, err := stmp.Query(exe.God.Args()...)
	if err != nil {
		return err
	}
	// Get column names
	columns, err := rows.Columns()
	if err != nil {
		return err
	}
	defer rows.Close() // finally close rows
	// Iterate over rows
	for rows.Next() {
		rse(&columns, rows)
	}
	return nil
}
func (exe *MySQLSMExecutor) Count() int64 {
	stmp, err := exe.MDB.Prepare(exe.God.Sql())
	if err != nil {
		panic(err)
	}
	rows, err := stmp.Query(exe.God.Args()...)
	if err != nil {
		panic(err)
	}
	defer rows.Close() // finally close rows
	var count int64
	for rows.Next() {
		err := rows.Scan(&count)
		if err != nil {
			panic(err)
		}
		return count //lint:ignore SA4004 just break
	}
	return 0
}
