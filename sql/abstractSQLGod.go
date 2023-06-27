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
	child        AbstractSQLGodChild
}

// 定义子类需要实现的接口
type AbstractSQLGodChild interface {
	_sql(sb *strings.Builder)
	_countSQL(sb *strings.Builder)
}

func NewAbstractSQLGod(child *AbstractSQLGodChild) *AbstractSQLGod {
	agod := &AbstractSQLGod{
		log:          log.Default(),
		logSql:       true,
		args:         make([]interface{}, 0),
		hasSubSelect: false,
	}
	agod.child = *child
	return agod
}

func (g *AbstractSQLGod) TableName(t reflect.Type) string {
	// 获取结构体名称
	name := t.Name()

	// 判断结构体中是否定义了TableName()方法
	method, ok := t.MethodByName("TableName")
	if ok {
		// 创建零值的结构体对象
		obj := reflect.Zero(t).Interface()

		// 调用TableName()方法获取表名
		values := method.Func.Call([]reflect.Value{reflect.ValueOf(obj)})
		name = values[0].String()
	}

	return name
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

func (g *AbstractSQLGod) Sql() string {
	var sb strings.Builder
	g.child._sql(&sb)
	if g.logSql {
		g.log.Println("sql:", sb.String(), g.args)
	}
	return sb.String()
}

func (g *AbstractSQLGod) CountSql() string {
	var sb strings.Builder
	g.child._countSQL(&sb)
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

func (g *AbstractSQLGod) toSQLGod() SQLGod {
	return g
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
