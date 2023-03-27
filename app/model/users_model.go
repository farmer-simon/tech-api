package model

import (
	"bytes"
	"fmt"
	"go.uber.org/zap"
	"goskeleton/app/data_type"
	"goskeleton/app/global/variable"
	"goskeleton/app/utils/md5_encrypt"
)

// 操作数据库喜欢使用gorm自带语法的开发者可以参考 GinSkeleton-Admin 系统相关代码
// Admin 项目地址：https://gitee.com/daitougege/gin-skeleton-admin-backend/
// gorm_v2 提供的语法+ ginskeleton 实践 ：  http://gitee.com/daitougege/gin-skeleton-admin-backend/blob/master/app/model/button_cn_en.go

// 创建 userFactory
// 参数说明： 传递空值，默认使用 配置文件选项：UseDbType（mysql）

func CreateUserFactory(sqlType string) *UsersModel {
	return &UsersModel{BaseModel: BaseModel{DB: UseDbConn(sqlType)}}
}

type UsersModel struct {
	BaseModel
	data_type.UsersBase
}

// TableName 表名
func (u *UsersModel) TableName() string {
	return "tech_users"
}

func (u *UsersModel) Index(name string, state, offset, limit int) (res []UsersModel, count int64) {
	var sql bytes.Buffer
	sql.WriteString(" FROM tech_users WHERE is_del=0 ")
	if name != "" {
		sql.WriteString(" AND name LIKE '%" + name + "%'")
	}
	if state != 0 {
		sql.WriteString(fmt.Sprintf(" AND state =%d", state))
	}
	u.Debug().Raw("SELECT COUNT(0) " + sql.String()).Count(&count)
	u.Debug().Raw("SELECT * "+sql.String()+" ORDER BY id DESC LIMIT ?,? ", offset, limit).Find(&res)
	return
}

// Login 用户登录,
func (u *UsersModel) Login(userName string, passwd string) *UsersModel {
	sql := "select *  from tech_users where  user_name=? and is_del=0 limit 1"
	result := u.Raw(sql, userName).First(&u)
	if result.Error == nil {
		// 账号密码验证成功
		if len(u.Passwd) > 0 && (u.Passwd == md5_encrypt.Base64Md5(passwd)) {
			return u
		}
	} else {
		variable.ZapLog.Error("根据账号查询单条记录出错:", zap.Error(result.Error))
	}
	return nil
}

//Create 新增
func (u *UsersModel) Create(userName, name, passwd, phone, avatar string, state int) bool {
	sql := "INSERT  INTO tech_users(user_name,name,passwd,phone,avatar, state) VALUES(?,?,?,?,?,?)"
	if u.Exec(sql, userName, name, passwd, phone, avatar, state).RowsAffected > 0 {
		return true
	}
	return false
}

//CheckRepeatUserName 更新前检查新的用户名是否已经存在（避免和别的账号重名）
func (u *UsersModel) CheckRepeatUserName(excludeId int, userName string) (exists int) {
	sql := "select count(*) as counts from tech_users where  id!=? AND is_del=0 AND user_name=?"
	_ = u.Raw(sql, excludeId, userName).First(&exists)
	return exists
}

//Update 更新
func (u *UsersModel) Update(id int, userName, name, passwd, phone, avatar string, state int) bool {
	if passwd == "" {
		sql := "update tech_users set user_name=?, name=?,phone=?,avatar=?, state=?  WHERE is_del=0 AND id=?"
		if u.Exec(sql, userName, name, phone, avatar, state, id).RowsAffected >= 0 {
			return true
		}
	} else {
		sql := "update tech_users set user_name=?, name=?,passwd=?,phone=?,avatar=?, state=?  WHERE is_del=0 AND id=?"
		if u.Exec(sql, userName, name, passwd, phone, avatar, state, id).RowsAffected >= 0 {
			return true
		}
	}
	return false
}

// ChangePasswd 改密码
func (u *UsersModel) ChangePasswd(id int, passwd string) bool {
	sql := "update tech_users set passwd=? WHERE is_del=0 AND id=?"
	if u.Exec(sql, passwd, id).RowsAffected >= 0 {
		return true
	}
	return false
}

// Destroy 删除
func (u *UsersModel) Destroy(id int) bool {
	sql := "update tech_users set is_del=1 WHERE id=?"
	if u.Exec(sql, id).RowsAffected >= 0 {
		return true
	}
	return false
}

//GetUserById 根据用户ID查询一条信息
func (u *UsersModel) GetUserById(id int) (*UsersModel, error) {
	sql := "SELECT  * FROM  `tech_users`  WHERE `is_del`=0 and   id=? LIMIT 1"
	result := u.Raw(sql, id).First(&u)
	if result.Error == nil && u.Id > 0 {
		return u, nil
	} else {
		return nil, result.Error
	}
}
