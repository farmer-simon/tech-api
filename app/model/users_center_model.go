package model

type UserCenter struct {
	BaseModel
	GH       string `json:"xm"`
	XM       string `json:"gh"`
	BMMC     string `json:"bmmc"`
	BMDM     string `json:"bmdm"`
	Uid      int64  `json:"uid"`
	UserName string `json:"user_name"`
	Phone    string `json:"phone"`
}
type userType struct {
	Id        int64  `gorm:"primarykey" json:"id"`
	Phone     string `json:"phone"`
	NickName  string `json:"nick_name"`
	RealName  string `json:"real_name"`
	Passwd    string `json:"passwd"`
	State     int    `json:"state"`
	SsoUserid string `json:"sso_userid"`
	Avatar    string `json:"avatar"`
}

func CreateUserCenterFactory(sqlType string) *UserCenter {
	return &UserCenter{BaseModel: BaseModel{DB: UseDbConn(sqlType)}}
}

func (u *UserCenter) TableName() string {
	return "tech_users_center"
}

func (u *UserCenter) AutoSyncUser(userId string) int64 {
	var user UserCenter
	u.Model(u).Debug().Where("gh = ?", userId).Find(&user)
	if user.Id == 0 {
		return 0
	}
	if user.BMDM == "" || user.BMMC == "" {
		return 0
	}

	tx := u.Begin()
	var tmpUser = userType{
		Phone:     user.Phone,
		NickName:  user.XM,
		RealName:  user.XM,
		Passwd:    "8sd9bb400c0634691f0e3bbbf1e2fd0d",
		State:     3,
		SsoUserid: userId,
		Avatar:    "/public/storage/resources/avatar.png",
	}
	if res := tx.Debug().Table("tech_members").Create(&tmpUser); res.Error != nil {
		tx.Rollback()
		return 0
	}
	tx.Commit()
	return tmpUser.Id
}
