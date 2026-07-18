package services

import (
	"fmt"
	"strconv"
	"time"

	"clash-admin/internal/base/model"
)

const timeLayout = "2006-01-02 15:04:05"

// PageResult 分页列表响应。
type PageResult[T any] struct {
	Current int   `json:"current"`
	Size    int   `json:"size"`
	Total   int64 `json:"total"`
	Records []T   `json:"records"`
}

func formatTime(t time.Time) string {
	if t.IsZero() {
		return ""
	}
	return t.Format(timeLayout)
}

func statusString(status int8) string {
	return strconv.Itoa(int(status))
}

func parseStatus(s string) (int8, error) {
	if s == "" {
		return 0, fmt.Errorf("status 不能为空")
	}
	v, err := strconv.Atoi(s)
	if err != nil {
		return 0, fmt.Errorf("status 无效")
	}
	if v != 1 && v != 2 {
		return 0, fmt.Errorf("status 仅支持 1 或 2")
	}
	return int8(v), nil
}

func parsePage(current, size int) (page, pageSize, offset int) {
	page = current
	if page < 1 {
		page = 1
	}
	pageSize = size
	if pageSize < 1 {
		pageSize = 10
	}
	if pageSize > 200 {
		pageSize = 200
	}
	offset = (page - 1) * pageSize
	return page, pageSize, offset
}

func stringSliceOrEmpty(s model.StringSlice) []string {
	if s == nil {
		return []string{}
	}
	return []string(s)
}
