package main

import (
	"fmt"
	"reflect"
	"sql-mak/mysql"
	"sql-mak/mysql/entity"
	"testing"
	"time"
)

func TestTableColumns(t *testing.T) {
	user := entity.User{
		Username: "johnDoe",
		Password: []byte("password123")}
	user.ID = []byte("1")
	user.InsertedAt = time.Now()
	user.UpdatedAt = time.Now()

	t.Log(getFields(user))
}

func getFields(o interface{}) []string {
	t := reflect.TypeOf(o)
	s := make([]string, 0)
	for i := 0; i < t.NumField(); i++ {
		fd := t.Field(i)
		if fd.Type.Kind() == reflect.Struct {
			if fd.Type.Kind() == reflect.Ptr {
				s = append(s, getFields(fd.Type.Elem())...)
			} else {
				s = append(s, getFields(fd.Type)...)
			}
			continue
		}
		if fd.Type.Kind() == reflect.Interface {
			continue
		}
		s = append(s, fd.Name)
	}
	return s
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
