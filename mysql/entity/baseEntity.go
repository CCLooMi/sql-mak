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

type ID []byte

func (id *ID) MarshalJSON() ([]byte, error) {
	return json.Marshal(hex.EncodeToString(*id))
}

func (id *ID) UnmarshalJSON(data []byte) error {
	var str string
	err := json.Unmarshal(data, &str)
	if err == nil {
		decoded, err := hex.DecodeString(str)
		if err != nil {
			*id = ID(data)
			return nil
		}
		*id = ID(decoded)
		return nil
	}
	var decoded []byte
	err = json.Unmarshal(data, &decoded)
	if err != nil {
		*id = ID(data)
		return nil
	}
	*id = ID(decoded)
	return nil
}

type BaseEntity interface {
	TableName() string
}

type IdEntity struct {
	BaseEntity
	Id *ID `orm:"type:binary(16); primaryKey; not null; comment:'主键ID'" column:"id"`
}

type TimeEntity struct {
	InsertedAt *DateTime `orm:"not null; comment:'插入时间'" column:"inserted_at"`
	UpdatedAt  *DateTime `orm:"not null; comment:'更新时间'" column:"updated_at"`
}
