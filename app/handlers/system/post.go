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
	handler.Api
}

// List GET /system/post/list
func (e PostHandler) List(c *gin.Context) {
	s := services.PostService{}
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

// Get GET /system/post/:postId
func (e PostHandler) Get(c *gin.Context) {
	s := services.PostService{}
	if err := e.MakeContext(c).MakeOrm().MakeService(&s.Service).Errors; err != nil {
		e.InternalError(c, err)
		return
	}
	id, _ := strconv.ParseUint(c.Param("postId"), 10, 64)
	data, err := s.Get(uint(id))
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			e.NotFound(c, "岗位不存在")
			return
		}
		e.InternalError(c, err)
		return
	}
	e.Success(c, data)
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
func (e PostHandler) Add(c *gin.Context) {
	s := services.PostService{}
	if err := e.MakeContext(c).MakeOrm().MakeService(&s.Service).Errors; err != nil {
		e.InternalError(c, err)
		return
	}
	var req postBody
	if err := c.ShouldBindJSON(&req); err != nil {
		e.BadRequest(c, "参数错误")
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
	if err := s.Add(post, services.ContextUserName(c)); err != nil {
		e.InternalError(c, err)
		return
	}
	e.SuccessMsg(c, "新增成功")
}

// Edit PUT /system/post
func (e PostHandler) Edit(c *gin.Context) {
	s := services.PostService{}
	if err := e.MakeContext(c).MakeOrm().MakeService(&s.Service).Errors; err != nil {
		e.InternalError(c, err)
		return
	}
	var req postBody
	if err := c.ShouldBindJSON(&req); err != nil {
		e.BadRequest(c, "参数错误")
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
	if err := s.Edit(post, services.ContextUserName(c)); err != nil {
		e.InternalError(c, err)
		return
	}
	e.SuccessMsg(c, "修改成功")
}

// Delete DELETE /system/post/:postId
func (e PostHandler) Delete(c *gin.Context) {
	s := services.PostService{}
	if err := e.MakeContext(c).MakeOrm().MakeService(&s.Service).Errors; err != nil {
		e.InternalError(c, err)
		return
	}
	ids := services.ParseUintIDs(c.Param("postId"))
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

// Export POST /system/post/export
func (e PostHandler) Export(c *gin.Context) {
	s := services.PostService{}
	if err := e.MakeContext(c).MakeOrm().MakeService(&s.Service).Errors; err != nil {
		e.InternalError(c, err)
		return
	}
	rows, err := s.Export(c)
	if err != nil {
		e.InternalError(c, err)
		return
	}
	services.WriteCSV(c, "post.csv", []string{"岗位编号", "岗位编码", "岗位名称", "岗位排序", "状态", "创建时间"}, rows)
}
