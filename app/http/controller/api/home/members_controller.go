package home

import (
	"github.com/gin-gonic/gin"
	"goskeleton/app/global/consts"
	"goskeleton/app/model"
	"goskeleton/app/service/members"
	"goskeleton/app/service/sms"
	"goskeleton/app/utils/cur_userinfo"
	"goskeleton/app/utils/response"
)

type Members struct {
}

// Info 当前用户资料返回
func (m *Members) Info(ctx *gin.Context) {
	membersId, _ := cur_userinfo.GetHomeCurrentUserId(ctx)
	info := model.CreateMemberFactory("").GetById(int(membersId))
	if info == nil {
		response.Fail(ctx, consts.CurdSelectFailCode, "查询资料失败", "")
		return
	}
	rs := []rune(info.Phone)
	phone := string(rs[0:3]) + "******" + string(rs[9:11])

	response.Success(ctx, "success", gin.H{
		"info": gin.H{
			"data": gin.H{
				"id":         info.Id,
				"phone":      phone,
				"nick_name":  info.NickName,
				"real_name":  info.RealName,
				"team_name":  info.TeamName,
				"team_major": info.TeamMajor,
				"team_intro": info.TeamIntro,
				"qq":         info.Qq,
				"wechat":     info.WeChat,
				"avatar":     info.Avatar,
				"state":      info.State,
			},
		},
	})
}

//SaveInfo 更新个人资料
func (m *Members) SaveInfo(ctx *gin.Context) {
	err := members.CreateMemberServiceFactory().UpdateExtendInfo(ctx)
	if err != nil {
		response.Fail(ctx, consts.CurdSelectFailCode, err.Error(), "")
		return
	}
	response.Success(ctx, "资料更新成功", gin.H{})
}

//GetPhoneCode 用户获取验证码
func (m *Members) GetPhoneCode(ctx *gin.Context) {

	member, _ := cur_userinfo.GetHomeCurrentUserInfo(ctx)
	smsClient := sms.CreatePhoneCodeFactory(member.Phone)
	err := smsClient.SendPhoneCode()
	if err != nil {
		response.Fail(ctx, consts.CurdSelectFailCode, "验证码发送失败，"+err.Error(), "")
		return
	}
	response.Success(ctx, "success", gin.H{})
}

//SetPasswd 设置密码
func (m *Members) SetPasswd(ctx *gin.Context) {
	code := ctx.GetString(consts.ValidatorPrefix + "code")
	passwd := ctx.GetString(consts.ValidatorPrefix + "passwd")
	memberId, _ := cur_userinfo.GetHomeCurrentUserId(ctx)
	err := members.CreateMemberServiceFactory().SetPasswd(int(memberId), code, passwd)
	if err != nil {
		response.Fail(ctx, consts.CurdUpdateFailCode, "设置密码失败，"+err.Error(), "")
		return
	}
	response.Success(ctx, "success", gin.H{})
}

//ChangePasswd 修改密码
func (m *Members) ChangePasswd(ctx *gin.Context) {
	OldPasswd := ctx.GetString(consts.ValidatorPrefix + "passwd")
	NewPasswd := ctx.GetString(consts.ValidatorPrefix + "new_passwd")
	memberId, _ := cur_userinfo.GetHomeCurrentUserId(ctx)
	err := members.CreateMemberServiceFactory().ChangePasswd(int(memberId), OldPasswd, NewPasswd)
	if err != nil {
		response.Fail(ctx, consts.CurdUpdateFailCode, "设置密码失败，"+err.Error(), "")
		return
	}
	response.Success(ctx, "success", gin.H{})
}

//BindPhone 绑定手机
func (m *Members) BindPhone(ctx *gin.Context) {
	phone := ctx.GetString(consts.ValidatorPrefix + "phone")
	code := ctx.GetString(consts.ValidatorPrefix + "code")
	memberId, _ := cur_userinfo.GetHomeCurrentUserId(ctx)
	err := members.CreateMemberServiceFactory().BindPhone(int(memberId), phone, code)
	if err != nil {
		response.Fail(ctx, -400101, "绑定失败，"+err.Error(), gin.H{})
		return
	}
	response.Success(ctx, "手机绑定成功", gin.H{})
}

//ChangePhone 更绑手机
func (m *Members) ChangePhone(ctx *gin.Context) {
	oldPhone := ctx.GetString(consts.ValidatorPrefix + "old_phone")
	phone := ctx.GetString(consts.ValidatorPrefix + "phone")
	code := ctx.GetString(consts.ValidatorPrefix + "code")
	memberId, _ := cur_userinfo.GetHomeCurrentUserId(ctx)
	err := members.CreateMemberServiceFactory().ChangePhone(int(memberId), oldPhone, phone, code)
	if err != nil {
		response.Fail(ctx, -400101, "绑定失败，"+err.Error(), gin.H{})
		return
	}
	response.Success(ctx, "手机绑定成功", gin.H{})
}
