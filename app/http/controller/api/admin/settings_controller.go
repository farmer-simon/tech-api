package admin

import (
	"github.com/gin-gonic/gin"
	"goskeleton/app/global/consts"
	"goskeleton/app/model"
	"goskeleton/app/utils/response"
)

type Settings struct {
}

//GetValues 获取配置
func (c *Settings) GetValues(ctx *gin.Context) {
	settings := model.CreateSettingsFactory("").GetSettings()
	response.Success(ctx, consts.CurdStatusOkMsg, gin.H{
		"settings": settings,
	})
}

//SetValues 保存配置
func (c *Settings) SetValues(ctx *gin.Context) {
	needsRecommendKeywords := ctx.GetString(consts.ValidatorPrefix + "NeedsRecommendKeywords")
	servicesRecommendKeywords := ctx.GetString(consts.ValidatorPrefix + "ServicesRecommendKeywords")
	userRegisterState := ctx.GetFloat64(consts.ValidatorPrefix + "UserRegisterState")
	contentPublishState := ctx.GetFloat64(consts.ValidatorPrefix + "ContentPublishState")

	model.CreateSettingsFactory("").SetSettings(needsRecommendKeywords, servicesRecommendKeywords, int(userRegisterState), int(contentPublishState))
	response.Success(ctx, consts.CurdStatusOkMsg, gin.H{})
}
