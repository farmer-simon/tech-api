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
	"time"
)

type Needs struct {
}

func (m *Needs) Index(ctx *gin.Context) {
	page := ctx.GetFloat64(consts.ValidatorPrefix + "page")
	limit := ctx.GetFloat64(consts.ValidatorPrefix + "limit")
	offset := math.Max((page-1)*limit, 0)
	var queryData data_type.NeedsQuery
	if err := data_bind.ShouldBindFormDataToModel(ctx, &queryData); err == nil {
		count, res := model.CreateNeedsFactory("").Index(&queryData, int(offset), int(limit))
		response.Success(ctx, "SUCCESS", gin.H{
			"list":  res,
			"count": count,
		})
		return
	}
	response.Fail(ctx, -400100, "参数错误", gin.H{})
}

func (m *Needs) Info(ctx *gin.Context) {
	query, _ := ctx.GetQuery("id")
	id := utils.String2Int(query)
	info := model.CreateNeedsFactory("").GetById(id)
	if info == nil {
		response.Fail(ctx, -400100, "需求不存在", gin.H{})
		return
	}
	if info.State != 3 {
		response.Fail(ctx, -400100, "需求不存在", gin.H{})
		return
	}

	//处理过期时间(返回剩余秒数)
	currentTime := time.Now().Unix()
	info.ExpireTime = int64(math.Max(0, float64(info.ExpireTime-currentTime)))

	author := model.CreateMemberFactory("").GetById(info.MembersId)
	response.Success(ctx, "SUCCESS", gin.H{
		"info": gin.H{
			"id":            info.Id,
			"category_id":   info.CategoryId,
			"category_name": info.CateName,
			"title":         info.Title,
			"price":         info.Price,
			"hits":          info.Hits,
			"expire_time":   info.ExpireTime,
			"content":       info.Content,
			"created_at":    info.CreatedAt,
		},
		"author": gin.H{
			"avatar":    author.Avatar,
			"nick_name": author.NickName,
			"id":        author.Id,
		},
	})
}

func (m *Needs) MyIndex(ctx *gin.Context) {
	page := ctx.GetFloat64(consts.ValidatorPrefix + "page")
	limit := ctx.GetFloat64(consts.ValidatorPrefix + "limit")
	offset := math.Max((page-1)*limit, 0)
	var queryData data_type.NeedsQuery
	if err := data_bind.ShouldBindFormDataToModel(ctx, &queryData); err == nil {
		membersId, _ := cur_userinfo.GetHomeCurrentUserId(ctx)
		queryData.MembersId = int(membersId)
		count, res := model.CreateNeedsFactory("").Index(&queryData, int(offset), int(limit))

		var ids = make([]int, 0)
		var list = make([]gin.H, 0)
		for _, item := range res {
			ids = append(ids, int(item.Id))
		}
		recordTotals := model.CreateNeedsRecordFactory("").GetRecordTotalByNeedsIds(ids)
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

func (m *Needs) MyInfo(ctx *gin.Context) {
	query, _ := ctx.GetQuery("id")
	id := utils.String2Int(query)
	info := model.CreateNeedsFactory("").GetById(id)
	if info == nil {
		response.Fail(ctx, -400100, "需求不存在", gin.H{})
		return
	}
	membersId, _ := cur_userinfo.GetHomeCurrentUserId(ctx)
	if info.MembersId != int(membersId) {
		response.Fail(ctx, -400100, "无权查看", gin.H{})
		return
	}
	response.Success(ctx, "SUCCESS", gin.H{
		"info": info,
	})
}

func (m *Needs) Create(ctx *gin.Context) {

	categoryId := ctx.GetFloat64(consts.ValidatorPrefix + "category_id")
	title := ctx.GetString(consts.ValidatorPrefix + "title")
	content := ctx.GetString(consts.ValidatorPrefix + "content")
	price := ctx.GetFloat64(consts.ValidatorPrefix + "price")
	phone := ctx.GetString(consts.ValidatorPrefix + "phone")
	wechat := ctx.GetString(consts.ValidatorPrefix + "wechat")
	qq := ctx.GetString(consts.ValidatorPrefix + "qq")
	expireTime := ctx.GetString(consts.ValidatorPrefix + "expire_time")

	info, err := cur_userinfo.GetHomeCurrentUserInfo(ctx)
	if err != nil {
		response.Fail(ctx, -400100, err.Error(), gin.H{})
		return
	}
	if info.State != 3 {
		response.Fail(ctx, -400100, "您的帐号未通过官方认证，暂不能发布信息！", gin.H{})
		return
	}
	err = model.CreateNeedsFactory("").Insert(int(info.Id), int(categoryId), price, title, content, phone, wechat, qq, expireTime)
	if err != nil {
		response.Fail(ctx, -400100, "需求发布失败,"+err.Error(), gin.H{})
		return
	}
	response.Success(ctx, "需求发布成功，请等待审核!", gin.H{})
}

func (m *Needs) Edit(ctx *gin.Context) {

	id := ctx.GetFloat64(consts.ValidatorPrefix + "id")
	categoryId := ctx.GetFloat64(consts.ValidatorPrefix + "category_id")
	title := ctx.GetString(consts.ValidatorPrefix + "title")
	content := ctx.GetString(consts.ValidatorPrefix + "content")
	price := ctx.GetFloat64(consts.ValidatorPrefix + "price")
	phone := ctx.GetString(consts.ValidatorPrefix + "phone")
	wechat := ctx.GetString(consts.ValidatorPrefix + "wechat")
	qq := ctx.GetString(consts.ValidatorPrefix + "qq")
	expireTime := ctx.GetString(consts.ValidatorPrefix + "expire_time")

	mod := model.CreateNeedsFactory("")
	membersId, _ := cur_userinfo.GetHomeCurrentUserId(ctx)
	needs := mod.GetById(int(id))
	if needs == nil || needs.MembersId != int(membersId) {
		response.Fail(ctx, -400100, "无权编辑此需求", gin.H{})
		return
	}
	err := mod.Edit(int(id), int(categoryId), price, title, content, phone, wechat, qq, expireTime)
	if err != nil {
		response.Fail(ctx, -400100, "需求编辑失败,"+err.Error(), gin.H{})
		return
	}
	response.Success(ctx, "需求编辑成功，请等待审核!", gin.H{})
}

func (m *Needs) Delete(ctx *gin.Context) {
	query, _ := ctx.GetQuery("id")
	id := utils.String2Int(query)
	mod := model.CreateNeedsFactory("")
	membersId, _ := cur_userinfo.GetHomeCurrentUserId(ctx)
	needs := mod.GetById(int(id))
	if needs == nil || needs.MembersId != int(membersId) {
		response.Fail(ctx, -400100, "无权此操作", gin.H{})
		return
	}
	mod.Delete(int(id))
	response.Success(ctx, "需求删除成功!", gin.H{})
}

func (m *Needs) Contact(ctx *gin.Context) {
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

	info := model.CreateNeedsFactory("").GetById(int(id))
	if info == nil || info.State != 3 {
		response.Fail(ctx, -400100, "需求不存在", gin.H{})
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

func (m *Needs) PostRecord(ctx *gin.Context) {
	id := ctx.GetFloat64(consts.ValidatorPrefix + "id")
	content := ctx.GetString(consts.ValidatorPrefix + "content")
	tenderPrice := ctx.GetString(consts.ValidatorPrefix + "price")
	member, _ := cur_userinfo.GetHomeCurrentUserInfo(ctx)
	err := model.CreateNeedsRecordFactory("").InsertRecord(int(member.Id), int(id), content, tenderPrice)
	if err != nil {
		response.Fail(ctx, -400100, "操作失败，请联系管理员", gin.H{})
	} else {
		response.Success(ctx, "您的报价已发送给需求方，请等待需求方联系您！", gin.H{})
	}
}

func (m *Needs) CloseRecord(ctx *gin.Context) {
	id := ctx.GetFloat64(consts.ValidatorPrefix + "id")
	recordMod := model.CreateNeedsRecordFactory("")
	record := recordMod.GetById(int(id))
	if record == nil {
		response.Fail(ctx, -400100, "操作失败，请联系管理员", gin.H{})
		return
	}
	member, _ := cur_userinfo.GetHomeCurrentUserInfo(ctx)
	service := model.CreateNeedsFactory("").GetById(record.NeedsId)
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
func (m *Needs) AcceptRecord(ctx *gin.Context) {
	id := ctx.GetFloat64(consts.ValidatorPrefix + "id")
	price := ctx.GetString(consts.ValidatorPrefix + "price")
	startTime := ctx.GetString(consts.ValidatorPrefix + "start_time")
	endTime := ctx.GetString(consts.ValidatorPrefix + "end_time")
	attach := ctx.GetString(consts.ValidatorPrefix + "attach")

	recordMod := model.CreateNeedsRecordFactory("")
	record := recordMod.GetById(int(id))
	if record == nil || record.State != 1 {
		response.Fail(ctx, -400100, "操作失败，请联系管理员", gin.H{})
		return
	}
	member, _ := cur_userinfo.GetHomeCurrentUserInfo(ctx)
	service := model.CreateNeedsFactory("").GetById(record.NeedsId)
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

func (m *Needs) RecordIndex(ctx *gin.Context) {
	id := ctx.GetFloat64(consts.ValidatorPrefix + "id")
	page := ctx.GetFloat64(consts.ValidatorPrefix + "page")
	limit := ctx.GetFloat64(consts.ValidatorPrefix + "limit")

	count, res := model.CreateNeedsRecordFactory("").Index(int(id), int(page), int(limit))

	response.Success(ctx, "SUCCESS", gin.H{
		"list":  res,
		"count": count,
	})
}

func (m *Needs) RecordList(ctx *gin.Context) {
	page := ctx.GetFloat64(consts.ValidatorPrefix + "page")
	limit := ctx.GetFloat64(consts.ValidatorPrefix + "limit")
	membersId, _ := cur_userinfo.GetHomeCurrentUserId(ctx)
	count, res := model.CreateNeedsRecordFactory("").MyList(int(membersId), int(page), int(limit))

	response.Success(ctx, "SUCCESS", gin.H{
		"list":  res,
		"count": count,
	})
}
