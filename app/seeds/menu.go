package seeds

import (
	"fmt"

	"clash-admin/app/models"
	"clash-admin/internal/base/model"

	"gorm.io/gorm"
)

// AutoMigrateMenus 迁移菜单表。
func AutoMigrateMenus(db *gorm.DB) error {
	return db.AutoMigrate(&models.Menu{})
}

func boolPtr(v bool) *bool    { return &v }
func intPtr(v int) *int       { return &v }
func strPtr(v string) *string { return &v }

// SeedMenus 幂等写入 manage 相关菜单树，并为 R_SUPER / R_ADMIN 绑定权限。
func SeedMenus(db *gorm.DB) error {
	var count int64
	if err := db.Model(&models.Menu{}).Count(&count).Error; err != nil {
		return err
	}
	if count > 0 {
		return bindRoleMenusIfNeeded(db)
	}

	homeOrder := 1
	manageOrder := 9
	userOrder := 1
	roleOrder := 2
	menuOrder := 3

	home := models.Menu{
		ParentID:  0,
		MenuType:  "2",
		MenuName:  "首页",
		RouteName: "home",
		RoutePath: "/home",
		Component: "layout.base$view.home",
		Icon:      "mdi:monitor-dashboard",
		IconType:  "1",
		Buttons:   model.JSONSlice[models.MenuButton]{},
		Query:     model.JSONSlice[models.MenuQuery]{},
		I18nKey:   strPtr("route.home"),
		KeepAlive: boolPtr(true),
		Order:     &homeOrder,
		Status:    1,
		CreateBy:  "system",
		UpdateBy:  "system",
	}
	if err := db.Create(&home).Error; err != nil {
		return fmt.Errorf("create home menu: %w", err)
	}

	manage := models.Menu{
		ParentID:  0,
		MenuType:  "1",
		MenuName:  "系统管理",
		RouteName: "manage",
		RoutePath: "/manage",
		Component: "layout.base",
		Icon:      "carbon:cloud-service-management",
		IconType:  "1",
		Buttons:   model.JSONSlice[models.MenuButton]{},
		Query:     model.JSONSlice[models.MenuQuery]{},
		I18nKey:   strPtr("route.manage"),
		Order:     &manageOrder,
		Status:    1,
		CreateBy:  "system",
		UpdateBy:  "system",
	}
	if err := db.Create(&manage).Error; err != nil {
		return fmt.Errorf("create manage menu: %w", err)
	}

	children := []models.Menu{
		{
			ParentID:  manage.ID,
			MenuType:  "2",
			MenuName:  "用户管理",
			RouteName: "manage_user",
			RoutePath: "/manage/user",
			Component: "view.manage_user",
			Icon:      "ic:round-manage-accounts",
			IconType:  "1",
			Buttons: model.JSONSlice[models.MenuButton]{
				{Code: "B_CODE_USER_ADD", Desc: "新增用户"},
				{Code: "B_CODE_USER_EDIT", Desc: "编辑用户"},
				{Code: "B_CODE_USER_DELETE", Desc: "删除用户"},
			},
			Query:     model.JSONSlice[models.MenuQuery]{},
			I18nKey:   strPtr("route.manage_user"),
			KeepAlive: boolPtr(false),
			Order:     &userOrder,
			Status:    1,
			CreateBy:  "system",
			UpdateBy:  "system",
		},
		{
			ParentID:  manage.ID,
			MenuType:  "2",
			MenuName:  "角色管理",
			RouteName: "manage_role",
			RoutePath: "/manage/role",
			Component: "view.manage_role",
			Icon:      "carbon:user-role",
			IconType:  "1",
			Buttons: model.JSONSlice[models.MenuButton]{
				{Code: "B_CODE_ROLE_ADD", Desc: "新增角色"},
				{Code: "B_CODE_ROLE_EDIT", Desc: "编辑角色"},
				{Code: "B_CODE_ROLE_DELETE", Desc: "删除角色"},
			},
			Query:    model.JSONSlice[models.MenuQuery]{},
			I18nKey:  strPtr("route.manage_role"),
			Order:    &roleOrder,
			Status:   1,
			CreateBy: "system",
			UpdateBy: "system",
		},
		{
			ParentID:  manage.ID,
			MenuType:  "2",
			MenuName:  "菜单管理",
			RouteName: "manage_menu",
			RoutePath: "/manage/menu",
			Component: "view.manage_menu",
			Icon:      "material-symbols:route",
			IconType:  "1",
			Buttons: model.JSONSlice[models.MenuButton]{
				{Code: "B_CODE_MENU_ADD", Desc: "新增菜单"},
				{Code: "B_CODE_MENU_EDIT", Desc: "编辑菜单"},
				{Code: "B_CODE_MENU_DELETE", Desc: "删除菜单"},
			},
			Query:     model.JSONSlice[models.MenuQuery]{},
			I18nKey:   strPtr("route.manage_menu"),
			KeepAlive: boolPtr(true),
			Order:     &menuOrder,
			Status:    1,
			CreateBy:  "system",
			UpdateBy:  "system",
		},
		{
			ParentID:   manage.ID,
			MenuType:   "2",
			MenuName:   "用户详情",
			RouteName:  "manage_user-detail",
			RoutePath:  "/manage/user-detail/:id",
			Component:  "view.manage_user-detail",
			Icon:       "",
			IconType:   "1",
			Buttons:    model.JSONSlice[models.MenuButton]{},
			Query:      model.JSONSlice[models.MenuQuery]{},
			I18nKey:    strPtr("route.manage_user-detail"),
			HideInMenu: boolPtr(true),
			ActiveMenu: strPtr("manage_user"),
			Order:      intPtr(4),
			Status:     1,
			CreateBy:   "system",
			UpdateBy:   "system",
		},
	}
	for i := range children {
		if err := db.Create(&children[i]).Error; err != nil {
			return fmt.Errorf("create child menu %s: %w", children[i].RouteName, err)
		}
	}

	return bindRoleMenusIfNeeded(db)
}

func bindRoleMenusIfNeeded(db *gorm.DB) error {
	var roleMenuCount int64
	if err := db.Model(&models.RoleMenu{}).Count(&roleMenuCount).Error; err != nil {
		return err
	}
	if roleMenuCount > 0 {
		return nil
	}

	var roles []models.Role
	if err := db.Where("role_code IN ?", []string{"R_SUPER", "R_ADMIN"}).Find(&roles).Error; err != nil {
		return err
	}
	var menus []models.Menu
	if err := db.Find(&menus).Error; err != nil {
		return err
	}
	if len(roles) == 0 || len(menus) == 0 {
		return nil
	}

	var roleMenus []models.RoleMenu
	var roleButtons []models.RoleButton
	for _, role := range roles {
		for _, menu := range menus {
			roleMenus = append(roleMenus, models.RoleMenu{RoleID: role.ID, MenuID: menu.ID})
			for _, btn := range menu.Buttons {
				roleButtons = append(roleButtons, models.RoleButton{RoleID: role.ID, ButtonCode: btn.Code})
			}
		}
	}
	if len(roleMenus) > 0 {
		if err := db.Create(&roleMenus).Error; err != nil {
			return fmt.Errorf("seed role menus: %w", err)
		}
	}
	if len(roleButtons) > 0 {
		if err := db.Create(&roleButtons).Error; err != nil {
			return fmt.Errorf("seed role buttons: %w", err)
		}
	}
	return nil
}
