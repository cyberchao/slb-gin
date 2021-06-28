package request

type PageInfo struct {
	Page     int `json:"page" form:"page"`
	PageSize int `json:"pageSize" form:"pageSize"`
}

type Data struct {
	Data     string `json:"data" form:"data"`
}
// Modify  user's auth structure
type SetUserAuth struct {
	ID        int `json:"id"`
	RoleId string    `json:"RoleId"`
}
