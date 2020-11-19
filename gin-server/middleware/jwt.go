package middleware

import (
	"github.com/gin-gonic/gin"
	"slb-admin/global/response"
	"slb-admin/service"
)

func JWTAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 我们这里jwt鉴权取头部信息 x-token 登录时回返回token信息 这里前端需要把token存储到cookie或者本地localStorage中 不过需要跟后端协商过期时间 可以约定刷新令牌或者重新登录
		token := c.Request.Header.Get("x-token")
		username := c.Request.Header.Get("x-user-name")
		if token == "" || username == "" {
			response.Result(response.ERROR, gin.H{
				"reload": true,
			}, "未登录或非法访问", c)
			return
		}
		err, RedisJwtToken := service.GetRedisJWT(username)
		if err != nil {
			response.Result(response.ERROR, gin.H{
				"reload": true,
			}, "未登录或非法访问", c)
			return
		} else if RedisJwtToken == token {
			// 刷新超时时间
			_ = service.SetRedisJWT(token, username)
			c.Next()
		}
	}
}
