package admin

import (
	"github.com/farmer-simon/go-utils"
	"github.com/gin-gonic/gin"
	"goskeleton/app/global/consts"
	"goskeleton/app/model"
	"goskeleton/app/service/members"
	"goskeleton/app/utils/response"
)

type Members struct {
}

func (c *Members) Verify(ctx *gin.Context) {
	id := ctx.GetFloat64(consts.ValidatorPrefix + "id")
	state := ctx.GetFloat64(consts.ValidatorPrefix + "state")
	rejectReason := ctx.GetString(consts.ValidatorPrefix + "reject_reason")
	pass := model.CreateMemberFactory("").Verify(int(id), int(state), rejectReason)
	if pass {
		response.Success(ctx, consts.CurdStatusOkMsg, "操作成功")
	} else {
		response.Fail(ctx, consts.CurdCreatFailCode, "操作失败", "")
	}
}
func (c *Members) Index(ctx *gin.Context) {
	count, list := members.CreateMemberServiceFactory().Index(ctx)

	response.Success(ctx, "Success", gin.H{
		"count": count,
		"list":  list,
	})
}

// Info 帐号详情
func (c *Members) Info(ctx *gin.Context) {
	id, _ := ctx.GetQuery("id")
	userId := utils.String2Int(id)
	if userId > 0 {
		user := model.CreateMemberFactory("").GetById(userId)
		if user != nil {
			user.Passwd = ""
			response.Success(ctx, "操作成功", user)
		} else {
			response.Fail(ctx, -400300, "用户信息查询失败", "")
		}
		return
	}
	response.Fail(ctx, -400300, "操作失败", "")
}
