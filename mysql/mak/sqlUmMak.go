package mak

import (
	"database/sql"
	"fmt"
	"log"
	"reflect"
	"strings"
)

type SQLUM struct {
	AbstractSQLMak
	AbstractSQLMakChild
	tables []string
	joins  []interface{}
	sets   []string
	where  string
	andor  []string
	limit  string
}

func NewSQLUM() *SQLUM {
	um := &SQLUM{
		where:  "WHERE 1=1",
		tables: make([]string, 0),
		sets:   make([]string, 0),
	}
	a := um.toAbstractSQLMakChild()
	um.AbstractSQLMak = *NewAbstractSQLMak(&a)
	return um
}

func (um *SQLUM) toAbstractSQLMakChild() AbstractSQLMakChild {
	return um
}
func (um *SQLUM) UPDATE(table interface{}, alias string) *SQLUM {
	switch tv := table.(type) {
	case reflect.Type:
		um.tables = append(um.tables, um.TableName(tv)+" "+alias)
	case string:
		um.tables = append(um.tables, tv+" "+alias)
	}
	return um
}
func (um *SQLUM) LEFT_JOIN(table interface{}, alias string, on string, args ...interface{}) *SQLUM {
	switch tv := table.(type) {
	case reflect.Type:
		if um.joins == nil {
			um.joins = make([]interface{}, 0)
		}
		um.joins = append(um.joins, "LEFT JOIN "+um.TableName(tv)+" "+alias+" ON "+on)
	case string:
		if um.joins == nil {
			um.joins = make([]interface{}, 0)
		}
		um.joins = append(um.joins, "LEFT JOIN "+tv+" "+alias+" ON "+on)
	case *SQLSM:
		um.hasSubArgs = true
		if um.joins == nil {
			um.joins = make([]interface{}, 0)
		}
		um.joins = append(um.joins, []interface{}{"LEFT JOIN ", A{tv, alias}, " ON " + on})
		um.args = append(um.args, tv.args...)
	}
	um.args = append(um.args, args...)
	return um
}

func (um *SQLUM) RIGHT_JOIN(table interface{}, alias string, on string, args ...interface{}) *SQLUM {
	switch tv := table.(type) {
	case reflect.Type:
		if um.joins == nil {
			um.joins = make([]interface{}, 0)
		}
		um.joins = append(um.joins, "RIGHT JOIN "+um.TableName(tv)+" "+alias+" ON "+on)
	case string:
		if um.joins == nil {
			um.joins = make([]interface{}, 0)
		}
		um.joins = append(um.joins, "RIGHT JOIN "+tv+" "+alias+" ON "+on)
	case *SQLSM:
		um.hasSubArgs = true
		if um.joins == nil {
			um.joins = make([]interface{}, 0)
		}
		um.joins = append(um.joins, []interface{}{"RIGHT JOIN ", A{tv, alias}, " ON " + on})
		um.args = append(um.args, tv.args...)
	}
	um.args = append(um.args, args...)
	return um
}

func (um *SQLUM) INNER_JOIN(table interface{}, alias string, on string, args ...interface{}) *SQLUM {
	switch tv := table.(type) {
	case reflect.Type:
		if um.joins == nil {
			um.joins = make([]interface{}, 0)
		}
		um.joins = append(um.joins, "INNER JOIN "+um.TableName(tv)+" "+alias+" ON "+on)
	case string:
		if um.joins == nil {
			um.joins = make([]interface{}, 0)
		}
		um.joins = append(um.joins, "INNER JOIN "+tv+" "+alias+" ON "+on)
	case *SQLSM:
		um.hasSubArgs = true
		if um.joins == nil {
			um.joins = make([]interface{}, 0)
		}
		um.joins = append(um.joins, []interface{}{"INNER JOIN ", A{tv, alias}, " ON " + on})
		um.args = append(um.args, tv.args...)
	}
	um.args = append(um.args, args...)
	return um
}

func (um *SQLUM) SET(set string, args ...interface{}) *SQLUM {
	um.sets = append(um.sets, set)
	um.args = append(um.args, args...)
	return um
}

func (um *SQLUM) WHERE(where string, args ...interface{}) *SQLUM {
	um.where = "WHERE " + where
	um.args = append(um.args, args...)
	return um
}

func (um *SQLUM) AND(and string, args ...interface{}) *SQLUM {
	if um.andor == nil {
		um.andor = make([]string, 0)
	}
	um.andor = append(um.andor, "AND "+and)
	um.args = append(um.args, args...)
	return um
}

func (um *SQLUM) OR(or string, args ...interface{}) *SQLUM {
	if um.andor == nil {
		um.andor = make([]string, 0)
	}
	um.andor = append(um.andor, "OR "+or)
	um.args = append(um.args, args...)
	return um
}

func (um *SQLUM) WHERE_IN(column string, inOrNotIn string, args ...interface{}) *SQLUM {
	if len(args) > 0 {
		sb := strings.Builder{}
		sb.WriteString("WHERE " + column + " " + inOrNotIn + " (")
		for i, arg := range args {
			um.args = append(um.args, arg)
			if i != len(args)-1 {
				sb.WriteString("?,")
				continue
			}
			sb.WriteRune('?')
		}
		sb.WriteString(")")
		um.where = sb.String()
	}
	return um
}

func (um *SQLUM) AND_IN(column string, inOrNotIn string, args ...interface{}) *SQLUM {
	if len(args) > 0 {
		sb := strings.Builder{}
		sb.WriteString("AND " + column + " " + inOrNotIn + " (")
		for i, arg := range args {
			um.args = append(um.args, arg)
			if i != len(args)-1 {
				sb.WriteString("?,")
				continue
			}
			sb.WriteRune('?')
		}
		sb.WriteString(")")
		if um.andor == nil {
			um.andor = make([]string, 0)
		}
		um.andor = append(um.andor, sb.String())
	}
	return um
}

func (um *SQLUM) OR_IN(column string, inOrNotIn string, args ...interface{}) *SQLUM {
	if len(args) > 0 {
		sb := strings.Builder{}
		sb.WriteString("OR " + column + " " + inOrNotIn + " (")
		for i, arg := range args {
			um.args = append(um.args, arg)
			if i != len(args)-1 {
				sb.WriteString("?,")
				continue
			}
			sb.WriteRune('?')
		}
		sb.WriteString(")")
		if um.andor == nil {
			um.andor = make([]string, 0)
		}
		um.andor = append(um.andor, sb.String())
	}
	return um
}

func (um *SQLUM) LIMIT(limits ...interface{}) *SQLUM {
	if len(limits) > 0 {
		if um.hasSubArgs {
			log.Fatal("非单表更新 SQL 不能使用 LIMIT，请在子查询中使用 LIMIT 查出需要更新的数据再使用 UPDATE 进行更新（无需 LIMIT）")
		}
		sb := strings.Builder{}
		sb.WriteString("LIMIT ")
		for i, lm := range limits {
			sb.WriteString(fmt.Sprint(lm))
			if i != len(limits)-1 {
				sb.WriteRune(',')
			}
		}
		sb.WriteString(" ")
		um.limit = sb.String()
	}
	return um
}

func (um *SQLUM) SetBatchArgs(batchArgs ...[]interface{}) *SQLUM {
	if um.batchArgs == nil {
		um.batchArgs = make([][]interface{}, 0)
	}
	um.batchArgs = append(um.batchArgs, batchArgs...)
	return um
}

func (um *SQLUM) Execute(mdb *sql.DB) *MySQLUMExecutor {
	return NewMySQLUMExecutor(um, mdb)
}
func (um *SQLUM) _sql(sb *strings.Builder) {
	sb.WriteString("UPDATE ")
	if len(um.tables) > 0 {
		for i, table := range um.tables {
			sb.WriteString(table)
			if i != len(um.tables)-1 {
				sb.WriteRune(',')
			}
		}
		sb.WriteString(" ")
	}
	if um.joins != nil {
		for _, join := range um.joins {
			if j, ok := join.(string); ok {
				sb.WriteString(j)
			} else {
				oa := join.([3]interface{})
				um := oa[1].(A)
				sb.WriteString(oa[0].(string) + "(")
				um.sm._sql(sb)
				sb.WriteString(")" + oa[2].(string))
			}
			sb.WriteString(" ")
		}
	}
	if len(um.sets) > 0 {
		sb.WriteString("SET ")
		for i, set := range um.sets {
			sb.WriteString(set)
			if i != len(um.sets)-1 {
				sb.WriteRune(',')
			}
		}
		sb.WriteString(" ")
	}
	sb.WriteString(um.where)
	sb.WriteString(" ")
	if um.andor != nil {
		for _, ao := range um.andor {
			sb.WriteString(ao)
			sb.WriteString(" ")
		}
	}
	if um.limit != "" {
		sb.WriteString(um.limit)
	}
}
