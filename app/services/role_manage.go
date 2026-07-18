package services

import (
	"context"
	"errors"
	"fmt"
	"sort"
	"strings"

	"clash-admin/app/models"
	"clash-admin/internal/base/service"

	"gorm.io/gorm"
)

// RoleVO 角色列表响应。
type RoleVO struct {
	ID         uint   `json:"id"`
	CreateBy   string `json:"createBy"`
	CreateTime string `json:"createTime"`
	UpdateBy   string `json:"updateBy"`
	UpdateTime string `json:"updateTime"`
	Status     string `json:"status"`
	RoleName   string `json:"roleName"`
	RoleCode   string `json:"roleCode"`
	RoleDesc   string `json:"roleDesc"`
}

// AllRoleVO 全部启用角色。
type AllRoleVO struct {
	ID       uint   `json:"id"`
	RoleName string `json:"roleName"`
	RoleCode string `json:"roleCode"`
}

// ButtonVO 按钮权限项。
type ButtonVO struct {
	ID    int    `json:"id"`
	Label string `json:"label"`
	Code  string `json:"code"`
}

// RoleListQuery 角色列表查询。
type RoleListQuery struct {
	Current  int
	Size     int
	RoleName string
	RoleCode string
	Status   string
}

// RoleUpsertInput 新增/更新角色入参。
type RoleUpsertInput struct {
	ID       uint
	RoleName string
	RoleCode string
	RoleDesc string
	Status   string
	Operator string
}

// RoleManageService 角色管理业务。
type RoleManageService struct {
	service.Service
}

// NewRoleManageService 创建角色管理服务。
func NewRoleManageService(db *gorm.DB) *RoleManageService {
	return &RoleManageService{Service: service.New(db)}
}

func (s *RoleManageService) toVO(r *models.Role) RoleVO {
	return RoleVO{
		ID:         r.ID,
		CreateBy:   r.CreateBy,
		CreateTime: formatTime(r.CreatedAt),
		UpdateBy:   r.UpdateBy,
		UpdateTime: formatTime(r.UpdatedAt),
		Status:     statusString(r.Status),
		RoleName:   r.RoleName,
		RoleCode:   r.RoleCode,
		RoleDesc:   r.RoleDesc,
	}
}

// GetRoleList 分页查询角色。
func (s *RoleManageService) GetRoleList(ctx context.Context, q RoleListQuery) (*PageResult[RoleVO], error) {
	page, size, offset := parsePage(q.Current, q.Size)
	db := s.DB.WithContext(ctx).Model(&models.Role{})
	if q.RoleName != "" {
		db = db.Where("role_name LIKE ?", "%"+q.RoleName+"%")
	}
	if q.RoleCode != "" {
		db = db.Where("role_code LIKE ?", "%"+q.RoleCode+"%")
	}
	if q.Status != "" {
		st, err := parseStatus(q.Status)
		if err != nil {
			return nil, err
		}
		db = db.Where("status = ?", st)
	}
	var total int64
	if err := db.Count(&total).Error; err != nil {
		return nil, err
	}
	var rows []models.Role
	if err := db.Order("id ASC").Offset(offset).Limit(size).Find(&rows).Error; err != nil {
		return nil, err
	}
	records := make([]RoleVO, 0, len(rows))
	for i := range rows {
		records = append(records, s.toVO(&rows[i]))
	}
	return &PageResult[RoleVO]{Current: page, Size: size, Total: total, Records: records}, nil
}

// GetAllRoles 获取全部启用角色。
func (s *RoleManageService) GetAllRoles(ctx context.Context) ([]AllRoleVO, error) {
	var rows []models.Role
	if err := s.DB.WithContext(ctx).Where("status = ?", 1).Order("id ASC").Find(&rows).Error; err != nil {
		return nil, err
	}
	out := make([]AllRoleVO, 0, len(rows))
	for _, r := range rows {
		out = append(out, AllRoleVO{ID: r.ID, RoleName: r.RoleName, RoleCode: r.RoleCode})
	}
	return out, nil
}

// AddRole 新增角色。
func (s *RoleManageService) AddRole(ctx context.Context, in RoleUpsertInput) error {
	name := strings.TrimSpace(in.RoleName)
	code := strings.TrimSpace(in.RoleCode)
	if name == "" || code == "" {
		return fmt.Errorf("%w: roleName/roleCode 不能为空", ErrManageInvalid)
	}
	st, err := parseStatus(in.Status)
	if err != nil {
		return err
	}
	var count int64
	if err := s.DB.WithContext(ctx).Model(&models.Role{}).Where("role_code = ?", code).Count(&count).Error; err != nil {
		return err
	}
	if count > 0 {
		return ErrManageConflict
	}
	role := models.Role{
		RoleName: name,
		RoleCode: code,
		RoleDesc: strings.TrimSpace(in.RoleDesc),
		Home:     "home",
		Status:   st,
		CreateBy: in.Operator,
		UpdateBy: in.Operator,
	}
	return s.DB.WithContext(ctx).Create(&role).Error
}

// UpdateRole 更新角色。
func (s *RoleManageService) UpdateRole(ctx context.Context, in RoleUpsertInput) error {
	if in.ID == 0 {
		return fmt.Errorf("%w: id 不能为空", ErrManageInvalid)
	}
	name := strings.TrimSpace(in.RoleName)
	code := strings.TrimSpace(in.RoleCode)
	if name == "" || code == "" {
		return fmt.Errorf("%w: roleName/roleCode 不能为空", ErrManageInvalid)
	}
	st, err := parseStatus(in.Status)
	if err != nil {
		return err
	}
	var role models.Role
	if err := s.DB.WithContext(ctx).First(&role, in.ID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return ErrManageNotFound
		}
		return err
	}
	var count int64
	if err := s.DB.WithContext(ctx).Model(&models.Role{}).Where("role_code = ? AND id <> ?", code, in.ID).Count(&count).Error; err != nil {
		return err
	}
	if count > 0 {
		return ErrManageConflict
	}
	return s.DB.WithContext(ctx).Model(&role).Updates(map[string]any{
		"role_name": name,
		"role_code": code,
		"role_desc": strings.TrimSpace(in.RoleDesc),
		"status":    st,
		"update_by": in.Operator,
	}).Error
}

// DeleteRoles 删除角色（若仍被用户引用则失败）。
func (s *RoleManageService) DeleteRoles(ctx context.Context, ids []uint) error {
	if len(ids) == 0 {
		return fmt.Errorf("%w: ids 不能为空", ErrManageInvalid)
	}
	var roles []models.Role
	if err := s.DB.WithContext(ctx).Where("id IN ?", ids).Find(&roles).Error; err != nil {
		return err
	}
	if len(roles) == 0 {
		return ErrManageNotFound
	}
	var users []models.User
	if err := s.DB.WithContext(ctx).Select("id", "roles").Find(&users).Error; err != nil {
		return err
	}
	for _, role := range roles {
		for _, u := range users {
			for _, code := range u.Roles {
				if code == role.RoleCode {
					return fmt.Errorf("%w: %s", ErrManageRoleInUse, role.RoleCode)
				}
			}
		}
	}
	return s.DB.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := tx.Where("role_id IN ?", ids).Delete(&models.RoleMenu{}).Error; err != nil {
			return err
		}
		if err := tx.Where("role_id IN ?", ids).Delete(&models.RoleButton{}).Error; err != nil {
			return err
		}
		return tx.Where("id IN ?", ids).Delete(&models.Role{}).Error
	})
}

// GetRoleMenus 获取角色已选菜单 ID。
func (s *RoleManageService) GetRoleMenus(ctx context.Context, roleID uint) ([]uint, error) {
	if roleID == 0 {
		return nil, fmt.Errorf("%w: roleId 不能为空", ErrManageInvalid)
	}
	var rows []models.RoleMenu
	if err := s.DB.WithContext(ctx).Where("role_id = ?", roleID).Find(&rows).Error; err != nil {
		return nil, err
	}
	ids := make([]uint, 0, len(rows))
	for _, r := range rows {
		ids = append(ids, r.MenuID)
	}
	return ids, nil
}

// UpdateRoleMenus 更新角色菜单权限。
func (s *RoleManageService) UpdateRoleMenus(ctx context.Context, roleID uint, menuIDs []uint) error {
	if roleID == 0 {
		return fmt.Errorf("%w: roleId 不能为空", ErrManageInvalid)
	}
	var count int64
	if err := s.DB.WithContext(ctx).Model(&models.Role{}).Where("id = ?", roleID).Count(&count).Error; err != nil {
		return err
	}
	if count == 0 {
		return ErrManageNotFound
	}
	return s.DB.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := tx.Where("role_id = ?", roleID).Delete(&models.RoleMenu{}).Error; err != nil {
			return err
		}
		if len(menuIDs) == 0 {
			return nil
		}
		rows := make([]models.RoleMenu, 0, len(menuIDs))
		for _, mid := range menuIDs {
			rows = append(rows, models.RoleMenu{RoleID: roleID, MenuID: mid})
		}
		return tx.Create(&rows).Error
	})
}

// GetRoleHome 获取角色首页。
func (s *RoleManageService) GetRoleHome(ctx context.Context, roleID uint) (string, error) {
	var role models.Role
	if err := s.DB.WithContext(ctx).Select("id", "home").First(&role, roleID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return "", ErrManageNotFound
		}
		return "", err
	}
	if role.Home == "" {
		return "home", nil
	}
	return role.Home, nil
}

// UpdateRoleHome 更新角色首页。
func (s *RoleManageService) UpdateRoleHome(ctx context.Context, roleID uint, home string) error {
	home = strings.TrimSpace(home)
	if roleID == 0 || home == "" {
		return fmt.Errorf("%w: roleId/home 不能为空", ErrManageInvalid)
	}
	res := s.DB.WithContext(ctx).Model(&models.Role{}).Where("id = ?", roleID).Update("home", home)
	if res.Error != nil {
		return res.Error
	}
	if res.RowsAffected == 0 {
		return ErrManageNotFound
	}
	return nil
}

// ListAllButtons 从菜单 buttons 展平生成稳定按钮目录。
func (s *RoleManageService) ListAllButtons(ctx context.Context) ([]ButtonVO, error) {
	catalog, err := s.buildButtonCatalog(ctx)
	if err != nil {
		return nil, err
	}
	return catalog, nil
}

// GetRoleButtons 获取角色已选按钮 ID。
func (s *RoleManageService) GetRoleButtons(ctx context.Context, roleID uint) ([]int, error) {
	catalog, err := s.buildButtonCatalog(ctx)
	if err != nil {
		return nil, err
	}
	codeToID := make(map[string]int, len(catalog))
	for _, b := range catalog {
		codeToID[b.Code] = b.ID
	}
	var rows []models.RoleButton
	if err := s.DB.WithContext(ctx).Where("role_id = ?", roleID).Find(&rows).Error; err != nil {
		return nil, err
	}
	ids := make([]int, 0, len(rows))
	for _, r := range rows {
		if id, ok := codeToID[r.ButtonCode]; ok {
			ids = append(ids, id)
		}
	}
	sort.Ints(ids)
	return ids, nil
}

// UpdateRoleButtons 更新角色按钮权限。
func (s *RoleManageService) UpdateRoleButtons(ctx context.Context, roleID uint, buttonIDs []int) error {
	if roleID == 0 {
		return fmt.Errorf("%w: roleId 不能为空", ErrManageInvalid)
	}
	var count int64
	if err := s.DB.WithContext(ctx).Model(&models.Role{}).Where("id = ?", roleID).Count(&count).Error; err != nil {
		return err
	}
	if count == 0 {
		return ErrManageNotFound
	}
	catalog, err := s.buildButtonCatalog(ctx)
	if err != nil {
		return err
	}
	idToCode := make(map[int]string, len(catalog))
	for _, b := range catalog {
		idToCode[b.ID] = b.Code
	}
	codes := make([]string, 0, len(buttonIDs))
	seen := map[string]struct{}{}
	for _, id := range buttonIDs {
		code, ok := idToCode[id]
		if !ok {
			continue
		}
		if _, exists := seen[code]; exists {
			continue
		}
		seen[code] = struct{}{}
		codes = append(codes, code)
	}
	return s.DB.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := tx.Where("role_id = ?", roleID).Delete(&models.RoleButton{}).Error; err != nil {
			return err
		}
		if len(codes) == 0 {
			return nil
		}
		rows := make([]models.RoleButton, 0, len(codes))
		for _, code := range codes {
			rows = append(rows, models.RoleButton{RoleID: roleID, ButtonCode: code})
		}
		return tx.Create(&rows).Error
	})
}

func (s *RoleManageService) buildButtonCatalog(ctx context.Context) ([]ButtonVO, error) {
	var menus []models.Menu
	if err := s.DB.WithContext(ctx).Select("id", "buttons").Find(&menus).Error; err != nil {
		return nil, err
	}
	type item struct {
		code  string
		label string
	}
	uniq := map[string]string{}
	for _, m := range menus {
		for _, b := range m.Buttons {
			code := strings.TrimSpace(b.Code)
			if code == "" {
				continue
			}
			if _, ok := uniq[code]; !ok {
				label := strings.TrimSpace(b.Desc)
				if label == "" {
					label = code
				}
				uniq[code] = label
			}
		}
	}
	codes := make([]string, 0, len(uniq))
	for code := range uniq {
		codes = append(codes, code)
	}
	sort.Strings(codes)
	out := make([]ButtonVO, 0, len(codes))
	for i, code := range codes {
		out = append(out, ButtonVO{ID: i + 1, Label: uniq[code], Code: code})
	}
	return out, nil
}

// ResolveButtonsByRoleCodes 按角色编码聚合按钮权限（供登录用户信息使用）。
func ResolveButtonsByRoleCodes(ctx context.Context, db *gorm.DB, roleCodes []string) ([]string, error) {
	for _, code := range roleCodes {
		if code == "R_SUPER" {
			return []string{"*"}, nil
		}
	}
	if len(roleCodes) == 0 {
		return []string{}, nil
	}
	var roles []models.Role
	if err := db.WithContext(ctx).Where("role_code IN ? AND status = ?", roleCodes, 1).Find(&roles).Error; err != nil {
		return nil, err
	}
	if len(roles) == 0 {
		return []string{}, nil
	}
	roleIDs := make([]uint, 0, len(roles))
	for _, r := range roles {
		roleIDs = append(roleIDs, r.ID)
	}
	var rows []models.RoleButton
	if err := db.WithContext(ctx).Where("role_id IN ?", roleIDs).Find(&rows).Error; err != nil {
		return nil, err
	}
	seen := map[string]struct{}{}
	out := make([]string, 0, len(rows))
	for _, r := range rows {
		if r.ButtonCode == "*" {
			return []string{"*"}, nil
		}
		if _, ok := seen[r.ButtonCode]; ok {
			continue
		}
		seen[r.ButtonCode] = struct{}{}
		out = append(out, r.ButtonCode)
	}
	sort.Strings(out)
	return out, nil
}
