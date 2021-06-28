package v1

//
//import (
//	"github.com/gin-gonic/gin"
//	"slb-admin/config"
//	"slb-admin/global/response"
//)
//
//func GetConfig(c *gin.Context) {
//
//	response.OkWithDetailed(config.Server, "获取成功", c)
//}
//
//// @Tags System
//// @Summary 设置配置文件内容
//// @Security ApiKeyAuth
//// @Produce  application/json
//// @Param data body model.System true "设置配置文件内容"
//// @Success 200 {string} string "{"success":true,"data":{},"msg":"设置成功"}"
//// @Router /system/setSystemConfig [post]
//func SetConfig(c *gin.Context) {
//	var sys model.System
//	_ = c.ShouldBindJSON(&sys)
//	if err := service.SetSystemConfig(sys); err != nil {
//		response.FailWithMessage("设置失败", c)
//	} else {
//		response.OkWithData("设置成功", c)
//	}
//}
