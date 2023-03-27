package model

import "errors"

//UserRole 后台角色
type UserRole struct {
	BaseModel
	Name string `gorm:"column:name" json:"name"`
}

//RoleUsers 后台角色下的用户
type RoleUsers struct {
	BaseModel
	RoleId int `gorm:"column:role_id" json:"role_id"`
	UserId int `gorm:"column:users_id" json:"user_id"`
}

func CreateUserRoleFactory(sqlType string) *UserRole {
	return &UserRole{BaseModel: BaseModel{DB: UseDbConn(sqlType)}}
}

// TableName 表名
func (u *UserRole) TableName() string {
	return "teach_auth_role"
}

// RoleUserTableName 表名
func (u *UserRole) RoleUserTableName() string {
	return "teach_auth_role_users"
}

// AddRole 添加角色
func (u *UserRole) AddRole(name string) (int64, error) {
	var newRole UserRole
	newRole.Name = name
	if res := u.Model(u).Create(&newRole); res != nil {
		return newRole.Id, nil
	}
	return 0, errors.New("添加角色失败")
}

//CheckRepeatRoleName 更新前检查新的角色名是否已经存在（避免和别的重名）
func (u *UserRole) CheckRepeatRoleName(excludeId int, name string) (exists int) {
	sql := "select count(*) as counts from teach_auth_role where  id!=?  AND name=?"
	_ = u.Raw(sql, excludeId, name).First(&exists)
	return exists
}

// EditRole 编辑角色
func (u *UserRole) EditRole(roleId int, name string) error {
	if res := u.Model(u).Where("id=?", roleId).Update("name", name); res != nil {
		return nil
	}
	return errors.New("编辑角色失败")
}

// DeleteRole 删除角色
func (u *UserRole) DeleteRole(roleId int) {
	//删除角色
	u.Model(u).Where("id=?", roleId).Delete(&UserRole{})
	////删除角色下的成员
	//u.Table(u.RoleUserTableName()).Where("role_id=?", roleId).Delete(&RoleUsers{})
	return
}

// GetRoleUsers 返回角色下的所有用户信息
func (u *UserRole) GetRoleUsers(roleId int) (users []RoleUsers) {
	u.Table(u.RoleUserTableName()).Where("role_id=?", roleId).Scan(&users)
	return
}

// GetRoleIdsByUserId 返回用户的所有角色关系
func (u *UserRole) GetRoleIdsByUserId(usersId int64) (roleIds []int) {
	var RoleUsers []RoleUsers
	u.Table(u.RoleUserTableName()).Where("users_id=?", usersId).Find(&RoleUsers)
	for _, role := range RoleUsers {
		roleIds = append(roleIds, role.RoleId)
	}
	return
}

//CheckRepeatRoleUser 检查角色下是否存在当前用户
func (u *UserRole) CheckRepeatRoleUser(roleId, usersId int) (exists int) {
	sql := "select count(*) as counts from teach_auth_role_users where  role_id=?  AND users_id=?"
	_ = u.Raw(sql, roleId, usersId).First(&exists)
	return exists
}

//AddRoleUser 添加角色成员
func (u *UserRole) AddRoleUser(roleId, usersId int) error {
	var newRoleUsers RoleUsers
	newRoleUsers.RoleId = roleId
	newRoleUsers.UserId = usersId
	if res := u.Table(u.RoleUserTableName()).Create(&newRoleUsers); res != nil {
		return errors.New("角色分配失败")
	}
	return nil

}

//DeleteRoleUser 删除角色成员
func (u *UserRole) DeleteRoleUser(roleId, usersId int) {
	u.Table(u.RoleUserTableName()).Where("role_id=? AND users_id=?", roleId, usersId).Delete(&RoleUsers{})
}
