package model

type Role struct {
	RoleId   string `json:"roleId" gorm:"not null;unique;primary_key;comment:角色ID;size:90"`
	RoleName string `json:"roleName" gorm:"comment:角色名"`
	Menus    []Menu `json:"menus" gorm:"many2many:sys_authority_menus;"`
}
