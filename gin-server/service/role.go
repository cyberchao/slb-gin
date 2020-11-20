package service

import (
	"errors"
	"gorm.io/gorm"
	"slb-admin/global"
	"slb-admin/model"
	"slb-admin/model/request"
)

// @title    CreateRole
// @description   创建一个角色
// @auth                     （2020/04/05  20:22）
// @param     auth            model.Role
// @return                    error
// @return    role       model.Role

func CreateRole(auth model.Role) (err error, role model.Role) {
	var roleBox model.Role
	if !errors.Is(global.DB.Where("role_id = ?", auth.RoleId).First(&roleBox).Error, gorm.ErrRecordNotFound) {
		return errors.New("存在相同角色id"), auth
	}
	err = global.DB.Create(&auth).Error
	return err, auth
}


// @title    UpdateRole
// @description   更改一个角色
// @auth                     （2020/04/05  20:22）
// @param     auth            model.Role
// @return                    error
// @return    role       model.Role

func UpdateRole(auth model.Role) (err error, role model.Role) {
	err = global.DB.Where("role_id = ?", auth.RoleId).First(&model.Role{}).Updates(&auth).Error
	return err, auth
}

// @title    DeleteRole
// @description   删除角色
// @auth                     （2020/04/05  20:22）
// @param     auth            model.Role
// @return                    error
// 删除角色

func DeleteRole(auth *model.Role) (err error) {
	if !errors.Is(global.DB.Where("role_id = ?", auth.RoleId).First(&model.User{}).Error, gorm.ErrRecordNotFound) {
		return errors.New("此角色有用户正在使用禁止删除")
	}
	if !errors.Is(global.DB.Where("parent_id = ?", auth.RoleId).First(&model.Role{}).Error, gorm.ErrRecordNotFound) {
		return errors.New("此角色存在子角色不允许删除")
	}
	db := global.DB.Preload("Menus").Where("role_id = ?", auth.RoleId).First(auth)
	err = db.Unscoped().Delete(auth).Error
	if len(auth.Menus) > 0 {
		err = global.DB.Model(auth).Association("Menus").Delete(auth.Menus)
		//err = db.Association("Menus").Delete(&auth)
	} else {
		err = db.Error
	}
	return err
}

// @title    GetInfoList
// @description   删除文件切片记录
// @auth                     （2020/04/05  20:22）
// @param     info            request.PaveInfo
// @return                    error
// 分页获取数据

func GetRoleInfoList(info request.PageInfo) (err error, list interface{}, total int64) {
	limit := info.PageSize
	offset := info.PageSize * (info.Page - 1)
	db := global.DB
	var role []model.Role
	err = db.Limit(limit).Offset(offset).Find(&role).Error

	return err, role, total
}

// @title    GetRoleInfo
// @description   获取所有角色信息
// @auth                     （2020/04/05  20:22）
// @param     auth            model.Role
// @return                    error
// @param     role       model.Role

func GetRoleInfo(auth model.Role) (err error, sa model.Role) {
	err = global.DB.Preload("DataRoleId").Where("role_id = ?", auth.RoleId).First(&sa).Error
	return err, sa
}


// @title    SetMenuRole
// @description   菜单与角色绑定
// @auth                     （2020/04/05  20:22）
// @param     auth            *model.Role
// @return                    error

func SetMenuRole(auth *model.Role) error {
	var s model.Role
	global.DB.Preload("Menus").First(&s, "role_id = ?", auth.RoleId)
	err := global.DB.Model(&s).Association("Menus").Replace(&auth.Menus)
	return err
}

