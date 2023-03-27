package members

import (
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"goskeleton/app/data_type"
	"goskeleton/app/global/consts"
	"goskeleton/app/model"
	"goskeleton/app/service/sms"
	"goskeleton/app/utils/cur_userinfo"
	"goskeleton/app/utils/data_bind"
	"goskeleton/app/utils/md5_encrypt"
	"math"
)

func CreateMemberServiceFactory() *MemberService {
	return &MemberService{model.CreateMemberFactory("")}
}

type MemberService struct {
	model *model.MemberModel
}

// Index 后台用户列表
func (serv *MemberService) Index(ctx *gin.Context) (count int64, res []model.MemberModel) {
	page := ctx.GetFloat64(consts.ValidatorPrefix + "page")
	limit := ctx.GetFloat64(consts.ValidatorPrefix + "limit")
	offset := math.Max((page-1)*limit, 0)
	var queryData data_type.MembersQuery
	if err := data_bind.ShouldBindFormDataToModel(ctx, &queryData); err == nil {
		count, res = serv.model.Index(&queryData, int(offset), int(limit))
		return
	}
	return
}

// Register 用户注册
func (serv *MemberService) Register(ctx *gin.Context) (*model.MemberModel, error) {
	phone := ctx.GetString(consts.ValidatorPrefix + "phone")
	nickName := ctx.GetString(consts.ValidatorPrefix + "nick_name")
	realName := ctx.GetString(consts.ValidatorPrefix + "real_name")
	Avatar := ctx.GetString(consts.ValidatorPrefix + "avatar")
	teamIntro := ctx.GetString(consts.ValidatorPrefix + "team_intro")
	teamName := ctx.GetString(consts.ValidatorPrefix + "team_name")
	teamMajor := ctx.GetString(consts.ValidatorPrefix + "team_major")
	qq := ctx.GetString(consts.ValidatorPrefix + "qq")
	wechat := ctx.GetString(consts.ValidatorPrefix + "wechat")

	//检查手机号是否存在
	member := serv.model.GetByPhone(phone)
	if member != nil {
		return nil, errors.New("手机号已存在")
	}
	err := serv.model.Insert(phone, nickName, realName, Avatar, teamIntro, teamName, teamMajor, qq, wechat)
	if err != nil {
		return nil, err
	}

	member = serv.model.GetByPhone(phone)
	if member != nil {
		return member, nil
	}

	return nil, errors.New("注册失败")
}

// PhoneLogin 手机验证码登录
func (serv *MemberService) PhoneLogin(phone, code string) (*model.MemberModel, error) {
	//检查验证码
	smsClient := sms.CreatePhoneCodeFactory(phone)
	pass := smsClient.CheckPhoneCode(code, true)
	if !pass {
		return nil, errors.New("验证码错误")
	}
	//检查手机号是否存在
	member := serv.model.GetByPhone(phone)
	if member != nil {
		return member, nil
	}
	return nil, errors.New("未注册帐号，请先注册！")
}

// PasswdLogin 手机号密码登录
func (serv *MemberService) PasswdLogin(phone, passwd string) (*model.MemberModel, error) {
	member := serv.model.GetByPhone(phone)
	if member == nil {
		return nil, errors.New("登录失败")
	}
	fmt.Println(md5_encrypt.Base64Md5(passwd))
	if member.Passwd == md5_encrypt.Base64Md5(passwd) {
		return member, nil
	}
	return nil, errors.New("登录失败")
}

// BindPhone 绑定手机号
func (serv *MemberService) BindPhone(memberId int, phone, code string) error {
	pass := sms.CreatePhoneCodeFactory(phone).CheckPhoneCode(code, true)
	if !pass {
		return errors.New("短信验证码错误")
	}

	member := serv.model.GetById(memberId)
	if member == nil {
		return errors.New("用户信息查询失败，请联系技术！")
	}
	if member.Phone != "" {
		return errors.New("账号已绑定手机号！")
	}
	member = serv.model.GetByPhone(phone)
	if member != nil {
		return errors.New("手机号已绑定在其他帐号")
	}

	serv.model.BindPhone(memberId, phone)
	return nil
}

// ChangePhone 更换绑定手机号
func (serv *MemberService) ChangePhone(memberId int, oldPhone, newPhone, code string) error {
	pass := sms.CreatePhoneCodeFactory(oldPhone).CheckPhoneCode(code, true)
	if !pass {
		return errors.New("短信验证码错误")
	}

	member := serv.model.GetById(memberId)
	if member == nil {
		return errors.New("用户信息查询失败，请联系技术！")
	}
	if member.Phone != oldPhone {
		return errors.New("旧手机号错误！")
	}
	member = serv.model.GetByPhone(newPhone)
	if member != nil && int(member.Id) != memberId {
		return errors.New("手机号已绑定在其他帐号")
	}

	serv.model.BindPhone(memberId, newPhone)
	return nil
}

//SetPasswd 设置密码
func (serv *MemberService) SetPasswd(membersId int, code, passwd string) error {
	member := serv.model.GetById(membersId)
	if member == nil {
		return errors.New("用户信息查询失败，请联系技术！")
	}
	if member.Phone == "" {
		return errors.New("请先绑定手机号！")
	}
	pass := sms.CreatePhoneCodeFactory(member.Phone).CheckPhoneCode(code, true)
	if !pass {
		return errors.New("短信验证码错误")
	}
	pass = serv.model.SetPasswd(int(member.Id), md5_encrypt.Base64Md5(passwd))
	if !pass {
		return errors.New("设置密码错误")
	}
	return nil
}

//ChangePasswd 修改密码
func (serv *MemberService) ChangePasswd(membersId int, OldPasswd, newPasswd string) error {
	member := serv.model.GetById(membersId)
	if member == nil {
		return errors.New("用户信息查询失败，请联系技术！")
	}
	if member.Phone == "" {
		return errors.New("请先绑定手机号！")
	}
	if member.Passwd != md5_encrypt.Base64Md5(OldPasswd) {
		return errors.New("旧密码错误！")
	}

	pass := serv.model.SetPasswd(int(member.Id), md5_encrypt.Base64Md5(newPasswd))
	if !pass {
		return errors.New("设置密码错误")
	}
	return nil
}

// UpdateExtendInfo 更新扩展资料
func (serv *MemberService) UpdateExtendInfo(ctx *gin.Context) error {

	var data data_type.MembersExtendBase
	if err := data_bind.ShouldBindFormDataToModel(ctx, &data); err == nil {
		membersId, isOk := cur_userinfo.GetHomeCurrentUserId(ctx)
		if !isOk || membersId == 0 {
			return errors.New("用户资料更新失败")
		}
		err = serv.model.UpdateExtendInfo(int(membersId), data)
		return err
	}
	return errors.New("用户资料更新失败")
}
