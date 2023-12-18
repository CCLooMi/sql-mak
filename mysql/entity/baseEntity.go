package entity

import (
	"encoding/hex"
	"encoding/json"
	"time"
)

type DateTime time.Time

func (dt *DateTime) MarshalJSON() ([]byte, error) {
	//判断dt是否为空
	if time.Time(*dt).IsZero() {
		return []byte("null"), nil
	}
	var stamp = `"` + time.Time(*dt).Format("2006-01-02 15:04:05") + `"`
	return []byte(stamp), nil
}

func (dt *DateTime) UnmarshalJSON(data []byte) error {
	var str string
	err := json.Unmarshal(data, &str)
	if err == nil {
		t, err := time.Parse("2006-01-02 15:04:05", str)
		if err != nil {
			return err
		}
		*dt = DateTime(t)
	} else {
		var timestamp int64
		err = json.Unmarshal(data, &timestamp)
		if err != nil {
			return err
		}
		*dt = DateTime(time.Unix(timestamp, 0))
	}

	return nil
}

type Bytes []byte

func (id *Bytes) MarshalJSON() ([]byte, error) {
	return json.Marshal(hex.EncodeToString(*id))
}
func (id *Bytes) UnmarshalJSON(data []byte) error {
	var str string
	err := json.Unmarshal(data, &str)
	if err == nil {
		decoded, err := hex.DecodeString(str)
		if err != nil {
			*id = data
			return nil
		}
		*id = decoded
		return nil
	}
	var decoded []byte
	err = json.Unmarshal(data, &decoded)
	if err != nil {
		*id = data
		return nil
	}
	*id = decoded
	return nil
}

type BaseEntity interface {
	TableName() string
}
type IdEntity struct {
	BaseEntity
	Id string `orm:"varchar(32) not null comment '主键ID'" column:"id" primaryKey:"true"`
}
type LongIdEntity struct {
	BaseEntity
	Id int64 `orm:"bigint auto_increment not null comment '主键ID'" column:"id" primaryKey:"true"`
}
type BidEntity struct {
	BaseEntity
	Id []byte `orm:"binary(16) not null comment '主键ID'" column:"id" primaryKey:"true"`
}
type TimeEntity struct {
	InsertedAt time.Time `orm:"datetime not null comment '插入时间'" column:"inserted_at"`
	UpdatedAt  time.Time `orm:"datetime not null comment '更新时间'" column:"updated_at"`
}
type StrTimeEntity struct {
	InsertedAt string `orm:"datetime not null comment '插入时间'" column:"inserted_at"`
	UpdatedAt  string `orm:"datetime not null comment '更新时间'" column:"updated_at"`
}
type TimestampEntity struct {
	InsertedAt int64 `orm:"datetime not null comment '插入时间'" column:"inserted_at"`
	UpdatedAt  int64 `orm:"datetime not null comment '更新时间'" column:"updated_at"`
}
type AuditEntity struct {
	CreatedBy string `orm:"varchar(32) not null comment '创建人'" column:"created_by"`
	UpdatedBy string `orm:"varchar(32) not null comment '更新人'" column:"updated_by"`
}
type BAuditEntity struct {
	CreatedBy []byte `orm:"binary(16) not null comment '创建人'" column:"created_by"`
	UpdatedBy []byte `orm:"binary(16) not null comment '更新人'" column:"updated_by"`
}
type LongAuditEntity struct {
	CreatedBy int64 `orm:"bigint not null comment '创建人'" column:"created_by"`
	UpdatedBy int64 `orm:"bigint not null comment '更新人'" column:"updated_by"`
}
