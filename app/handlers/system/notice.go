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

// NoticeHandler 公告 HTTP 处理器。
type NoticeHandler struct {
	handler.Handler
	svc *services.NoticeService
}

// NewNoticeHandler 创建公告 Handler。
func NewNoticeHandler(svc *services.NoticeService) *NoticeHandler {
	return &NoticeHandler{svc: svc}
}

// List GET /system/notice/list
func (h *NoticeHandler) List(ctx *gin.Context) {
	total, rows, err := h.svc.List(ctx)
	if err != nil {
		h.InternalError(ctx, err)
		return
	}
	h.Table(ctx, total, rows)
}

// Get GET /system/notice/:noticeId
func (h *NoticeHandler) Get(ctx *gin.Context) {
	id, _ := strconv.ParseUint(ctx.Param("noticeId"), 10, 64)
	data, err := h.svc.Get(uint(id))
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			h.NotFound(ctx, "公告不存在")
			return
		}
		h.InternalError(ctx, err)
		return
	}
	h.Success(ctx, data)
}

type noticeBody struct {
	NoticeID      uint   `json:"noticeId"`
	NoticeTitle   string `json:"noticeTitle"`
	NoticeType    string `json:"noticeType"`
	NoticeContent string `json:"noticeContent"`
	Status        string `json:"status"`
	Remark        string `json:"remark"`
}

// Add POST /system/notice
func (h *NoticeHandler) Add(ctx *gin.Context) {
	var req noticeBody
	if err := ctx.ShouldBindJSON(&req); err != nil {
		h.BadRequest(ctx, "参数错误")
		return
	}
	notice := &models.Notice{
		NoticeTitle:   req.NoticeTitle,
		NoticeType:    req.NoticeType,
		NoticeContent: []byte(req.NoticeContent),
		Status:        req.Status,
		Remark:        req.Remark,
	}
	if notice.Status == "" {
		notice.Status = "0"
	}
	if err := h.svc.Add(notice, services.ContextUserName(ctx)); err != nil {
		h.InternalError(ctx, err)
		return
	}
	h.SuccessMsg(ctx, "新增成功")
}

// Edit PUT /system/notice
func (h *NoticeHandler) Edit(ctx *gin.Context) {
	var req noticeBody
	if err := ctx.ShouldBindJSON(&req); err != nil {
		h.BadRequest(ctx, "参数错误")
		return
	}
	notice := &models.Notice{
		NoticeTitle:   req.NoticeTitle,
		NoticeType:    req.NoticeType,
		NoticeContent: []byte(req.NoticeContent),
		Status:        req.Status,
		Remark:        req.Remark,
	}
	notice.ID = req.NoticeID
	if err := h.svc.Edit(notice, services.ContextUserName(ctx)); err != nil {
		h.InternalError(ctx, err)
		return
	}
	h.SuccessMsg(ctx, "修改成功")
}

// Delete DELETE /system/notice/:noticeId
func (h *NoticeHandler) Delete(ctx *gin.Context) {
	ids := services.ParseUintIDs(ctx.Param("noticeId"))
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

// ListTop GET /system/notice/listTop
func (h *NoticeHandler) ListTop(ctx *gin.Context) {
	rows, unread, err := h.svc.ListTop(services.ContextUserID(ctx))
	if err != nil {
		h.InternalError(ctx, err)
		return
	}
	h.OK(ctx, gin.H{
		"data":        rows,
		"unreadCount": unread,
	})
}

// MarkRead POST /system/notice/markRead
func (h *NoticeHandler) MarkRead(ctx *gin.Context) {
	noticeID, _ := strconv.ParseUint(ctx.Query("noticeId"), 10, 64)
	if noticeID == 0 {
		h.BadRequest(ctx, "参数错误")
		return
	}
	if err := h.svc.MarkRead(services.ContextUserID(ctx), uint(noticeID)); err != nil {
		h.InternalError(ctx, err)
		return
	}
	h.SuccessMsg(ctx, "操作成功")
}

// MarkReadAll POST /system/notice/markReadAll
func (h *NoticeHandler) MarkReadAll(ctx *gin.Context) {
	ids := services.ParseUintIDs(ctx.Query("ids"))
	if len(ids) == 0 {
		h.BadRequest(ctx, "参数错误")
		return
	}
	if err := h.svc.MarkReadAll(services.ContextUserID(ctx), ids); err != nil {
		h.InternalError(ctx, err)
		return
	}
	h.SuccessMsg(ctx, "操作成功")
}

// ReadUsersList GET /system/notice/readUsers/list
func (h *NoticeHandler) ReadUsersList(ctx *gin.Context) {
	total, rows, err := h.svc.ReadUsersList(ctx)
	if err != nil {
		h.InternalError(ctx, err)
		return
	}
	h.Table(ctx, total, rows)
}
