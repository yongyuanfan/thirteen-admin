package handlers

import (
	"net/url"
	"os"
	"path/filepath"
	"strconv"

	"thirteen-admin/app/runtime"
	"thirteen-admin/internal/base/handler"

	"github.com/gin-gonic/gin"
)

// CommonHandler 公共上传下载。
type CommonHandler struct {
	handler.Api
}

// Upload POST /common/upload
func (e CommonHandler) Upload(c *gin.Context) {
	files := runtime.File()
	file, err := c.FormFile("file")
	if err != nil {
		e.BadRequest(c, "请选择文件")
		return
	}
	f, err := file.Open()
	if err != nil {
		e.InternalError(c, err)
		return
	}
	defer f.Close()
	fileName, newFileName, urlPath, err := files.SaveUpload(file.Filename, f)
	if err != nil {
		e.InternalError(c, err)
		return
	}
	e.OK(c, gin.H{
		"fileName":         fileName,
		"newFileName":      newFileName,
		"url":              urlPath,
		"originalFilename": file.Filename,
	})
}

// Download GET /common/download
func (e CommonHandler) Download(c *gin.Context) {
	files := runtime.File()
	name := c.Query("fileName")
	del := c.Query("delete") == "true"
	abs, err := files.Resolve(name)
	if err != nil {
		e.BadRequest(c, err.Error())
		return
	}
	base := filepath.Base(abs)
	c.Header("download-filename", url.QueryEscape(base))
	c.FileAttachment(abs, base)
	if del {
		_ = os.Remove(abs)
	}
}

// DownloadResource GET /common/download/resource
func (e CommonHandler) DownloadResource(c *gin.Context) {
	files := runtime.File()
	resource := c.Query("resource")
	abs, err := files.Resolve(resource)
	if err != nil {
		e.BadRequest(c, err.Error())
		return
	}
	base := filepath.Base(abs)
	c.Header("download-filename", url.QueryEscape(base))
	c.FileAttachment(abs, base)
}

// ParseBoolQuery 解析布尔 query。
func ParseBoolQuery(c *gin.Context, key string) bool {
	v := c.Query(key)
	b, _ := strconv.ParseBool(v)
	return b
}
