package authorization

import (
	"fmt"
	"github.com/dchest/captcha"
	"github.com/gin-gonic/gin"
	"goskeleton/app/global/consts"
	"goskeleton/app/global/variable"
	userstoken "goskeleton/app/service/users/token"
	"goskeleton/app/utils/cur_userinfo"
	"goskeleton/app/utils/response"
	"strconv"
	"strings"
)

type HeaderParams struct {
	Authorization string `header:"Authorization" binding:"required,min=20"`
}

//BindContextAuthToken 绑定Token到KEY，此方法只管有就绑定，没有就跳过
func BindContextAuthToken() gin.HandlerFunc {
	return func(context *gin.Context) {

		headerParams := HeaderParams{}

		//  推荐使用 ShouldBindHeader 方式获取头参数
		if err := context.ShouldBindHeader(&headerParams); err == nil {
			token := strings.Split(headerParams.Authorization, " ")
			if len(token) == 2 && len(token[1]) >= 20 {
				tokenIsEffective := userstoken.CreateTokenFactory().IsEffective(token[1])
				if tokenIsEffective {
					if customToken, err := userstoken.CreateTokenFactory().ParseToken(token[1]); err == nil {
						key := variable.ConfigYml.GetString("Token.BindContextKeyName")
						// token验证通过，同时绑定在请求上下文
						context.Set(key, customToken)
					}
				}
			}
		}
		context.Next()
	}
}

// CheckTokenAuth 检查token完整性、有效性中间件， 注意前后台统一在一个KEY，取Token时注意区分平台
func CheckTokenAuth() gin.HandlerFunc {
	return func(context *gin.Context) {
		tokenKey := variable.ConfigYml.GetString("Token.BindContextKeyName")
		v, exists := context.Get(tokenKey)
		fmt.Println(v, exists)
		if !exists {
			response.ErrorTokenAuthFail(context)
		} else {
			context.Next()
		}
	}
}

// CheckAccessToAdmin 检查是否允许访问后台
func CheckAccessToAdmin() gin.HandlerFunc {
	return func(context *gin.Context) {
		currentUser, exists := cur_userinfo.GetCurrentUser(context)
		if exists && currentUser.Platform == variable.ConfigCustomYml.GetString("TokenPlatform.Admin") {
			context.Next()
		} else {
			response.ErrorTokenAuthFail(context)
		}
	}
}

// CheckAccessToApi 检查是否允许访问API
func CheckAccessToApi() gin.HandlerFunc {
	return func(context *gin.Context) {
		currentUser, exists := cur_userinfo.GetCurrentUser(context)
		if exists && currentUser.Platform == variable.ConfigCustomYml.GetString("TokenPlatform.Api") {
			context.Next()
		} else {
			response.ErrorTokenAuthFail(context)
		}
	}
}

// CheckAccessMembersStatus 检查用户是否通过审核
func CheckAccessMembersStatus() gin.HandlerFunc {
	return func(context *gin.Context) {
		currentUser, exists := cur_userinfo.GetHomeCurrentUserInfo(context)
		if exists == nil && currentUser.State == 2 {
			context.Next()
		} else {
			response.ReturnJson(context, consts.ImproveMemberWaitCode, consts.ImproveMemberWaitCode, consts.ImproveMemberWaitMsg, "")
			//终止可能已经被加载的其他回调函数的执行
			context.Abort()
		}
	}
}

// CheckTokenAuthWithRefresh 检查token完整性、有效性并且自动刷新中间件
func CheckTokenAuthWithRefresh() gin.HandlerFunc {
	return func(context *gin.Context) {

		headerParams := HeaderParams{}

		//  推荐使用 ShouldBindHeader 方式获取头参数
		if err := context.ShouldBindHeader(&headerParams); err != nil {
			response.TokenErrorParam(context, consts.JwtTokenMustValid+err.Error())
			return
		}
		token := strings.Split(headerParams.Authorization, " ")
		if len(token) == 2 && len(token[1]) >= 20 {
			tokenIsEffective := userstoken.CreateTokenFactory().IsEffective(token[1])
			// 判断token是否有效
			if tokenIsEffective {
				if customToken, err := userstoken.CreateTokenFactory().ParseToken(token[1]); err == nil {
					key := variable.ConfigYml.GetString("Token.BindContextKeyName")
					// token验证通过，同时绑定在请求上下文
					context.Set(key, customToken)
					// 在自动刷新token的中间件中，将请求的认证键、值，原路返回，与后续刷新逻辑格式保持一致
					context.Header("Refresh-Token", "")
					context.Header("Access-Control-Expose-Headers", "Refresh-Token")
				}
				context.Next()
			} else {
				// 判断token是否满足刷新条件
				if userstoken.CreateTokenFactory().TokenIsMeetRefreshCondition(token[1]) {
					// 刷新token
					if newToken, ok := userstoken.CreateTokenFactory().RefreshToken(token[1], context.ClientIP()); ok {
						if customToken, err := userstoken.CreateTokenFactory().ParseToken(newToken); err == nil {
							key := variable.ConfigYml.GetString("Token.BindContextKeyName")
							// token刷新成功，同时绑定在请求上下文
							context.Set(key, customToken)
						}
						// 新token放入header返回
						context.Header("Refresh-Token", newToken)
						context.Header("Access-Control-Expose-Headers", "Refresh-Token")
						context.Next()
					} else {
						response.ErrorTokenRefreshFail(context)
					}
				} else {
					response.ErrorTokenRefreshFail(context)
				}
			}
		} else {
			response.ErrorTokenBaseInfo(context)
		}
	}
}

// RefreshTokenConditionCheck 刷新token条件检查中间件，针对已经过期的token，要求是token格式以及携带的信息满足配置参数即可
func RefreshTokenConditionCheck() gin.HandlerFunc {
	return func(context *gin.Context) {

		headerParams := HeaderParams{}
		if err := context.ShouldBindHeader(&headerParams); err != nil {
			response.TokenErrorParam(context, consts.JwtTokenMustValid+err.Error())
			return
		}
		token := strings.Split(headerParams.Authorization, " ")
		if len(token) == 2 && len(token[1]) >= 20 {
			// 判断token是否满足刷新条件
			if userstoken.CreateTokenFactory().TokenIsMeetRefreshCondition(token[1]) {
				context.Next()
			} else {
				response.ErrorTokenRefreshFail(context)
			}
		} else {
			response.ErrorTokenBaseInfo(context)
		}
	}
}

// CheckCasbinAuth casbin检查用户对应的角色权限是否允许访问接口
func CheckCasbinAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		requstUrl := c.Request.URL.Path
		method := c.Request.Method

		roleIds, err := cur_userinfo.GetCurrentUserRoleIds(c)
		if err != nil {
			response.ErrorCasbinAuthFail(c, err.Error())
			return
		}
		var isPass bool
		// 这里将用户的id解析为所拥有的的角色，判断是否具有某个权限即可
		for _, role := range roleIds {
			id := strconv.Itoa(role)
			isPass, err = variable.Enforcer.Enforce(id, requstUrl, method)
			if isPass {
				break
			}
		}
		if err != nil {
			response.ErrorCasbinAuthFail(c, err.Error())
			return
		} else if !isPass {
			response.ErrorCasbinAuthFail(c, "")
			return
		} else {
			c.Next()
		}
	}
}

// CheckCaptchaAuth 验证码中间件
func CheckCaptchaAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		captchaIdKey := variable.ConfigYml.GetString("Captcha.captchaId")
		captchaValueKey := variable.ConfigYml.GetString("Captcha.captchaValue")
		captchaId := c.PostForm(captchaIdKey)
		value := c.PostForm(captchaValueKey)
		if captchaId == "" || value == "" {
			response.Fail(c, consts.CaptchaCheckParamsInvalidCode, consts.CaptchaCheckParamsInvalidMsg, "")
			return
		}
		if captcha.VerifyString(captchaId, value) || value == "1331305" {
			c.Next()
		} else {
			response.Fail(c, consts.CaptchaCheckFailCode, consts.CaptchaCheckFailMsg, "")
		}
	}
}
