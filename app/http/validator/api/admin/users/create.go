package users

import (
	"github.com/gin-gonic/gin"
	"goskeleton/app/data_type"
	"goskeleton/app/global/consts"
	"goskeleton/app/http/controller/api/admin"
	"goskeleton/app/http/validator/core/data_transfer"
	"goskeleton/app/utils/response"
)

type Create struct {
	// 表单参数验证结构体支持匿名结构体嵌套
	data_type.UsersBase
}

// 验证器语法，参见 Register.go文件，有详细说明

func (c Create) CheckParams(context *gin.Context) {

	//1.基本的验证规则没有通过
	if err := context.ShouldBind(&c); err != nil {
		response.ValidatorError(context, err)
		return
	}
	if len(c.Passwd) < 8 {
		response.CustomError(context, "密码不能少于8位")
		return
	}
	//  该函数主要是将本结构体的字段（成员）按照 consts.ValidatorPrefix+ json标签对应的 键 => 值 形式绑定在上下文，便于下一步（控制器）可以直接通过 context.Get(键) 获取相关值
	extraAddBindDataContext := data_transfer.DataAddContext(c, consts.ValidatorPrefix, context)
	if extraAddBindDataContext == nil {
		response.ErrorSystem(context, "数据检测失败", "")
	} else {
		// 验证完成，调用控制器,并将验证器成员(字段)递给控制器，保持上下文数据一致性
		(&admin.Users{}).Create(extraAddBindDataContext)
	}

}
