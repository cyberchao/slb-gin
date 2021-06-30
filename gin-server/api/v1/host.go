package v1

import (
	"encoding/json"
	"errors"
	"fmt"
	"slb-admin/global"
	"slb-admin/global/response"
	"slb-admin/model"
	resp "slb-admin/model/response"
	"slb-admin/service"
	"strings"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// @Tags nginx_host
// @Summary 新建主机
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body model.Host true "新建主机"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"新增host成功"}"
// @Router /role/deleteRole [post]
func CreateHost(c *gin.Context) {
	var host model.Host
	_ = c.ShouldBindJSON(&host)

	if !errors.Is(global.DB.Where("ip = ?", host.Ip).First(&host).Error, gorm.ErrRecordNotFound) {
		response.FailWithMessage("存在相同ip", c)
	}
	err := global.DB.Create(&host).Error
	if err != nil {
		response.FailWithMessage(fmt.Sprintf("新增host失败，%v", err), c)
	} else {
		response.OkWithMessage("新增host成功", c)
	}
}

// @Tags nginx_host
// @Summary nginx语法检查
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Success 200 {string} string "{"success":true,"data":{},"msg":"nginx语法检查成功"}"
// @Router /role/deleteRole [post]
func CheckHost(c *gin.Context) {
	var postData []string
	if err := c.ShouldBindJSON(&postData); err != nil {
		global.Logger.Errorf("c.ShouldBindJSON failed. err: [%s]", err.Error())
		response.FailWithMessage("请求数据异常", c)
	}
	res, err := service.SshWoker(postData, "/wls/openresty/nginx/sbin/nginx -t")
	if err != nil {
		global.Logger.Errorf("ssh err. err: [%s]", err.Error())
		response.FailWithMessage("内部错误", c)
	}
	for ip, out := range res {
		if strings.Contains(out, "syntax is ok") {
			delete(res, ip)
		}
	}
	if len(res) == 0 {

		global.Logger.Infof("ssh res: [%s]", res)
		response.OkWithMessage("检查通过", c)
	} else {
		jsonString, _ := json.Marshal(res)
		global.Logger.Infof("ssh res: [%s]", res)
		response.FailWithMessage(string(jsonString), c)
	}
}

// @Tags nginx_host
// @Summary nginx reload
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Success 200 {string} string "{"success":true,"data":{},"msg":"nginx reload成功"}"
// @Router /role/deleteRole [post]
func ReloadHost(c *gin.Context) {
	var postData []string
	if err := c.ShouldBindJSON(&postData); err != nil {
		global.Logger.Errorf("c.ShouldBindJSON failed. err: [%s]", err.Error())
		response.FailWithMessage("请求数据异常", c)
	}
	res, err := service.SshWoker(postData, "/wls/openresty/nginx/sbin/nginx -s reload")
	if err != nil {
		global.Logger.Errorf("ssh err. err: [%s]", err.Error())
		response.FailWithMessage("内部错误", c)
	}
	for ip, out := range res {
		if out == "" {
			delete(res, ip)
		}
	}
	if len(res) == 0 {
		global.Logger.Infof("ssh res: [%s]", res)
		response.OkWithMessage("Reload 成功", c)
	} else {
		jsonString, _ := json.Marshal(res)
		global.Logger.Infof("ssh res: [%s]", res)
		response.FailWithMessage(string(jsonString), c)
	}
}

// @Tags 删除主机
// @Summary nginx reload
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Success 200 {string} string "{"success":true,"data":{},"msg":"删除主机成功"}"
// @Router /role/deleteRole [post]
func DeleteHost(c *gin.Context) {
	type requestData struct {
		Id int `json:"id"`
	}
	var postData requestData
	if err := c.ShouldBindJSON(&postData); err != nil {
		response.FailWithMessage("请求数据异常", c)
	}
	db := global.DB
	db.Delete(&model.Host{}, postData.Id)
	response.OkWithMessage("删除成功", c)
}

// @Tags nginx_host
// @Summary 获取主机列表
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Success 200 {string} string "{"success":true,"data":{},"msg":"获取主机列表成功"}"
// @Router /role/deleteRole [post]
func GetHostList(c *gin.Context) {
	type requestData struct {
		Page     int      `json:"page"`
		PageSize int      `json:"pageSize"`
		Env      []string `json:"env"`
		Cluster  []string `json:"cluster"`
		Ip       string   `json:"ip"`
	}
	var postData requestData
	_ = c.ShouldBindJSON(&postData)

	limit := postData.PageSize
	offset := postData.PageSize * (postData.Page - 1)
	db := global.DB

	//根据查询条件生成where参数
	query := ""
	queryArgs := make([]interface{}, 0)
	if len(postData.Cluster) != 0 {
		query += "cluster IN ?"
		queryArgs = append(queryArgs, postData.Cluster)
	}
	if len(postData.Env) != 0 {
		if query == "" {
			query += "env IN ?"
		} else {
			query += " AND env IN ?"
		}
		queryArgs = append(queryArgs, postData.Env)
	}
	if postData.Ip != "" {
		if query == "" {
			query += "ip LIKE ?"
		} else {
			query += " AND ip LIKE ?"
		}
		queryArgs = append(queryArgs, "%"+postData.Ip+"%")
	}
	var count int64
	var hosts []model.Host
	if query != "" {
		db.Limit(limit).Offset(offset).Where(query, queryArgs...).Find(&hosts).Count(&count)
		_ = db.Limit(limit).Offset(offset).Where(query, queryArgs...).Find(&hosts).Error
	} else {
		db.Find(&hosts).Count(&count)
		_ = db.Limit(limit).Offset(offset).Find(&hosts).Error
	}

	response.OkWithData(resp.PageResult{
		List:     hosts,
		Total:    count,
		Page:     postData.Page,
		PageSize: postData.PageSize,
	}, c)
}
