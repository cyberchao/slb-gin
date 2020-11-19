package v1

import (
	"fmt"
	"slb-admin/global/response"
	"slb-admin/model"
	"slb-admin/model/request"
	resp "slb-admin/model/response"
	"slb-admin/service"
	"slb-admin/utils"
	"github.com/gin-gonic/gin"
)

// @Tags authority
// @Summary 创建角色
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body model.SysRole true "创建角色"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"获取成功"}"
// @Router /authority/createRole [post]
func CreateRole(c *gin.Context) {
	var auth model.Role
	_ = c.ShouldBindJSON(&auth)

	err, authBack := service.CreateRole(auth)
	if err != nil {
		response.FailWithMessage(fmt.Sprintf("创建失败，%v", err), c)
	} else {
		response.OkWithData(resp.SysRoleResponse{Role: authBack}, c)
	}
}

// @Tags authority
// @Summary 拷贝角色
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body response.SysRoleCopyResponse true "拷贝角色"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"拷贝成功"}"
// @Router /authority/copyRole [post]
func CopyRole(c *gin.Context) {
	var copyInfo resp.SysRoleCopyResponse
	_ = c.ShouldBindJSON(&copyInfo)
	OldRoleVerify := utils.Rules{
		"OldRoleId": {utils.NotEmpty()},
	}
	OldRoleVerifyErr := utils.Verify(copyInfo, OldRoleVerify)
	if OldRoleVerifyErr != nil {
		response.FailWithMessage(OldRoleVerifyErr.Error(), c)
		return
	}
	RoleVerify := utils.Rules{
		"RoleId":   {utils.NotEmpty()},
		"RoleName": {utils.NotEmpty()},
		"ParentId":      {utils.NotEmpty()},
	}
	RoleVerifyErr := utils.Verify(copyInfo.Role, RoleVerify)
	if RoleVerifyErr != nil {
		response.FailWithMessage(RoleVerifyErr.Error(), c)
		return
	}
	err, authBack := service.CopyRole(copyInfo)
	if err != nil {
		response.FailWithMessage(fmt.Sprintf("拷贝失败，%v", err), c)
	} else {
		response.OkWithData(resp.SysRoleResponse{Role: authBack}, c)
	}
}

// @Tags authority
// @Summary 删除角色
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body model.SysRole true "删除角色"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"获取成功"}"
// @Router /authority/deleteRole [post]
func DeleteRole(c *gin.Context) {
	var a model.SysRole
	_ = c.ShouldBindJSON(&a)
	RoleIdVerifyErr := utils.Verify(a, utils.CustomizeMap["RoleIdVerify"])
	if RoleIdVerifyErr != nil {
		response.FailWithMessage(RoleIdVerifyErr.Error(), c)
		return
	}
	// 删除角色之前需要判断是否有用户正在使用此角色
	err := service.DeleteRole(&a)
	if err != nil {
		response.FailWithMessage(fmt.Sprintf("删除失败，%v", err), c)
	} else {
		response.OkWithMessage("删除成功", c)
	}
}

// @Tags authority
// @Summary 设置角色资源权限
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body model.SysRole true "设置角色资源权限"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"设置成功"}"
// @Router /authority/updateRole [post]
func UpdateRole(c *gin.Context) {
	var auth model.SysRole
	_ = c.ShouldBindJSON(&auth)
	RoleVerify := utils.Rules{
		"RoleId":   {utils.NotEmpty()},
		"RoleName": {utils.NotEmpty()},
		"ParentId":      {utils.NotEmpty()},
	}
	RoleVerifyErr := utils.Verify(auth, RoleVerify)
	if RoleVerifyErr != nil {
		response.FailWithMessage(RoleVerifyErr.Error(), c)
		return
	}
	err, authority := service.UpdateRole(auth)
	if err != nil {
		response.FailWithMessage(fmt.Sprintf("更新失败，%v", err), c)
	} else {
		response.OkWithData(resp.SysRoleResponse{Role: authority}, c)
	}
}

// @Tags authority
// @Summary 分页获取角色列表
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body request.PageInfo true "分页获取用户列表"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"获取成功"}"
// @Router /authority/getRoleList [post]
func GetRoleList(c *gin.Context) {
	var pageInfo request.PageInfo
	_ = c.ShouldBindJSON(&pageInfo)
	PageVerifyErr := utils.Verify(pageInfo, utils.CustomizeMap["PageVerify"])
	if PageVerifyErr != nil {
		response.FailWithMessage(PageVerifyErr.Error(), c)
		return
	}
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

// @Tags authority
// @Summary 设置角色资源权
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body model.SysRole true "设置角色资源权限"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"设置成功"}"
// @Router /authority/setDataRole [post]
func SetDataRole(c *gin.Context) {
	var auth model.SysRole
	_ = c.ShouldBindJSON(&auth)
	RoleIdVerifyErr := utils.Verify(auth, utils.CustomizeMap["RoleIdVerify"])
	if RoleIdVerifyErr != nil {
		response.FailWithMessage(RoleIdVerifyErr.Error(), c)
		return
	}
	err := service.SetDataRole(auth)
	if err != nil {
		response.FailWithMessage(fmt.Sprintf("设置关联失败，%v", err), c)
	} else {
		response.Ok(c)
	}
}
