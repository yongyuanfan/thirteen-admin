package services

import "testing"

func TestToSnake(t *testing.T) {
	tests := []struct {
		in, want string
	}{
		{"createdAt", "created_at"},
		{"updatedAt", "updated_at"},
		{"createdBy", "created_by"},
		{"updatedBy", "updated_by"},
		{"tableName", "table_name"},
		{"operTime", "oper_time"},
	}
	for _, tt := range tests {
		if got := toSnake(tt.in); got != tt.want {
			t.Errorf("toSnake(%q) = %q, want %q", tt.in, got, tt.want)
		}
	}
}
