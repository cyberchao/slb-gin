import service from '@/utils/request'

// @Summary 用户登录
// @Produce  application/json
// @Param {
//  page     int
//	pageSize int
// }
// @Router /role/getroleList [post]
export const getRoleList = (data) => {
    return service({
        url: "/role/getRoleList",
        method: 'post',
        data
    })
}


// @Summary 删除角色
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body {roleId uint} true "删除角色"
// @Success 200 {string} json "{"success":true,"data":{},"msg":"获取成功"}"
// @Router /role/deleterole [post]
export const deleteRole = (data) => {
    return service({
        url: "/role/deleteRole",
        method: 'post',
        data
    })
}

// @Summary 创建角色
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body api.CreaterolePatams true "创建角色"
// @Success 200 {string} json "{"success":true,"data":{},"msg":"获取成功"}"
// @Router /role/createrole [post]
export const createRole = (data) => {
    return service({
        url: "/role/createRole",
        method: 'post',
        data
    })
}


// @Summary 设置角色资源权限
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body sysModel.Sysrole true "设置角色资源权限"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"设置成功"}"
// @Router /role/setDatarole [post]
export const setDataRole = (data) => {
    return service({
        url: "/role/setDataRole",
        method: 'post',
        data
    })
}

// @Summary 修改角色
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body model.Sysrole true "修改角色"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"设置成功"}"
// @Router /role/setDatarole [post]
export const updateRole = (data) => {
    return service({
        url: "/role/updateRole",
        method: 'put',
        data
    })
}