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
	handler.Api
}

// List GET /system/notice/list
func (e NoticeHandler) List(c *gin.Context) {
	s := services.NoticeService{}
	if err := e.MakeContext(c).MakeOrm().MakeService(&s.Service).Errors; err != nil {
		e.InternalError(c, err)
		return
	}
	total, rows, err := s.List(c)
	if err != nil {
		e.InternalError(c, err)
		return
	}
	e.Table(c, total, rows)
}

// Get GET /system/notice/:noticeId
func (e NoticeHandler) Get(c *gin.Context) {
	s := services.NoticeService{}
	if err := e.MakeContext(c).MakeOrm().MakeService(&s.Service).Errors; err != nil {
		e.InternalError(c, err)
		return
	}
	id, _ := strconv.ParseUint(c.Param("noticeId"), 10, 64)
	data, err := s.Get(uint(id))
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			e.NotFound(c, "公告不存在")
			return
		}
		e.InternalError(c, err)
		return
	}
	e.Success(c, data)
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
func (e NoticeHandler) Add(c *gin.Context) {
	s := services.NoticeService{}
	if err := e.MakeContext(c).MakeOrm().MakeService(&s.Service).Errors; err != nil {
		e.InternalError(c, err)
		return
	}
	var req noticeBody
	if err := c.ShouldBindJSON(&req); err != nil {
		e.BadRequest(c, "参数错误")
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
	if err := s.Add(notice, services.ContextUserName(c)); err != nil {
		e.InternalError(c, err)
		return
	}
	e.SuccessMsg(c, "新增成功")
}

// Edit PUT /system/notice
func (e NoticeHandler) Edit(c *gin.Context) {
	s := services.NoticeService{}
	if err := e.MakeContext(c).MakeOrm().MakeService(&s.Service).Errors; err != nil {
		e.InternalError(c, err)
		return
	}
	var req noticeBody
	if err := c.ShouldBindJSON(&req); err != nil {
		e.BadRequest(c, "参数错误")
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
	if err := s.Edit(notice, services.ContextUserName(c)); err != nil {
		e.InternalError(c, err)
		return
	}
	e.SuccessMsg(c, "修改成功")
}

// Delete DELETE /system/notice/:noticeId
func (e NoticeHandler) Delete(c *gin.Context) {
	s := services.NoticeService{}
	if err := e.MakeContext(c).MakeOrm().MakeService(&s.Service).Errors; err != nil {
		e.InternalError(c, err)
		return
	}
	ids := services.ParseUintIDs(c.Param("noticeId"))
	if len(ids) == 0 {
		e.BadRequest(c, "参数错误")
		return
	}
	if err := s.Delete(ids); err != nil {
		e.Fail(c, handler.CodeError, err.Error())
		return
	}
	e.SuccessMsg(c, "删除成功")
}

// ListTop GET /system/notice/listTop
func (e NoticeHandler) ListTop(c *gin.Context) {
	s := services.NoticeService{}
	if err := e.MakeContext(c).MakeOrm().MakeService(&s.Service).Errors; err != nil {
		e.InternalError(c, err)
		return
	}
	rows, unread, err := s.ListTop(services.ContextUserID(c))
	if err != nil {
		e.InternalError(c, err)
		return
	}
	e.OK(c, gin.H{
		"data":        rows,
		"unreadCount": unread,
	})
}

// MarkRead POST /system/notice/markRead
func (e NoticeHandler) MarkRead(c *gin.Context) {
	s := services.NoticeService{}
	if err := e.MakeContext(c).MakeOrm().MakeService(&s.Service).Errors; err != nil {
		e.InternalError(c, err)
		return
	}
	noticeID, _ := strconv.ParseUint(c.Query("noticeId"), 10, 64)
	if noticeID == 0 {
		e.BadRequest(c, "参数错误")
		return
	}
	if err := s.MarkRead(services.ContextUserID(c), uint(noticeID)); err != nil {
		e.InternalError(c, err)
		return
	}
	e.SuccessMsg(c, "操作成功")
}

// MarkReadAll POST /system/notice/markReadAll
func (e NoticeHandler) MarkReadAll(c *gin.Context) {
	s := services.NoticeService{}
	if err := e.MakeContext(c).MakeOrm().MakeService(&s.Service).Errors; err != nil {
		e.InternalError(c, err)
		return
	}
	ids := services.ParseUintIDs(c.Query("ids"))
	if len(ids) == 0 {
		e.BadRequest(c, "参数错误")
		return
	}
	if err := s.MarkReadAll(services.ContextUserID(c), ids); err != nil {
		e.InternalError(c, err)
		return
	}
	e.SuccessMsg(c, "操作成功")
}

// ReadUsersList GET /system/notice/readUsers/list
func (e NoticeHandler) ReadUsersList(c *gin.Context) {
	s := services.NoticeService{}
	if err := e.MakeContext(c).MakeOrm().MakeService(&s.Service).Errors; err != nil {
		e.InternalError(c, err)
		return
	}
	total, rows, err := s.ReadUsersList(c)
	if err != nil {
		e.InternalError(c, err)
		return
	}
	e.Table(c, total, rows)
}
