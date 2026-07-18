package services

import (
	"context"
	"errors"
	"fmt"
	"sort"
	"strings"

	"clash-admin/app/models"
	"clash-admin/internal/base/model"
	"clash-admin/internal/base/service"

	"gorm.io/gorm"
)

// MenuVO 菜单响应（含树形 children）。
type MenuVO struct {
	ID              uint                `json:"id"`
	CreateBy        string              `json:"createBy"`
	CreateTime      string              `json:"createTime"`
	UpdateBy        string              `json:"updateBy"`
	UpdateTime      string              `json:"updateTime"`
	Status          string              `json:"status"`
	ParentID        uint                `json:"parentId"`
	MenuType        string              `json:"menuType"`
	MenuName        string              `json:"menuName"`
	RouteName       string              `json:"routeName"`
	RoutePath       string              `json:"routePath"`
	Component       string              `json:"component,omitempty"`
	Icon            string              `json:"icon"`
	IconType        string              `json:"iconType"`
	Buttons         []models.MenuButton `json:"buttons"`
	Children        []MenuVO            `json:"children"`
	I18nKey         *string             `json:"i18nKey"`
	KeepAlive       *bool               `json:"keepAlive"`
	Constant        *bool               `json:"constant"`
	Order           *int                `json:"order"`
	Href            *string             `json:"href"`
	HideInMenu      *bool               `json:"hideInMenu"`
	ActiveMenu      *string             `json:"activeMenu"`
	MultiTab        *bool               `json:"multiTab"`
	FixedIndexInTab *int                `json:"fixedIndexInTab"`
	Query           []models.MenuQuery  `json:"query"`
}

// MenuTreeVO 角色菜单树节点。
type MenuTreeVO struct {
	ID       uint         `json:"id"`
	Label    string       `json:"label"`
	PID      uint         `json:"pId"`
	Children []MenuTreeVO `json:"children,omitempty"`
}

// MenuUpsertInput 新增/更新菜单入参。
type MenuUpsertInput struct {
	ID              uint
	ParentID        uint
	MenuType        string
	MenuName        string
	RouteName       string
	RoutePath       string
	Component       string
	Icon            string
	IconType        string
	Buttons         []models.MenuButton
	Query           []models.MenuQuery
	I18nKey         *string
	KeepAlive       *bool
	Constant        *bool
	Order           *int
	Href            *string
	HideInMenu      *bool
	ActiveMenu      *string
	MultiTab        *bool
	FixedIndexInTab *int
	Status          string
	Operator        string
}

// MenuManageService 菜单管理业务。
type MenuManageService struct {
	service.Service
}

// NewMenuManageService 创建菜单管理服务。
func NewMenuManageService(db *gorm.DB) *MenuManageService {
	return &MenuManageService{Service: service.New(db)}
}

func (s *MenuManageService) toVO(m *models.Menu) MenuVO {
	buttons := []models.MenuButton(m.Buttons)
	if buttons == nil {
		buttons = []models.MenuButton{}
	}
	query := []models.MenuQuery(m.Query)
	if query == nil {
		query = []models.MenuQuery{}
	}
	return MenuVO{
		ID:              m.ID,
		CreateBy:        m.CreateBy,
		CreateTime:      formatTime(m.CreatedAt),
		UpdateBy:        m.UpdateBy,
		UpdateTime:      formatTime(m.UpdatedAt),
		Status:          statusString(m.Status),
		ParentID:        m.ParentID,
		MenuType:        m.MenuType,
		MenuName:        m.MenuName,
		RouteName:       m.RouteName,
		RoutePath:       m.RoutePath,
		Component:       m.Component,
		Icon:            m.Icon,
		IconType:        m.IconType,
		Buttons:         buttons,
		Children:        []MenuVO{},
		I18nKey:         m.I18nKey,
		KeepAlive:       m.KeepAlive,
		Constant:        m.Constant,
		Order:           m.Order,
		Href:            m.Href,
		HideInMenu:      m.HideInMenu,
		ActiveMenu:      m.ActiveMenu,
		MultiTab:        m.MultiTab,
		FixedIndexInTab: m.FixedIndexInTab,
		Query:           query,
	}
}

func menuOrderValue(m *models.Menu) int {
	if m.Order == nil {
		return 0
	}
	return *m.Order
}

// GetMenuList 返回菜单树列表。
func (s *MenuManageService) GetMenuList(ctx context.Context) (*PageResult[MenuVO], error) {
	var rows []models.Menu
	if err := s.DB.WithContext(ctx).Find(&rows).Error; err != nil {
		return nil, err
	}
	sort.Slice(rows, func(i, j int) bool {
		oi, oj := menuOrderValue(&rows[i]), menuOrderValue(&rows[j])
		if oi != oj {
			return oi < oj
		}
		return rows[i].ID < rows[j].ID
	})
	nodes := make(map[uint]MenuVO, len(rows))
	childrenMap := map[uint][]uint{}
	for i := range rows {
		vo := s.toVO(&rows[i])
		nodes[rows[i].ID] = vo
		childrenMap[rows[i].ParentID] = append(childrenMap[rows[i].ParentID], rows[i].ID)
	}
	var build func(id uint) MenuVO
	build = func(id uint) MenuVO {
		vo := nodes[id]
		kids := childrenMap[id]
		vo.Children = make([]MenuVO, 0, len(kids))
		for _, kid := range kids {
			vo.Children = append(vo.Children, build(kid))
		}
		return vo
	}
	rootIDs := childrenMap[0]
	records := make([]MenuVO, 0, len(rootIDs))
	for _, id := range rootIDs {
		records = append(records, build(id))
	}
	total := int64(len(rows))
	return &PageResult[MenuVO]{
		Current: 1,
		Size:    int(total),
		Total:   total,
		Records: records,
	}, nil
}

// GetMenuTree 获取精简菜单树。
func (s *MenuManageService) GetMenuTree(ctx context.Context) ([]MenuTreeVO, error) {
	var rows []models.Menu
	if err := s.DB.WithContext(ctx).Find(&rows).Error; err != nil {
		return nil, err
	}
	sort.Slice(rows, func(i, j int) bool {
		oi, oj := menuOrderValue(&rows[i]), menuOrderValue(&rows[j])
		if oi != oj {
			return oi < oj
		}
		return rows[i].ID < rows[j].ID
	})
	childrenMap := map[uint][]uint{}
	labelMap := map[uint]string{}
	pidMap := map[uint]uint{}
	for _, r := range rows {
		childrenMap[r.ParentID] = append(childrenMap[r.ParentID], r.ID)
		labelMap[r.ID] = r.MenuName
		pidMap[r.ID] = r.ParentID
	}
	var build func(id uint) MenuTreeVO
	build = func(id uint) MenuTreeVO {
		node := MenuTreeVO{ID: id, Label: labelMap[id], PID: pidMap[id]}
		for _, kid := range childrenMap[id] {
			node.Children = append(node.Children, build(kid))
		}
		return node
	}
	roots := make([]MenuTreeVO, 0, len(childrenMap[0]))
	for _, id := range childrenMap[0] {
		roots = append(roots, build(id))
	}
	return roots, nil
}

// GetAllPages 返回全部 routeName。
func (s *MenuManageService) GetAllPages(ctx context.Context) ([]string, error) {
	var rows []models.Menu
	if err := s.DB.WithContext(ctx).Select("route_name").Order("id ASC").Find(&rows).Error; err != nil {
		return nil, err
	}
	out := make([]string, 0, len(rows))
	seen := map[string]struct{}{}
	for _, r := range rows {
		name := strings.TrimSpace(r.RouteName)
		if name == "" {
			continue
		}
		if _, ok := seen[name]; ok {
			continue
		}
		seen[name] = struct{}{}
		out = append(out, name)
	}
	return out, nil
}

// AddMenu 新增菜单。
func (s *MenuManageService) AddMenu(ctx context.Context, in MenuUpsertInput) error {
	menu, err := s.buildMenuModel(in, true)
	if err != nil {
		return err
	}
	return s.DB.WithContext(ctx).Create(menu).Error
}

// UpdateMenu 更新菜单。
func (s *MenuManageService) UpdateMenu(ctx context.Context, in MenuUpsertInput) error {
	if in.ID == 0 {
		return fmt.Errorf("%w: id 不能为空", ErrManageInvalid)
	}
	var existing models.Menu
	if err := s.DB.WithContext(ctx).First(&existing, in.ID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return ErrManageNotFound
		}
		return err
	}
	if in.ParentID == in.ID {
		return fmt.Errorf("%w: parentId 不能等于自身", ErrManageInvalid)
	}
	menu, err := s.buildMenuModel(in, false)
	if err != nil {
		return err
	}
	menu.ID = in.ID
	menu.CreatedAt = existing.CreatedAt
	menu.CreateBy = existing.CreateBy
	return s.DB.WithContext(ctx).Save(menu).Error
}

func (s *MenuManageService) buildMenuModel(in MenuUpsertInput, isCreate bool) (*models.Menu, error) {
	menuType := strings.TrimSpace(in.MenuType)
	menuName := strings.TrimSpace(in.MenuName)
	routeName := strings.TrimSpace(in.RouteName)
	if menuType == "" || menuName == "" || routeName == "" {
		return nil, fmt.Errorf("%w: menuType/menuName/routeName 不能为空", ErrManageInvalid)
	}
	if menuType != "1" && menuType != "2" {
		return nil, fmt.Errorf("%w: menuType 仅支持 1 或 2", ErrManageInvalid)
	}
	st, err := parseStatus(in.Status)
	if err != nil {
		return nil, err
	}
	iconType := strings.TrimSpace(in.IconType)
	if iconType == "" {
		iconType = "1"
	}
	buttons := in.Buttons
	if buttons == nil {
		buttons = []models.MenuButton{}
	}
	query := in.Query
	if query == nil {
		query = []models.MenuQuery{}
	}
	menu := &models.Menu{
		ParentID:        in.ParentID,
		MenuType:        menuType,
		MenuName:        menuName,
		RouteName:       routeName,
		RoutePath:       strings.TrimSpace(in.RoutePath),
		Component:       strings.TrimSpace(in.Component),
		Icon:            strings.TrimSpace(in.Icon),
		IconType:        iconType,
		Buttons:         model.JSONSlice[models.MenuButton](buttons),
		Query:           model.JSONSlice[models.MenuQuery](query),
		I18nKey:         in.I18nKey,
		KeepAlive:       in.KeepAlive,
		Constant:        in.Constant,
		Order:           in.Order,
		Href:            in.Href,
		HideInMenu:      in.HideInMenu,
		ActiveMenu:      in.ActiveMenu,
		MultiTab:        in.MultiTab,
		FixedIndexInTab: in.FixedIndexInTab,
		Status:          st,
		UpdateBy:        in.Operator,
	}
	if isCreate {
		menu.CreateBy = in.Operator
	}
	return menu, nil
}

// DeleteMenus 删除菜单及其子孙，并清理 role_menus。
func (s *MenuManageService) DeleteMenus(ctx context.Context, ids []uint) error {
	if len(ids) == 0 {
		return fmt.Errorf("%w: ids 不能为空", ErrManageInvalid)
	}
	var all []models.Menu
	if err := s.DB.WithContext(ctx).Select("id", "parent_id").Find(&all).Error; err != nil {
		return err
	}
	childrenMap := map[uint][]uint{}
	for _, m := range all {
		childrenMap[m.ParentID] = append(childrenMap[m.ParentID], m.ID)
	}
	toDelete := map[uint]struct{}{}
	var walk func(id uint)
	walk = func(id uint) {
		if _, ok := toDelete[id]; ok {
			return
		}
		toDelete[id] = struct{}{}
		for _, kid := range childrenMap[id] {
			walk(kid)
		}
	}
	for _, id := range ids {
		walk(id)
	}
	if len(toDelete) == 0 {
		return ErrManageNotFound
	}
	deleteIDs := make([]uint, 0, len(toDelete))
	for id := range toDelete {
		deleteIDs = append(deleteIDs, id)
	}
	return s.DB.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := tx.Where("menu_id IN ?", deleteIDs).Delete(&models.RoleMenu{}).Error; err != nil {
			return err
		}
		return tx.Where("id IN ?", deleteIDs).Delete(&models.Menu{}).Error
	})
}
