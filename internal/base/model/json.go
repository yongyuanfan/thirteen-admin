package model

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
)

// JSONSlice 通用 JSON 数组字段。
type JSONSlice[T any] []T

func (s JSONSlice[T]) Value() (driver.Value, error) {
	if s == nil {
		return "[]", nil
	}
	b, err := json.Marshal(s)
	if err != nil {
		return nil, err
	}
	return string(b), nil
}

func (s *JSONSlice[T]) Scan(value any) error {
	if value == nil {
		*s = JSONSlice[T]{}
		return nil
	}
	var raw []byte
	switch v := value.(type) {
	case []byte:
		raw = v
	case string:
		raw = []byte(v)
	default:
		return errors.New("JSONSlice: unsupported scan type")
	}
	if len(raw) == 0 {
		*s = JSONSlice[T]{}
		return nil
	}
	return json.Unmarshal(raw, s)
}
