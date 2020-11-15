package model

type User struct {
	ID       string `json:"id" gorm:"autoIncrement;comment:用户ID"`
	Username string `json:"userName" gorm:"comment:用户登录名"`
	Role     Role   `json:"role" gorm:"foreignKey:RoleId;references:RoleId;comment:用户角色"`
	RoleId   string `json:"roleId" gorm:"default:888;comment:用户角色ID"`
}
