package monitor

import (
	"thirteen-admin/app/runtime"
	"thirteen-admin/app/services"
	"thirteen-admin/internal/base/handler"

	"github.com/gin-gonic/gin"
)

// LogininforHandler 登录日志 HTTP 处理器。
type LogininforHandler struct {
	handler.Api
}

func (e *LogininforHandler) svc(c *gin.Context) (*services.LogininforService, error) {
	if err := e.MakeContext(c).MakeOrm().Errors; err != nil {
		return nil, err
	}
	return services.NewLogininforService(e.Orm, runtime.Cache()), nil
}

// List GET /monitor/logininfor/list
func (e LogininforHandler) List(c *gin.Context) {
	s, err := e.svc(c)
	if err != nil {
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

// Delete DELETE /monitor/logininfor/:infoId
func (e LogininforHandler) Delete(c *gin.Context) {
	s, err := e.svc(c)
	if err != nil {
		e.InternalError(c, err)
		return
	}
	ids := services.ParseUintIDs(c.Param("infoId"))
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

// Unlock GET /monitor/logininfor/unlock/:userName
func (e LogininforHandler) Unlock(c *gin.Context) {
	s, err := e.svc(c)
	if err != nil {
		e.InternalError(c, err)
		return
	}
	userName := c.Param("userName")
	if err := s.Unlock(c.Request.Context(), userName); err != nil {
		e.InternalError(c, err)
		return
	}
	e.SuccessMsg(c, "解锁成功")
}

// Clean DELETE /monitor/logininfor/clean
func (e LogininforHandler) Clean(c *gin.Context) {
	s, err := e.svc(c)
	if err != nil {
		e.InternalError(c, err)
		return
	}
	if err := s.Clean(); err != nil {
		e.InternalError(c, err)
		return
	}
	e.SuccessMsg(c, "清空成功")
}

// Export POST /monitor/logininfor/export
func (e LogininforHandler) Export(c *gin.Context) {
	s, err := e.svc(c)
	if err != nil {
		e.InternalError(c, err)
		return
	}
	rows, err := s.Export(c)
	if err != nil {
		e.InternalError(c, err)
		return
	}
	services.WriteCSV(c, "logininfor.csv", []string{"访问编号", "用户名称", "登录地址", "登录状态", "描述", "访问时间"}, rows)
}
