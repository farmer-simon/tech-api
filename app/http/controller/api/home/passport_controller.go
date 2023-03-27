package home

import (
	"github.com/gin-gonic/gin"
	"goskeleton/app/global/consts"
	"goskeleton/app/global/variable"
	"goskeleton/app/model"
	"goskeleton/app/service"
	"goskeleton/app/service/members"
	"goskeleton/app/service/sms"
	"goskeleton/app/service/users/token"
	"goskeleton/app/utils/response"
	"goskeleton/app/utils/teach_tool"
)

type Passport struct {
}

//Register 注册
func (c *Passport) Register(ctx *gin.Context) {
	phone := ctx.GetString(consts.ValidatorPrefix + "phone")
	code := ctx.GetString(consts.ValidatorPrefix + "code")
	//检查验证码
	smsClient := sms.CreatePhoneCodeFactory(phone)
	pass := smsClient.CheckPhoneCode(code, true)
	if !pass {
		response.Fail(ctx, -400100, "验证码错误", gin.H{})
		return
	}
	member, err := members.CreateMemberServiceFactory().Register(ctx)
	if err != nil {
		response.Fail(ctx, -400100, "注册失败,"+err.Error(), gin.H{})
		return
	}

	tokenData := c.memberTokenData(member)
	if tokenData != nil {
		response.Success(ctx, consts.CurdStatusOkMsg, tokenData)
		return
	}
	response.Fail(ctx, -400100, "注册失败", gin.H{})
}

func (c *Passport) PasswdLogin(ctx *gin.Context) {
	phone := ctx.GetString(consts.ValidatorPrefix + "phone")
	passwd := ctx.GetString(consts.ValidatorPrefix + "passwd")
	member, err := members.CreateMemberServiceFactory().PasswdLogin(phone, passwd)
	if err != nil {
		response.Fail(ctx, -400101, "登录失败"+err.Error(), gin.H{})
		return
	}
	tokenData := c.memberTokenData(member)
	if tokenData != nil {
		response.Success(ctx, consts.CurdStatusOkMsg, tokenData)
		return
	}
	response.Fail(ctx, -400101, "登录失败,Token生成出错啦", gin.H{})
}

// PhoneLogin 手机验证码登录
func (c *Passport) PhoneLogin(ctx *gin.Context) {
	phone := ctx.GetString(consts.ValidatorPrefix + "phone")
	code := ctx.GetString(consts.ValidatorPrefix + "code")

	member, err := members.CreateMemberServiceFactory().PhoneLogin(phone, code)
	if err != nil {
		response.Fail(ctx, -400100, "登录失败"+err.Error(), gin.H{})
		return
	}
	tokenData := c.memberTokenData(member)
	if tokenData != nil {
		response.Success(ctx, consts.CurdStatusOkMsg, tokenData)
		return
	}
	response.Fail(ctx, -400100, "登录失败,Token生成出错啦", gin.H{})
}

//Logout 退出
func (c *Passport) Logout(ctx *gin.Context) {

}

// memberTokenData 返回用户Token
func (c *Passport) memberTokenData(member *model.MemberModel) gin.H {
	userToken, err := token.CreateTokenFactory().GenerateToken(variable.ConfigCustomYml.GetString("TokenPlatform.Api"), member.Id, "", variable.ConfigYml.GetInt64("Token.JwtTokenCreatedExpireAt"))
	if err == nil {
		token.CreateTokenFactory().RecordLoginToken(userToken)

		rs := []rune(member.Phone)
		member.Phone = string(rs[0:3]) + "******" + string(rs[9:11])

		data := gin.H{
			"token": userToken,
			"data": gin.H{
				"id":        member.Id,
				"phone":     member.Phone,
				"nick_name": member.NickName,
				"real_name": member.RealName,
				"state":     member.State,
			},
		}
		return data
	}
	return nil
}

//PhoneCode 获取手机验证码
func (c *Passport) PhoneCode(ctx *gin.Context) {
	phone, isOk := ctx.GetPostForm("phone")
	if !isOk {
		response.Fail(ctx, -400101, "请输入手机号码", gin.H{})
		return
	}

	if !teach_tool.CheckPhone(phone) {
		response.Fail(ctx, -400101, "手机号输入错误", gin.H{})
		return
	}
	err := sms.CreatePhoneCodeFactory(phone).SendPhoneCode()
	if err != nil {
		response.Fail(ctx, -400101, err.Error(), gin.H{})
		return
	}
	response.Success(ctx, "验证码发送成功", gin.H{})
}

//RefreshToken 刷新Token
func (c *Passport) RefreshToken(ctx *gin.Context) {

}

// Sso 单点登录处理
func (c *Passport) Sso(context *gin.Context) {
	ticket, _ := context.GetQuery("ticket")
	if ticket == "" {
		response.Fail(context, consts.CurdCreatFailCode, "参数错误", "")
		return
	}

	userId, err := (&service.Client{}).CheckTicket(ticket)
	if err != nil {
		response.Fail(context, consts.CurdCreatFailCode, err.Error(), "")
		return
	}
	//userId := "cs10000"
	// TODO 为和原系统做最小的更改，针对第一次登录回来的用户，自动在user表中登录一条数据
	uid := model.CreateMemberFactory("").SSOAutoLogin(userId)
	if uid == 0 {
		response.Fail(context, consts.CurdCreatFailCode, "帐号不存在，请联系信息中心同步账号", "")
		return
	}
	member := model.CreateMemberFactory("").GetById(int(uid))
	tokenData := c.memberTokenData(member)
	if tokenData != nil {
		response.Success(context, consts.CurdStatusOkMsg, tokenData)
		return
	}
	response.Fail(context, consts.CurdCreatFailCode, "帐号不存在，请联系信息中心同步账号", "")
	return
}
