package router

import (
	"github.com/gin-gonic/gin"
	"slb-admin/api/v1"
	"slb-admin/middleware"
)

func InitRootRouter(Router *gin.RouterGroup) {
	UserRouter := Router.Group("/").
		Use(middleware.JWTAuth())
	{
		UserRouter.GET("main", v1.Root)
	}
}
