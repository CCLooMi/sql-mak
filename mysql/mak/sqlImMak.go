package mak

import (
	"database/sql"
	"strings"

	"github.com/CCLooMi/sql-mak/utils"
)

type SQLIM struct {
	AbstractSQLMak
	AbstractSQLMakChild
	table    string
	columns  []string
	valuesSM *SQLSM
	sets     []string
	setArgs  []interface{}
}

func NewSQLIM() *SQLIM {
	im := &SQLIM{
		columns: make([]string, 0),
	}
	a := im.toAbstractSQLMakChild()
	im.AbstractSQLMak = *NewAbstractSQLMak(&a)
	return im
}
func (im *SQLIM) toAbstractSQLMakChild() AbstractSQLMakChild {
	return im
}
func (im *SQLIM) INSERT_INTO(table interface{}, columns ...string) *SQLIM {
	switch t := table.(type) {
	case string:
		im.table = t
	default:
		im.table = im.TableName(table)
	}
	ei := utils.GetEntityInfo(table)
	if len(columns) == 0 {
		im.columns = append(im.columns, ei.Columns...)
		fvs := utils.GetFieldValues(table, ei.Fields)
		for i, vi := range fvs {
			exp := ei.IExpMap[ei.Columns[i]]
			if exp != "" {
				fvs[i] = ExpStr(exp, vi)
			}
		}
		im.args = append(im.args, fvs...)
	} else {
		im.columns = append(im.columns, columns...)
	}
	return im
}

func (im *SQLIM) INTO_COLUMNS(columns ...string) *SQLIM {
	im.columns = append(im.columns, columns...)
	return im
}

func (im *SQLIM) VALUES(values ...interface{}) *SQLIM {
	im.args = append(im.args, values...)
	return im
}

func (im *SQLIM) VALUES_SM(sm *SQLSM) *SQLIM {
	im.hasSubArgs = true
	im.valuesSM = sm
	if len(im.columns) == 0 {
		im.columns = append(im.columns, sm.ColumnAlias()...)
	}
	im.args = append(im.args, sm.args...)
	return im
}

func (im *SQLIM) ON_DUPLICATE_KEY_UPDATE() *SQLIM {
	if im.sets == nil {
		im.sets = make([]string, 0)
	}
	return im
}

func (im *SQLIM) SET(set string, args ...interface{}) *SQLIM {
	im.sets = append(im.sets, set)
	if len(args) > 0 {
		if im.setArgs == nil {
			im.setArgs = make([]interface{}, 0)
			im.args = append(im.args, &im.setArgs)
			im.hasSubArgs = true
		}
		im.setArgs = append(im.setArgs, args...)
	}
	return im
}

func (im *SQLIM) BatchArgs(batchArgs ...[]interface{}) *SQLIM {
	if im.batchArgs == nil {
		im.batchArgs = make([][]interface{}, 0)
	}
	im.batchArgs = append(im.batchArgs, batchArgs...)
	return im
}
func (im *SQLIM) Execute(mdb *sql.DB) *MySQLIMExecutor {
	return NewMySQLIMExecutor(im, mdb)
}

func (im *SQLIM) _sql(sb *strings.Builder) {
	sb.WriteString("INSERT INTO ")
	sb.WriteString(im.table)
	sb.WriteString(" ")
	cL := len(im.columns)
	if cL > 0 {
		sb.WriteString("(")
		for i, column := range im.columns {
			if column[0] != '`' {
				sb.WriteString("`" + column + "`")
				if i != cL-1 {
					sb.WriteRune(',')
				}
				continue
			}
			sb.WriteString(column)
			if i != cL-1 {
				sb.WriteRune(',')
			}
		}
		sb.WriteString(")")
	}
	if im.valuesSM != nil {
		sb.WriteString("(")
		im.valuesSM._sql(sb)
		sb.WriteString(")")
	} else if len(im.args) > 0 {
		sb.WriteString("VALUES (")
		ags := im.args
		l := len(ags)
		if im.setArgs != nil {
			l -= len(im.setArgs)
		}
		for i, ag := range ags {
			if i < l {
				if e, ok := ag.(*EXP); ok {
					sb.WriteString(e.Exp())
					im.args = append(im.args[:i], append(e.Args(), im.args[i+1:]...)...)
				} else {
					sb.WriteRune('?')
				}
				if i != len(ags)-1 {
					sb.WriteString(",")
				}
			}
		}
		sb.WriteRune(')')
	} else if im.batchArgs != nil {
		sb.WriteString("VALUES (")
		l := len(im.batchArgs[0])
		L := len(im.columns)
		for i := 0; i < l && i < L; i++ {
			if i != l-1 && i != L-1 {
				sb.WriteString("?,")
				continue
			}
			sb.WriteRune('?')
		}
		sb.WriteRune(')')
	}
	if im.sets != nil && len(im.sets) > 0 {
		sb.WriteString("ON DUPLICATE KEY UPDATE ")
		for i, set := range im.sets {
			sb.WriteString(set)
			if i != len(im.sets)-1 {
				sb.WriteRune(',')
			}
		}
	}
}
