package model

import (
	"github.com/farmer-simon/go-utils"
	"strings"
)

func CreateAttrsFactory(sqlType string) *AttrsModel {
	return &AttrsModel{BaseModel: BaseModel{DB: UseDbConn(sqlType)}}
}

type AttrsModel struct {
	BaseModel
	Name      string `json:"name"  gorm:"column:name"`
	Type      string `json:"type"  gorm:"column:type"`
	Src       string `json:"src"  gorm:"column:src"`
	TargetId  int    `json:"target_id"  gorm:"column:target_id"`
	MembersId int    `json:"members_id"  gorm:"column:members_id"`
}

// TableName 表名
func (mod *AttrsModel) TableName() string {
	return "tech_attachment"
}

func (mod *AttrsModel) Insert(src, name, srcType string, targetId, membersId int) (*AttrsModel, error) {
	mod.Name = name
	mod.Src = src
	mod.Type = srcType
	mod.TargetId = targetId
	mod.MembersId = membersId
	if res := mod.Omit("CreatedAt", "UpdatedAt").Create(&mod); res.Error != nil {
		return nil, res.Error
	}
	return mod, nil
}

func (mod *AttrsModel) GetByTargetId(id int, srcType string) (res []AttrsModel) {
	mod.Table(mod.TableName()).Where("target_id=? AND type=?", id, srcType).Find(&res)
	return
}

func (mod *AttrsModel) UpdateTargetIdByIds(ids string, targetId int) bool {
	var attrIds []int
	for _, id := range strings.Split(ids, ",") {
		attrIds = append(attrIds, utils.String2Int(id))
	}
	//清除历史的
	mod.Table(mod.TableName()).Where("target_id =?", targetId).Update("target_id", 0)
	//设置新的
	mod.Table(mod.TableName()).Where("id IN(?)", attrIds).Update("target_id", targetId)
	return true
}

func (mod *AttrsModel) GetCoverByTargetIds(targetIds []int, srcType string) map[int]string {
	var ids []int
	var cover = make(map[int]string, 0)
	for _, id := range targetIds {
		ids = append(ids, id)
		cover[id] = ""
	}
	var res []AttrsModel
	mod.Table(mod.TableName()).Where("target_id IN ? AND type=?", ids, srcType).Find(&res)
	for _, item := range res {
		if cover[item.TargetId] == "" {
			cover[item.TargetId] = item.Src
		}
	}
	return cover
}
