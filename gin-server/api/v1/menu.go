package v1

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"slb-admin/global"
	"slb-admin/global/response"
	"slb-admin/model"
	"slb-admin/service"
)

func GetMenu(c *gin.Context) {
	type MenusResponse struct {
		Menus []model.Menu `json:"menus"`
	}

	username := c.Request.Header.Get("x-user-name")
	var user model.User
	global.DB.Where("username = ?", username).First(&user)
	err, menus := service.GetMenuTree(user.RoleId)
	if err != nil {
		response.FailWithMessage(fmt.Sprintf("获取失败qqq，%v", err), c)
	} else {
		response.OkWithData(MenusResponse{Menus: menus}, c)
	}

}
