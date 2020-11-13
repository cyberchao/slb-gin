package model

type SysUser struct {
	ID        string `json:"id" gorm:"comment:用户ID"`
	Username  string `json:"userName" gorm:"comment:用户登录名"`
	Password  string `json:"-"  gorm:"comment:用户登录密码"`
	HeaderImg string `json:"headerImg" gorm:"default:http://qmplusimg.henrongyi.top/head.png;comment:用户头像"`
	Role      Role   `json:"role" gorm:"foreignKey:AuthorityId;references:AuthorityId;comment:用户角色"`
	RoleId    string `json:"roleId" gorm:"default:888;comment:用户角色ID"`
}
