package curd

import (
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"goskeleton/app/global/consts"
	"goskeleton/app/global/variable"
	"goskeleton/app/model"
	"goskeleton/app/utils/cur_userinfo"
	"goskeleton/app/utils/data_bind"
	"goskeleton/app/utils/md5_encrypt"
)

func CreateUserCurdFactory() *UsersCurd {
	return &UsersCurd{model.CreateUserFactory("")}
}

type UsersCurd struct {
	model *model.UsersModel
}

// Create 添加用户
func (u *UsersCurd) Create(ctx *gin.Context) error {
	var data model.UsersModel
	if err := data_bind.ShouldBindFormDataToModel(ctx, &data); err == nil {
		data.Passwd = md5_encrypt.Base64Md5(data.Passwd) // 预先处理密码加密，然后存储在数据库
		exists := u.model.CheckRepeatUserName(0, data.UserName)
		if exists > 0 {
			return errors.New(fmt.Sprintf("用户名 %s 已存在", data.UserName))
		}
		res := u.model.Create(data.UserName, data.Name, data.Passwd, data.Phone, data.Avatar, data.State)
		if res {
			return nil
		}
		return errors.New("帐号添加失败")
	}
	return nil
}

func (u *UsersCurd) Update(ctx *gin.Context) error {
	var data model.UsersModel
	if err := data_bind.ShouldBindFormDataToModel(ctx, &data); err == nil {
		// 密码不输入则不修改密码
		if data.Passwd != "" {
			data.Passwd = md5_encrypt.Base64Md5(data.Passwd) // 预先处理密码加密，然后存储在数据库
		}
		exists := u.model.CheckRepeatUserName(int(data.Id), data.UserName)
		if exists > 0 {
			return errors.New(fmt.Sprintf("用户名 %s 已存在", data.UserName))
		}
		res := u.model.Update(int(data.Id), data.UserName, data.Name, data.Passwd, data.Phone, data.Avatar, data.State)
		if res {
			if data.Passwd != "" {
				// TODO 将原TOKEN失效
				model.CreateOauthFactory("").Destroy(int(data.Id), variable.ConfigCustomYml.GetString("TokenPlatform.Admin"))
			}
			return nil
		}
		return errors.New("帐号编辑失败")
	}
	return nil
}

func (u *UsersCurd) ChangePasswd(ctx *gin.Context) error {

	passwd := ctx.GetString(consts.ValidatorPrefix + "passwd")
	newPasswd := ctx.GetString(consts.ValidatorPrefix + "new_passwd")
	rePasswd := ctx.GetString(consts.ValidatorPrefix + "re_passwd")
	if newPasswd != rePasswd {
		return errors.New("两次输入的密码不一致")
	}
	passwd = md5_encrypt.Base64Md5(passwd)
	newPasswd = md5_encrypt.Base64Md5(newPasswd)

	uid, isOk := cur_userinfo.GetAdminCurrentUserId(ctx)
	if isOk {
		uInfo, err := u.model.GetUserById(int(uid))
		if err != nil {
			return errors.New("帐号信息查询失败")
		}
		if uInfo.Passwd != passwd {
			return errors.New("原密码校验失败")
		}
		isOk = u.model.ChangePasswd(int(uInfo.Id), newPasswd)
		if isOk {
			// TODO 将原TOKEN失效
			model.CreateOauthFactory("").Destroy(int(uInfo.Id), variable.ConfigCustomYml.GetString("TokenPlatform.Admin"))
			return nil
		}
	}
	return errors.New("更新密码失败，请联系技术")
}

//Destroy 删除
func (u *UsersCurd) Destroy(id int) error {
	isOk := u.model.Destroy(id)
	if isOk {
		// TODO 将原TOKEN失效
		model.CreateOauthFactory("").Destroy(id, variable.ConfigCustomYml.GetString("TokenPlatform.Admin"))
		return nil
	}
	return errors.New("删除失败")
}
