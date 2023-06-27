package god

import (
	"fmt"
	"log"
	"reflect"
	"strings"
)

type SQLUM struct {
	AbstractSQLGod
	AbstractSQLGodChild
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
	a := um.toAbstractSQLGodChild()
	um.AbstractSQLGod = *NewAbstractSQLGod(&a)
	return um
}

func (sm *SQLUM) toAbstractSQLGodChild() AbstractSQLGodChild {
	return sm
}
func (sm *SQLUM) UPDATE(table interface{}, alias string) *SQLUM {
	switch tv := table.(type) {
	case reflect.Type:
		sm.tables = append(sm.tables, sm.TableName(tv)+" "+alias)
	case string:
		sm.tables = append(sm.tables, tv+" "+alias)
	}
	return sm
}
func (sm *SQLUM) LEFT_JOIN(table interface{}, alias string, on string, args ...interface{}) *SQLUM {
	switch tv := table.(type) {
	case reflect.Type:
		if sm.joins == nil {
			sm.joins = make([]interface{}, 0)
		}
		sm.joins = append(sm.joins, "LEFT JOIN "+sm.TableName(tv)+" "+alias+" ON "+on)
	case string:
		if sm.joins == nil {
			sm.joins = make([]interface{}, 0)
		}
		sm.joins = append(sm.joins, "LEFT JOIN "+tv+" "+alias+" ON "+on)
	case *SQLSM:
		sm.hasSubSelect = true
		if sm.joins == nil {
			sm.joins = make([]interface{}, 0)
		}
		sm.joins = append(sm.joins, []interface{}{"LEFT JOIN ", A{tv, alias}, " ON " + on})
		sm.args = append(sm.args, tv.args...)
	}
	sm.args = append(sm.args, args...)
	return sm
}

func (sm *SQLUM) RIGHT_JOIN(table interface{}, alias string, on string, args ...interface{}) *SQLUM {
	switch tv := table.(type) {
	case reflect.Type:
		if sm.joins == nil {
			sm.joins = make([]interface{}, 0)
		}
		sm.joins = append(sm.joins, "RIGHT JOIN "+sm.TableName(tv)+" "+alias+" ON "+on)
	case string:
		if sm.joins == nil {
			sm.joins = make([]interface{}, 0)
		}
		sm.joins = append(sm.joins, "RIGHT JOIN "+tv+" "+alias+" ON "+on)
	case *SQLSM:
		sm.hasSubSelect = true
		if sm.joins == nil {
			sm.joins = make([]interface{}, 0)
		}
		sm.joins = append(sm.joins, []interface{}{"RIGHT JOIN ", A{tv, alias}, " ON " + on})
		sm.args = append(sm.args, tv.args...)
	}
	sm.args = append(sm.args, args...)
	return sm
}

func (sm *SQLUM) INNER_JOIN(table interface{}, alias string, on string, args ...interface{}) *SQLUM {
	switch tv := table.(type) {
	case reflect.Type:
		if sm.joins == nil {
			sm.joins = make([]interface{}, 0)
		}
		sm.joins = append(sm.joins, "INNER JOIN "+sm.TableName(tv)+" "+alias+" ON "+on)
	case string:
		if sm.joins == nil {
			sm.joins = make([]interface{}, 0)
		}
		sm.joins = append(sm.joins, "INNER JOIN "+tv+" "+alias+" ON "+on)
	case *SQLSM:
		sm.hasSubSelect = true
		if sm.joins == nil {
			sm.joins = make([]interface{}, 0)
		}
		sm.joins = append(sm.joins, []interface{}{"INNER JOIN ", A{tv, alias}, " ON " + on})
		sm.args = append(sm.args, tv.args...)
	}
	sm.args = append(sm.args, args...)
	return sm
}

func (sm *SQLUM) SET(set string, args ...interface{}) *SQLUM {
	sm.sets = append(sm.sets, set)
	sm.args = append(sm.args, args...)
	return sm
}

func (sm *SQLUM) WHERE(where string, args ...interface{}) *SQLUM {
	sm.where = "WHERE " + where
	sm.args = append(sm.args, args...)
	return sm
}

func (sm *SQLUM) AND(and string, args ...interface{}) *SQLUM {
	if sm.andor == nil {
		sm.andor = make([]string, 0)
	}
	sm.andor = append(sm.andor, "AND "+and)
	sm.args = append(sm.args, args...)
	return sm
}

func (sm *SQLUM) OR(or string, args ...interface{}) *SQLUM {
	if sm.andor == nil {
		sm.andor = make([]string, 0)
	}
	sm.andor = append(sm.andor, "OR "+or)
	sm.args = append(sm.args, args...)
	return sm
}

func (sm *SQLUM) WHERE_IN(column string, inOrNotIn string, args ...interface{}) *SQLUM {
	if len(args) > 0 {
		sb := strings.Builder{}
		sb.WriteString("WHERE " + column + " " + inOrNotIn + " (")
		for i, arg := range args {
			sm.args = append(sm.args, arg)
			if i != len(args)-1 {
				sb.WriteString("?,")
			}
		}
		sb.WriteString(")")
		sm.where = sb.String()
	}
	return sm
}

func (sm *SQLUM) AND_IN(column string, inOrNotIn string, args ...interface{}) *SQLUM {
	if len(args) > 0 {
		sb := strings.Builder{}
		sb.WriteString("AND " + column + " " + inOrNotIn + " (")
		for i, arg := range args {
			sm.args = append(sm.args, arg)
			if i != len(args)-1 {
				sb.WriteString("?,")
			}
		}
		sb.WriteString(")")
		if sm.andor == nil {
			sm.andor = make([]string, 0)
		}
		sm.andor = append(sm.andor, sb.String())
	}
	return sm
}

func (sm *SQLUM) OR_IN(column string, inOrNotIn string, args ...interface{}) *SQLUM {
	if len(args) > 0 {
		sb := strings.Builder{}
		sb.WriteString("OR " + column + " " + inOrNotIn + " (")
		for i, arg := range args {
			sm.args = append(sm.args, arg)
			if i != len(args)-1 {
				sb.WriteString("?,")
			}
		}
		sb.WriteString(")")
		if sm.andor == nil {
			sm.andor = make([]string, 0)
		}
		sm.andor = append(sm.andor, sb.String())
	}
	return sm
}

func (sm *SQLUM) LIMIT(limits ...interface{}) *SQLUM {
	if len(limits) > 0 {
		if sm.hasSubSelect {
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
		sm.limit = sb.String()
	}
	return sm
}

func (sm *SQLUM) BatchArgs(batchArgs ...[]interface{}) *SQLUM {
	if sm.batchArgs == nil {
		sm.batchArgs = make([][]interface{}, 0)
	}
	sm.batchArgs = append(sm.batchArgs, batchArgs...)
	return sm
}

func (sm *SQLUM) SetBatchArgs(batchArgs [][]interface{}) *SQLUM {
	sm.batchArgs = batchArgs
	return sm
}

func (sm *SQLUM) Execute() *SQLUMExecutor {
	god := sm.toSQLGod()
	executor, _ := GetExecutor(god, reflect.TypeOf(SQLUMExecutor{})).(SQLUMExecutor)
	return &executor
}
func (sm *SQLUM) _sql(sb *strings.Builder) {
	sb.WriteString("UPDATE ")
	if len(sm.tables) > 0 {
		for i, table := range sm.tables {
			sb.WriteString(table)
			if i != len(sm.tables)-1 {
				sb.WriteRune(',')
			}
		}
		sb.WriteString(" ")
	}
	if sm.joins != nil {
		for _, join := range sm.joins {
			if j, ok := join.(string); ok {
				sb.WriteString(j)
			} else {
				oa := join.([3]interface{})
				sm := oa[1].(A)
				sb.WriteString(oa[0].(string) + "(")
				sm.sm._sql(sb)
				sb.WriteString(")" + oa[2].(string))
			}
			sb.WriteString(" ")
		}
	}
	if len(sm.sets) > 0 {
		sb.WriteString("SET ")
		for i, set := range sm.sets {
			sb.WriteString(set)
			if i != len(sm.sets)-1 {
				sb.WriteRune(',')
			}
		}
		sb.WriteString(" ")
	}
	sb.WriteString(sm.where)
	sb.WriteString(" ")
	if sm.andor != nil {
		for _, ao := range sm.andor {
			sb.WriteString(ao)
			sb.WriteString(" ")
		}
	}
	if sm.limit != "" {
		sb.WriteString(sm.limit)
	}
}
