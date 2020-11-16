package model

type User struct {
	ID       int `json:"id" gorm:"AUTO_INCREMENT"`
	Username string `json:"userName" gorm:"comment:用户登录名"`
	Role     Role   `json:"role" gorm:"foreignKey:RoleId;references:RoleId;comment:用户角色"`
	RoleId   string `json:"roleId" gorm:"default:888;comment:用户角色ID"`
}
