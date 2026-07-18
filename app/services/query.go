package services

import (
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// PageQuery 分页参数。
type PageQuery struct {
	PageNum  int
	PageSize int
}

// ParsePage 从 query 解析分页。
func ParsePage(ctx *gin.Context) PageQuery {
	pageNum, _ := strconv.Atoi(ctx.DefaultQuery("pageNum", "1"))
	pageSize, _ := strconv.Atoi(ctx.DefaultQuery("pageSize", "10"))
	if pageNum < 1 {
		pageNum = 1
	}
	if pageSize < 1 {
		pageSize = 10
	}
	if pageSize > 500 {
		pageSize = 500
	}
	return PageQuery{PageNum: pageNum, PageSize: pageSize}
}

// ApplyPage 应用分页。
func ApplyPage(db *gorm.DB, page PageQuery) *gorm.DB {
	offset := (page.PageNum - 1) * page.PageSize
	return db.Offset(offset).Limit(page.PageSize)
}

// ApplyOrder 应用排序（orderByColumn 为属性名或列名）。
func ApplyOrder(db *gorm.DB, ctx *gin.Context, defaultOrder string) *gorm.DB {
	col := ctx.Query("orderByColumn")
	asc := strings.ToLower(ctx.Query("isAsc"))
	if col == "" {
		return db.Order(defaultOrder)
	}
	col = toSnake(col)
	if asc == "asc" {
		return db.Order(col + " asc")
	}
	return db.Order(col + " desc")
}

// ApplyTimeRange 应用 params[beginTime]/params[endTime]。
func ApplyTimeRange(db *gorm.DB, ctx *gin.Context, column string) *gorm.DB {
	begin := ctx.Query("params[beginTime]")
	end := ctx.Query("params[endTime]")
	if begin != "" {
		db = db.Where(column+" >= ?", begin)
	}
	if end != "" {
		db = db.Where(column+" <= ?", end+" 23:59:59")
	}
	return db
}

func toSnake(s string) string {
	var b strings.Builder
	for i, r := range s {
		if r >= 'A' && r <= 'Z' {
			if i > 0 {
				b.WriteByte('_')
			}
			b.WriteRune(r + ('a' - 'A'))
			continue
		}
		b.WriteRune(r)
	}
	return b.String()
}

// FormatTime 格式化时间为前端字符串。
func FormatTime(t *time.Time) string {
	if t == nil || t.IsZero() {
		return ""
	}
	return t.Format("2006-01-02 15:04:05")
}

// FormatTimeVal 格式化 time.Time。
func FormatTimeVal(t time.Time) string {
	if t.IsZero() {
		return ""
	}
	return t.Format("2006-01-02 15:04:05")
}

// ParseUintIDs 解析 "1,2,3" 为 uint 切片。
func ParseUintIDs(s string) []uint {
	parts := strings.Split(s, ",")
	ids := make([]uint, 0, len(parts))
	for _, p := range parts {
		p = strings.TrimSpace(p)
		if p == "" {
			continue
		}
		n, err := strconv.ParseUint(p, 10, 64)
		if err != nil {
			continue
		}
		ids = append(ids, uint(n))
	}
	return ids
}

// ContextUserID 从 gin context 取用户 ID。
func ContextUserID(ctx *gin.Context) uint {
	v, ok := ctx.Get("authUserID")
	if !ok {
		return 0
	}
	switch t := v.(type) {
	case string:
		n, _ := strconv.ParseUint(t, 10, 64)
		return uint(n)
	case uint:
		return t
	default:
		return 0
	}
}

// ContextUserName 从 gin context 取用户名。
func ContextUserName(ctx *gin.Context) string {
	v, ok := ctx.Get("authUserName")
	if !ok {
		return ""
	}
	s, _ := v.(string)
	return s
}
