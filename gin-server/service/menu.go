package service

import (
	"fmt"
	"slb-admin/global"
	"slb-admin/model"
)

func getMenuTreeMap(roleId string) (err error, treeMap map[string][]model.Menu) {
	var allMenus []model.Menu
	treeMap = make(map[string][]model.Menu)

	var role []model.Role
	global.DB.Where("role_id= ?", roleId).First(&role)
	fmt.Println(role)
	global.DB.Model(&role).Association("Menus").Find(&allMenus)
	//err = global.DB.Where("role_id= ?", "3").Order("sort").Find(&allMenus).Error
	for _, v := range allMenus {
		treeMap[v.ParentId] = append(treeMap[v.ParentId], v)
	}
	return err, treeMap
}

func getChildrenList(menu *model.Menu, treeMap map[string][]model.Menu) (err error) {
	menu.Children = treeMap[menu.ID]
	for i := 0; i < len(menu.Children); i++ {
		err = getChildrenList(&menu.Children[i], treeMap)
	}
	return err
}

func GetMenuTree(roleId string) (err error, menus []model.Menu) {
	err, menuTree := getMenuTreeMap(roleId)
	menus = menuTree["0"]
	for i := 0; i < len(menus); i++ {
		err = getChildrenList(&menus[i], menuTree)
	}
	return err, menus
}

func AddMenuRole(menus []model.Menu, roleId string) (err error) {
	var auth model.Role
	auth.RoleId = roleId
	auth.Menus = menus
	err = SetMenuRole(&auth)
	return err
}

func GetMenuRole(roleId string) (err error, menus []model.Menu) {
	var role []model.Role
	var allMenus []model.Menu
	global.DB.Where("role_id= ?", roleId).First(&role)
	global.DB.Model(&role).Association("Menus").Find(&allMenus)
	//err = global.DB.Raw(sql, roleId).Scan(&menus).Error
	return err, allMenus
}

func getBaseMenuTreeMap() (err error, treeMap map[string][]model.Menu) {
	var allMenus []model.Menu
	treeMap = make(map[string][]model.Menu)
	err = global.DB.Order("sort").Find(&allMenus).Error
	for _, v := range allMenus {
		treeMap[v.ParentId] = append(treeMap[v.ParentId], v)
	}
	return err, treeMap
}

func getBaseChildrenList(menu *model.Menu, treeMap map[string][]model.Menu) (err error) {
	menu.Children = treeMap[menu.ID]
	for i := 0; i < len(menu.Children); i++ {
		err = getBaseChildrenList(&menu.Children[i], treeMap)
	}
	return err
}

func GetBaseMenuTree() (err error, menus []model.Menu) {
	err, treeMap := getBaseMenuTreeMap()
	menus = treeMap["0"]
	for i := 0; i < len(menus); i++ {
		err = getBaseChildrenList(&menus[i], treeMap)
	}
	return err, menus
}