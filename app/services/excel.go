package services

import (
	"bytes"
	"encoding/csv"
	"fmt"
	"net/url"

	"github.com/gin-gonic/gin"
)

// WriteCSV 写出 CSV 下载（兼容前端 blob 下载）。
func WriteCSV(ctx *gin.Context, filename string, headers []string, rows [][]string) {
	var buf bytes.Buffer
	// UTF-8 BOM，便于 Excel 打开
	buf.Write([]byte{0xEF, 0xBB, 0xBF})
	w := csv.NewWriter(&buf)
	_ = w.Write(headers)
	for _, row := range rows {
		_ = w.Write(row)
	}
	w.Flush()

	name := filename
	if name == "" {
		name = "export.csv"
	}
	ctx.Header("Content-Type", "application/vnd.ms-excel;charset=utf-8")
	ctx.Header("Content-Disposition", fmt.Sprintf("attachment; filename*=UTF-8''%s", url.QueryEscape(name)))
	ctx.Header("download-filename", url.QueryEscape(name))
	ctx.Data(200, "application/vnd.ms-excel;charset=utf-8", buf.Bytes())
}
