package service

import (
	"slb-admin/global"
	"slb-admin/model"
)

func getMenuTreeMap(authorityId string) (err error, treeMap map[string][]model.Menu) {
	var allMenus []model.Menu
	treeMap = make(map[string][]model.Menu)
	err = global.GVA_DB.Where("role_id = ?", authorityId).Order("sort").Preload("Parameters").Find(&allMenus).Error
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

func GetMenuTree(authorityId string) (err error, menus []model.Menu) {
	err, menuTree := getMenuTreeMap(authorityId)
	menus = menuTree["0"]
	for i := 0; i < len(menus); i++ {
		err = getChildrenList(&menus[i], menuTree)
	}
	return err, menus
}
