package services

import (
	"errors"

	"thirteen-admin/app/models"
	"thirteen-admin/internal/base/service"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// MenuService 菜单业务。
type MenuService struct {
	service.Service
}

// NewMenuService 创建菜单服务。
func NewMenuService(db *gorm.DB) *MenuService {
	return &MenuService{Service: service.New(db)}
}

// List 查询菜单列表。
func (s *MenuService) List(ctx *gin.Context) ([]map[string]any, error) {
	var menus []models.Menu
	q := s.DB.Model(&models.Menu{})
	if name := ctx.Query("menuName"); name != "" {
		q = q.Where("menu_name LIKE ?", "%"+name+"%")
	}
	if status := ctx.Query("status"); status != "" {
		q = q.Where("status = ?", status)
	}
	if err := q.Order("parent_id asc, order_num asc").Find(&menus).Error; err != nil {
		return nil, err
	}
	rows := make([]map[string]any, 0, len(menus))
	for i := range menus {
		rows = append(rows, MenuToMap(&menus[i]))
	}
	return rows, nil
}

// Get 获取菜单详情。
func (s *MenuService) Get(id uint) (map[string]any, error) {
	var menu models.Menu
	if err := s.DB.First(&menu, id).Error; err != nil {
		return nil, err
	}
	return MenuToMap(&menu), nil
}

// TreeSelect 菜单下拉树。
func (s *MenuService) TreeSelect() ([]map[string]any, error) {
	var menus []models.Menu
	if err := s.DB.Where("status = ?", "0").Order("parent_id asc, order_num asc").Find(&menus).Error; err != nil {
		return nil, err
	}
	return BuildMenuTreeSelect(menus, 0), nil
}

// RoleMenuTreeSelect 角色菜单树及已选 keys。
func (s *MenuService) RoleMenuTreeSelect(roleID uint) (menus []map[string]any, checkedKeys []uint, err error) {
	menus, err = s.TreeSelect()
	if err != nil {
		return nil, nil, err
	}
	err = s.DB.Model(&models.RoleMenu{}).Where("role_id = ?", roleID).Pluck("menu_id", &checkedKeys).Error
	if checkedKeys == nil {
		checkedKeys = []uint{}
	}
	return menus, checkedKeys, err
}

// Add 新增菜单。
func (s *MenuService) Add(menu *models.Menu, operator string) error {
	menu.CreatedBy = operator
	return s.DB.Create(menu).Error
}

// Edit 修改菜单。
func (s *MenuService) Edit(menu *models.Menu, operator string) error {
	if menu.ParentID == menu.ID {
		return errors.New("上级菜单不能选择自己")
	}
	menu.UpdatedBy = operator
	return s.DB.Model(menu).Updates(map[string]any{
		"menu_name":  menu.MenuName,
		"parent_id":  menu.ParentID,
		"order_num":  menu.OrderNum,
		"path":       menu.Path,
		"component":  menu.Component,
		"query":      menu.Query,
		"route_name": menu.RouteName,
		"is_frame":   menu.IsFrame,
		"is_cache":   menu.IsCache,
		"menu_type":  menu.MenuType,
		"visible":    menu.Visible,
		"status":     menu.Status,
		"perms":      menu.Perms,
		"icon":       menu.Icon,
		"remark":     menu.Remark,
		"updated_by": menu.UpdatedBy,
	}).Error
}

// UpdateSort 批量更新排序。
func (s *MenuService) UpdateSort(idsStr, numsStr string) error {
	ids, nums, err := ParseSortUpdate(idsStr, numsStr)
	if err != nil {
		return err
	}
	return s.DB.Transaction(func(tx *gorm.DB) error {
		for i, id := range ids {
			if err := tx.Model(&models.Menu{}).Where("id = ?", id).Update("order_num", nums[i]).Error; err != nil {
				return err
			}
		}
		return nil
	})
}

// Delete 删除菜单。
func (s *MenuService) Delete(id uint) error {
	var childCount int64
	if err := s.DB.Model(&models.Menu{}).Where("parent_id = ?", id).Count(&childCount).Error; err != nil {
		return err
	}
	if childCount > 0 {
		return errors.New("存在子菜单，不允许删除")
	}
	return s.DB.Transaction(func(tx *gorm.DB) error {
		if err := tx.Delete(&models.Menu{}, id).Error; err != nil {
			return err
		}
		return tx.Where("menu_id = ?", id).Delete(&models.RoleMenu{}).Error
	})
}
