package router

import (
	"github.com/gin-gonic/gin"
	"slb-admin/api/v1"
	"slb-admin/middleware"
)

func InitLoginRouter(Router *gin.RouterGroup) {
	UserRouter := Router.Group("user")
	{
		UserRouter.POST("login", v1.Login)
	}
}

func InitUserRouter(Router *gin.RouterGroup) {
	UserRouter := Router.Group("user").
		Use(middleware.JWTAuth())
	{
		UserRouter.POST("getUserList", v1.GetUserList) // 分页获取用户列表
		UserRouter.POST("setUserRole", v1.SetUserRole) // 设置用户权限
	}
}

func InitMenuRouter(Router *gin.RouterGroup) (R gin.IRoutes) {
	MenuRouter := Router.Group("menu").
		Use(middleware.JWTAuth())
	{
		MenuRouter.POST("getMenu", v1.GetMenu) // 获取菜单树
	}
	return MenuRouter
}
