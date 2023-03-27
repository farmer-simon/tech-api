package model

import "github.com/farmer-simon/go-utils"

func CreateSettingsFactory(sqlType string) *SettingsModel {
	return &SettingsModel{BaseModel: BaseModel{DB: UseDbConn(sqlType)}}
}

type SettingsModel struct {
	BaseModel
	Param string `json:"param"  gorm:"column:param"`
	Value string `json:"value"  gorm:"column:value"`
}

// TableName 表名
func (mod *SettingsModel) TableName() string {
	return "tech_settings"
}

func (mod *SettingsModel) SetSettings(needsRecommendKeywords, servicesRecommendKeywords string, userRegisterState, contentPublishState int) bool {

	mod.Table(mod.TableName()).Where("param =?", "needs_recommend_keywords").Update("value", needsRecommendKeywords)
	mod.Table(mod.TableName()).Where("param =?", "services_recommend_keywords").Update("value", servicesRecommendKeywords)
	mod.Table(mod.TableName()).Where("param =?", "user_register_state").Update("value", userRegisterState)
	mod.Table(mod.TableName()).Where("param =?", "content_publish_state").Update("value", contentPublishState)

	return true
}

func (mod *SettingsModel) GetSettings() map[string]interface{} {
	var res []SettingsModel
	var results = make(map[string]interface{}, 0)
	mod.Table(mod.TableName()).Find(&res)
	for _, item := range res {
		results[item.Param] = item.Value
	}
	return results
}

func (mod *SettingsModel) SetTokenCache(token, refreshToken string, expires int64) bool {
	tx := mod.Table(mod.TableName()).Begin()
	tx.Debug().Where("param=?", "token").Update("value", token)
	tx.Debug().Where("param=?", "refresh_token").Update("value", refreshToken)
	tx.Debug().Where("param=?", "expires").Update("value", expires)
	tx.Commit()
	return true
}

func (mod *SettingsModel) GetTokenCache() (token, refreshToken string, expires int64) {
	settings := mod.GetSettings()
	token = settings["token"].(string)
	refreshToken = settings["refresh_token"].(string)
	expires = utils.String2Int64(settings["expires"].(string))
	return
}
