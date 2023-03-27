package data_type

// UsersBase 后台操作员
type UsersBase struct {
	UserName string `json:"user_name" gorm:"column:user_name" form:"user_name" binding:"required"`
	Name     string `json:"name" gorm:"column:name" form:"name" binding:"required"`
	Phone    string `json:"phone" gorm:"column:phone" form:"phone" binding:"required"`
	Passwd   string `gorm:"column:passwd" json:"passwd" form:"passwd"`
	Avatar   string `json:"avatar" gorm:"column:avatar" form:"avatar" binding:"required"`
	State    int    `json:"state" gorm:"column:state" form:"state" binding:"required,min=1,max=2"`
}
