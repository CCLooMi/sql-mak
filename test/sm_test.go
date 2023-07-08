package test

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"sql-mak/mysql"
	"sql-mak/mysql/entity"
	"testing"

	_ "github.com/go-sql-driver/mysql"
)

var MYDB *sql.DB

func init() {
	_db, err := sql.Open("mysql", "root:apple@tcp(127.0.0.1:3308)/wios?charset=utf8mb4&parseTime=True&loc=Local")
	if err != nil {
		panic(err)
	}
	if err = _db.Ping(); err != nil {
		_db, _ := sql.Open("mysql", "root:123456@tcp(127.0.0.1:3306)/test?charset=utf8mb4&parseTime=True&loc=Local")
		if err = _db.Ping(); err != nil {
			panic(err)
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
	sm.Execute(MYDB).
		ExtractorResultTo(us)
	fmt.Println(toJSONString(us))
}

func toJSONString(v interface{}) string {
	b, _ := json.Marshal(v)
	return string(b)
}

type User struct {
	entity.IdEntity
	entity.TimeEntity
	Username string `orm:"type:varchar(255); comment:'用户名'" column:"username"`
	Password []byte `orm:"type:varbinary(32); comment:'用户密码'" column:"password"`
}

func (*User) TableName() string {
	return "users"
}
