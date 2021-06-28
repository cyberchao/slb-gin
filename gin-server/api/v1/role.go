package v1

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"slb-admin/global/response"
	"slb-admin/model"
	"slb-admin/model/request"
	resp "slb-admin/model/response"
	"slb-admin/service"
)

// @Tags role
// @Summary 创建角色
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body model.Role true "创建角色"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"获取成功"}"
// @Router /role/createRole [post]
func CreateRole(c *gin.Context) {
	var auth model.Role
	_ = c.ShouldBindJSON(&auth)
	type SysRoleResponse struct {
		role model.Role `json:"role"`
	}

	err, authBack := service.CreateRole(auth)
	if err != nil {
		response.FailWithMessage(fmt.Sprintf("创建失败，%v", err), c)
	} else {
		response.OkWithData(SysRoleResponse{role: authBack}, c)
	}
}

// @Tags role
// @Summary 删除角色
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body model.Role true "删除角色"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"获取成功"}"
// @Router /role/deleteRole [post]
func DeleteRole(c *gin.Context) {
	var a model.Role
	_ = c.ShouldBindJSON(&a)
	// 删除角色之前需要判断是否有用户正在使用此角色
	err := service.DeleteRole(&a)
	if err != nil {
		response.FailWithMessage(fmt.Sprintf("删除失败，%v", err), c)
	} else {
		response.OkWithMessage("删除成功", c)
	}
}

// @Tags role
// @Summary 设置角色资源权限
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body model.Role true "设置角色资源权限"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"设置成功"}"
// @Router /role/updateRole [post]
func UpdateRole(c *gin.Context) {
	var auth model.Role
	_ = c.ShouldBindJSON(&auth)
	type SysRoleResponse struct {
		role model.Role `json:"role"`
	}

	err, role := service.UpdateRole(auth)
	if err != nil {
		response.FailWithMessage(fmt.Sprintf("更新失败，%v", err), c)
	} else {
		response.OkWithData(SysRoleResponse{role: role}, c)
	}
}

// @Tags role
// @Summary 分页获取角色列表
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body request.PageInfo true "分页获取用户列表"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"获取成功"}"
// @Router /role/getRoleList [post]
func GetRoleList(c *gin.Context) {
	var pageInfo request.PageInfo
	_ = c.ShouldBindJSON(&pageInfo)

	err, list, total := service.GetRoleInfoList(pageInfo)
	if err != nil {
		response.FailWithMessage(fmt.Sprintf("获取数据失败，%v", err), c)
	} else {
		response.OkWithData(resp.PageResult{
			List:     list,
			Total:    total,
			Page:     pageInfo.Page,
			PageSize: pageInfo.PageSize,
		}, c)
	}
}
