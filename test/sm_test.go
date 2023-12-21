package test

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"testing"
	"time"

	"github.com/CCLooMi/sql-mak/mysql"
	"github.com/CCLooMi/sql-mak/mysql/entity"
	"github.com/CCLooMi/sql-mak/utils"

	_ "github.com/go-sql-driver/mysql"
)

var MYDB *sql.DB

func init() {
	_db, err := sql.Open("mysql", "root:apple@tcp(127.0.0.1:3306)/wios?charset=utf8mb4&parseTime=True&loc=Local")
	if err != nil {
		panic(err)
	}
	if err = _db.Ping(); err != nil {
		_db, _ := sql.Open("mysql", "root:123456@tcp(127.0.0.1:3306)/test?charset=utf8mb4&parseTime=True&loc=Local")
		if err = _db.Ping(); err != nil {
			fmt.Errorf("Failed to connect to database: %v", err)
			return
		}
		MYDB = _db
	} else {
		MYDB = _db
	}

	MYDB.SetMaxOpenConns(10)
	MYDB.SetMaxIdleConns(5)
}

func TestSelectExtract(t *testing.T) {
	us := &[]*User{}
	sm := mysql.SELECT("*").FROM("users", "u")
	sm.Execute(MYDB).ExtractorResultTo(us)
	fmt.Println(toJSONString(us))

	u := &User{}
	sm.LIMIT(1).Execute(MYDB).ExtractorResultTo(u)
	fmt.Println(toJSONString(u))
}
func TestInsert(t *testing.T) {
	u := &User{Username: "Joy", Password: []byte("123456")}
	id := "123456"
	time := time.Now()
	u.Id = &id
	u.UpdatedAt = &time
	u.InsertedAt = &time

	im := mysql.INSERT_INTO(u).
		ON_DUPLICATE_KEY_UPDATE().SET("username=?", "JoyNew")
	fmt.Println("sql", im.Sql())
	c, err := im.Execute(MYDB).Update().RowsAffected()
	fmt.Println(c, err)
}

func TestEntityInfo(t *testing.T) {
	ei := utils.GetEntityInfo(User{})
	if ei.PrimaryKey != "id" {
		t.Errorf("Expected primary key to be 'id', but got '%s'", ei.PrimaryKey)
	}
}

func toJSONString(v interface{}) string {
	b, _ := json.Marshal(v)
	return string(b)
}

type User struct {
	entity.IdEntity
	entity.TimeEntity
	Username string `orm:"varchar(255) comment '用户名'" column:"username"`
	Password []byte `orm:"varbinary(32) comment '用户密码'" column:"password"`
}

func (*User) TableName() string {
	return "users"
}
