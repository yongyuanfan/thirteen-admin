package services

import (
	"errors"

	"thirteen-admin/app/models"
	"thirteen-admin/internal/base/service"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// DictService 字典业务。
type DictService struct {
	service.Service
}

// NewDictService 创建字典服务。
func NewDictService(db *gorm.DB) *DictService {
	return &DictService{Service: service.New(db)}
}

// --- 字典类型 ---

// TypeList 分页查询字典类型。
func (s *DictService) TypeList(ctx *gin.Context) (int64, []map[string]any, error) {
	page := ParsePage(ctx)
	q := s.DB.Model(&models.DictType{})
	if name := ctx.Query("dictName"); name != "" {
		q = q.Where("dict_name LIKE ?", "%"+name+"%")
	}
	if typ := ctx.Query("dictType"); typ != "" {
		q = q.Where("dict_type LIKE ?", "%"+typ+"%")
	}
	if status := ctx.Query("status"); status != "" {
		q = q.Where("status = ?", status)
	}
	q = ApplyTimeRange(q, ctx, "created_at")
	var total int64
	if err := q.Count(&total).Error; err != nil {
		return 0, nil, err
	}
	q = ApplyOrder(q, ctx, "id asc")
	var items []models.DictType
	if err := ApplyPage(q, page).Find(&items).Error; err != nil {
		return 0, nil, err
	}
	rows := make([]map[string]any, 0, len(items))
	for i := range items {
		rows = append(rows, DictTypeToMap(&items[i]))
	}
	return total, rows, nil
}

// TypeGet 获取字典类型。
func (s *DictService) TypeGet(id uint) (map[string]any, error) {
	var item models.DictType
	if err := s.DB.First(&item, id).Error; err != nil {
		return nil, err
	}
	return DictTypeToMap(&item), nil
}

// TypeAdd 新增字典类型。
func (s *DictService) TypeAdd(item *models.DictType, operator string) error {
	var count int64
	if err := s.DB.Model(&models.DictType{}).Where("dict_type = ?", item.DictType).Count(&count).Error; err != nil {
		return err
	}
	if count > 0 {
		return errors.New("字典类型已存在")
	}
	item.CreatedBy = operator
	return s.DB.Create(item).Error
}

// TypeEdit 修改字典类型。
func (s *DictService) TypeEdit(item *models.DictType, operator string) error {
	item.UpdatedBy = operator
	return s.DB.Transaction(func(tx *gorm.DB) error {
		var old models.DictType
		if err := tx.First(&old, item.ID).Error; err != nil {
			return err
		}
		if err := tx.Model(item).Updates(map[string]any{
			"dict_name":  item.DictName,
			"dict_type":  item.DictType,
			"status":     item.Status,
			"remark":     item.Remark,
			"updated_by": item.UpdatedBy,
		}).Error; err != nil {
			return err
		}
		if old.DictType != item.DictType {
			return tx.Model(&models.DictData{}).Where("dict_type = ?", old.DictType).
				Update("dict_type", item.DictType).Error
		}
		return nil
	})
}

// TypeDelete 删除字典类型。
func (s *DictService) TypeDelete(ids []uint) error {
	if len(ids) == 0 {
		return errors.New("请选择要删除的数据")
	}
	return s.DB.Transaction(func(tx *gorm.DB) error {
		var types []models.DictType
		if err := tx.Where("id IN ?", ids).Find(&types).Error; err != nil {
			return err
		}
		for _, t := range types {
			if err := tx.Where("dict_type = ?", t.DictType).Delete(&models.DictData{}).Error; err != nil {
				return err
			}
		}
		return tx.Where("id IN ?", ids).Delete(&models.DictType{}).Error
	})
}

// TypeOptionSelect 字典类型下拉。
func (s *DictService) TypeOptionSelect() ([]map[string]any, error) {
	var items []models.DictType
	if err := s.DB.Where("status = ?", "0").Order("id asc").Find(&items).Error; err != nil {
		return nil, err
	}
	rows := make([]map[string]any, 0, len(items))
	for i := range items {
		rows = append(rows, DictTypeToMap(&items[i]))
	}
	return rows, nil
}

// TypeExport 导出字典类型。
func (s *DictService) TypeExport(ctx *gin.Context) ([][]string, error) {
	_, rows, err := s.TypeList(ctx)
	if err != nil {
		return nil, err
	}
	out := make([][]string, 0, len(rows))
	for _, row := range rows {
		out = append(out, []string{
			toString(row["dictId"]),
			toString(row["dictName"]),
			toString(row["dictType"]),
			toString(row["status"]),
			toString(row["createdAt"]),
		})
	}
	return out, nil
}

// --- 字典数据 ---

// DataList 分页查询字典数据。
func (s *DictService) DataList(ctx *gin.Context) (int64, []map[string]any, error) {
	page := ParsePage(ctx)
	q := s.DB.Model(&models.DictData{})
	if typ := ctx.Query("dictType"); typ != "" {
		q = q.Where("dict_type = ?", typ)
	}
	if label := ctx.Query("dictLabel"); label != "" {
		q = q.Where("dict_label LIKE ?", "%"+label+"%")
	}
	if status := ctx.Query("status"); status != "" {
		q = q.Where("status = ?", status)
	}
	var total int64
	if err := q.Count(&total).Error; err != nil {
		return 0, nil, err
	}
	q = ApplyOrder(q, ctx, "dict_sort asc")
	var items []models.DictData
	if err := ApplyPage(q, page).Find(&items).Error; err != nil {
		return 0, nil, err
	}
	rows := make([]map[string]any, 0, len(items))
	for i := range items {
		rows = append(rows, DictDataToMap(&items[i]))
	}
	return total, rows, nil
}

// DataGet 获取字典数据。
func (s *DictService) DataGet(id uint) (map[string]any, error) {
	var item models.DictData
	if err := s.DB.First(&item, id).Error; err != nil {
		return nil, err
	}
	return DictDataToMap(&item), nil
}

// DataByType 按类型获取字典数据。
func (s *DictService) DataByType(dictType string) ([]map[string]any, error) {
	var items []models.DictData
	if err := s.DB.Where("dict_type = ? AND status = ?", dictType, "0").
		Order("dict_sort asc").Find(&items).Error; err != nil {
		return nil, err
	}
	rows := make([]map[string]any, 0, len(items))
	for i := range items {
		rows = append(rows, DictDataToMap(&items[i]))
	}
	return rows, nil
}

// DataAdd 新增字典数据。
func (s *DictService) DataAdd(item *models.DictData, operator string) error {
	item.CreatedBy = operator
	return s.DB.Create(item).Error
}

// DataEdit 修改字典数据。
func (s *DictService) DataEdit(item *models.DictData, operator string) error {
	item.UpdatedBy = operator
	return s.DB.Model(item).Updates(map[string]any{
		"dict_sort":  item.DictSort,
		"dict_label": item.DictLabel,
		"dict_value": item.DictValue,
		"dict_type":  item.DictType,
		"css_class":  item.CSSClass,
		"list_class": item.ListClass,
		"is_default": item.IsDefault,
		"status":     item.Status,
		"remark":     item.Remark,
		"updated_by": item.UpdatedBy,
	}).Error
}

// DataDelete 删除字典数据。
func (s *DictService) DataDelete(ids []uint) error {
	if len(ids) == 0 {
		return errors.New("请选择要删除的数据")
	}
	return s.DB.Where("id IN ?", ids).Delete(&models.DictData{}).Error
}

// DataExport 导出字典数据。
func (s *DictService) DataExport(ctx *gin.Context) ([][]string, error) {
	_, rows, err := s.DataList(ctx)
	if err != nil {
		return nil, err
	}
	out := make([][]string, 0, len(rows))
	for _, row := range rows {
		out = append(out, []string{
			toString(row["dictCode"]),
			toString(row["dictLabel"]),
			toString(row["dictValue"]),
			toString(row["dictType"]),
			toString(row["dictSort"]),
			toString(row["status"]),
		})
	}
	return out, nil
}
