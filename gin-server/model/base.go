package model

type Role struct {
	RoleId   string `json:"roleId" gorm:"not null;unique;primary_key;comment:角色ID;size:90"`
	RoleName string `json:"roleName" gorm:"comment:角色名"`
	Menus    []Menu `json:"menus" gorm:"many2many:role_menus;"`
}

type User struct {
	ID       int    `json:"id" gorm:"AUTO_INCREMENT"`
	Username string `json:"userName" gorm:"comment:用户登录名"`
	Role     Role   `json:"role" gorm:"foreignKey:RoleId;references:RoleId;comment:用户角色"`
	RoleId   string `json:"roleId" gorm:"default:888;comment:用户角色ID"`
}

type Menu struct {
	ID        string `gorm:"primarykey"`
	MenuLevel uint   `json:"-"`
	ParentId  string `json:"parentId" gorm:"comment:父菜单ID"`
	Path      string `json:"path" gorm:"comment:路由path"`
	Name      string `json:"name" gorm:"comment:路由name"`
	Hidden    bool   `json:"hidden" gorm:"comment:是否在列表隐藏"`
	Component string `json:"component" gorm:"comment:对应前端文件路径"`
	Sort      int    `json:"sort" gorm:"comment:排序标记"`
	KeepAlive bool   `json:"keepAlive" gorm:"comment:是否缓存"`
	Title     string `json:"title" gorm:"comment:菜单名"`
	Icon      string `json:"icon" gorm:"comment:菜单图标"`
	Meta      `json:"meta" gorm:"comment:附加属性"`
	Roles     []Role `json:"roles" gorm:"many2many:role_menus;"`
	Children  []Menu `json:"children" gorm:"-"`
}
type Meta struct {
	Titletab string `json:"title" gorm:"comment:菜单名"`
}
