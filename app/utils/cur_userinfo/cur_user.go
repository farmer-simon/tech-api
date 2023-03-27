package cur_userinfo

import (
	"errors"
	"github.com/gin-gonic/gin"
	"goskeleton/app/global/variable"
	"goskeleton/app/http/middleware/my_jwt"
	"goskeleton/app/model"
)

// GetHomeCurrentUserId 获取API当前用户的id
// @context 请求上下文
func GetHomeCurrentUserId(context *gin.Context) (int64, bool) {
	currentUser, exists := GetCurrentUser(context)
	if exists && currentUser.Platform == variable.ConfigCustomYml.GetString("TokenPlatform.Api") {
		return currentUser.UserId, true
	}
	return 0, false
}

func GetHomeCurrentUserInfo(context *gin.Context) (member *model.MemberModel, err error) {
	id, isExists := GetHomeCurrentUserId(context)
	if !isExists {
		return nil, errors.New("请重新登录")
	}
	member = model.CreateMemberFactory("").GetById(int(id))
	if member == nil {
		err = errors.New("用户查询失败")
	}
	//为了安全处理一下
	member.Passwd = ""
	return
}

// GetAdminCurrentUserId 获取ADMIN当前用户的id
// @context 请求上下文
func GetAdminCurrentUserId(context *gin.Context) (int64, bool) {
	currentUser, exists := GetCurrentUser(context)
	if exists && currentUser.Platform == variable.ConfigCustomYml.GetString("TokenPlatform.Admin") {
		return currentUser.UserId, true
	}
	return 0, false
}

// GetCurrentUser 获取当前用户
// @context 请求上下文
func GetCurrentUser(context *gin.Context) (currentUser my_jwt.CustomClaims, exists bool) {
	tokenKey := variable.ConfigYml.GetString("Token.BindContextKeyName")
	value, exists := context.Get(tokenKey)
	if !exists {
		return
	}
	currentUser = value.(my_jwt.CustomClaims)
	return
}

func GetCurrentUserRoleIds(context *gin.Context) (roleIds []int, err error) {
	userId, exists := GetAdminCurrentUserId(context)
	if exists {
		roleIds = model.CreateUserRoleFactory("").GetRoleIdsByUserId(userId)
		if len(roleIds) == 0 {
			err = errors.New("无任何角色，请联系管理员！")
		}
		return
	}
	return nil, errors.New("未知角色，请联系管理员！")
}
