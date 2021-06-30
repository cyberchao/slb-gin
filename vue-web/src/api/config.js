import service from "@/utils/request";

// @Tags systrm
// @Summary 获取配置文件内容
// @Security ApiKeyAuth
// @Produce  application/json
// @Success 200 {string} string "{"success":true,"data":{},"msg":"返回成功"}"
// @Router /system/getSystemConfig [post]
export const getConfig = () => {
    return service({
        url: "/getconfig",
        method: "post",
    });
};

// @Tags system
// @Summary 设置配置文件内容
// @Security ApiKeyAuth
// @Produce  application/json
// @Param data body sysModel.System true
// @Success 200 {string} string "{"success":true,"data":{},"msg":"返回成功"}"
// @Router /system/setSystemConfig [post]
export const setSystemConfig = (data) => {
    return service({
        url: "/setconfig",
        method: "post",
        data,
    });
};

export const getSystemConfig = (data) => {
    return service({
        url: "/getconfig",
        method: "get",
        data,
    });
};