package model

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
	"time"
)

// BaseModel 通用模型字段。
type BaseModel struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}

// StringSlice 用于 GORM JSON 数组字段。
type StringSlice []string

func (s StringSlice) Value() (driver.Value, error) {
	if s == nil {
		return "[]", nil
	}
	b, err := json.Marshal(s)
	if err != nil {
		return nil, err
	}
	return string(b), nil
}

func (s *StringSlice) Scan(value any) error {
	if value == nil {
		*s = StringSlice{}
		return nil
	}
	var raw []byte
	switch v := value.(type) {
	case []byte:
		raw = v
	case string:
		raw = []byte(v)
	default:
		return errors.New("StringSlice: unsupported scan type")
	}
	if len(raw) == 0 {
		*s = StringSlice{}
		return nil
	}
	return json.Unmarshal(raw, s)
}
