package mak

import (
	"github.com/CCLooMi/sql-mak/utils"
	"log"
	"strings"
)

type AbstractSQLMak struct {
	SQLMak
	logSql     bool
	args       []interface{}
	batchArgs  [][]interface{}
	hasSubArgs bool
	child      AbstractSQLMakChild
}

// 定义子类需要实现的接口
type AbstractSQLMakChild interface {
	_sql(sb *strings.Builder)
	_countSQL(sb *strings.Builder)
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

func NewAbstractSQLMak(child *AbstractSQLMakChild) *AbstractSQLMak {
	agod := &AbstractSQLMak{
		logSql:     true,
		args:       make([]interface{}, 0),
		hasSubArgs: false,
	}
	agod.child = *child
	return agod
}

func (g *AbstractSQLMak) TableName(table interface{}) string {
	return utils.TableName(table)
}

func (g *AbstractSQLMak) IF(condition bool, fs ...func(g *AbstractSQLMak)) *AbstractSQLMak {
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

func (g *AbstractSQLMak) SWITCH(i int, fs ...func(g *AbstractSQLMak)) *AbstractSQLMak {
	if i < 0 || i > len(fs)-1 {
		return g
	}
	fs[i](g)
	return g
}

func (g *AbstractSQLMak) LOGSQL(b bool) *AbstractSQLMak {
	g.logSql = b
	return g
}

func (g *AbstractSQLMak) Sql() string {
	var sb strings.Builder
	g.child._sql(&sb)
	if g.logSql {
		log.Println("sql:", sb.String(), g.args)
	}
	return sb.String()
}

func (g *AbstractSQLMak) CountSql() string {
	var sb strings.Builder
	g.child._countSQL(&sb)
	if g.logSql {
		log.Println("sql:", sb.String(), g.args)
	}
	return sb.String()
}

func (g *AbstractSQLMak) Args() []interface{} {
	if !g.hasSubArgs {
		return g.args
	} else {
		ags := make([]interface{}, 0)
		g.flat(g.args, &ags)
		return ags
	}
}

func (g *AbstractSQLMak) BatchArgs() [][]interface{} {
	return g.batchArgs
}

func (g *AbstractSQLMak) flat(args []interface{}, result *[]interface{}) {
	for _, o := range args {
		if _, ok := o.([]interface{}); ok {
			g.flat(o.([]interface{}), result)
			continue
		}
		if _, ok := o.(*[]interface{}); ok {
			g.flat(*o.(*[]interface{}), result)
			continue
		}
		*result = append(*result, o)
	}
}

func (g *AbstractSQLMak) toSQLMak() SQLMak {
	return g
}
