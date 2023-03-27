package users

import (
	"errors"
	"github.com/gin-gonic/gin"
	"goskeleton/app/global/consts"
	"goskeleton/app/model"
	"goskeleton/app/utils/data_bind"
)

func CreateRoleServiceFactory() *RoleService {
	return &RoleService{model.CreateUserRoleFactory("")}
}

type RoleService struct {
	model *model.UserRole
}

// AddRole 添加角色
func (serv *RoleService) AddRole(ctx *gin.Context) error {
	name := ctx.GetString(consts.ValidatorPrefix + "name")
	exists := serv.model.CheckRepeatRoleName(0, name)
	if exists > 0 {
		return errors.New("重复的角色名称")
	}
	_, err := serv.model.AddRole(name)
	if err != nil {
		return err
	}
	return nil
}

func (serv *RoleService) EditRole(ctx *gin.Context) error {
	id := ctx.GetFloat64(consts.ValidatorPrefix + "id")
	name := ctx.GetString(consts.ValidatorPrefix + "name")
	exists := serv.model.CheckRepeatRoleName(int(id), name)
	if exists > 0 {
		return errors.New("重复的角色名称")
	}
	err := serv.model.EditRole(int(id), name)
	if err != nil {
		return err
	}
	return nil
}

//DestroyRole 删除角色
func (serv *RoleService) DestroyRole(id int) error {
	if consts.FounderRoleId == id {
		return errors.New("为了保障系统有效运行，此角色不允许删除")
	}

	users := serv.model.GetRoleUsers(id)
	if len(users) > 0 {
		return errors.New("请先移除角色下的成员")
	}
	serv.model.DeleteRole(id)
	return nil
}

//AddRoleUser 添加角色成员
func (serv *RoleService) AddRoleUser(ctx *gin.Context) error {
	var data model.RoleUsers
	if err := data_bind.ShouldBindFormDataToModel(ctx, &data); err == nil {
		exists := serv.model.CheckRepeatRoleUser(data.RoleId, data.UserId)
		if exists > 0 {
			return errors.New("已添加的成员")
		}
		return serv.model.AddRoleUser(data.RoleId, data.UserId)
	}
	return errors.New("参数错误")
}

//DeleteRoleUser 删除角色成员
func (serv *RoleService) DeleteRoleUser(ctx *gin.Context) error {
	var data model.RoleUsers
	if err := data_bind.ShouldBindFormDataToModel(ctx, &data); err == nil {
		if consts.FounderRoleId == data.RoleId && consts.FounderUserId == data.UserId {
			return errors.New("为了保障系统有效运行，此数据不允许删除")
		}
		serv.model.DeleteRoleUser(data.RoleId, data.UserId)
		return nil
	}
	return errors.New("参数错误")
}
