package sql

type SQLGod interface {
	Sql() string
	CountSql() string
	Args() []interface{}
	BatchArgs() [][]interface{}
	Execute() SQLExecutor
}

type IN string

const (
	INValue    IN = "IN"
	NOTINValue IN = "NOT IN"
)

func (i IN) value() string {
	return string(i)
}
