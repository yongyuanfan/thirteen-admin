package services

import (
	"errors"
	"strconv"

	"thirteen-admin/app/models"
	"thirteen-admin/internal/base/service"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// RoleService 角色业务。
type RoleService struct {
	service.Service
	dept *DeptService
}

// NewRoleService 创建角色服务。
func NewRoleService(db *gorm.DB, dept *DeptService) *RoleService {
	return &RoleService{
		Service: service.New(db),
		dept:    dept,
	}
}

// List 分页查询角色。
func (s *RoleService) List(ctx *gin.Context) (int64, []map[string]any, error) {
	page := ParsePage(ctx)
	q := s.DB.Model(&models.Role{}).Where("del_flag = ?", "0")
	if name := ctx.Query("roleName"); name != "" {
		q = q.Where("role_name LIKE ?", "%"+name+"%")
	}
	if key := ctx.Query("roleKey"); key != "" {
		q = q.Where("role_key LIKE ?", "%"+key+"%")
	}
	if status := ctx.Query("status"); status != "" {
		q = q.Where("status = ?", status)
	}
	q = ApplyTimeRange(q, ctx, "created_at")
	var total int64
	if err := q.Count(&total).Error; err != nil {
		return 0, nil, err
	}
	q = ApplyOrder(q, ctx, "role_sort asc")
	var roles []models.Role
	if err := ApplyPage(q, page).Find(&roles).Error; err != nil {
		return 0, nil, err
	}
	rows := make([]map[string]any, 0, len(roles))
	for i := range roles {
		rows = append(rows, RoleToMap(&roles[i]))
	}
	return total, rows, nil
}

// Get 获取角色详情。
func (s *RoleService) Get(id uint) (map[string]any, error) {
	var role models.Role
	if err := s.DB.Where("id = ? AND del_flag = ?", id, "0").First(&role).Error; err != nil {
		return nil, err
	}
	return RoleToMap(&role), nil
}

// Add 新增角色。
func (s *RoleService) Add(role *models.Role, menuIDs []uint, operator string) error {
	var count int64
	if err := s.DB.Model(&models.Role{}).Where("role_key = ? AND del_flag = ?", role.RoleKey, "0").Count(&count).Error; err != nil {
		return err
	}
	if count > 0 {
		return errors.New("新增角色'" + role.RoleName + "'失败，角色权限已存在")
	}
	role.DelFlag = "0"
	role.CreatedBy = operator
	return s.DB.Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(role).Error; err != nil {
			return err
		}
		return s.saveRoleMenus(tx, role.ID, menuIDs)
	})
}

// Edit 修改角色。
func (s *RoleService) Edit(role *models.Role, menuIDs []uint, operator string) error {
	if role.ID == 1 {
		return errors.New("不允许操作超级管理员角色")
	}
	role.UpdatedBy = operator
	return s.DB.Transaction(func(tx *gorm.DB) error {
		if err := tx.Model(role).Updates(map[string]any{
			"role_name":           role.RoleName,
			"role_key":            role.RoleKey,
			"role_sort":           role.RoleSort,
			"data_scope":          role.DataScope,
			"menu_check_strictly": role.MenuCheckStrictly,
			"dept_check_strictly": role.DeptCheckStrictly,
			"status":              role.Status,
			"remark":              role.Remark,
			"updated_by":          role.UpdatedBy,
		}).Error; err != nil {
			return err
		}
		return s.saveRoleMenus(tx, role.ID, menuIDs)
	})
}

// DataScope 保存数据权限。
func (s *RoleService) DataScope(role *models.Role, deptIDs []uint, operator string) error {
	if role.ID == 1 {
		return errors.New("不允许操作超级管理员角色")
	}
	role.UpdatedBy = operator
	return s.DB.Transaction(func(tx *gorm.DB) error {
		if err := tx.Model(role).Updates(map[string]any{
			"data_scope":          role.DataScope,
			"dept_check_strictly": role.DeptCheckStrictly,
			"updated_by":          role.UpdatedBy,
		}).Error; err != nil {
			return err
		}
		if err := tx.Where("role_id = ?", role.ID).Delete(&models.RoleDept{}).Error; err != nil {
			return err
		}
		for _, deptID := range deptIDs {
			if err := tx.Create(&models.RoleDept{RoleID: role.ID, DeptID: deptID}).Error; err != nil {
				return err
			}
		}
		return nil
	})
}

// ChangeStatus 修改状态。
func (s *RoleService) ChangeStatus(roleID uint, status, operator string) error {
	if roleID == 1 {
		return errors.New("不允许操作超级管理员角色")
	}
	return s.DB.Model(&models.Role{}).Where("id = ?", roleID).Updates(map[string]any{
		"status":     status,
		"updated_by": operator,
	}).Error
}

// Delete 删除角色。
func (s *RoleService) Delete(ids []uint) error {
	if len(ids) == 0 {
		return errors.New("请选择要删除的数据")
	}
	for _, id := range ids {
		if id == 1 {
			return errors.New("不允许操作超级管理员角色")
		}
	}
	return s.DB.Transaction(func(tx *gorm.DB) error {
		if err := tx.Model(&models.Role{}).Where("id IN ?", ids).Update("del_flag", "2").Error; err != nil {
			return err
		}
		if err := tx.Where("role_id IN ?", ids).Delete(&models.RoleMenu{}).Error; err != nil {
			return err
		}
		if err := tx.Where("role_id IN ?", ids).Delete(&models.RoleDept{}).Error; err != nil {
			return err
		}
		return tx.Where("role_id IN ?", ids).Delete(&models.UserRole{}).Error
	})
}

// AllocatedList 已分配用户。
func (s *RoleService) AllocatedList(ctx *gin.Context) (int64, []map[string]any, error) {
	return s.authUserList(ctx, true)
}

// UnallocatedList 未分配用户。
func (s *RoleService) UnallocatedList(ctx *gin.Context) (int64, []map[string]any, error) {
	return s.authUserList(ctx, false)
}

// CancelAuthUser 取消单个授权。
func (s *RoleService) CancelAuthUser(userID, roleID uint) error {
	return s.DB.Where("user_id = ? AND role_id = ?", userID, roleID).Delete(&models.UserRole{}).Error
}

// CancelAuthUserAll 批量取消授权。
func (s *RoleService) CancelAuthUserAll(roleID uint, userIDs []uint) error {
	if len(userIDs) == 0 {
		return errors.New("请选择要取消授权的用户")
	}
	return s.DB.Where("role_id = ? AND user_id IN ?", roleID, userIDs).Delete(&models.UserRole{}).Error
}

// SelectAuthUserAll 批量授权。
func (s *RoleService) SelectAuthUserAll(roleID uint, userIDs []uint) error {
	if len(userIDs) == 0 {
		return errors.New("请选择要授权的用户")
	}
	return s.DB.Transaction(func(tx *gorm.DB) error {
		for _, userID := range userIDs {
			var count int64
			if err := tx.Model(&models.UserRole{}).Where("user_id = ? AND role_id = ?", userID, roleID).Count(&count).Error; err != nil {
				return err
			}
			if count > 0 {
				continue
			}
			if err := tx.Create(&models.UserRole{UserID: userID, RoleID: roleID}).Error; err != nil {
				return err
			}
		}
		return nil
	})
}

// DeptTree 数据权限部门树。
func (s *RoleService) DeptTree(roleID uint) ([]map[string]any, []uint, error) {
	depts, err := s.dept.TreeSelect()
	if err != nil {
		return nil, nil, err
	}
	var checked []uint
	if err := s.DB.Model(&models.RoleDept{}).Where("role_id = ?", roleID).Pluck("dept_id", &checked).Error; err != nil {
		return nil, nil, err
	}
	if checked == nil {
		checked = []uint{}
	}
	return depts, checked, nil
}

// Export 导出角色。
func (s *RoleService) Export(ctx *gin.Context) ([][]string, error) {
	_, rows, err := s.List(ctx)
	if err != nil {
		return nil, err
	}
	out := make([][]string, 0, len(rows))
	for _, row := range rows {
		out = append(out, []string{
			toString(row["roleId"]),
			toString(row["roleName"]),
			toString(row["roleKey"]),
			toString(row["roleSort"]),
			toString(row["status"]),
			toString(row["createdAt"]),
		})
	}
	return out, nil
}

func (s *RoleService) saveRoleMenus(tx *gorm.DB, roleID uint, menuIDs []uint) error {
	if err := tx.Where("role_id = ?", roleID).Delete(&models.RoleMenu{}).Error; err != nil {
		return err
	}
	for _, menuID := range menuIDs {
		if err := tx.Create(&models.RoleMenu{RoleID: roleID, MenuID: menuID}).Error; err != nil {
			return err
		}
	}
	return nil
}

func (s *RoleService) authUserList(ctx *gin.Context, allocated bool) (int64, []map[string]any, error) {
	page := ParsePage(ctx)
	roleIDStr := ctx.Query("roleId")
	roleID, _ := strconv.ParseUint(roleIDStr, 10, 64)
	q := s.DB.Model(&models.User{}).Where("del_flag = ?", "0")
	if allocated {
		q = q.Where("id IN (SELECT user_id FROM sys_user_role WHERE role_id = ?)", roleID)
	} else {
		q = q.Where("id NOT IN (SELECT user_id FROM sys_user_role WHERE role_id = ?)", roleID)
	}
	if name := ctx.Query("userName"); name != "" {
		q = q.Where("user_name LIKE ?", "%"+name+"%")
	}
	if phone := ctx.Query("phonenumber"); phone != "" {
		q = q.Where("phonenumber LIKE ?", "%"+phone+"%")
	}
	var total int64
	if err := q.Count(&total).Error; err != nil {
		return 0, nil, err
	}
	var users []models.User
	if err := ApplyPage(q.Order("id asc"), page).Find(&users).Error; err != nil {
		return 0, nil, err
	}
	rows := make([]map[string]any, 0, len(users))
	for i := range users {
		rows = append(rows, UserToMap(&users[i]))
	}
	return total, rows, nil
}
