package admin

import (
	"github.com/farmer-simon/go-utils"
	"github.com/gin-gonic/gin"
	"goskeleton/app/global/consts"
	"goskeleton/app/model"
	"goskeleton/app/service/category"
	"goskeleton/app/utils/response"
)

type Category struct {
}

//Index 分类列表
func (c *Category) Index(ctx *gin.Context) {
	count, list := category.CreateCateServiceFactory().Index(ctx)
	response.Success(ctx, consts.CurdStatusOkMsg, gin.H{
		"count": count,
		"list":  list,
	})
}

// Create 添加分类
func (c *Category) Create(ctx *gin.Context) {
	err := category.CreateCateServiceFactory().Create(ctx)
	if err == nil {
		response.Success(ctx, consts.CurdStatusOkMsg, "")
	} else {
		response.Fail(ctx, consts.CurdCreatFailCode, err.Error(), "")
	}
}

// Edit 添加分类
func (c *Category) Edit(ctx *gin.Context) {
	err := category.CreateCateServiceFactory().Update(ctx)
	if err == nil {
		response.Success(ctx, consts.CurdStatusOkMsg, "")
	} else {
		response.Fail(ctx, consts.CurdUpdateFailCode, err.Error(), "")
	}
}

// Info 分类详情
func (c *Category) Info(ctx *gin.Context) {
	id, _ := ctx.GetQuery("id")
	cateId := utils.String2Int(id)
	cate := model.CreateCateModelFactory("").GetById(cateId)
	if cate != nil {
		if cate.IsDel == 1 {
			response.Fail(ctx, consts.CurdSelectFailCode, "分类已删除", "")
			return
		}
		response.Success(ctx, consts.CurdStatusOkMsg, cate)
	} else {
		response.Fail(ctx, consts.CurdSelectFailCode, "分类信息未找到", "")
	}
}

//Destroy 删除分类
func (c *Category) Destroy(ctx *gin.Context) {
	id, _ := ctx.GetQuery("id")
	advId := utils.String2Int(id)
	if advId > 0 {
		err := category.CreateCateServiceFactory().Destroy(advId)
		if err == nil {
			response.Success(ctx, consts.CurdStatusOkMsg, "")
		} else {
			response.Fail(ctx, consts.CurdDeleteFailCode, err.Error(), "")
		}
		return
	}
	response.Fail(ctx, consts.CurdDeleteFailCode, consts.CurdDeleteFailMsg, "")
}
