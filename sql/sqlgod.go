package sql

import (
	"log"
	"reflect"
	"strings"
)

type AbstractSQLGod struct {
	SQLGod
	log          *log.Logger
	logSql       bool
	args         []interface{}
	batchArgs    [][]interface{}
	hasSubSelect bool
}

func NewAbstractSQLGod() *AbstractSQLGod {
	return &AbstractSQLGod{
		log:          log.Default(),
		logSql:       true,
		args:         make([]interface{}, 0),
		hasSubSelect: false,
	}
}

func (g *AbstractSQLGod) tableName(c reflect.Type) string {
	tableName := ""
	field, ok := c.FieldByName("Table")
	if ok {
		tableName = field.Tag.Get("name")
	}

	if len(tableName) == 0 {
		// don't rely on hibernate-jpa to get the Table annotation name value
		jpaTable, ok := reflect.TypeOf((*interface{})(nil)).Elem().FieldByName("Table")
		if ok {
			a := field.Tag.Get(jpaTable.Tag.Get("name"))
			if len(a) > 0 {
				tableName = a
			}
		}
	}
	if len(tableName) == 0 {
		tableName = c.Name()
	}
	return tableName
}

func (g *AbstractSQLGod) IF(condition bool, fs ...func(g SQLGod)) SQLGod {
	if len(fs) == 0 {
		return g
	}
	if condition {
		fs[0](g)
		return g
	}
	if len(fs) > 1 {
		fs[0](g)
	}
	return g
}

func (g *AbstractSQLGod) SWITCH(i int, fs ...func(g SQLGod)) SQLGod {
	if i < 0 || i > len(fs)-1 {
		return g
	}
	fs[i](g)
	return g
}

func (g *AbstractSQLGod) LOGSQL(b bool) SQLGod {
	g.logSql = b
	return g
}

func (g *AbstractSQLGod) sql() string {
	var sb strings.Builder
	g._sql(&sb)
	if g.logSql {
		g.log.Println("sql:", sb.String(), g.args)
	}
	return sb.String()
}

func (g *AbstractSQLGod) countSql() string {
	var sb strings.Builder
	g._countSQL(&sb)
	if g.logSql {
		g.log.Println("sql:", sb.String(), g.args)
	}
	return sb.String()
}

func (g *AbstractSQLGod) Args() []interface{} {
	if !g.hasSubSelect {
		return g.args
	} else {
		ags := make([]interface{}, 0)
		g.flat(g.args, &ags)
		return ags
	}
}

func (g *AbstractSQLGod) BatchArgs() [][]interface{} {
	return g.batchArgs
}

func (g *AbstractSQLGod) flat(args []interface{}, result *[]interface{}) {
	for _, o := range args {
		if reflect.TypeOf(o).Kind() != reflect.Slice {
			*result = append(*result, o)
		} else {
			g.flat(o.([]interface{}), result)
		}
	}
}
func (g *AbstractSQLGod) _sql(sb *strings.Builder) {
	panic("sql not support!")
}

func (g *AbstractSQLGod) _countSQL(sb *strings.Builder) {
	panic("count sql not support!")
}

// A struct for binding SQLSM objects to aliases
type A struct {
	sm    *SQLSM
	alias string
}

// A struct for binding EXP objects to aliases
type B struct {
	exp   *EXP
	alias string
}
