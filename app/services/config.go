package services

import (
	"errors"

	"thirteen-admin/app/models"
	"thirteen-admin/internal/base/service"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// ConfigService 参数配置业务。
type ConfigService struct {
	service.Service
}

// NewConfigService 创建参数服务。
func NewConfigService(db *gorm.DB) *ConfigService {
	return &ConfigService{Service: service.New(db)}
}

// List 分页查询参数。
func (s *ConfigService) List(ctx *gin.Context) (int64, []map[string]any, error) {
	page := ParsePage(ctx)
	q := s.DB.Model(&models.Config{})
	if name := ctx.Query("configName"); name != "" {
		q = q.Where("config_name LIKE ?", "%"+name+"%")
	}
	if key := ctx.Query("configKey"); key != "" {
		q = q.Where("config_key LIKE ?", "%"+key+"%")
	}
	if typ := ctx.Query("configType"); typ != "" {
		q = q.Where("config_type = ?", typ)
	}
	q = ApplyTimeRange(q, ctx, "created_at")
	var total int64
	if err := q.Count(&total).Error; err != nil {
		return 0, nil, err
	}
	q = ApplyOrder(q, ctx, "id asc")
	var configs []models.Config
	if err := ApplyPage(q, page).Find(&configs).Error; err != nil {
		return 0, nil, err
	}
	rows := make([]map[string]any, 0, len(configs))
	for i := range configs {
		rows = append(rows, ConfigToMap(&configs[i]))
	}
	return total, rows, nil
}

// Get 获取参数详情。
func (s *ConfigService) Get(id uint) (map[string]any, error) {
	var cfg models.Config
	if err := s.DB.First(&cfg, id).Error; err != nil {
		return nil, err
	}
	return ConfigToMap(&cfg), nil
}

// GetByKey 按键名获取参数值。
func (s *ConfigService) GetByKey(key string) (string, error) {
	var cfg models.Config
	if err := s.DB.Where("config_key = ?", key).First(&cfg).Error; err != nil {
		return "", err
	}
	return cfg.ConfigValue, nil
}

// Add 新增参数。
func (s *ConfigService) Add(cfg *models.Config, operator string) error {
	var count int64
	if err := s.DB.Model(&models.Config{}).Where("config_key = ?", cfg.ConfigKey).Count(&count).Error; err != nil {
		return err
	}
	if count > 0 {
		return errors.New("参数键名已存在")
	}
	cfg.CreatedBy = operator
	return s.DB.Create(cfg).Error
}

// Edit 修改参数。
func (s *ConfigService) Edit(cfg *models.Config, operator string) error {
	cfg.UpdatedBy = operator
	return s.DB.Model(cfg).Updates(map[string]any{
		"config_name":  cfg.ConfigName,
		"config_key":   cfg.ConfigKey,
		"config_value": cfg.ConfigValue,
		"config_type":  cfg.ConfigType,
		"remark":       cfg.Remark,
		"updated_by":   cfg.UpdatedBy,
	}).Error
}

// Delete 删除参数。
func (s *ConfigService) Delete(ids []uint) error {
	if len(ids) == 0 {
		return errors.New("请选择要删除的数据")
	}
	return s.DB.Where("id IN ?", ids).Delete(&models.Config{}).Error
}

// Export 导出参数 CSV 行。
func (s *ConfigService) Export(ctx *gin.Context) ([][]string, error) {
	_, rows, err := s.List(ctx)
	if err != nil {
		return nil, err
	}
	out := make([][]string, 0, len(rows))
	for _, row := range rows {
		out = append(out, []string{
			toString(row["configId"]),
			toString(row["configName"]),
			toString(row["configKey"]),
			toString(row["configValue"]),
			toString(row["configType"]),
			toString(row["createTime"]),
		})
	}
	return out, nil
}
