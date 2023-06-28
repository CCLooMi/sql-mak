package god

import (
	"fmt"
	"reflect"
	"strings"
)

type SQLIM struct {
	AbstractSQLGod
	AbstractSQLGodChild
	table    string
	columns  []string
	valuesSM *SQLSM
	sets     []string
	setArgs  []interface{}
}

func NewSQLIM() *SQLIM {
	return &SQLIM{
		columns: make([]string, 0),
	}
}
func (im *SQLIM) INSERT_INTO(table interface{}, columns ...string) *SQLIM {
	switch t := table.(type) {
	case string:
		im.table = t
	case reflect.Type:
		im.table = im.TableName(t)
	default:
		im.table = fmt.Sprintf("%s", table)
	}
	im.columns = append(im.columns, columns...)
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
	im.hasSubSelect = true
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
func (sm *SQLIM) Execute() *SQLIMExecutor {
	god := sm.toSQLGod()
	executor, _ := GetExecutor(god, reflect.TypeOf(SQLIMExecutor{})).(SQLIMExecutor)
	return &executor
}

func (im *SQLIM) _sql(sb *strings.Builder) {
	sb.WriteString("INSERT INTO ")
	sb.WriteString(im.table)
	sb.WriteString(" ")
	if len(im.columns) > 0 {
		sb.WriteString("(")
		for i, column := range im.columns {
			if column[0] != '`' {
				sb.WriteString("`" + column + "`,")
				continue
			}
			sb.WriteString(column)
			if i != len(im.columns)-1 {
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
		l := len(ags) - 1
		if im.setArgs != nil {
			l--
		}
		L := len(im.columns)
		for i, idx := 0, 0; i < L; i++ {
			if _, ok := ags[i].(EXP); ok {
				ri := i - idx
				args := im.args
				args = append(args[:ri], args[ri+1:]...)
				e := ags[i].(EXP)
				sb.WriteString(e.Exp())
				for _, arg := range e.Args() {
					args = append(args[:ri], append([]interface{}{arg}, args[ri:]...)...)
					idx--
					ri++
				}
				im.args = args
			} else {
				if i != L-1 {
					sb.WriteString("?,")
				} else {
					sb.WriteRune('?')
				}
			}
		}
		if L <= 0 {
			sb.WriteRune(')')
		}
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
