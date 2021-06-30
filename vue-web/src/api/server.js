import service from '@/utils/request'
// @Tags api
// @Summary 分页获取角色列表
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body modelInterface.PageInfo true "分页获取用户列表"
// @Success 200 {string} json "{"success":true,"data":{},"msg":"获取成功"}"
// @Router /server/getApiList [post]
// {
//  page     int
//	pageSize int
// }
export const getServerList = (data) => {
    console.log('request.data:',data)
    return service({
        url: "/api/getServerList",
        method: 'post',
        data
    })
}


// @Tags Api
// @Summary 创建基础api
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body api.CreateApiParams true "创建api"
// @Success 200 {string} json "{"success":true,"data":{},"msg":"获取成功"}"
// @Router /createApi [post]
export const createServer = (data) => {
    return service({
        url: "/api/createServer",
        method: 'post',
        data
    })
}

// @Tags menu
// @Summary 根据id获取菜单
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body api.GetById true "根据id获取菜单"
// @Success 200 {string} json "{"success":true,"data":{},"msg":"获取成功"}"
// @Router /menu/getApiById [post]
export const getServerById = (data) => {
    return service({
        url: "/api/getServerById",
        method: 'post',
        data
    })
}



// @Tags Api
// @Summary 更新api
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body api.CreateApiParams true "更新api"
// @Success 200 {string} json "{"success":true,"data":{},"msg":"更新成功"}"
// @Router /updateApi [post]
export const updateServer = (data) => {
    return service({
        url: "/api/updateServer",
        method: 'post',
        data
    })
}


// @Tags Api
// @Summary 删除指定api
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body dbModel.Api true "删除api"
// @Success 200 {string} json "{"success":true,"data":{},"msg":"获取成功"}"
// @Router /deleteApi [post]
export const deleteServer = (data) => {
    console.log('request.data:',data)
    return service({
        url: "/api/deleteServer",
        method: 'post',
        data
    })
}

export const publishServer = (data) => {
    console.log('request.data:',data)
    return service({
        url: "/api/publishServer",
        method: 'post',
        data
    })
}