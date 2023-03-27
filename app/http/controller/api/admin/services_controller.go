package admin

import (
	"github.com/farmer-simon/go-utils"
	"github.com/gin-gonic/gin"
	"goskeleton/app/data_type"
	"goskeleton/app/global/consts"
	"goskeleton/app/model"
	"goskeleton/app/utils/data_bind"
	"goskeleton/app/utils/response"
	"math"
)

type Services struct {
}

func (m *Services) Index(ctx *gin.Context) {
	page := ctx.GetFloat64(consts.ValidatorPrefix + "page")
	limit := ctx.GetFloat64(consts.ValidatorPrefix + "limit")
	offset := math.Max((page-1)*limit, 0)
	var queryData data_type.ServicesQuery
	if err := data_bind.ShouldBindFormDataToModel(ctx, &queryData); err == nil {
		count, res := model.CreateServicesFactory("").Index(&queryData, int(offset), int(limit))
		response.Success(ctx, "SUCCESS", gin.H{
			"list":  res,
			"count": count,
		})
		return
	}
	response.Fail(ctx, -400100, "参数错误", gin.H{})
}

func (m *Services) Info(ctx *gin.Context) {
	query, _ := ctx.GetQuery("id")
	id := utils.String2Int(query)
	info := model.CreateServicesFactory("").GetById(id)
	if info == nil {
		response.Fail(ctx, -400100, "服务不存在", gin.H{})
		return
	}
	author := model.CreateMemberFactory("").GetById(info.MembersId)
	attrs := model.CreateAttrsFactory("").GetByTargetId(int(info.Id), "services")
	response.Success(ctx, "SUCCESS", gin.H{
		"info":  info,
		"attrs": attrs,
		"author": gin.H{
			"avatar":    author.Avatar,
			"nick_name": author.NickName,
			"id":        author.Id,
		},
	})
}

func (m *Services) Verify(ctx *gin.Context) {
	id := ctx.GetFloat64(consts.ValidatorPrefix + "id")
	state := ctx.GetFloat64(consts.ValidatorPrefix + "state")
	rejectReason := ctx.GetString(consts.ValidatorPrefix + "reject_reason")
	pass := model.CreateServicesFactory("").Verify(int(id), int(state), rejectReason)
	if pass {
		response.Success(ctx, consts.CurdStatusOkMsg, "操作成功")
	} else {
		response.Fail(ctx, consts.CurdCreatFailCode, "操作失败", "")
	}
}

func (m *Services) RecordIndex(ctx *gin.Context) {
	servicesId := ctx.GetFloat64(consts.ValidatorPrefix + "services_id")
	state := ctx.GetFloat64(consts.ValidatorPrefix + "state")
	page := ctx.GetFloat64(consts.ValidatorPrefix + "page")
	limit := ctx.GetFloat64(consts.ValidatorPrefix + "limit")

	count, res := model.CreateServicesRecordFactory("").AdminIndex(int(servicesId), int(state), int(page), int(limit))
	response.Success(ctx, "SUCCESS", gin.H{
		"list":  res,
		"count": count,
	})
}
