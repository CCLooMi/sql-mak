package test

import (
	"encoding/hex"
	"encoding/json"
	"github.com/CCLooMi/sql-mak/mysql/entity"
	"testing"
	"time"
)

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

type Obj struct {
	ID       *ID       `json:"id,omitempty"`
	UpdateAt *DateTime `json:"update_at,omitempty"`
}

func TestJSON(ts *testing.T) {
	now := DateTime(time.Now())
	data := Obj{
		ID:       &ID{1, 2, 3, 4, 5},
		UpdateAt: &now,
	}

	jsonData, err := json.Marshal(data)
	if err != nil {
		ts.Log("Failed to marshal struct to JSON:", err)
		return
	}

	ts.Log("Serialized JSON:", string(jsonData))

	// 示例2：将JSON反序列化为结构体实例
	jsonStr := `{"id": [6, 7, 8, 9, 10]}`

	var result Obj
	err = json.Unmarshal([]byte(jsonStr), &result)
	if err != nil {
		ts.Log("Failed to unmarshal JSON to struct:", err)
		return
	}

	ts.Log("Deserialized Struct:", toJSON(result))

	jsonStr = `{"id": "0102030405"}`
	err = json.Unmarshal([]byte(jsonStr), &result)
	if err != nil {
		ts.Log("Failed to unmarshal JSON to struct:", err)
		return
	}
	ts.Log("Deserialized Struct:", toJSON(result))
}

func toJSON(v any) string {
	b, err := json.MarshalIndent(v, "", "  ")
	if err != nil {
		panic(err)
	}
	return string(b)
}

func TestBaseEntity(ts *testing.T) {
	type U struct {
		entity.BidEntity
		entity.TimeEntity
		entity.BAuditEntity
	}
	t := time.Now()
	u := U{
		BidEntity: entity.BidEntity{
			Id: &[]byte{1, 2, 3, 4, 5},
		},
		TimeEntity: entity.TimeEntity{
			InsertedAt: &t,
			UpdatedAt:  &t,
		},
		BAuditEntity: entity.BAuditEntity{
			CreatedBy: &[]byte{1, 2, 3, 4, 5},
			UpdatedBy: &[]byte{1, 2, 3, 4, 5},
		},
	}
	r, err := json.Marshal(u)
	if err != nil {
		ts.Log("Failed to marshal struct to JSON:", err)
		return
	}
	ts.Log("Serialized JSON:", string(r))
}
