package handlers

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
	handler.Handler
	svc *services.JobService
}

// NewJobHandler 创建定时任务 Handler。
func NewJobHandler(svc *services.JobService) *JobHandler {
	return &JobHandler{svc: svc}
}

// List GET /monitor/job/list
func (h *JobHandler) List(ctx *gin.Context) {
	total, rows, err := h.svc.List(ctx)
	if err != nil {
		h.InternalError(ctx, err)
		return
	}
	h.Table(ctx, total, rows)
}

// Get GET /monitor/job/:jobId
func (h *JobHandler) Get(ctx *gin.Context) {
	id, _ := strconv.ParseUint(ctx.Param("jobId"), 10, 64)
	data, err := h.svc.Get(uint(id))
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			h.NotFound(ctx, "任务不存在")
			return
		}
		h.InternalError(ctx, err)
		return
	}
	h.Success(ctx, data)
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
func (h *JobHandler) Add(ctx *gin.Context) {
	var req jobBody
	if err := ctx.ShouldBindJSON(&req); err != nil {
		h.BadRequest(ctx, "参数错误")
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
	if err := h.svc.Add(job, services.ContextUserName(ctx)); err != nil {
		h.Fail(ctx, handler.CodeError, err.Error())
		return
	}
	h.SuccessMsg(ctx, "新增成功")
}

// Edit PUT /monitor/job
func (h *JobHandler) Edit(ctx *gin.Context) {
	var req jobBody
	if err := ctx.ShouldBindJSON(&req); err != nil {
		h.BadRequest(ctx, "参数错误")
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
	if err := h.svc.Edit(job, services.ContextUserName(ctx)); err != nil {
		h.InternalError(ctx, err)
		return
	}
	h.SuccessMsg(ctx, "修改成功")
}

// Delete DELETE /monitor/job/:jobId
func (h *JobHandler) Delete(ctx *gin.Context) {
	ids := services.ParseUintIDs(ctx.Param("jobId"))
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

type jobStatusBody struct {
	JobID  uint   `json:"jobId"`
	Status string `json:"status"`
}

// ChangeStatus PUT /monitor/job/changeStatus
func (h *JobHandler) ChangeStatus(ctx *gin.Context) {
	var req jobStatusBody
	if err := ctx.ShouldBindJSON(&req); err != nil {
		h.BadRequest(ctx, "参数错误")
		return
	}
	if err := h.svc.ChangeStatus(req.JobID, req.Status); err != nil {
		h.InternalError(ctx, err)
		return
	}
	h.SuccessMsg(ctx, "修改成功")
}

type jobRunBody struct {
	JobID    uint   `json:"jobId"`
	JobGroup string `json:"jobGroup"`
}

// Run PUT /monitor/job/run
func (h *JobHandler) Run(ctx *gin.Context) {
	var req jobRunBody
	if err := ctx.ShouldBindJSON(&req); err != nil {
		h.BadRequest(ctx, "参数错误")
		return
	}
	if err := h.svc.Run(req.JobID); err != nil {
		h.InternalError(ctx, err)
		return
	}
	h.SuccessMsg(ctx, "执行成功")
}

// Export POST /monitor/job/export
func (h *JobHandler) Export(ctx *gin.Context) {
	rows, err := h.svc.Export(ctx)
	if err != nil {
		h.InternalError(ctx, err)
		return
	}
	services.WriteCSV(ctx, "job.csv", []string{"任务编号", "任务名称", "任务组名", "调用目标", "cron表达式", "状态"}, rows)
}

// JobLogHandler 任务日志 HTTP 处理器。
type JobLogHandler struct {
	handler.Handler
	svc *services.JobLogService
}

// NewJobLogHandler 创建任务日志 Handler。
func NewJobLogHandler(svc *services.JobLogService) *JobLogHandler {
	return &JobLogHandler{svc: svc}
}

// List GET /monitor/jobLog/list
func (h *JobLogHandler) List(ctx *gin.Context) {
	total, rows, err := h.svc.List(ctx)
	if err != nil {
		h.InternalError(ctx, err)
		return
	}
	h.Table(ctx, total, rows)
}

// Delete DELETE /monitor/jobLog/:jobLogId
func (h *JobLogHandler) Delete(ctx *gin.Context) {
	ids := services.ParseUintIDs(ctx.Param("jobLogId"))
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

// Clean DELETE /monitor/jobLog/clean
func (h *JobLogHandler) Clean(ctx *gin.Context) {
	if err := h.svc.Clean(); err != nil {
		h.InternalError(ctx, err)
		return
	}
	h.SuccessMsg(ctx, "清空成功")
}

// Export POST /monitor/jobLog/export
func (h *JobLogHandler) Export(ctx *gin.Context) {
	rows, err := h.svc.Export(ctx)
	if err != nil {
		h.InternalError(ctx, err)
		return
	}
	services.WriteCSV(ctx, "jobLog.csv", []string{"日志编号", "任务名称", "任务组名", "执行状态", "日志信息", "执行时间"}, rows)
}
