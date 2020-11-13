package v1

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"slb-admin/model/request"
	"slb-admin/service"
)
type Response struct {
	Code int         `json:"code"`
	Data interface{} `json:"data"`
	Msg  string      `json:"msg"`
}
func response(code int, data interface{}, msg string, c *gin.Context) {
	// 开始时间
	c.JSON(http.StatusOK, Response{
		code,
		data,
		msg,
	})
}
func GetMenu(c *gin.Context) {
	claims, _ := c.Get("claims")
	waitUse := claims.(*request.CustomClaims)
	err, menus := service.GetMenuTree(waitUse.AuthorityId)
	if err != nil {
		response(1, "", "fail",c)
	} else {
		response(0, menus, "操作成功",c)
	}


}
