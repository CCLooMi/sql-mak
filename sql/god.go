package sql

type SQLGod interface {
	sql() string
	countSql() string
	Args() []interface{}
	BatchArgs() [][]interface{}
	execute() SQLExecutor
}

type IN string

const (
	INValue    IN = "IN"
	NOTINValue IN = "NOT IN"
)

func (i IN) value() string {
	return string(i)
}
