package data_type

type CategoryBase struct {
	Name      string  `gorm:"column:name" json:"name" form:"name" binding:"required"`
	Desc      string  `gorm:"column:desc" json:"desc" form:"desc" binding:"required"`
	ParentId  float64 `gorm:"column:parent_id" json:"parent_id" form:"parent_id" binding:"min=0"`
	State     float64 `gorm:"column:state" json:"state" form:"state" binding:"min=1,max=2"`
	Sort      int     `gorm:"column:sort" json:"sort" form:"sort" binding:"min=0,max=999"`
	Recommend int64   `gorm:"column:recommend" json:"recommend" form:"recommend"`
	Src       string  `gorm:"column:src" json:"src" form:"src"`
	IsDel     int     `gorm:"column:is_del" json:"is_del" form:"is_del"`
}

type CategoryItem struct {
	Id       int64  `json:"id"`
	ParentId int64  `json:"parent_id"`
	Name     string `json:"name"`
}

type CategoryTree struct {
	CategoryItem
	SubTree []CategoryItem `json:"sub_tree"`
}
