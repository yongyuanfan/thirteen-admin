package services

import (
	"context"
	"errors"
	"fmt"

	"thirteen-admin/app/models"
	"thirteen-admin/internal/base/service"

	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

const loginLockKeyPref = "login:lock:"

// LogininforService 登录日志业务。
type LogininforService struct {
	service.Service
	rdb *redis.Client
}

// NewLogininforService 创建登录日志服务。
func NewLogininforService(db *gorm.DB, rdb *redis.Client) *LogininforService {
	return &LogininforService{Service: service.New(db), rdb: rdb}
}

func (s *LogininforService) buildQuery(ctx *gin.Context) *gorm.DB {
	q := s.DB.Model(&models.Logininfor{})
	if userName := ctx.Query("userName"); userName != "" {
		q = q.Where("user_name LIKE ?", "%"+userName+"%")
	}
	if ip := ctx.Query("ipaddr"); ip != "" {
		q = q.Where("ipaddr LIKE ?", "%"+ip+"%")
	}
	if st := ctx.Query("status"); st != "" {
		q = q.Where("status = ?", st)
	}
	return ApplyTimeRange(q, ctx, "login_time")
}

// List 分页查询登录日志。
func (s *LogininforService) List(ctx *gin.Context) (int64, []map[string]any, error) {
	page := ParsePage(ctx)
	q := s.buildQuery(ctx)
	var total int64
	if err := q.Count(&total).Error; err != nil {
		return 0, nil, err
	}
	q = ApplyOrder(q, ctx, "login_time desc")
	var logs []models.Logininfor
	if err := ApplyPage(q, page).Find(&logs).Error; err != nil {
		return 0, nil, err
	}
	rows := make([]map[string]any, 0, len(logs))
	for i := range logs {
		rows = append(rows, LogininforToMap(&logs[i]))
	}
	return total, rows, nil
}

// Delete 删除登录日志。
func (s *LogininforService) Delete(ids []uint) error {
	if len(ids) == 0 {
		return errors.New("请选择要删除的数据")
	}
	return s.DB.Where("id IN ?", ids).Delete(&models.Logininfor{}).Error
}

// Unlock 解锁登录失败锁定的账号。
func (s *LogininforService) Unlock(ctx context.Context, userName string) error {
	if userName == "" {
		return errors.New("用户名不能为空")
	}
	if s.rdb == nil {
		return nil
	}
	return s.rdb.Del(ctx, loginLockKey(userName)).Err()
}

// Clean 清空登录日志。
func (s *LogininforService) Clean() error {
	return s.DB.Where("1 = 1").Delete(&models.Logininfor{}).Error
}

// Export 导出登录日志 CSV 行。
func (s *LogininforService) Export(ctx *gin.Context) ([][]string, error) {
	_, rows, err := s.List(ctx)
	if err != nil {
		return nil, err
	}
	out := make([][]string, 0, len(rows))
	for _, row := range rows {
		out = append(out, []string{
			toString(row["infoId"]),
			toString(row["userName"]),
			toString(row["ipaddr"]),
			toString(row["status"]),
			toString(row["msg"]),
			toString(row["loginTime"]),
		})
	}
	return out, nil
}

func loginLockKey(userName string) string {
	return fmt.Sprintf("%s%s", loginLockKeyPref, userName)
}
