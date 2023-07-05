package god

import (
	"database/sql"
	"fmt"
	"reflect"
	"sql-mak/utils"
)

type MySQLSMExecutor struct {
	SQLSMExecutor
	MDB *sql.DB
}

func NewMySQLSMExecutor(sm *SQLSM, mdb *sql.DB) *MySQLSMExecutor {
	exe := &MySQLSMExecutor{MDB: mdb}
	exe.SQLSMExecutor = *NewSQLSMExecutor(sm)
	return exe
}

func (exe *MySQLSMExecutor) INSERT_INTO_TABLExtractorResultSet(rse ResultSetExtractor) interface{} {
	stmp, err := exe.MDB.Prepare(exe.God.Sql())
	if err != nil {
		panic(err)
	}
	rows, err := stmp.Query(exe.God.Args()...)
	if err != nil {
		panic(err)
	}
	return rse(rows)
}
func (exe *MySQLSMExecutor) GetResultAsMap() map[string]interface{} {
	stmp, err := exe.MDB.Prepare(exe.God.Sql())
	if err != nil {
		panic(err)
	}
	rows, err := stmp.Query(exe.God.Args()...)
	if err != nil {
		panic(err)
	}
	m, err := exe.RowsToMap(rows)
	if err != nil {
		panic(err)
	}
	return m
}
func (exe *MySQLSMExecutor) GetResultAsMapList() []map[string]interface{} {
	stmp, err := exe.MDB.Prepare(exe.God.Sql())
	if err != nil {
		panic(err)
	}
	rows, err := stmp.Query(exe.God.Args()...)
	if err != nil {
		panic(err)
	}
	list, err := exe.RowsToMaps(rows)
	if err != nil {
		panic(err)
	}
	return list
}
func (exe *MySQLSMExecutor) ExtractorResultSet(rse ResultSetExtractor) interface{} {
	stmp, err := exe.MDB.Prepare(exe.God.Sql())
	if err != nil {
		panic(err)
	}
	rows, err := stmp.Query(exe.God.Args()...)
	if err != nil {
		panic(err)
	}
	return rse(rows)
}

func (e *MySQLSMExecutor) ExtractorResultTo(out interface{}) {
	e.ExtractorResultSet(func(rs *sql.Rows) interface{} {
		t := reflect.TypeOf(out)
		fmt.Println(t)

		//判断out是不是指针,*xxx
		if t.Kind() == reflect.Ptr {
			t = t.Elem() //xxx
			fmt.Println(t)
		}
		//判断out是不是数组//[]xxx
		if t.Kind() == reflect.Slice {
			//获取数组中元素类型
			t = t.Elem()
			fmt.Println(t)

			//创建一个切片，用于存储结果
			slice := reflect.MakeSlice(reflect.SliceOf(t), 0, 0)
			for rs.Next() {
				//创建一个新的元素实例
				elem := reflect.New(t).Elem()
				fmt.Println(elem.Type())
				//将rs的结果扫描到elem中
				rowsTo(rs, elem)
				//将elem添加到切片中
				slice = reflect.Append(slice, elem)
			}
			//将切片赋值给out
			reflect.ValueOf(out).Elem().Set(slice)
		} else {
			//将rs的结果扫描到out中
			rowsTo(rs, reflect.ValueOf(out).Elem())
		}
		return nil
	})
}

func rowsTo(rs *sql.Rows, out reflect.Value) {
	ei := utils.GetEntityInfo(out.Type())
	columns, _ := rs.Columns()
	cL := len(columns)
	//判断out是不是指针
	if out.Kind() == reflect.Ptr {
		out = out.Elem()
		fmt.Println(out.Type())
	}
	// 创建用于存储结果的切片
	values := make([]interface{}, cL)
	for i := range values {
		fi := ei.CFMap[columns[i]]
		fv := out.FieldByName(fi)
		if fv.CanAddr() {
			values[i] = fv.Addr().Interface()
		} else {
			values[i] = nil
		}
	}
	rs.Scan(values...)
}

func (exe *MySQLSMExecutor) Count() int64 {
	stmp, err := exe.MDB.Prepare(exe.God.Sql())
	if err != nil {
		panic(err)
	}
	rows, err := stmp.Query(exe.God.Args()...)
	if err != nil {
		panic(err)
	}
	defer rows.Close() // finally close rows
	var count int64
	for rows.Next() {
		err := rows.Scan(&count)
		if err != nil {
			panic(err)
		}
		return count //lint:ignore SA4004 just break
	}
	return 0
}
