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
func GetBaseMenuTree(c *gin.Context) {
	type SysBaseMenusResponse struct {
		Menus []model.Menu `json:"menus"`
	}
	err, menus := service.GetBaseMenuTree()
	if err != nil {
		response.FailWithMessage(fmt.Sprintf("获取失败，%v", err), c)
	} else {
		response.OkWithData(SysBaseMenusResponse{Menus: menus}, c)
	}
}

func AddMenuRole(c *gin.Context) {
	type AddMenuRoleInfo struct {
		Menus       []model.Menu
		RoleId string
	}
	var addMenuRoleInfo AddMenuRoleInfo
	_ = c.ShouldBindJSON(&addMenuRoleInfo)

	err := service.AddMenuRole(addMenuRoleInfo.Menus, addMenuRoleInfo.RoleId)
	if err != nil {
		response.FailWithMessage(fmt.Sprintf("添加失败，%v", err), c)
	} else {
		response.OkWithMessage("添加成功", c)
	}
}

func GetMenuRole(c *gin.Context) {
	type RoleIdInfo struct {
		RoleId string
	}
	type SysMenusResponse struct {
		Menus []model.Menu `json:"menus"`
	}
	var roleIdInfo RoleIdInfo
	_ = c.ShouldBindJSON(&roleIdInfo)

	err, menus := service.GetMenuRole(roleIdInfo.RoleId)
	if err != nil {
		response.FailWithDetailed(response.ERROR, SysMenusResponse{Menus: menus}, fmt.Sprintf("添加失败，%v", err), c)
	} else {
		response.Result(response.SUCCESS, gin.H{"menus": menus}, "获取成功", c)
	}
}
