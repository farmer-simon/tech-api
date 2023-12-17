package model

import (
	"bytes"
	"errors"
	"fmt"
	"github.com/farmer-simon/go-utils"
	"goskeleton/app/data_type"
)

func CreateMemberFactory(sqlType string) *MemberModel {
	return &MemberModel{BaseModel: BaseModel{DB: UseDbConn(sqlType)}}
}

type MemberModel struct {
	BaseModel
	data_type.MembersBase
	data_type.MembersExtendBase
}

// TableName 表名
func (mod *MemberModel) TableName() string {
	return "tech_members"
}

// Index 后台列表
func (mod *MemberModel) Index(q *data_type.MembersQuery, offset, limit int) (count int64, res []MemberModel) {
	var sql bytes.Buffer
	sql.WriteString(" FROM tech_members WHERE is_del=0 ")
	if q.Phone != "" {
		sql.WriteString(fmt.Sprintf(" AND phone = %s", q.Phone))
	}

	if q.StartTime != "" {
		sql.WriteString(fmt.Sprintf(" AND created_at>=%s", q.StartTime))
	}
	if q.EndTime != "" {
		sql.WriteString(fmt.Sprintf(" AND created_at<%s", q.EndTime))
	}
	if q.NickName != "" {
		sql.WriteString(fmt.Sprintf(" AND nick_name = %s", q.NickName))
	}
	if q.RealName != "" {
		sql.WriteString(fmt.Sprintf(" AND real_name = %s", q.RealName))
	}

	if q.TeamName != "" {
		sql.WriteString(" AND team_name LIKE '%" + q.TeamName + "%'")
	}
	if q.TeamMajor != "" {
		sql.WriteString(" AND team_majors LIKE '%" + q.TeamMajor + "%'")
	}
	if q.State > 0 {
		sql.WriteString(fmt.Sprintf(" AND state = %d", int(q.State)))
	}

	mod.Raw("SELECT COUNT(0) " + sql.String()).Count(&count)
	mod.Raw("SELECT * "+sql.String()+" ORDER BY  id DESC LIMIT ?,? ", offset, limit).Find(&res)
	return
}

// Insert 注册
func (mod *MemberModel) Insert(phone, nickName, realName, avatar, teamIntro, teamName, teamMajor, qq, wechat string) error {
	state := 2
	settings := CreateSettingsFactory("").GetSettings()
	if utils.String2Int(settings["user_register_state"].(string)) == 2 {
		state = 3
	}

	sql := "INSERT INTO tech_members(phone,passwd, nick_name, real_name, avatar,team_intro, team_name, team_major, qq, wechat, state) VALUES(?,?,?,?,?,?,?,?,?,?,?)"
	if res := mod.Exec(sql, phone, "", nickName, realName, avatar, teamIntro, teamName, teamMajor, qq, wechat, state); res.Error != nil {
		return res.Error
	}
	return nil
}

// UpdateExtendInfo 更新扩展资料
func (mod *MemberModel) UpdateExtendInfo(membersId int, data data_type.MembersExtendBase) error {
	sql := "UPDATE tech_members SET nick_name=?, real_name=?, avatar=?, team_intro=?, team_name=?,team_major=?, qq=?, wechat=? WHERE id=?"
	if res := mod.Exec(sql, data.NickName, data.RealName, data.Avatar, data.TeamIntro, data.TeamName, data.TeamMajor, data.Qq, data.WeChat, membersId); res.Error != nil {
		return res.Error
	}
	return nil
}

// GetByPhone 根据电话查信息
func (mod *MemberModel) GetByPhone(phone string) *MemberModel {
	sql := "select *  from tech_members where  phone=? and is_del=0 limit 1"
	result := mod.Raw(sql, phone).First(&mod)
	if result.Error == nil && mod.Id > 0 {
		return mod
	}
	return nil
}

// GetById 根据用户ID查询一条信息
func (mod *MemberModel) GetById(id int) *MemberModel {
	sql := "SELECT  * FROM  `tech_members`  WHERE `is_del`=0 and   id=? LIMIT 1"
	result := mod.Raw(sql, id).First(mod)
	if result.Error == nil && mod.Id > 0 {
		return mod
	}
	return nil
}

// QuickCreateByPhone 根据手机号快速创建用户
func (mod *MemberModel) QuickCreateByPhone(phone string) (*MemberModel, error) {
	sql := "INSERT  INTO tech_members(phone, nick_name, passwd, state,avatar) VALUES(?,?,?,?,?)"
	if mod.Exec(sql, phone, fmt.Sprintf("手机用户%s", phone[7:11]), "", 1, "/public/storage/resources/avatar.png").RowsAffected > 0 {
		return mod.GetByPhone(phone), nil
	}
	return nil, errors.New("帐号信息写入失败")
}

// BindPhone 绑定手机号
func (mod *MemberModel) BindPhone(id int, phone string) bool {
	sql := "UPDATE tech_members SET phone=? WHERE id=? LIMIT 1"
	if mod.Exec(sql, phone, id).RowsAffected > 0 {
		return true
	}
	return false
}

func (mod *MemberModel) SetPasswd(id int, passwd string) bool {
	sql := "UPDATE tech_members SET passwd=? WHERE id=? LIMIT 1"
	if mod.Exec(sql, passwd, id).RowsAffected > 0 {
		return true
	}
	return false
}

// Verify 更新状态
func (mod *MemberModel) Verify(id, state int, rejectReason string) bool {
	sql := "UPDATE tech_members SET state=? ,reject_reason=?  WHERE id=? LIMIT 1"
	if mod.Exec(sql, state, rejectReason, id).RowsAffected > 0 {
		return true
	}
	return false
}

func (mod *MemberModel) SSOAutoLogin(userId string) int64 {
	mod.Model(mod).Where("sso_userid=?", userId).Find(&mod)
	if mod.Id != 0 {
		return mod.Id
	}
	//自动从清洗的表中摘取数据
	autoInsertId := CreateUserCenterFactory("").AutoSyncUser(userId)
	if autoInsertId != 0 {
		return autoInsertId
	}
	return 0
}
