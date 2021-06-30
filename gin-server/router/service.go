package router

import (
	"github.com/gin-gonic/gin"
	"slb-admin/api/v1"
	"slb-admin/middleware"
)

func InitServiceRouter(Router *gin.RouterGroup) {
	UserRouter := Router.Group("/api").
		Use(middleware.JWTAuth())
	{
		//UserRouter.POST("getconfig", v1.GetConfig)
		//UserRouter.POST("setconfig", v1.SetConfig)
		UserRouter.POST("cross_parse", v1.CrossParse)
		UserRouter.POST("createServer", v1.CreateServer)
		UserRouter.POST("updateServer", v1.UpdateServer)
		UserRouter.POST("getServerList", v1.GetServerList)
		UserRouter.POST("deleteServer", v1.DeleteServer)
		UserRouter.POST("publishServer", v1.PublishServer)
		UserRouter.POST("createHost", v1.CreateHost)
		UserRouter.POST("checkHost", v1.CheckHost)
		UserRouter.POST("reloadHost", v1.ReloadHost)
		UserRouter.POST("getHostList", v1.GetHostList)
		UserRouter.POST("deleteHost", v1.DeleteHost)
		UserRouter.POST("getUpstreamList", v1.GetUpstreamList)
		UserRouter.POST("createUpstream", v1.CreateUpstream)
		UserRouter.POST("updateUpstream", v1.UpdateUpstream)
		UserRouter.POST("deleteUpstream", v1.DeleteUpstream)
	}
}
