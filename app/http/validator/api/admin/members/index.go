package members

import (
	"github.com/farmer-simon/go-utils"
	"github.com/gin-gonic/gin"
	"goskeleton/app/data_type"
	"goskeleton/app/global/consts"
	"goskeleton/app/http/controller/api/admin"
	commonDataType "goskeleton/app/http/validator/common/data_type"
	"goskeleton/app/http/validator/core/data_transfer"
	"goskeleton/app/utils/response"
	"time"
)

type Index struct {
	// 表单参数验证结构体支持匿名结构体嵌套
	data_type.MembersQuery
	commonDataType.Page
}

// 验证器语法，参见 Register.go文件，有详细说明

func (validVal Index) CheckParams(context *gin.Context) {

	//1.基本的验证规则没有通过
	if err := context.ShouldBind(&validVal); err != nil {
		response.ValidatorError(context, err)
		return
	}

	if validVal.StartTime != "" {
		tmpTime, err := utils.FormatDateTimeStringToDateTime(validVal.StartTime)
		if err == nil {
			validVal.StartTime = tmpTime.Format(utils.TimeFormat)
		} else {
			validVal.StartTime = time.Now().Format(utils.TimeFormat)
		}
	}
	if validVal.EndTime != "" {
		tmpTime, err := utils.FormatDateTimeStringToDateTime(validVal.EndTime)
		if err == nil {
			validVal.EndTime = tmpTime.Format(utils.TimeFormat)
		} else {
			validVal.EndTime = time.Now().Format(utils.TimeFormat)
		}
	}

	//  该函数主要是将本结构体的字段（成员）按照 consts.ValidatorPrefix+ json标签对应的 键 => 值 形式绑定在上下文，便于下一步（控制器）可以直接通过 context.Get(键) 获取相关值
	extraAddBindDataContext := data_transfer.DataAddContext(validVal, consts.ValidatorPrefix, context)
	if extraAddBindDataContext == nil {
		response.ErrorSystem(context, "数据检测失败", "")
	} else {
		// 验证完成，调用控制器,并将验证器成员(字段)递给控制器，保持上下文数据一致性
		(&admin.Members{}).Index(extraAddBindDataContext)
	}

}
