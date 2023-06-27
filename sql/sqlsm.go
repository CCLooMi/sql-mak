package sql

import (
	"fmt"
	"reflect"
	"strings"
)

type SQLSM struct {
	AbstractSQLGod
	AbstractSQLGodChild
	columns      []interface{}
	columnAlias  []string
	tables       []interface{}
	unions       []interface{}
	joins        []interface{}
	where        string
	andor        []string
	orderBy      []string
	groupBy      []string
	limit        string
	hasSubSelect bool
}

func NewSQLSM() *SQLSM {
	sm := &SQLSM{
		where:       "WHERE 1=1",
		columns:     make([]interface{}, 0),
		columnAlias: make([]string, 0),
		tables:      make([]interface{}, 0),
	}
	a := sm.toAbstractSQLGodChild()
	sm.AbstractSQLGod = *NewAbstractSQLGod(&a)
	return sm
}
func (sm *SQLSM) toAbstractSQLGodChild() AbstractSQLGodChild {
	return sm
}

func (sm *SQLSM) ColumnAlias() []string {
	return sm.columnAlias
}

func (s *SQLSM) SELECT(cols ...interface{}) *SQLSM {
	for _, col := range cols {
		switch colV := col.(type) {
		case string:
			s.columns = append(s.columns, colV)
			s.columnAlias = append(s.columnAlias, colV[strings.LastIndex(colV, ".")+1:])
		case EXP:
			exp := col.(EXP)
			s.columns = append(s.columns, exp)
			s.args = append(s.args, exp.args...)
			s.hasSubSelect = true
		}
	}
	return s
}

func (s *SQLSM) SELECT_AS(column, alias string) *SQLSM {
	s.columns = append(s.columns, column+" AS '"+alias+"'")
	s.columnAlias = append(s.columnAlias, alias)
	return s
}

func (s *SQLSM) SELECT_AS_SM(column *SQLSM, alias string) *SQLSM {
	s.hasSubSelect = true
	s.columns = append(s.columns, A{sm: column, alias: "'" + alias + "'"})
	s.columnAlias = append(s.columnAlias, alias)
	s.args = append(s.args, column.args...)
	return s
}

func (s *SQLSM) SELECT_AS_EXP(exp *EXP, alias string) *SQLSM {
	s.hasSubSelect = true
	s.columns = append(s.columns, B{exp: exp, alias: "'" + alias + "'"})
	s.columnAlias = append(s.columnAlias, alias)
	s.args = append(s.args, exp.args...)
	return s
}
func (s *SQLSM) FROM(table interface{}, alias string) *SQLSM {
	switch v := table.(type) {
	case reflect.Type:
		s.tables = append(s.tables, s.TableName(v)+" "+alias)
	case string:
		s.tables = append(s.tables, table.(string)+" "+alias)
	case *SQLSM:
		s.hasSubSelect = true
		s.tables = append(s.tables, A{sm: table.(*SQLSM), alias: alias})
		s.args = append(s.args, table.(*SQLSM).args...)
	}
	return s
}

func (s *SQLSM) JOIN(table interface{}, alias string, args ...interface{}) *SQLSM {
	var join string
	switch v := table.(type) {
	case reflect.Type:
		join = "JOIN " + s.TableName(v) + " " + alias
	case string:
		join = "JOIN " + table.(string) + " " + alias
	case *SQLSM:
		s.hasSubSelect = true
		s.joins = append(s.joins, []interface{}{"JOIN", A{sm: table.(*SQLSM), alias: alias}})
		s.args = append(s.args, table.(*SQLSM).args...)
	}

	if len(args) > 0 {
		s.args = append(s.args, args...)
	}

	s.joins = append(s.joins, join)
	return s
}
func (s *SQLSM) JOIN_ON(table interface{}, alias string, on string, args ...interface{}) *SQLSM {
	var join string
	switch v := table.(type) {
	case reflect.Type:
		join = "JOIN " + s.TableName(v) + " " + alias + " ON " + on
	case string:
		join = "JOIN " + table.(string) + " " + alias + " ON " + on
	case *SQLSM:
		s.hasSubSelect = true
		s.joins = append(s.joins, []interface{}{"JOIN", A{sm: table.(*SQLSM), alias: alias}, " ON " + on})
		s.args = append(s.args, table.(*SQLSM).args...)
	}

	if len(args) > 0 {
		s.args = append(s.args, args...)
	}

	s.joins = append(s.joins, join)
	return s
}

func (s *SQLSM) LEFT_JOIN(table interface{}, alias string, on string, args ...interface{}) *SQLSM {
	var join string
	switch v := table.(type) {
	case reflect.Type:
		join = "LEFT JOIN " + s.TableName(v) + " " + alias + " ON " + on
	case string:
		join = "LEFT JOIN " + table.(string) + " " + alias + " ON " + on
	case *SQLSM:
		s.hasSubSelect = true
		s.joins = append(s.joins, []interface{}{"LEFT JOIN ", A{sm: table.(*SQLSM), alias: alias}, " ON " + on})
		s.args = append(s.args, table.(*SQLSM).args...)
	}

	if len(args) > 0 {
		s.args = append(s.args, args...)
	}

	s.joins = append(s.joins, join)
	return s
}

func (s *SQLSM) RIGHT_JOIN(table interface{}, alias string, on string, args ...interface{}) *SQLSM {
	var join string
	switch v := table.(type) {
	case reflect.Type:
		join = "RIGHT JOIN " + s.TableName(v) + " " + alias + " ON " + on
	case string:
		join = "RIGHT JOIN " + table.(string) + " " + alias + " ON " + on
	case *SQLSM:
		s.hasSubSelect = true
		s.joins = append(s.joins, []interface{}{"RIGHT JOIN ", A{sm: table.(*SQLSM), alias: alias}, " ON " + on})
		s.args = append(s.args, table.(*SQLSM).args...)
	}

	if len(args) > 0 {
		s.args = append(s.args, args...)
	}

	s.joins = append(s.joins, join)
	return s
}

func (s *SQLSM) INNER_JOIN(table interface{}, alias string, on string, args ...interface{}) *SQLSM {
	var join string
	switch v := table.(type) {
	case reflect.Type:
		join = "INNER JOIN " + s.TableName(v) + " " + alias + " ON " + on
	case string:
		join = "INNER JOIN " + table.(string) + " " + alias + " ON " + on
	case *SQLSM:
		s.hasSubSelect = true
		s.joins = append(s.joins, []interface{}{"INNER JOIN ", A{sm: table.(*SQLSM), alias: alias}, " ON " + on})
		s.args = append(s.args, table.(*SQLSM).args...)
	}

	if len(args) > 0 {
		s.args = append(s.args, args...)
	}

	s.joins = append(s.joins, join)
	return s
}

func (s *SQLSM) UNION(sm *SQLSM, args ...interface{}) *SQLSM {
	s.hasSubSelect = true
	if s.unions == nil {
		s.unions = make([]interface{}, 0)
	}
	s.unions = append(s.unions, [2]interface{}{"UNION ", A{sm: sm}})

	s.args = append(s.args, sm.args...)
	s.args = append(s.args, args...)

	return s
}

func (s *SQLSM) UNION_ALIAS(sm *SQLSM, alias string, args ...interface{}) *SQLSM {
	s.hasSubSelect = true
	if s.unions == nil {
		s.unions = make([]interface{}, 0)
	}
	s.unions = append(s.unions, [2]interface{}{"UNION SELECT * FROM ", A{sm: sm, alias: alias}})

	s.args = append(s.args, sm.args...)
	s.args = append(s.args, args...)

	return s
}

func (s *SQLSM) UNION_ALL(sm *SQLSM, args ...interface{}) *SQLSM {
	s.hasSubSelect = true
	if s.unions == nil {
		s.unions = make([]interface{}, 0)
	}
	s.unions = append(s.unions, [2]interface{}{"UNION ALL ", A{sm: sm}})

	s.args = append(s.args, sm.args...)
	s.args = append(s.args, args...)

	return s
}

func (s *SQLSM) UNION_ALL_ALIAS(sm *SQLSM, alias string, args ...interface{}) *SQLSM {
	s.hasSubSelect = true
	if s.unions == nil {
		s.unions = make([]interface{}, 0)
	}
	s.unions = append(s.unions, [2]interface{}{"UNION ALL SELECT * FROM ", A{sm: sm, alias: alias}})

	s.args = append(s.args, sm.args...)
	s.args = append(s.args, args...)

	return s
}
func (s *SQLSM) WHERE(where string, args ...interface{}) *SQLSM {
	s.where = "WHERE " + where
	s.args = append(s.args, args...)
	return s
}

func (s *SQLSM) WHERE_IN(column string, inOrNotIn IN, args ...interface{}) *SQLSM {
	if len(args) > 0 {
		var sb strings.Builder
		sb.WriteString("WHERE ")
		sb.WriteString(column)
		sb.WriteRune(' ')
		sb.WriteString(inOrNotIn.value())
		sb.WriteRune('(')
		for i, arg := range args {
			s.args = append(s.args, arg)
			sb.WriteRune('?')
			if i < len(args)-1 {
				sb.WriteRune(',')
			}
		}
		sb.WriteRune(')')
		s.where = sb.String()
	}
	return s
}

func (s *SQLSM) WHERE_SUBQUERY(column string, inOrNotIn IN, subquery *SQLSM) *SQLSM {
	var sb strings.Builder
	sb.WriteString("WHERE ")
	sb.WriteString(column)
	sb.WriteRune(' ')
	sb.WriteString(inOrNotIn.value())
	sb.WriteRune('(')
	sb.WriteString(subquery.Sql())
	sb.WriteRune(')')
	s.args = append(s.args, subquery.args...)
	s.where = sb.String()
	return s
}

func (s *SQLSM) AND(and string, args ...interface{}) *SQLSM {
	if s.andor == nil {
		s.andor = make([]string, 0)
	}
	s.andor = append(s.andor, "AND "+and)
	s.args = append(s.args, args...)
	return s
}

func (s *SQLSM) AND_IN(column string, inOrNotIn IN, args ...interface{}) *SQLSM {
	if len(args) > 0 {
		var sb strings.Builder
		sb.WriteString("AND ")
		sb.WriteString(column)
		sb.WriteRune(' ')
		sb.WriteString(inOrNotIn.value())
		sb.WriteRune('(')
		for i, arg := range args {
			s.args = append(s.args, arg)
			sb.WriteRune('?')
			if i != len(args)-1 {
				sb.WriteRune(',')
			}
		}
		sb.WriteRune(')')
		if s.andor == nil {
			s.andor = make([]string, 0)
		}
		s.andor = append(s.andor, sb.String())
	}
	return s
}

func (s *SQLSM) AND_SUBQUERY(column string, inOrNotIn IN, subquery *SQLSM) *SQLSM {
	var sb strings.Builder
	sb.WriteString("AND ")
	sb.WriteString(column)
	sb.WriteRune(' ')
	sb.WriteString(inOrNotIn.value())
	sb.WriteRune('(')
	sb.WriteString(subquery.Sql())
	sb.WriteRune(')')
	s.args = append(s.args, subquery.args...)
	if s.andor == nil {
		s.andor = make([]string, 0)
	}
	s.andor = append(s.andor, sb.String())
	return s
}
func (s *SQLSM) OR(or string, args ...interface{}) *SQLSM {
	if s.andor == nil {
		s.andor = make([]string, 0)
	}
	s.andor = append(s.andor, "OR "+or)
	s.args = append(s.args, args...)
	return s
}

func (s *SQLSM) OR_IN(column string, inOrNotIn IN, args ...interface{}) *SQLSM {
	if len(args) > 0 {
		var sb strings.Builder
		sb.WriteString("OR ")
		sb.WriteString(column)
		sb.WriteRune(' ')
		sb.WriteString(inOrNotIn.value())
		sb.WriteRune('(')
		for i, arg := range args {
			s.args = append(s.args, arg)
			sb.WriteRune('?')
			if i != len(args)-1 {
				sb.WriteRune(',')
			}
		}
		sb.WriteRune(')')
		if s.andor == nil {
			s.andor = make([]string, 0)
		}
		s.andor = append(s.andor, sb.String())
	}
	return s
}

func (s *SQLSM) OR_SUBQUERY(column string, inOrNotIn IN, subquery *SQLSM) *SQLSM {
	var sb strings.Builder
	sb.WriteString("OR ")
	sb.WriteString(column)
	sb.WriteRune(' ')
	sb.WriteString(inOrNotIn.value())
	sb.WriteRune('(')
	sb.WriteString(subquery.Sql())
	sb.WriteRune(')')
	s.args = append(s.args, subquery.args...)
	if s.andor == nil {
		s.andor = make([]string, 0)
	}
	s.andor = append(s.andor, sb.String())
	return s
}

func (s *SQLSM) GROUP_BY(gbs ...string) *SQLSM {
	if s.groupBy == nil {
		s.groupBy = make([]string, 0)
	}
	s.groupBy = append(s.groupBy, gbs...)
	return s
}

func (s *SQLSM) ORDER_BY(obs ...string) *SQLSM {
	if s.orderBy == nil {
		s.orderBy = make([]string, 0)
	}
	s.orderBy = append(s.orderBy, obs...)
	return s
}

func (s *SQLSM) LIMIT(limits ...interface{}) *SQLSM {
	if len(limits) > 0 {
		var sb strings.Builder
		sb.WriteString("LIMIT ")
		for i, lm := range limits {
			sb.WriteString(fmt.Sprintf("%v", lm))
			if i != len(limits)-1 {
				sb.WriteRune(',')
			}
		}
		sb.WriteRune(' ')
		s.limit = sb.String()
	}
	return s
}
func (s *SQLSM) ExpSQL() string {
	var sb strings.Builder
	sb.WriteRune('(')
	s._sql(&sb)
	sb.WriteRune(')')
	return sb.String()
}

func (s *SQLSM) Execute() SQLSMExecutor {
	god := s.toSQLGod()
	executor, _ := GetExecutor(&god, reflect.TypeOf(SQLSMExecutor{})).(SQLSMExecutor)
	return executor
}

func (s *SQLSM) _sql(sb *strings.Builder) {
	sb.WriteString("SELECT ")
	if len(s.columns) > 0 {
		for i, col := range s.columns {
			switch colV := col.(type) {
			case string:
				sb.WriteString(colV)
			case A:
				sb.WriteRune('(')
				colV.sm._sql(sb)
				sb.WriteString(") AS ")
				sb.WriteString(colV.alias)
			case B:
				exp := colV.exp
				if colV.alias != "" {
					sb.WriteString(exp.Exp())
					sb.WriteString(" AS ")
					sb.WriteString(colV.alias)
				} else {
					sb.WriteString(exp.Exp())
				}
			}
			if i != len(s.columns)-1 {
				sb.WriteRune(',')
			}
		}
		sb.WriteRune(' ')
	} else {
		sb.WriteString("* ")
	}
	if len(s.tables) > 0 {
		sb.WriteString("FROM ")
		for i, table := range s.tables {
			switch tv := table.(type) {
			case string:
				sb.WriteString(tv)
			case A:
				sm := tv.sm
				sb.WriteRune('(')
				sm._sql(sb)
				sb.WriteRune(')')
				sb.WriteString(tv.alias)
			}
			if i != len(s.tables) {
				sb.WriteRune(',')
			}
		}
		sb.WriteRune(' ')
	}
	if s.joins != nil {
		for _, join := range s.joins {
			switch jv := join.(type) {
			case string:
				sb.WriteString(jv)
			default:
				oa := join.([3]interface{})
				sb.WriteString(oa[0].(string))
				oa1 := oa[1].(A)
				sm := oa1.sm
				if oa1.alias != "" {
					sb.WriteRune('(')
					sm._sql(sb)
					sb.WriteRune(')')
					sb.WriteString(oa1.alias)
				} else {
					sm._sql(sb)
				}
				if len(oa) > 2 {
					sb.WriteString(oa[2].(string))
				}
			}
			sb.WriteRune(' ')
		}
	}
	sb.WriteString(s.where)
	sb.WriteRune(' ')
	if s.andor != nil {
		for _, cd := range s.andor {
			sb.WriteString(cd)
			sb.WriteRune(' ')
		}
	}
	if s.groupBy != nil {
		sb.WriteString("GROUP BY ")
		for i, gb := range s.groupBy {
			sb.WriteString(gb)
			if i != len(s.groupBy) {
				sb.WriteRune(',')
			}
		}
		sb.WriteRune(' ')
	}
	if s.orderBy != nil {
		sb.WriteString("ORDER BY ")
		for i, ob := range s.orderBy {
			sb.WriteString(ob)
			if i != len(s.orderBy) {
				sb.WriteRune(',')
			}
		}
		sb.WriteRune(' ')
	}
	if s.unions != nil {
		for _, union := range s.unions {
			switch uv := union.(type) {
			case string:
				sb.WriteString(uv)
			default:
				oa := union.([3]interface{})
				sb.WriteString(oa[0].(string))
				oa1 := oa[1].(A)
				sm := oa1.sm
				if oa1.alias != "" {
					sb.WriteRune('(')
					sm._sql(sb)
					sb.WriteRune(')')
					sb.WriteString(oa1.alias)
				} else {
					sm._sql(sb)
				}
				if len(oa) > 2 {
					sb.WriteString(oa[2].(string))
				}
			}
			sb.WriteRune(' ')
		}
	}
	if s.limit != "" {
		sb.WriteString(s.limit)
	}
}
func (s *SQLSM) _countSQL(sb *strings.Builder) {
	sb.WriteString("SELECT COUNT(1) ")
	if s.groupBy != nil {
		sb.WriteString("FROM (SELECT ")
		for i, gb := range s.groupBy {
			sb.WriteString(gb)
			if i != len(s.groupBy)-1 {
				sb.WriteRune(',')
			}
		}
		sb.WriteRune(' ')
	}
	if len(s.tables) > 0 {
		sb.WriteString("FROM ")
		for i, table := range s.tables {
			switch tv := table.(type) {
			case string:
				sb.WriteString(tv)
			case A:
				sm := tv.sm
				sb.WriteRune('(')
				sm._sql(sb)
				sb.WriteRune(')')
				sb.WriteString(tv.alias)
			}
			if i != len(s.tables)-1 {
				sb.WriteRune(',')
			}
		}
		sb.WriteRune(' ')
	}
	if s.joins != nil {
		for _, join := range s.joins {
			switch jv := join.(type) {
			case string:
				sb.WriteString(jv)
			default:
				oa := join.([3]interface{})
				sb.WriteString(oa[0].(string))
				oa1 := oa[1].(A)
				sm := oa1.sm
				if oa1.alias != "" {
					sb.WriteRune('(')
					sm._sql(sb)
					sb.WriteRune(')')
					sb.WriteString(oa1.alias)
				} else {
					sm._sql(sb)
				}
				if len(oa) > 2 {
					sb.WriteString(oa[2].(string))
				}
			}
			sb.WriteRune(' ')
		}
	}
	sb.WriteString(s.where)
	sb.WriteRune(' ')
	if s.andor != nil {
		for _, cd := range s.andor {
			sb.WriteString(cd)
			sb.WriteRune(' ')
		}
	}
	if s.groupBy != nil {
		sb.WriteString("GROUP BY ")
		for i, gb := range s.groupBy {
			sb.WriteString(gb)
			if i != len(s.groupBy)-1 {
				sb.WriteRune(',')
			}
		}
		sb.WriteRune(')')
		sb.WriteRune('a')
	}
}
