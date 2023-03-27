package data_type

type MembersQuery struct {
	Phone     string  `json:"phone" form:"phone"`
	NickName  string  `json:"nick_name" form:"nick_name"`
	RealName  string  `json:"real_name" form:"real_name"`
	TeamName  string  `json:"team_name" form:"team_name"`
	TeamMajor string  `json:"team_major" form:"team_major"`
	State     float64 `json:"state" form:"state" binding:"min=0,max=2"`
	StartTime string  `json:"start_time" form:"start_time"`
	EndTime   string  `json:"end_time" form:"end_time"`
}

type MembersBase struct {
	Phone        string  `json:"phone" gorm:"column:phone"  form:"phone" binding:"required"`
	Passwd       string  `json:"passwd" gorm:"column:passwd"  form:"passwd"`
	State        float64 `json:"state" gorm:"column:state" form:"state" binding:"min=0,max=2"`
	RejectReason string  `json:"reject_reason" gorm:"column:reject_reason"  form:"reject_reason"`
}

type MembersExtendBase struct {
	Avatar    string `json:"avatar" gorm:"column:avatar" form:"avatar"`
	NickName  string `json:"nick_name" gorm:"column:nick_name" form:"nick_name" binding:"required"`
	RealName  string `json:"real_name" gorm:"column:real_name" form:"real_name" binding:"required"`
	TeamIntro string `json:"team_intro" gorm:"column:team_intro"  form:"team_intro" binding:"required"`
	TeamName  string `json:"team_name" form:"team_name" gorm:"column:team_name" binding:"required"`
	TeamMajor string `json:"team_major" form:"team_major" gorm:"column:team_major" binding:"required"`
	Qq        string `json:"qq" gorm:"column:qq"  form:"qq" binding:"required"`
	WeChat    string `json:"wechat" gorm:"column:wechat"  form:"wechat" binding:"required"`
	SsoUserid string `json:"sso_userid" gorm:"column:sso_userid"`
}
