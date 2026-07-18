package monitor

import (
	"errors"
	"strconv"

	"thirteen-admin/app/models"
	"thirteen-admin/app/services"
	"thirteen-admin/internal/base/handler"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// JobHandler 定时任务 HTTP 处理器。
type JobHandler struct {
	handler.Api
}

// List GET /monitor/job/list
func (e JobHandler) List(c *gin.Context) {
	s := services.JobService{}
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

// Get GET /monitor/job/:jobId
func (e JobHandler) Get(c *gin.Context) {
	s := services.JobService{}
	if err := e.MakeContext(c).MakeOrm().MakeService(&s.Service).Errors; err != nil {
		e.InternalError(c, err)
		return
	}
	id, _ := strconv.ParseUint(c.Param("jobId"), 10, 64)
	data, err := s.Get(uint(id))
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			e.NotFound(c, "任务不存在")
			return
		}
		e.InternalError(c, err)
		return
	}
	e.Success(c, data)
}

type jobBody struct {
	JobID          uint   `json:"jobId"`
	JobName        string `json:"jobName"`
	JobGroup       string `json:"jobGroup"`
	InvokeTarget   string `json:"invokeTarget"`
	CronExpression string `json:"cronExpression"`
	MisfirePolicy  string `json:"misfirePolicy"`
	Concurrent     string `json:"concurrent"`
	Status         string `json:"status"`
	Remark         string `json:"remark"`
}

// Add POST /monitor/job
func (e JobHandler) Add(c *gin.Context) {
	s := services.JobService{}
	if err := e.MakeContext(c).MakeOrm().MakeService(&s.Service).Errors; err != nil {
		e.InternalError(c, err)
		return
	}
	var req jobBody
	if err := c.ShouldBindJSON(&req); err != nil {
		e.BadRequest(c, "参数错误")
		return
	}
	job := &models.Job{
		JobName:        req.JobName,
		JobGroup:       req.JobGroup,
		InvokeTarget:   req.InvokeTarget,
		CronExpression: req.CronExpression,
		MisfirePolicy:  req.MisfirePolicy,
		Concurrent:     req.Concurrent,
		Status:         req.Status,
		Remark:         req.Remark,
	}
	if err := s.Add(job, services.ContextUserName(c)); err != nil {
		e.Fail(c, handler.CodeError, err.Error())
		return
	}
	e.SuccessMsg(c, "新增成功")
}

// Edit PUT /monitor/job
func (e JobHandler) Edit(c *gin.Context) {
	s := services.JobService{}
	if err := e.MakeContext(c).MakeOrm().MakeService(&s.Service).Errors; err != nil {
		e.InternalError(c, err)
		return
	}
	var req jobBody
	if err := c.ShouldBindJSON(&req); err != nil {
		e.BadRequest(c, "参数错误")
		return
	}
	job := &models.Job{
		JobName:        req.JobName,
		JobGroup:       req.JobGroup,
		InvokeTarget:   req.InvokeTarget,
		CronExpression: req.CronExpression,
		MisfirePolicy:  req.MisfirePolicy,
		Concurrent:     req.Concurrent,
		Status:         req.Status,
		Remark:         req.Remark,
	}
	job.ID = req.JobID
	if err := s.Edit(job, services.ContextUserName(c)); err != nil {
		e.InternalError(c, err)
		return
	}
	e.SuccessMsg(c, "修改成功")
}

// Delete DELETE /monitor/job/:jobId
func (e JobHandler) Delete(c *gin.Context) {
	s := services.JobService{}
	if err := e.MakeContext(c).MakeOrm().MakeService(&s.Service).Errors; err != nil {
		e.InternalError(c, err)
		return
	}
	ids := services.ParseUintIDs(c.Param("jobId"))
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

type jobStatusBody struct {
	JobID  uint   `json:"jobId"`
	Status string `json:"status"`
}

// ChangeStatus PUT /monitor/job/changeStatus
func (e JobHandler) ChangeStatus(c *gin.Context) {
	s := services.JobService{}
	if err := e.MakeContext(c).MakeOrm().MakeService(&s.Service).Errors; err != nil {
		e.InternalError(c, err)
		return
	}
	var req jobStatusBody
	if err := c.ShouldBindJSON(&req); err != nil {
		e.BadRequest(c, "参数错误")
		return
	}
	if err := s.ChangeStatus(req.JobID, req.Status); err != nil {
		e.InternalError(c, err)
		return
	}
	e.SuccessMsg(c, "修改成功")
}

type jobRunBody struct {
	JobID    uint   `json:"jobId"`
	JobGroup string `json:"jobGroup"`
}

// Run PUT /monitor/job/run
func (e JobHandler) Run(c *gin.Context) {
	s := services.JobService{}
	if err := e.MakeContext(c).MakeOrm().MakeService(&s.Service).Errors; err != nil {
		e.InternalError(c, err)
		return
	}
	var req jobRunBody
	if err := c.ShouldBindJSON(&req); err != nil {
		e.BadRequest(c, "参数错误")
		return
	}
	if err := s.Run(req.JobID); err != nil {
		e.InternalError(c, err)
		return
	}
	e.SuccessMsg(c, "执行成功")
}

// Export POST /monitor/job/export
func (e JobHandler) Export(c *gin.Context) {
	s := services.JobService{}
	if err := e.MakeContext(c).MakeOrm().MakeService(&s.Service).Errors; err != nil {
		e.InternalError(c, err)
		return
	}
	rows, err := s.Export(c)
	if err != nil {
		e.InternalError(c, err)
		return
	}
	services.WriteCSV(c, "job.csv", []string{"任务编号", "任务名称", "任务组名", "调用目标", "cron表达式", "状态"}, rows)
}

// JobLogHandler 任务日志 HTTP 处理器。
type JobLogHandler struct {
	handler.Api
}

// List GET /monitor/jobLog/list
func (e JobLogHandler) List(c *gin.Context) {
	s := services.JobLogService{}
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

// Delete DELETE /monitor/jobLog/:jobLogId
func (e JobLogHandler) Delete(c *gin.Context) {
	s := services.JobLogService{}
	if err := e.MakeContext(c).MakeOrm().MakeService(&s.Service).Errors; err != nil {
		e.InternalError(c, err)
		return
	}
	ids := services.ParseUintIDs(c.Param("jobLogId"))
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

// Clean DELETE /monitor/jobLog/clean
func (e JobLogHandler) Clean(c *gin.Context) {
	s := services.JobLogService{}
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

// Export POST /monitor/jobLog/export
func (e JobLogHandler) Export(c *gin.Context) {
	s := services.JobLogService{}
	if err := e.MakeContext(c).MakeOrm().MakeService(&s.Service).Errors; err != nil {
		e.InternalError(c, err)
		return
	}
	rows, err := s.Export(c)
	if err != nil {
		e.InternalError(c, err)
		return
	}
	services.WriteCSV(c, "jobLog.csv", []string{"日志编号", "任务名称", "任务组名", "执行状态", "日志信息", "执行时间"}, rows)
}
