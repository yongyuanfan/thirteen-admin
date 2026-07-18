package handlers

import (
	"net/url"
	"os"
	"path/filepath"
	"strconv"

	"thirteen-admin/app/services"
	"thirteen-admin/internal/base/handler"

	"github.com/gin-gonic/gin"
)

// CommonHandler 公共上传下载。
type CommonHandler struct {
	handler.Handler
	files *services.FileService
}

// NewCommonHandler 创建。
func NewCommonHandler(files *services.FileService) *CommonHandler {
	return &CommonHandler{files: files}
}

// Upload POST /common/upload
func (h *CommonHandler) Upload(ctx *gin.Context) {
	file, err := ctx.FormFile("file")
	if err != nil {
		h.BadRequest(ctx, "请选择文件")
		return
	}
	f, err := file.Open()
	if err != nil {
		h.InternalError(ctx, err)
		return
	}
	defer f.Close()
	fileName, newFileName, urlPath, err := h.files.SaveUpload(file.Filename, f)
	if err != nil {
		h.InternalError(ctx, err)
		return
	}
	h.OK(ctx, gin.H{
		"fileName":         fileName,
		"newFileName":      newFileName,
		"url":              urlPath,
		"originalFilename": file.Filename,
	})
}

// Download GET /common/download
func (h *CommonHandler) Download(ctx *gin.Context) {
	name := ctx.Query("fileName")
	del := ctx.Query("delete") == "true"
	abs, err := h.files.Resolve(name)
	if err != nil {
		h.BadRequest(ctx, err.Error())
		return
	}
	base := filepath.Base(abs)
	ctx.Header("download-filename", url.QueryEscape(base))
	ctx.FileAttachment(abs, base)
	if del {
		_ = os.Remove(abs)
	}
}

// DownloadResource GET /common/download/resource
func (h *CommonHandler) DownloadResource(ctx *gin.Context) {
	resource := ctx.Query("resource")
	abs, err := h.files.Resolve(resource)
	if err != nil {
		h.BadRequest(ctx, err.Error())
		return
	}
	base := filepath.Base(abs)
	ctx.Header("download-filename", url.QueryEscape(base))
	ctx.FileAttachment(abs, base)
}

// ParseBoolQuery 解析布尔 query。
func ParseBoolQuery(ctx *gin.Context, key string) bool {
	v := ctx.Query(key)
	b, _ := strconv.ParseBool(v)
	return b
}
