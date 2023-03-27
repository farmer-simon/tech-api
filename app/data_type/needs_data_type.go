package data_type

type NeedsQuery struct {
	QueryType string  `json:"query_type"`
	MembersId int     `json:"members_id"`
	Cid       float64 `json:"cid" form:"cid"`
	Sort      string  `json:"sort" form:"sort"`

	Phone     string  `json:"phone" form:"phone"`
	Name      string  `json:"name" form:"name"`
	Title     string  `json:"title" form:"title"`
	Recommend int     `json:"recommend" form:"recommend"`
	State     float64 `json:"state" form:"state" binding:"min=0,max=3"`
	StartTime string  `json:"start_time" form:"start_time"`
	EndTime   string  `json:"end_time" form:"end_time"`
}

type NeedsParam struct {
	CategoryId int     `json:"category_id" form:"category_id"  binding:"required"`
	Title      string  `json:"title" form:"title"  binding:"required"`
	Content    string  `json:"content" form:"content"  binding:"required"`
	Price      float64 `json:"price" form:"price" `
	Phone      string  `json:"phone" form:"phone"  binding:"required"`
	Wechat     string  `json:"wechat" form:"wechat"`
	Qq         string  `json:"qq" form:"qq"`
	ExpireTime string  `json:"expire_time" form:"expire_time"  binding:"required"`
}

type NeedsMod struct {
	MembersId    int     `json:"members_id"  gorm:"column:members_id" `
	CategoryId   int     `json:"category_id" gorm:"column:category_id"`
	Title        string  `json:"title" gorm:"column:title"`
	Content      string  `json:"content" gorm:"column:content"`
	Price        float64 `json:"price" gorm:"column:price" `
	Phone        string  `json:"phone" gorm:"column:phone"`
	Qq           string  `json:"qq" gorm:"column:qq"`
	Wechat       string  `json:"wechat" gorm:"column:wechat"`
	ExpireTime   int64   `json:"expire_time" gorm:"column:expire_time"`
	Recommend    int     `json:"recommend" gorm:"column:recommend"`
	State        float64 `json:"state" gorm:"column:state"`
	RejectReason string  `json:"reject_reason" gorm:"column:reject_reason"`
	Hits         float64 `json:"hits" gorm:"column:hits"`
	NickName     string  `json:"nick_name" gorm:"column:nick_name"`
	Avatar       string  `json:"avatar" gorm:"column:avatar"`
	CateName     string  `json:"cate_name" gorm:"column:cate_name"`
}
