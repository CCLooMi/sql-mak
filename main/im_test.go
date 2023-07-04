package main

import (
	"fmt"
	"sql-mak/mysql"
	"sql-mak/mysql/entity"
	"sql-mak/utils"
	"testing"
	"time"
)

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
