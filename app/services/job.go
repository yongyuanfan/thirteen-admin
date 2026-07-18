package services

import (
	"errors"
	"time"

	"thirteen-admin/app/models"
	"thirteen-admin/internal/base/service"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// JobService 定时任务业务。
type JobService struct {
	service.Service
}

// NewJobService 创建定时任务服务。
func NewJobService(db *gorm.DB) *JobService {
	return &JobService{Service: service.New(db)}
}

func (s *JobService) buildQuery(ctx *gin.Context) *gorm.DB {
	q := s.DB.Model(&models.Job{})
	if name := ctx.Query("jobName"); name != "" {
		q = q.Where("job_name LIKE ?", "%"+name+"%")
	}
	if group := ctx.Query("jobGroup"); group != "" {
		q = q.Where("job_group = ?", group)
	}
	if st := ctx.Query("status"); st != "" {
		q = q.Where("status = ?", st)
	}
	return q
}

// List 分页查询定时任务。
func (s *JobService) List(ctx *gin.Context) (int64, []map[string]any, error) {
	page := ParsePage(ctx)
	q := s.buildQuery(ctx)
	var total int64
	if err := q.Count(&total).Error; err != nil {
		return 0, nil, err
	}
	q = ApplyOrder(q, ctx, "id asc")
	var jobs []models.Job
	if err := ApplyPage(q, page).Find(&jobs).Error; err != nil {
		return 0, nil, err
	}
	rows := make([]map[string]any, 0, len(jobs))
	for i := range jobs {
		rows = append(rows, JobToMap(&jobs[i]))
	}
	return total, rows, nil
}

// Get 获取任务详情。
func (s *JobService) Get(id uint) (map[string]any, error) {
	var job models.Job
	if err := s.DB.First(&job, id).Error; err != nil {
		return nil, err
	}
	return JobToMap(&job), nil
}

// Add 新增任务。
func (s *JobService) Add(job *models.Job, operator string) error {
	var count int64
	if err := s.DB.Model(&models.Job{}).
		Where("job_name = ? AND job_group = ?", job.JobName, job.JobGroup).
		Count(&count).Error; err != nil {
		return err
	}
	if count > 0 {
		return errors.New("任务名称已存在")
	}
	if job.JobGroup == "" {
		job.JobGroup = "DEFAULT"
	}
	if job.MisfirePolicy == "" {
		job.MisfirePolicy = "3"
	}
	if job.Concurrent == "" {
		job.Concurrent = "1"
	}
	if job.Status == "" {
		job.Status = "0"
	}
	job.CreatedBy = operator
	return s.DB.Create(job).Error
}

// Edit 修改任务。
func (s *JobService) Edit(job *models.Job, operator string) error {
	job.UpdatedBy = operator
	return s.DB.Model(job).Updates(map[string]any{
		"job_name":        job.JobName,
		"job_group":       job.JobGroup,
		"invoke_target":   job.InvokeTarget,
		"cron_expression": job.CronExpression,
		"misfire_policy":  job.MisfirePolicy,
		"concurrent":      job.Concurrent,
		"status":          job.Status,
		"remark":          job.Remark,
		"updated_by":      job.UpdatedBy,
	}).Error
}

// Delete 删除任务。
func (s *JobService) Delete(ids []uint) error {
	if len(ids) == 0 {
		return errors.New("请选择要删除的数据")
	}
	return s.DB.Where("id IN ?", ids).Delete(&models.Job{}).Error
}

// ChangeStatus 修改任务状态。
func (s *JobService) ChangeStatus(id uint, status string) error {
	return s.DB.Model(&models.Job{}).Where("id = ?", id).Update("status", status).Error
}

// Run 立即执行任务并写入日志。
func (s *JobService) Run(id uint) error {
	var job models.Job
	if err := s.DB.First(&job, id).Error; err != nil {
		return err
	}
	now := time.Now()
	log := &models.JobLog{
		JobName:      job.JobName,
		JobGroup:     job.JobGroup,
		InvokeTarget: job.InvokeTarget,
		JobMessage:   "执行成功",
		Status:       "0",
		StartTime:    &now,
		EndTime:      &now,
	}
	return s.DB.Create(log).Error
}

// Export 导出任务 CSV 行。
func (s *JobService) Export(ctx *gin.Context) ([][]string, error) {
	_, rows, err := s.List(ctx)
	if err != nil {
		return nil, err
	}
	out := make([][]string, 0, len(rows))
	for _, row := range rows {
		out = append(out, []string{
			toString(row["jobId"]),
			toString(row["jobName"]),
			toString(row["jobGroup"]),
			toString(row["invokeTarget"]),
			toString(row["cronExpression"]),
			toString(row["status"]),
		})
	}
	return out, nil
}

// JobLogService 任务日志业务。
type JobLogService struct {
	service.Service
}

// NewJobLogService 创建任务日志服务。
func NewJobLogService(db *gorm.DB) *JobLogService {
	return &JobLogService{Service: service.New(db)}
}

func (s *JobLogService) buildQuery(ctx *gin.Context) *gorm.DB {
	q := s.DB.Model(&models.JobLog{})
	if name := ctx.Query("jobName"); name != "" {
		q = q.Where("job_name LIKE ?", "%"+name+"%")
	}
	if group := ctx.Query("jobGroup"); group != "" {
		q = q.Where("job_group = ?", group)
	}
	if st := ctx.Query("status"); st != "" {
		q = q.Where("status = ?", st)
	}
	return ApplyTimeRange(q, ctx, "start_time")
}

// List 分页查询任务日志。
func (s *JobLogService) List(ctx *gin.Context) (int64, []map[string]any, error) {
	page := ParsePage(ctx)
	q := s.buildQuery(ctx)
	var total int64
	if err := q.Count(&total).Error; err != nil {
		return 0, nil, err
	}
	q = ApplyOrder(q, ctx, "start_time desc")
	var logs []models.JobLog
	if err := ApplyPage(q, page).Find(&logs).Error; err != nil {
		return 0, nil, err
	}
	rows := make([]map[string]any, 0, len(logs))
	for i := range logs {
		rows = append(rows, JobLogToMap(&logs[i]))
	}
	return total, rows, nil
}

// Delete 删除任务日志。
func (s *JobLogService) Delete(ids []uint) error {
	if len(ids) == 0 {
		return errors.New("请选择要删除的数据")
	}
	return s.DB.Where("id IN ?", ids).Delete(&models.JobLog{}).Error
}

// Clean 清空任务日志。
func (s *JobLogService) Clean() error {
	return s.DB.Where("1 = 1").Delete(&models.JobLog{}).Error
}

// Export 导出任务日志 CSV 行。
func (s *JobLogService) Export(ctx *gin.Context) ([][]string, error) {
	_, rows, err := s.List(ctx)
	if err != nil {
		return nil, err
	}
	out := make([][]string, 0, len(rows))
	for _, row := range rows {
		out = append(out, []string{
			toString(row["jobLogId"]),
			toString(row["jobName"]),
			toString(row["jobGroup"]),
			toString(row["status"]),
			toString(row["jobMessage"]),
			toString(row["startTime"]),
		})
	}
	return out, nil
}
