package test

import (
	"encoding/json"
	"fmt"
	"sql-mak/mysql"
	"sql-mak/mysql/entity"
	"testing"
)

func TestSelectExtract(t *testing.T) {
	us := &[]*entity.User{}
	sm := mysql.SELECT("*").FROM("users", "u")
	sm.Execute(mysql.MYDB).
		ExtractorResultTo(us)
	fmt.Println(toJSONString(us))
}

func toJSONString(v interface{}) string {
	b, _ := json.Marshal(v)
	return string(b)
}
