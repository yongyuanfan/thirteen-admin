package services

import (
	"errors"

	"thirteen-admin/app/models"
	"thirteen-admin/internal/base/service"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// OperLogService 操作日志业务。
type OperLogService struct {
	service.Service
}

// NewOperLogService 创建操作日志服务。
func NewOperLogService(db *gorm.DB) *OperLogService {
	return &OperLogService{Service: service.New(db)}
}

func (s *OperLogService) buildQuery(ctx *gin.Context) *gorm.DB {
	q := s.DB.Model(&models.OperLog{})
	if title := ctx.Query("title"); title != "" {
		q = q.Where("title LIKE ?", "%"+title+"%")
	}
	if operName := ctx.Query("operName"); operName != "" {
		q = q.Where("oper_name LIKE ?", "%"+operName+"%")
	}
	if bt := ctx.Query("businessType"); bt != "" {
		q = q.Where("business_type = ?", bt)
	}
	if st := ctx.Query("status"); st != "" {
		q = q.Where("status = ?", st)
	}
	if ip := ctx.Query("operIp"); ip != "" {
		q = q.Where("oper_ip LIKE ?", "%"+ip+"%")
	}
	return ApplyTimeRange(q, ctx, "oper_time")
}

// List 分页查询操作日志。
func (s *OperLogService) List(ctx *gin.Context) (int64, []map[string]any, error) {
	page := ParsePage(ctx)
	q := s.buildQuery(ctx)
	var total int64
	if err := q.Count(&total).Error; err != nil {
		return 0, nil, err
	}
	q = ApplyOrder(q, ctx, "oper_time desc")
	var logs []models.OperLog
	if err := ApplyPage(q, page).Find(&logs).Error; err != nil {
		return 0, nil, err
	}
	rows := make([]map[string]any, 0, len(logs))
	for i := range logs {
		rows = append(rows, OperLogToMap(&logs[i]))
	}
	return total, rows, nil
}

// Delete 删除操作日志。
func (s *OperLogService) Delete(ids []uint) error {
	if len(ids) == 0 {
		return errors.New("请选择要删除的数据")
	}
	return s.DB.Where("id IN ?", ids).Delete(&models.OperLog{}).Error
}

// Clean 清空操作日志。
func (s *OperLogService) Clean() error {
	return s.DB.Where("1 = 1").Delete(&models.OperLog{}).Error
}

// Export 导出操作日志 CSV 行。
func (s *OperLogService) Export(ctx *gin.Context) ([][]string, error) {
	_, rows, err := s.List(ctx)
	if err != nil {
		return nil, err
	}
	out := make([][]string, 0, len(rows))
	for _, row := range rows {
		out = append(out, []string{
			toString(row["operId"]),
			toString(row["title"]),
			toString(row["businessType"]),
			toString(row["operName"]),
			toString(row["operIp"]),
			toString(row["status"]),
			toString(row["operTime"]),
		})
	}
	return out, nil
}
