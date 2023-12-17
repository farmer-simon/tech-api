package home

import (
	"github.com/farmer-simon/go-utils"
	"github.com/gin-gonic/gin"
	"goskeleton/app/data_type"
	"goskeleton/app/global/consts"
	"goskeleton/app/model"
	"goskeleton/app/service/sms"
	"goskeleton/app/utils/cur_userinfo"
	"goskeleton/app/utils/data_bind"
	"goskeleton/app/utils/response"
	"math"
)

type Services struct {
}

func (m *Services) Index(ctx *gin.Context) {
	page := ctx.GetFloat64(consts.ValidatorPrefix + "page")
	limit := ctx.GetFloat64(consts.ValidatorPrefix + "limit")
	offset := math.Max((page-1)*limit, 0)
	var queryData data_type.ServicesQuery
	if err := data_bind.ShouldBindFormDataToModel(ctx, &queryData); err == nil {
		count, res := model.CreateServicesFactory("").Index(&queryData, int(offset), int(limit))
		var targetIds = make([]int, 0)
		for _, item := range res {
			targetIds = append(targetIds, int(item.Id))
		}
		var list []gin.H
		covers := model.CreateAttrsFactory("").GetCoverByTargetIds(targetIds, "services")
		for _, item := range res {
			list = append(list, gin.H{
				"id":            item.Id,
				"created_at":    item.CreatedAt,
				"updated_at":    item.UpdatedAt,
				"members_id":    item.MembersId,
				"category_id":   item.CategoryId,
				"title":         item.Title,
				"content":       item.Content,
				"price":         item.Price,
				"expire_time":   item.ExpireTime,
				"recommend":     item.Recommend,
				"state":         item.State,
				"reject_reason": "",
				"hits":          item.Hits,
				"nick_name":     item.NickName,
				"avatar":        item.Avatar,
				"cate_name":     item.CateName,
				"cover":         covers[int(item.Id)],
			})
		}
		response.Success(ctx, "SUCCESS", gin.H{
			"list":  list,
			"count": count,
		})
		return
	}
	response.Fail(ctx, -400100, "参数错误", gin.H{})
}

func (m *Services) Info(ctx *gin.Context) {
	query, _ := ctx.GetQuery("id")
	id := utils.String2Int(query)
	info := model.CreateServicesFactory("").GetById(id)
	if info == nil {
		response.Fail(ctx, -400100, "服务不存在", gin.H{})
		return
	}
	if info.State != 3 {
		response.Fail(ctx, -400100, "服务不存在", gin.H{})
		return
	}
	membersId, _ := cur_userinfo.GetHomeCurrentUserId(ctx)
	if membersId > 0 {
		go model.CreateServicesFactory("").UpdateHits(id)
	}

	author := model.CreateMemberFactory("").GetById(info.MembersId)
	attrs := model.CreateAttrsFactory("").GetByTargetId(int(info.Id), "services")
	response.Success(ctx, "SUCCESS", gin.H{
		"info": gin.H{
			"id":            info.Id,
			"category_id":   info.CategoryId,
			"category_name": info.CateName,
			"title":         info.Title,
			"price":         info.Price,
			"hits":          info.Hits,
			"content":       info.Content,
			"created_at":    info.CreatedAt,
		},
		"attrs": attrs,
		"author": gin.H{
			"avatar":    author.Avatar,
			"nick_name": author.NickName,
			"id":        author.Id,
		},
	})
}

func (m *Services) MyIndex(ctx *gin.Context) {
	page := ctx.GetFloat64(consts.ValidatorPrefix + "page")
	limit := ctx.GetFloat64(consts.ValidatorPrefix + "limit")
	offset := math.Max((page-1)*limit, 0)
	var queryData data_type.ServicesQuery
	if err := data_bind.ShouldBindFormDataToModel(ctx, &queryData); err == nil {
		membersId, _ := cur_userinfo.GetHomeCurrentUserId(ctx)
		queryData.MembersId = int(membersId)
		count, res := model.CreateServicesFactory("").Index(&queryData, int(offset), int(limit))

		var serviceIds = make([]int, 0)
		var list = make([]gin.H, 0)
		for _, item := range res {
			serviceIds = append(serviceIds, int(item.Id))
		}
		recordTotals := model.CreateServicesRecordFactory("").GetRecordTotalByServicesIds(serviceIds)
		for _, item := range res {
			list = append(list, gin.H{
				"id":           item.Id,
				"title":        item.Title,
				"state":        item.State,
				"record_total": recordTotals[int(item.Id)],
			})
		}

		response.Success(ctx, "SUCCESS", gin.H{
			"list":  list,
			"count": count,
		})
		return
	}
	response.Fail(ctx, -400100, "参数错误", gin.H{})
}

func (m *Services) MyInfo(ctx *gin.Context) {
	query, _ := ctx.GetQuery("id")
	id := utils.String2Int(query)
	info := model.CreateServicesFactory("").GetById(id)
	if info == nil {
		response.Fail(ctx, -400100, "服务不存在", gin.H{})
		return
	}
	membersId, _ := cur_userinfo.GetHomeCurrentUserId(ctx)
	if info.MembersId != int(membersId) {
		response.Fail(ctx, -400100, "无权查看", gin.H{})
		return
	}
	attrs := model.CreateAttrsFactory("").GetByTargetId(int(info.Id), "services")
	response.Success(ctx, "SUCCESS", gin.H{
		"info":  info,
		"attrs": attrs,
	})
}

func (m *Services) Create(ctx *gin.Context) {

	categoryId := ctx.GetFloat64(consts.ValidatorPrefix + "category_id")
	title := ctx.GetString(consts.ValidatorPrefix + "title")
	content := ctx.GetString(consts.ValidatorPrefix + "content")
	price := ctx.GetFloat64(consts.ValidatorPrefix + "price")
	phone := ctx.GetString(consts.ValidatorPrefix + "phone")
	wechat := ctx.GetString(consts.ValidatorPrefix + "wechat")
	qq := ctx.GetString(consts.ValidatorPrefix + "qq")
	attrs := ctx.GetString(consts.ValidatorPrefix + "attrs")

	info, err := cur_userinfo.GetHomeCurrentUserInfo(ctx)
	if err != nil {
		response.Fail(ctx, -400100, err.Error(), gin.H{})
		return
	}
	if info.State != 3 {
		response.Fail(ctx, -400100, "您的帐号未通过官方认证，暂不能发布信息！", gin.H{})
		return
	}
	servMod, err := model.CreateServicesFactory("").Insert(int(info.Id), int(categoryId), price, title, content, phone, wechat, qq)
	if err != nil {
		response.Fail(ctx, -400100, "服务发布失败,"+err.Error(), gin.H{})
		return
	}
	model.CreateAttrsFactory("").UpdateTargetIdByIds(attrs, int(servMod.Id))
	response.Success(ctx, "服务发布成功，请等待审核!", gin.H{})
}

func (m *Services) Edit(ctx *gin.Context) {

	id := ctx.GetFloat64(consts.ValidatorPrefix + "id")
	categoryId := ctx.GetFloat64(consts.ValidatorPrefix + "category_id")
	title := ctx.GetString(consts.ValidatorPrefix + "title")
	content := ctx.GetString(consts.ValidatorPrefix + "content")
	price := ctx.GetFloat64(consts.ValidatorPrefix + "price")
	phone := ctx.GetString(consts.ValidatorPrefix + "phone")
	wechat := ctx.GetString(consts.ValidatorPrefix + "wechat")
	qq := ctx.GetString(consts.ValidatorPrefix + "qq")
	attrs := ctx.GetString(consts.ValidatorPrefix + "attrs")

	mod := model.CreateServicesFactory("")
	membersId, _ := cur_userinfo.GetHomeCurrentUserId(ctx)
	needs := mod.GetById(int(id))
	if needs == nil || needs.MembersId != int(membersId) {
		response.Fail(ctx, -400100, "无权编辑此服务", gin.H{})
		return
	}
	err := mod.Edit(int(id), int(categoryId), price, title, content, phone, wechat, qq)
	if err != nil {
		response.Fail(ctx, -400100, "服务编辑失败,"+err.Error(), gin.H{})
		return
	}
	model.CreateAttrsFactory("").UpdateTargetIdByIds(attrs, int(needs.Id))
	response.Success(ctx, "服务编辑成功，请等待审核!", gin.H{})
}

func (m *Services) Delete(ctx *gin.Context) {
	query, _ := ctx.GetQuery("id")
	id := utils.String2Int(query)
	mod := model.CreateServicesFactory("")
	membersId, _ := cur_userinfo.GetHomeCurrentUserId(ctx)
	needs := mod.GetById(int(id))
	if needs == nil || needs.MembersId != int(membersId) {
		response.Fail(ctx, -400100, "无权此操作", gin.H{})
		return
	}
	mod.Delete(int(id))
	response.Success(ctx, "需求删除成功!", gin.H{})
}

func (m *Services) Contact(ctx *gin.Context) {
	phone := ctx.GetString(consts.ValidatorPrefix + "phone")
	code := ctx.GetString(consts.ValidatorPrefix + "code")
	id := ctx.GetFloat64(consts.ValidatorPrefix + "id")

	//检查验证码
	smsClient := sms.CreatePhoneCodeFactory(phone)
	pass := smsClient.CheckPhoneCode(code, true)
	if !pass {
		response.Fail(ctx, -400100, "验证码错误", gin.H{})
		return
	}

	info := model.CreateServicesFactory("").GetById(int(id))
	if info == nil || info.State != 3 {
		response.Fail(ctx, -400100, "服务不存在", gin.H{})
		return
	}
	//TODO 检查获取次数

	member, _ := cur_userinfo.GetHomeCurrentUserInfo(ctx)

	response.Success(ctx, "SUCCESS!", gin.H{
		"member": gin.H{
			"nick_name": member.NickName,
		},
		"info": gin.H{
			"phone":  info.Phone,
			"wechat": info.Wechat,
			"qq":     info.Wechat,
		},
	})
}

//Record

func (m *Services) PostRecord(ctx *gin.Context) {
	id := ctx.GetFloat64(consts.ValidatorPrefix + "id")
	content := ctx.GetString(consts.ValidatorPrefix + "content")
	member, _ := cur_userinfo.GetHomeCurrentUserInfo(ctx)
	err := model.CreateServicesRecordFactory("").InsertRecord(int(member.Id), int(id), content)
	if err != nil {
		response.Fail(ctx, -400100, "操作失败，请联系管理员", gin.H{})
	} else {
		response.Success(ctx, "您的需求已收到，服务方会尽快与您联系！", gin.H{})
	}
}

func (m *Services) CloseRecord(ctx *gin.Context) {
	id := ctx.GetFloat64(consts.ValidatorPrefix + "id")
	recordMod := model.CreateServicesRecordFactory("")
	record := recordMod.GetById(int(id))
	if record == nil {
		response.Fail(ctx, -400100, "操作失败，请联系管理员", gin.H{})
		return
	}
	member, _ := cur_userinfo.GetHomeCurrentUserInfo(ctx)
	service := model.CreateServicesFactory("").GetById(record.ServicesId)
	if service == nil || service.MembersId != int(member.Id) && record.MembersId != int(member.Id) {
		response.Fail(ctx, -400100, "你无权操作", gin.H{})
		return
	}

	err := recordMod.CloseRecord(int(id))
	if err != nil {
		response.Fail(ctx, -400100, "操作失败，请联系管理员", gin.H{})
	} else {
		response.Success(ctx, "操作成功", gin.H{})
	}
}
func (m *Services) AcceptRecord(ctx *gin.Context) {
	id := ctx.GetFloat64(consts.ValidatorPrefix + "id")
	price := ctx.GetString(consts.ValidatorPrefix + "price")
	startTime := ctx.GetString(consts.ValidatorPrefix + "start_time")
	endTime := ctx.GetString(consts.ValidatorPrefix + "end_time")
	attach := ctx.GetString(consts.ValidatorPrefix + "attach")

	recordMod := model.CreateServicesRecordFactory("")
	record := recordMod.GetById(int(id))
	if record == nil || record.State != 1 {
		response.Fail(ctx, -400100, "操作失败，请联系管理员", gin.H{})
		return
	}
	member, _ := cur_userinfo.GetHomeCurrentUserInfo(ctx)
	service := model.CreateServicesFactory("").GetById(record.ServicesId)
	if service == nil || service.MembersId != int(member.Id) {
		response.Fail(ctx, -400100, "你无权操作", gin.H{})
		return
	}

	err := recordMod.AcceptRecord(int(id), price, startTime, endTime, attach)
	if err != nil {
		response.Fail(ctx, -400100, "操作失败，请联系管理员", gin.H{})
	} else {
		response.Success(ctx, "操作成功", gin.H{})
	}
}

func (m *Services) RecordIndex(ctx *gin.Context) {
	id := ctx.GetFloat64(consts.ValidatorPrefix + "id")
	page := ctx.GetFloat64(consts.ValidatorPrefix + "page")
	limit := ctx.GetFloat64(consts.ValidatorPrefix + "limit")

	count, res := model.CreateServicesRecordFactory("").Index(int(id), int(page), int(limit))

	response.Success(ctx, "SUCCESS", gin.H{
		"list":  res,
		"count": count,
	})
}

func (m *Services) RecordList(ctx *gin.Context) {
	page := ctx.GetFloat64(consts.ValidatorPrefix + "page")
	limit := ctx.GetFloat64(consts.ValidatorPrefix + "limit")
	membersId, _ := cur_userinfo.GetHomeCurrentUserId(ctx)
	count, res := model.CreateServicesRecordFactory("").MyList(int(membersId), int(page), int(limit))

	response.Success(ctx, "SUCCESS", gin.H{
		"list":  res,
		"count": count,
	})
}
