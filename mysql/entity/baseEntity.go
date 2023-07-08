package entity

import (
	"time"
)

// type JsonTime time.Time

// func (jt JsonTime) MarshalJSON() ([]byte, error) {
// 	var stamp = fmt.Sprintf("\"%s\"", time.Time(jt).Format("2006-01-02 15:04:05"))
// 	return []byte(stamp), nil
// }

type BaseEntity interface {
	TableName() string
}

type IdEntity struct {
	BaseEntity
	ID []byte `orm:"type:binary(16); primaryKey; not null; comment:'主键ID'" column:"id"`
}

type TimeEntity struct {
	InsertedAt time.Time `orm:"not null; comment:'插入时间'" column:"inserted_at"`
	UpdatedAt  time.Time `orm:"not null; comment:'更新时间'" column:"updated_at"`
}
