package mak

type SQLMak interface {
	Sql() string
	CountSql() string
	Args() []interface{}
	BatchArgs() [][]interface{}
}
