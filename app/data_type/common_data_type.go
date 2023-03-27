package data_type

type Id struct {
	Id float64 `form:"id" json:"id" binding:"required"` // 必填，页面值>0
}
type ParentId struct {
	ParentId float64 `form:"parent_id" json:"parent_id"` // 必填，页面值>0
}

type State struct {
	State float64 `form:"state" json:"state"`
}

type IsDel struct {
	IsDel int64 `gorm:"column:is_del" json:"is_del" form:"is_del"`
}

type Sort struct {
	Sort int `gorm:"column:sort" json:"sort" form:"sort" binding:"min=0,max=999"`
}

type TimeSlot struct {
	StartTime string `gorm:"column:start_time" json:"start_time" form:"start_time" binding:"required"`
	EndTime   string `gorm:"column:end_time" json:"end_time" form:"end_time"`
}

type RecordTotal struct {
	Id    int `gorm:"column:id"`
	Total int `gorm:"column:total"`
}
