package mak

import (
	"database/sql"
	"strings"
)

type SQLDM struct {
	AbstractSQLMak
	AbstractSQLMakChild
	table string
	where string
	andor []string
}

func NewSQLDM() *SQLDM {
	dm := &SQLDM{
		where: "WHERE 1=1",
	}
	a := dm.toAbstractSQLMakChild()
	dm.AbstractSQLMak = *NewAbstractSQLMak(&a)
	return dm
}

func (dm *SQLDM) toAbstractSQLMakChild() AbstractSQLMakChild {
	return dm
}
func (dm *SQLDM) FROM(table interface{}) *SQLDM {
	switch v := table.(type) {
	case string:
		dm.table = v
	default:
		dm.table = dm.TableName(v)
	}
	return dm
}
func (dm *SQLDM) WHERE(where string, args ...interface{}) *SQLDM {
	dm.where = "WHERE " + where
	dm.args = append(dm.args, args...)
	return dm
}

func (dm *SQLDM) AND(and string, args ...interface{}) *SQLDM {
	if dm.andor == nil {
		dm.andor = make([]string, 0)
	}
	dm.andor = append(dm.andor, "AND "+and)
	dm.args = append(dm.args, args...)
	return dm
}

func (dm *SQLDM) OR(or string, args ...interface{}) *SQLDM {
	if dm.andor == nil {
		dm.andor = make([]string, 0)
	}
	dm.andor = append(dm.andor, "OR "+or)
	dm.args = append(dm.args, args...)
	return dm
}

func (dm *SQLDM) WHERE_IN(column string, args ...interface{}) *SQLDM {
	if len(args) > 0 {
		var sb strings.Builder
		sb.WriteString("WHERE ")
		sb.WriteString(column)
		sb.WriteString(" IN(")
		for i, arg := range args {
			dm.args = append(dm.args, arg)
			if i != len(args)-1 {
				sb.WriteString("?,")
				continue
			}
			sb.WriteRune('?')
		}
		sb.WriteByte(')')
		dm.where = sb.String()
	}
	return dm
}

func (dm *SQLDM) WHERE_SUBQUERY(column string, inOrNotIn IN, subquery *SQLSM) *SQLDM {
	var sb strings.Builder
	sb.WriteString("WHERE ")
	sb.WriteString(column)
	sb.WriteRune(' ')
	sb.WriteString(inOrNotIn.value())
	sb.WriteRune('(')
	sb.WriteString(subquery.Sql())
	sb.WriteRune(')')
	dm.args = append(dm.args, subquery.args...)
	dm.where = sb.String()
	return dm
}
func (dm *SQLDM) AND_IN(column string, args ...interface{}) *SQLDM {
	if len(args) > 0 {
		var sb strings.Builder
		sb.WriteString("AND ")
		sb.WriteString(column)
		sb.WriteString(" IN(")
		for i, arg := range args {
			dm.args = append(dm.args, arg)
			if i != len(args)-1 {
				sb.WriteString("?,")
				continue
			}
			sb.WriteRune('?')
		}
		sb.WriteByte(')')
		if dm.andor == nil {
			dm.andor = make([]string, 0)
		}
		dm.andor = append(dm.andor, sb.String())
	}
	return dm
}

func (dm *SQLDM) AND_SUBQUERY(column string, inOrNotIn IN, subquery *SQLSM) *SQLDM {
	var sb strings.Builder
	sb.WriteString("AND ")
	sb.WriteString(column)
	sb.WriteRune(' ')
	sb.WriteString(inOrNotIn.value())
	sb.WriteRune('(')
	sb.WriteString(subquery.Sql())
	sb.WriteRune(')')
	dm.args = append(dm.args, subquery.args...)
	dm.where = sb.String()
	return dm
}
func (dm *SQLDM) OR_IN(column string, args ...interface{}) *SQLDM {
	if len(args) > 0 {
		var sb strings.Builder
		sb.WriteString("OR ")
		sb.WriteString(column)
		sb.WriteString(" IN(")
		for i, arg := range args {
			dm.args = append(dm.args, arg)
			if i != len(args)-1 {
				sb.WriteString("?,")
				continue
			}
			sb.WriteRune('?')
		}
		sb.WriteByte(')')
		if dm.andor == nil {
			dm.andor = make([]string, 0)
		}
		dm.andor = append(dm.andor, sb.String())
	}
	return dm
}

func (dm *SQLDM) OR_SUBQUERY(column string, inOrNotIn IN, subquery *SQLSM) *SQLDM {
	var sb strings.Builder
	sb.WriteString("OR ")
	sb.WriteString(column)
	sb.WriteRune(' ')
	sb.WriteString(inOrNotIn.value())
	sb.WriteRune('(')
	sb.WriteString(subquery.Sql())
	sb.WriteRune(')')
	dm.args = append(dm.args, subquery.args...)
	dm.where = sb.String()
	return dm
}
func (dm *SQLDM) SetBatchArgs(batchArgs ...[]interface{}) *SQLDM {
	if dm.batchArgs == nil {
		dm.batchArgs = make([][]interface{}, 0)
	}
	dm.batchArgs = append(dm.batchArgs, batchArgs...)
	return dm
}
func (dm *SQLDM) Execute(mdb *sql.DB) *MySQLDMExecutor {
	return NewMySQLDMExecutor(dm, mdb)
}

func (dm *SQLDM) _sql(sb *strings.Builder) {
	sb.WriteString("DELETE FROM ")
	sb.WriteString(dm.table)
	sb.WriteString(" ")
	sb.WriteString(dm.where)
	sb.WriteString(" ")
	if dm.andor != nil {
		for _, cd := range dm.andor {
			sb.WriteString(cd)
			sb.WriteString(" ")
		}
	}
}
