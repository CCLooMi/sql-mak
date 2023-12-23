package mysql

import (
	"database/sql"
	"github.com/CCLooMi/sql-mak/mysql/mak"
)

func SELECT(columns ...interface{}) *mak.SQLSM {
	return mak.NewSQLSM().SELECT(columns...)
}

func SELECT_EXP(exp *mak.EXP) *mak.SQLSM {
	return mak.NewSQLSM().SELECT(exp)
}

func SELECT_AS(column, alias string) *mak.SQLSM {
	return mak.NewSQLSM().SELECT_AS(column, alias)
}

func SELECT_SM_AS(column *mak.SQLSM, alias string) *mak.SQLSM {
	return mak.NewSQLSM().SELECT_SM(column, alias)
}

func SELECT_EXP_AS(exp *mak.EXP, alias string) *mak.SQLSM {
	return mak.NewSQLSM().SELECT_EXP(exp, alias)
}

func INSERT_INTO(table interface{}, columns ...string) *mak.SQLIM {
	return mak.NewSQLIM().INSERT_INTO(table, columns...)
}

func UPDATE(table interface{}, alias string) *mak.SQLUM {
	return mak.NewSQLUM().UPDATE(table, alias)
}
func DELETE() *mak.SQLDM {
	return mak.NewSQLDM()
}
func TxExecute(tx *sql.Tx, maks ...mak.SQLMak) []sql.Result {
	var rsList []sql.Result
	for _, v := range maks {
		stmt, err := tx.Prepare(v.Sql())
		if err != nil {
			tx.Rollback()
			panic(err.Error())
		}
		defer stmt.Close()
		rs, err := stmt.Exec(v.Args()...)
		if err != nil {
			tx.Rollback()
			panic(err.Error())
		}
		rsList = append(rsList, rs)
	}
	err := tx.Commit()
	if err != nil {
		panic(err.Error())
	}
	return rsList
}
