package services

import (
	"errors"
	"fmt"
	"strconv"

	"thirteen-admin/app/models"
	"thirteen-admin/internal/base/service"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// DeptService 部门业务。
type DeptService struct {
	service.Service
}

// NewDeptService 创建部门服务。
func NewDeptService(db *gorm.DB) *DeptService {
	return &DeptService{Service: service.New(db)}
}

// List 查询部门列表。
func (s *DeptService) List(ctx *gin.Context) ([]map[string]any, error) {
	var depts []models.Dept
	q := s.DB.Model(&models.Dept{}).Where("del_flag = ?", "0")
	if name := ctx.Query("deptName"); name != "" {
		q = q.Where("dept_name LIKE ?", "%"+name+"%")
	}
	if status := ctx.Query("status"); status != "" {
		q = q.Where("status = ?", status)
	}
	if err := q.Order("parent_id asc, order_num asc").Find(&depts).Error; err != nil {
		return nil, err
	}
	rows := make([]map[string]any, 0, len(depts))
	for i := range depts {
		rows = append(rows, DeptToMap(&depts[i]))
	}
	return rows, nil
}

// ListExclude 排除指定部门及其子孙。
func (s *DeptService) ListExclude(excludeID uint) ([]map[string]any, error) {
	excludeIDs, err := s.ChildDeptIDs(excludeID)
	if err != nil {
		return nil, err
	}
	var depts []models.Dept
	q := s.DB.Where("del_flag = ?", "0")
	if len(excludeIDs) > 0 {
		q = q.Where("id NOT IN ?", excludeIDs)
	}
	if err := q.Order("parent_id asc, order_num asc").Find(&depts).Error; err != nil {
		return nil, err
	}
	rows := make([]map[string]any, 0, len(depts))
	for i := range depts {
		rows = append(rows, DeptToMap(&depts[i]))
	}
	return rows, nil
}

// Get 获取部门详情。
func (s *DeptService) Get(id uint) (map[string]any, error) {
	var dept models.Dept
	if err := s.DB.Where("id = ? AND del_flag = ?", id, "0").First(&dept).Error; err != nil {
		return nil, err
	}
	return DeptToMap(&dept), nil
}

// Add 新增部门。
func (s *DeptService) Add(dept *models.Dept, operator string) error {
	ancestors, err := s.buildAncestors(dept.ParentID)
	if err != nil {
		return err
	}
	dept.Ancestors = ancestors
	dept.DelFlag = "0"
	dept.CreatedBy = operator
	return s.DB.Create(dept).Error
}

// Edit 修改部门。
func (s *DeptService) Edit(dept *models.Dept, operator string) error {
	var old models.Dept
	if err := s.DB.Where("id = ? AND del_flag = ?", dept.ID, "0").First(&old).Error; err != nil {
		return err
	}
	if dept.ParentID == dept.ID {
		return errors.New("上级部门不能是自己")
	}
	newAncestors, err := s.buildAncestors(dept.ParentID)
	if err != nil {
		return err
	}
	oldPrefix := old.Ancestors + "," + strconv.FormatUint(uint64(old.ID), 10)
	newPrefix := newAncestors + "," + strconv.FormatUint(uint64(old.ID), 10)
	dept.Ancestors = newAncestors
	dept.UpdatedBy = operator
	return s.DB.Transaction(func(tx *gorm.DB) error {
		if err := tx.Model(dept).Updates(map[string]any{
			"parent_id":  dept.ParentID,
			"ancestors":  dept.Ancestors,
			"dept_name":  dept.DeptName,
			"order_num":  dept.OrderNum,
			"leader":     dept.Leader,
			"phone":      dept.Phone,
			"email":      dept.Email,
			"status":     dept.Status,
			"updated_by": dept.UpdatedBy,
		}).Error; err != nil {
			return err
		}
		if old.ParentID != dept.ParentID {
			return tx.Model(&models.Dept{}).
				Where("ancestors LIKE ? OR ancestors LIKE ?", oldPrefix+"%", oldPrefix).
				Update("ancestors", gorm.Expr("REPLACE(ancestors, ?, ?)", oldPrefix, newPrefix)).Error
		}
		return nil
	})
}

// UpdateSort 批量更新排序。
func (s *DeptService) UpdateSort(idsStr, numsStr string) error {
	ids, nums, err := ParseSortUpdate(idsStr, numsStr)
	if err != nil {
		return err
	}
	return s.DB.Transaction(func(tx *gorm.DB) error {
		for i, id := range ids {
			if err := tx.Model(&models.Dept{}).Where("id = ?", id).Update("order_num", nums[i]).Error; err != nil {
				return err
			}
		}
		return nil
	})
}

// Delete 删除部门。
func (s *DeptService) Delete(id uint) error {
	var childCount int64
	if err := s.DB.Model(&models.Dept{}).Where("parent_id = ? AND del_flag = ?", id, "0").Count(&childCount).Error; err != nil {
		return err
	}
	if childCount > 0 {
		return errors.New("存在下级部门，不允许删除")
	}
	var userCount int64
	if err := s.DB.Model(&models.User{}).Where("dept_id = ? AND del_flag = ?", id, "0").Count(&userCount).Error; err != nil {
		return err
	}
	if userCount > 0 {
		return errors.New("部门存在用户，不允许删除")
	}
	return s.DB.Model(&models.Dept{}).Where("id = ?", id).Update("del_flag", "2").Error
}

// TreeSelect 部门下拉树。
func (s *DeptService) TreeSelect() ([]map[string]any, error) {
	var depts []models.Dept
	if err := s.DB.Where("del_flag = ?", "0").Order("parent_id asc, order_num asc").Find(&depts).Error; err != nil {
		return nil, err
	}
	return BuildDeptTreeSelect(depts, 0), nil
}

// ChildDeptIDs 获取部门及所有子部门 ID。
func (s *DeptService) ChildDeptIDs(deptID uint) ([]uint, error) {
	if deptID == 0 {
		return nil, nil
	}
	var ids []uint
	prefix := fmt.Sprintf("%%,%d,%%", deptID)
	if err := s.DB.Model(&models.Dept{}).Where("del_flag = ? AND (id = ? OR ancestors LIKE ? OR ancestors LIKE ?)",
		"0", deptID, prefix, fmt.Sprintf("%%,%d", deptID)).Pluck("id", &ids).Error; err != nil {
		return nil, err
	}
	return ids, nil
}

func (s *DeptService) buildAncestors(parentID uint) (string, error) {
	if parentID == 0 {
		return "0", nil
	}
	var parent models.Dept
	if err := s.DB.Where("id = ? AND del_flag = ?", parentID, "0").First(&parent).Error; err != nil {
		return "", errors.New("上级部门不存在")
	}
	return parent.Ancestors + "," + strconv.FormatUint(uint64(parentID), 10), nil
}
