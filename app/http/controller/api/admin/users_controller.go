package admin

import (
	"github.com/farmer-simon/go-utils"
	"github.com/gin-gonic/gin"
	"goskeleton/app/global/consts"
	"goskeleton/app/global/variable"
	"goskeleton/app/model"
	"goskeleton/app/service/users/curd"
	"goskeleton/app/service/users/token"
	"goskeleton/app/utils/cur_userinfo"
	"goskeleton/app/utils/response"
	"math"
	"time"
)

type Users struct {
}

func (c *Users) Index(ctx *gin.Context) {
	name := ctx.GetString(consts.ValidatorPrefix + "name")
	state := ctx.GetFloat64(consts.ValidatorPrefix + "state")
	page := ctx.GetFloat64(consts.ValidatorPrefix + "page")
	limit := ctx.GetFloat64(consts.ValidatorPrefix + "limit")
	offset := math.Max(0, (page-1)*limit)

	list, count := model.CreateUserFactory("").Index(name, int(state), int(offset), int(limit))
	for i, item := range list {
		item.Passwd = ""
		list[i] = item
	}
	response.Success(ctx, consts.CurdStatusOkMsg, gin.H{
		"list":  list,
		"count": count,
	})
}

func (c *Users) Login(ctx *gin.Context) {
	userName := ctx.GetString(consts.ValidatorPrefix + "user_name")
	passwd := ctx.GetString(consts.ValidatorPrefix + "passwd")
	userModel := model.CreateUserFactory("").Login(userName, passwd)
	if userModel != nil {
		if userModel.State != 2 {
			response.Fail(ctx, -400300, "登录失败,您的帐号已被停用！", "")
			return
		}
		userToken, err := token.CreateTokenFactory().GenerateToken(variable.ConfigCustomYml.GetString("TokenPlatform.Admin"), userModel.Id, "", variable.ConfigYml.GetInt64("Token.JwtTokenCreatedExpireAt"))
		if err == nil {
			token.CreateTokenFactory().RecordLoginToken(userToken)
			data := gin.H{
				"id":         userModel.Id,
				"user_name":  userName,
				"name":       userModel.Name,
				"avatar":     userModel.Avatar,
				"phone":      userModel.Phone,
				"token":      userToken,
				"updated_at": time.Now().Format(variable.DateFormat),
			}
			response.Success(ctx, consts.CurdStatusOkMsg, data)
			return
		}
	}
	response.Fail(ctx, -400300, "登录失败,用户名或密码错误！", "")
}

// LoginOut 退出登录
func (c *Users) LoginOut(ctx *gin.Context) {
	//退出登录，直接删除Token
	uid, _ := cur_userinfo.GetAdminCurrentUserId(ctx)
	model.CreateOauthFactory("").Destroy(int(uid), variable.ConfigCustomYml.GetString("TokenPlatform.Admin"))
	response.Success(ctx, consts.CurdStatusOkMsg, gin.H{})
}

//Create 添加帐号
func (c *Users) Create(ctx *gin.Context) {
	err := curd.CreateUserCurdFactory().Create(ctx)
	if err != nil {
		response.Fail(ctx, -400300, err.Error(), "")
		return
	}
	response.Success(ctx, "操作成功", gin.H{})
}

//Edit 修改帐号
func (c *Users) Edit(ctx *gin.Context) {
	err := curd.CreateUserCurdFactory().Update(ctx)
	if err != nil {
		response.Fail(ctx, -400300, err.Error(), "")
		return
	}
	response.Success(ctx, "操作成功", gin.H{})
}

// Destroy 删除帐号
func (c *Users) Destroy(ctx *gin.Context) {
	id, _ := ctx.GetQuery("id")
	userId := utils.String2Int(id)
	currentUid, _ := cur_userinfo.GetAdminCurrentUserId(ctx)
	if userId == int(currentUid) {
		response.Fail(ctx, -400300, "为了保障系统有效运行，不允许删除自己", "")
		return
	}
	err := curd.CreateUserCurdFactory().Destroy(userId)
	if err == nil {
		response.Success(ctx, "操作成功", "")
	} else {
		response.Fail(ctx, -400300, err.Error(), err)
	}
}

// Info 帐号详情
func (c *Users) Info(ctx *gin.Context) {
	id, _ := ctx.GetQuery("id")
	userId := utils.String2Int(id)
	if userId > 0 {
		user, err := model.CreateUserFactory("").GetUserById(userId)
		if err == nil {
			user.Passwd = ""
			response.Success(ctx, "操作成功", user)
		} else {
			response.Fail(ctx, -400300, err.Error(), err)
		}
		return
	}
	response.Fail(ctx, -400300, "操作失败", "")
}

// ChangePasswd 修改密码
func (c *Users) ChangePasswd(ctx *gin.Context) {
	err := curd.CreateUserCurdFactory().ChangePasswd(ctx)
	if err != nil {
		response.Fail(ctx, -400306, err.Error(), "")
		return
	}
	response.Success(ctx, "操作成功", gin.H{})
}
