package register_validator

import (
	"goskeleton/app/core/container"
	"goskeleton/app/global/consts"
	"goskeleton/app/http/validator/api/admin/category"
	"goskeleton/app/http/validator/api/admin/members"
	"goskeleton/app/http/validator/api/admin/needs"
	"goskeleton/app/http/validator/api/admin/services"
	"goskeleton/app/http/validator/api/admin/settings"
	"goskeleton/app/http/validator/api/admin/users"
	homeMembers "goskeleton/app/http/validator/api/home/members"
	homeNeeds "goskeleton/app/http/validator/api/home/needs"
	homeServices "goskeleton/app/http/validator/api/home/services"
	"goskeleton/app/http/validator/common/upload_files"
)

// ApiRegisterValidator 各个业务模块验证器必须进行注册（初始化），程序启动时会自动加载到容器
func ApiRegisterValidator() {
	//创建容器
	containers := container.CreateContainersFactory()

	// ADMIN API部分
	containers.Set(consts.ValidatorPrefix+"AdminLogin", users.Login{})
	containers.Set(consts.ValidatorPrefix+"Settings", settings.Set{})

	containers.Set(consts.ValidatorPrefix+"CategoryCreate", category.Create{})
	containers.Set(consts.ValidatorPrefix+"CategoryEdit", category.Edit{})
	containers.Set(consts.ValidatorPrefix+"CategoryIndex", category.Index{})

	containers.Set(consts.ValidatorPrefix+"UsersCreate", users.Create{})
	containers.Set(consts.ValidatorPrefix+"UsersEdit", users.Edit{})
	containers.Set(consts.ValidatorPrefix+"UsersIndex", users.Index{})
	containers.Set(consts.ValidatorPrefix+"UsersChangePasswd", users.ChangePasswd{})

	containers.Set(consts.ValidatorPrefix+"MembersIndex", members.Index{})
	containers.Set(consts.ValidatorPrefix+"MembersVerify", members.Verify{})

	containers.Set(consts.ValidatorPrefix+"AdminNeedIndex", needs.Index{})
	containers.Set(consts.ValidatorPrefix+"AdminNeedVerify", needs.Verify{})
	containers.Set(consts.ValidatorPrefix+"AdminNeedRecord", needs.RecordIndex{})

	containers.Set(consts.ValidatorPrefix+"AdminServicesIndex", services.Index{})
	containers.Set(consts.ValidatorPrefix+"AdminServicesVerify", services.Verify{})
	containers.Set(consts.ValidatorPrefix+"AdminServicesRecord", services.RecordIndex{})
	//HOME
	containers.Set(consts.ValidatorPrefix+"PhoneLogin", homeMembers.PhoneLogin{})
	containers.Set(consts.ValidatorPrefix+"Register", homeMembers.Register{})

	containers.Set(consts.ValidatorPrefix+"MyNeedIndex", homeNeeds.MyIndex{})
	containers.Set(consts.ValidatorPrefix+"MyNeedCreate", homeNeeds.Create{})
	containers.Set(consts.ValidatorPrefix+"MyNeedEdit", homeNeeds.Edit{})

	containers.Set(consts.ValidatorPrefix+"NeedIndex", homeNeeds.Index{})

	containers.Set(consts.ValidatorPrefix+"MyServicesIndex", homeServices.MyIndex{})
	containers.Set(consts.ValidatorPrefix+"MyServicesCreate", homeServices.Create{})
	containers.Set(consts.ValidatorPrefix+"MyServicesEdit", homeServices.Edit{})

	containers.Set(consts.ValidatorPrefix+"ServicesIndex", homeServices.Index{})

	containers.Set(consts.ValidatorPrefix+"UpdateInfo", homeMembers.UpdateInfo{})
	containers.Set(consts.ValidatorPrefix+"ChangePhone", homeMembers.ChangePhone{})
	// 文件上传
	containers.Set(consts.ValidatorPrefix+"UploadFiles", upload_files.UpFiles{})
	containers.Set(consts.ValidatorPrefix+"ContactServices", homeServices.Contact{})
	containers.Set(consts.ValidatorPrefix+"ContactNeeds", homeNeeds.Contact{})

	containers.Set(consts.ValidatorPrefix+"ServicesRecordPost", homeServices.RecordPost{})
	containers.Set(consts.ValidatorPrefix+"MyServicesRecordAccept", homeServices.RecordAccept{})
	containers.Set(consts.ValidatorPrefix+"MyServicesRecordClose", homeServices.RecordClose{})
	containers.Set(consts.ValidatorPrefix+"MyServicesRecordIndex", homeServices.RecordIndex{})
	containers.Set(consts.ValidatorPrefix+"MyServicesRecordList", homeServices.RecordList{})

	containers.Set(consts.ValidatorPrefix+"NeedsRecordPost", homeNeeds.RecordPost{})
	containers.Set(consts.ValidatorPrefix+"MyNeedsRecordAccept", homeNeeds.RecordAccept{})
	containers.Set(consts.ValidatorPrefix+"MyNeedsRecordClose", homeNeeds.RecordClose{})
	containers.Set(consts.ValidatorPrefix+"MyNeedsRecordIndex", homeNeeds.RecordIndex{})
	containers.Set(consts.ValidatorPrefix+"MyNeedsRecordList", homeNeeds.RecordList{})

}
