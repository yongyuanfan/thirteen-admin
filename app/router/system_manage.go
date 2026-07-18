package router

import (
	"clash-admin/app/handlers"

	"github.com/gin-gonic/gin"
)

// RegisterSystemManageRoutes 注册系统管理路由（均需登录）。
func RegisterSystemManageRoutes(h *handlers.SystemManageHandler) {
	RegisterAuthRouter(func(g *gin.RouterGroup, authMW gin.HandlerFunc) {
		group := g.Group("/systemManage")
		if authMW != nil {
			group.Use(authMW)
		}
		{
			group.GET("/getUserList", h.GetUserList)
			group.POST("/addUser", h.AddUser)
			group.POST("/updateUser", h.UpdateUser)
			group.DELETE("/deleteUser", h.DeleteUser)
			group.DELETE("/batchDeleteUser", h.BatchDeleteUser)

			group.GET("/getRoleList", h.GetRoleList)
			group.GET("/getAllRoles", h.GetAllRoles)
			group.POST("/addRole", h.AddRole)
			group.POST("/updateRole", h.UpdateRole)
			group.DELETE("/deleteRole", h.DeleteRole)
			group.DELETE("/batchDeleteRole", h.BatchDeleteRole)
			group.GET("/getRoleMenus", h.GetRoleMenus)
			group.POST("/updateRoleMenus", h.UpdateRoleMenus)
			group.GET("/getRoleHome", h.GetRoleHome)
			group.POST("/updateRoleHome", h.UpdateRoleHome)
			group.GET("/getAllButtons", h.GetAllButtons)
			group.GET("/getRoleButtons", h.GetRoleButtons)
			group.POST("/updateRoleButtons", h.UpdateRoleButtons)

			group.GET("/getMenuList/v2", h.GetMenuList)
			group.GET("/getMenuTree", h.GetMenuTree)
			group.GET("/getAllPages", h.GetAllPages)
			group.POST("/addMenu", h.AddMenu)
			group.POST("/updateMenu", h.UpdateMenu)
			group.DELETE("/deleteMenu", h.DeleteMenu)
			group.DELETE("/batchDeleteMenu", h.BatchDeleteMenu)
		}
	})
}
