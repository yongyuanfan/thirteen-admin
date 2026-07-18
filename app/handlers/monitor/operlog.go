package monitor

import (
	"thirteen-admin/app/services"
	"thirteen-admin/internal/base/handler"

	"github.com/gin-gonic/gin"
)

// OperLogHandler 操作日志 HTTP 处理器。
type OperLogHandler struct {
	handler.Api
}

// List GET /monitor/operlog/list
func (e OperLogHandler) List(c *gin.Context) {
	s := services.OperLogService{}
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

// Delete DELETE /monitor/operlog/:operId
func (e OperLogHandler) Delete(c *gin.Context) {
	s := services.OperLogService{}
	if err := e.MakeContext(c).MakeOrm().MakeService(&s.Service).Errors; err != nil {
		e.InternalError(c, err)
		return
	}
	ids := services.ParseUintIDs(c.Param("operId"))
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

// Clean DELETE /monitor/operlog/clean
func (e OperLogHandler) Clean(c *gin.Context) {
	s := services.OperLogService{}
	if err := e.MakeContext(c).MakeOrm().MakeService(&s.Service).Errors; err != nil {
		e.InternalError(c, err)
		return
	}
	if err := s.Clean(); err != nil {
		e.InternalError(c, err)
		return
	}
	e.SuccessMsg(c, "清空成功")
}

// Export POST /monitor/operlog/export
func (e OperLogHandler) Export(c *gin.Context) {
	s := services.OperLogService{}
	if err := e.MakeContext(c).MakeOrm().MakeService(&s.Service).Errors; err != nil {
		e.InternalError(c, err)
		return
	}
	rows, err := s.Export(c)
	if err != nil {
		e.InternalError(c, err)
		return
	}
	services.WriteCSV(c, "operlog.csv", []string{"日志编号", "系统模块", "业务类型", "操作人员", "主机", "操作状态", "操作时间"}, rows)
}
