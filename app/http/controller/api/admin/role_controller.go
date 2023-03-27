package admin

import (
	"github.com/farmer-simon/go-utils"
	"github.com/gin-gonic/gin"
	"goskeleton/app/service/users"
	"goskeleton/app/utils/response"
)

type Role struct {
}

func (r *Role) Index(ctx *gin.Context) {
	//users.CreateRoleServiceFactory()
}

func (r *Role) Create(ctx *gin.Context) {
	err := users.CreateRoleServiceFactory().AddRole(ctx)
	if err != nil {
		response.Fail(ctx, -400401, err.Error(), "")
		return
	}
	response.Success(ctx, "操作成功", gin.H{})
}

func (r *Role) Edit(ctx *gin.Context) {
	err := users.CreateRoleServiceFactory().EditRole(ctx)
	if err != nil {
		response.Fail(ctx, -400402, err.Error(), "")
		return
	}
	response.Success(ctx, "操作成功", gin.H{})
}

func (r *Role) Destroy(ctx *gin.Context) {
	id, _ := ctx.GetQuery("id")
	roleId := utils.String2Int(id)
	if roleId > 0 {
		err := users.CreateRoleServiceFactory().DestroyRole(roleId)
		if err == nil {
			response.Success(ctx, "操作成功", gin.H{})
			return
		}
		response.Fail(ctx, -400403, err.Error(), "")
		return
	}
	response.Fail(ctx, -400404, "参数错误", "")

}

type RoleUsers struct {
}

func (ru *RoleUsers) Index(ctx *gin.Context) {

}
func (ru *RoleUsers) Create(ctx *gin.Context) {
	err := users.CreateRoleServiceFactory().AddRoleUser(ctx)
	if err != nil {
		response.Fail(ctx, -400405, err.Error(), "")
		return
	}
	response.Success(ctx, "操作成功", gin.H{})
}

func (ru *RoleUsers) Destroy(ctx *gin.Context) {
	err := users.CreateRoleServiceFactory().DeleteRoleUser(ctx)
	if err != nil {
		response.Fail(ctx, -400406, err.Error(), "")
		return
	}
	response.Success(ctx, "操作成功", gin.H{})
}
