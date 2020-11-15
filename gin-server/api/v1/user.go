package v1

import (
	"errors"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"slb-admin/global"
	"slb-admin/global/response"
	"slb-admin/model"
	"slb-admin/model/request"
	resp "slb-admin/model/response"
	"slb-admin/service"
)

func Login(c *gin.Context) {
	var L request.LoginStruct
	_ = c.ShouldBindJSON(&L)
	err, ssores := service.LoginSSO(L)
	if err != nil {
		response.FailWithMessage("sso连接失败", c)
	}
	returnCode := ssores["returnCode"].(string)
	token := ssores["token"].(string)
	if returnCode == "1" {
		response.FailWithMessage("sso连接成功，获取token失败", c)
	} else {
		if err, u := tokenNext(token, L); err != nil {
			response.FailWithMessage("保存token失败", c)
		} else {
			response.OkWithData(resp.LoginResponse{
				User:  u,
				Token: token,
			}, c)
		}
	}
}

// 保存token，创建用户
func tokenNext(token string, L request.LoginStruct) (err error, us model.User) {
	user := &model.User{Username: L.Username, RoleId: "666"}
	var u model.User
	if err := service.SetRedisJWT(token, user.Username); err != nil {
		return err, u
	}

	// 用户不存在
	if errors.Is(global.DB.Where("username = ?", L.Username).First(&u).Error, gorm.ErrRecordNotFound) {
		err = global.DB.Create(&user).Error
		return err, u
	}
	return nil, u
}
