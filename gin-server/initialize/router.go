package initialize

import (
	"github.com/gin-gonic/gin"
	ginSwagger "github.com/swaggo/gin-swagger"
	"github.com/swaggo/gin-swagger/swaggerFiles"
	"slb-admin/router"
)

// 初始化总路由

func Routers() *gin.Engine {
	var Router = gin.Default()
	Router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	ApiGroup := Router.Group("")
	router.InitRootRouter(ApiGroup) // 注册用户路由
	router.InitUserRouter(ApiGroup) // 注册基础功能路由 不做鉴权

	return Router
}
