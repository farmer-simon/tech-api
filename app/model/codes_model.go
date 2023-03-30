package model

import "github.com/farmer-simon/go-utils"

func CreateCodesFactory(sqlType string) *CodesModel {
	return &CodesModel{BaseModel: BaseModel{DB: UseDbConn(sqlType)}}
}

type CodesModel struct {
	BaseModel
	Phone      string `json:"phone"  gorm:"column:phone"`
	Code       string `json:"code"  gorm:"column:code"`
	ExpireTime string `json:"expire_time"  gorm:"column:expire_time"`
}

// TableName 表名
func (mod *CodesModel) TableName() string {
	return "tech_codes"
}

//GetCodeByPhone 根据手机号获取验证码
func (mod *CodesModel) GetCodeByPhone(phone string) (code string) {
	currentTime := utils.GetTimeStamp()
	sql := "SELECT code FROM  `tech_codes` WHERE `phone`=? and expire_time>? ORDER BY id DESC LIMIT 1"
	result := mod.Raw(sql, phone, currentTime).First(&code)
	if result.Error == nil {
		return
	}
	return ""
}

func (mod *CodesModel) SetPhoneCode(phone, code string) error {
	expireTime := utils.GetTimeStamp() + 5*60
	sql := "INSERT INTO tech_codes(phone, code, expire_time) VALUES(?,?,?)"
	if res := mod.Exec(sql, phone, code, expireTime); res.Error != nil {
		return res.Error
	}
	return nil
}

func (mod *CodesModel) DeleteCodeByPhone(phone string) error {
	sql := "DELETE FROM tech_codes WHERE phone=?"
	mod.Exec(sql, phone)
	return nil
}

func (mod *CodesModel) DeleteExpireCodes() error {
	currentTime := utils.GetTimeStamp()
	sql := "DELETE FROM tech_codes WHERE expire_time<?"
	if res := mod.Exec(sql, currentTime); res.Error != nil {
		return res.Error
	}
	return nil
}
