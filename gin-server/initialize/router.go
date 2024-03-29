package initialize

import (
	v1 "slb-admin/api/v1"
	_ "slb-admin/docs"
	"slb-admin/middleware"
	"slb-admin/router"

	"github.com/gin-gonic/gin"
	ginSwagger "github.com/swaggo/gin-swagger"
	"github.com/swaggo/gin-swagger/swaggerFiles"
)

// 初始化总路由

func Routers() *gin.Engine {
	var Router = gin.Default()
	Router.Use(middleware.Cors())
	Router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	Router.GET("ws", v1.RemoteTail)
	ApiGroup := Router.Group("")
	router.InitRootRouter(ApiGroup)  //
	router.InitLoginRouter(ApiGroup) // 注册登录路由 不做鉴权 其它路由全部鉴权
	router.InitUserRouter(ApiGroup)
	router.InitRoleRouter(ApiGroup)
	router.InitMenuRouter(ApiGroup)
	router.InitServiceRouter(ApiGroup)
	return Router
}
