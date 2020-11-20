package router

import (
	"slb-admin/api/v1"
	"slb-admin/middleware"
	"github.com/gin-gonic/gin"
)

func InitRoleRouter(Router *gin.RouterGroup) {
	RoleRouter := Router.Group("role").
		Use(middleware.JWTAuth())
	{
		RoleRouter.POST("createRole", v1.CreateRole)   // 创建角色
		RoleRouter.POST("deleteRole", v1.DeleteRole)   // 删除角色
		RoleRouter.PUT("updateRole", v1.UpdateRole)    // 更新角色
		RoleRouter.POST("getRoleList", v1.GetRoleList) // 获取角色列表
		//RoleRouter.POST("setDataRole", v1.SetDataRole) // 设置角色资源权限
	}
}
