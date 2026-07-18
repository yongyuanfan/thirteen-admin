package system

import (
	"errors"
	"strconv"

	"thirteen-admin/app/models"
	"thirteen-admin/app/services"
	"thirteen-admin/internal/base/handler"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// PostHandler 岗位 HTTP 处理器。
type PostHandler struct {
	handler.Handler
	svc *services.PostService
}

// NewPostHandler 创建岗位 Handler。
func NewPostHandler(svc *services.PostService) *PostHandler {
	return &PostHandler{svc: svc}
}

// List GET /system/post/list
func (h *PostHandler) List(ctx *gin.Context) {
	total, rows, err := h.svc.List(ctx)
	if err != nil {
		h.InternalError(ctx, err)
		return
	}
	h.Table(ctx, total, rows)
}

// Get GET /system/post/:postId
func (h *PostHandler) Get(ctx *gin.Context) {
	id, _ := strconv.ParseUint(ctx.Param("postId"), 10, 64)
	data, err := h.svc.Get(uint(id))
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			h.NotFound(ctx, "岗位不存在")
			return
		}
		h.InternalError(ctx, err)
		return
	}
	h.Success(ctx, data)
}

type postBody struct {
	PostID   uint   `json:"postId"`
	PostCode string `json:"postCode"`
	PostName string `json:"postName"`
	PostSort int    `json:"postSort"`
	Status   string `json:"status"`
	Remark   string `json:"remark"`
}

// Add POST /system/post
func (h *PostHandler) Add(ctx *gin.Context) {
	var req postBody
	if err := ctx.ShouldBindJSON(&req); err != nil {
		h.BadRequest(ctx, "参数错误")
		return
	}
	post := &models.Post{
		PostCode: req.PostCode,
		PostName: req.PostName,
		PostSort: req.PostSort,
		Status:   req.Status,
		Remark:   req.Remark,
	}
	if post.Status == "" {
		post.Status = "0"
	}
	if err := h.svc.Add(post, services.ContextUserName(ctx)); err != nil {
		h.InternalError(ctx, err)
		return
	}
	h.SuccessMsg(ctx, "新增成功")
}

// Edit PUT /system/post
func (h *PostHandler) Edit(ctx *gin.Context) {
	var req postBody
	if err := ctx.ShouldBindJSON(&req); err != nil {
		h.BadRequest(ctx, "参数错误")
		return
	}
	post := &models.Post{
		PostCode: req.PostCode,
		PostName: req.PostName,
		PostSort: req.PostSort,
		Status:   req.Status,
		Remark:   req.Remark,
	}
	post.ID = req.PostID
	if err := h.svc.Edit(post, services.ContextUserName(ctx)); err != nil {
		h.InternalError(ctx, err)
		return
	}
	h.SuccessMsg(ctx, "修改成功")
}

// Delete DELETE /system/post/:postId
func (h *PostHandler) Delete(ctx *gin.Context) {
	ids := services.ParseUintIDs(ctx.Param("postId"))
	if len(ids) == 0 {
		h.BadRequest(ctx, "参数错误")
		return
	}
	if err := h.svc.Delete(ids); err != nil {
		h.Fail(ctx, handler.CodeError, err.Error())
		return
	}
	h.SuccessMsg(ctx, "删除成功")
}

// Export POST /system/post/export
func (h *PostHandler) Export(ctx *gin.Context) {
	rows, err := h.svc.Export(ctx)
	if err != nil {
		h.InternalError(ctx, err)
		return
	}
	services.WriteCSV(ctx, "post.csv", []string{"岗位编号", "岗位编码", "岗位名称", "岗位排序", "状态", "创建时间"}, rows)
}
