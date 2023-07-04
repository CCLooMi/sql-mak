package main

import (
	"database/sql"
	"fmt"
	"reflect"
	"sql-mak/mysql"
	"sql-mak/mysql/entity"
	"sql-mak/mysql/god"
	"sql-mak/utils"
	"testing"
	"time"
)

func TestSelectStruct(t *testing.T) {
	sm := mysql.SELECT("*").
		FROM("users", "u").
		Execute(mysql.MYDB)
	ei := utils.GetEntityInfo(&entity.User{})
	us := make([]*entity.User, 0)
	sm.ExtractorRows(god.RowsExtractor(func(columns *[]string, rs *sql.Rows) {
		u := &entity.User{}
		cL := len(*columns)
		// 创建用于存储结果的切片
		values := make([]interface{}, cL)
		fields := reflect.ValueOf(u).Elem()
		for i := range values {
			fi := ei.CFMap[(*columns)[i]]
			t.Log(fi)
			fv := fields.FieldByName(fi)
			if fv.CanAddr() {
				values[i] = fv.Addr().Interface()
			} else {
				values[i] = nil
			}
		}
		rs.Scan(values...)
		us = append(us, u)
	}))

	t.Log(us)
}
func TestEntityInfo(t *testing.T) {
	ei := utils.GetEntityInfo(&entity.User{})
	t.Log(ei.TableName, ei.Columns, ei.Fields, ei.FCMap, ei.CFMap, ei.Tags)
}

func TestTableName(t *testing.T) {
	u1 := utils.TableName(entity.User{})
	u2 := utils.TableName(&entity.User{})

	if u1 != u2 {
		t.Fail()
	}
}

func TestTableValues(t *testing.T) {
	user := entity.User{
		Username: "johnDoe",
		Password: []byte("password123")}
	user.ID = []byte("1")
	user.InsertedAt = time.Now()
	user.UpdatedAt = time.Now()
	t.Log(utils.GetEntityInfo(user))
	t.Log(utils.GetEntityInfo(&user))
}

func TestInsert(t *testing.T) {
	user := entity.User{
		Username: "johnDoe",
		Password: []byte("password123")}
	user.ID = []byte("1")
	user.InsertedAt = time.Now()
	user.UpdatedAt = time.Now()

	fmt.Print(user.TableName())
	im := mysql.INSERT_INTO(user)

	t.Log(im.Execute(mysql.MYDB).Update().RowsAffected())
}
func TestInsertSQL(t *testing.T) {
	sm := mysql.SELECT("u.id", "u.name").FROM("users", "u")
	sm.LOGSQL(false)
	sql := mysql.
		INSERT_INTO("users").
		INTO_COLUMNS("u.id", "u.name").
		VALUES_SM(sm).
		LOGSQL(false).
		Sql()
	t.Log(sql)
}

func TestUpdateSQL(t *testing.T) {
	um := mysql.UPDATE("users", "u").SET("u.name = ?", "123456").WHERE("u.id = 1")
	um.LOGSQL(false)
	sql := um.Sql()
	t.Log(sql)
}

func TestDeleteSQL(t *testing.T) {
	dm := mysql.DELETE().FROM("users").WHERE("id = 1")
	dm.LOGSQL(false)
	sql := dm.Sql()
	t.Log(sql)
}
