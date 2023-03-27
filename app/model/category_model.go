package model

import (
	"fmt"
	"goskeleton/app/data_type"
)

func CreateCateModelFactory(sqlType string) *CateModel {
	return &CateModel{BaseModel: BaseModel{DB: UseDbConn(sqlType)}}
}

type CateModel struct {
	BaseModel
	data_type.CategoryBase
}

// TableName 表名
func (mod *CateModel) TableName() string {
	return "tech_category"
}

// Insert 插入分类
func (mod *CateModel) Insert(data CateModel) error {
	if res := mod.Model(mod).Create(&data); res.Error != nil {
		return res.Error
	}
	return nil
}

// Update 编辑分类
func (mod *CateModel) Update(data *CateModel) error {
	if res := mod.Model(mod).Where("id=?", data.Id).Omit("CreatedAt").Save(data); res.Error != nil {
		return res.Error
	}
	return nil
}

//Destroy 删除分类
func (mod *CateModel) Destroy(id int) bool {
	mod.Model(mod).Where("id=?", id).Update("is_del", 1)
	return true
}

// CheckRepeatName  检查重复名称
func (mod *CateModel) CheckRepeatName(excludeId, parentId int, name string) bool {
	var count int64
	mod.Model(mod).Where("parent_id=? AND name=? AND id<>? AND is_del=0", parentId, name, excludeId).Count(&count)
	if count > 0 {
		return true
	}
	return false

}

func (mod *CateModel) GetById(id int) *CateModel {
	var res *CateModel
	mod.Model(mod).Where("id=?", id).Find(&res)
	if res.Id > 0 {
		return res
	}
	return nil
}

// GetList 分类列表
func (mod *CateModel) GetList(parentId, state, offset, limit int) (count int64, res []CateModel) {
	sql := `SELECT %s FROM tech_category WHERE is_del=0`
	if state == 1 || state == 2 {
		sql = sql + fmt.Sprintf(" AND state = %d", state)
	}
	if parentId > 0 {
		sql = sql + fmt.Sprintf(" AND parent_id = %d", parentId)
	}
	mod.Raw(fmt.Sprintf(sql, "COUNT(0)")).Count(&count)
	sql = fmt.Sprintf(sql, "*") + " ORDER BY sort ASC, id ASC"
	//当limit=0的时候返回全部
	if limit > 0 {
		sql = sql + fmt.Sprintf(" limit %d, %d", offset, limit)
	}
	mod.Debug().Raw(sql).Find(&res)
	return
}

//GetRecommendList 获取推荐分类
func (mod *CateModel) GetRecommendList(parentId, limit int) (res []CateModel) {
	sql := fmt.Sprintf("SELECT * FROM tech_category WHERE is_del=0 AND state=2 AND parent_id = %d ORDER BY sort ASC, id ASC LIMIT %d", parentId, limit)
	mod.Raw(sql).Find(&res)
	return
}

func (mod *CateModel) GetAllOnlineCategory() (res []CateModel) {
	mod.Raw("SELECT * FROM tech_category WHERE is_del=0 AND state=1 ORDER BY sort ASC, id ASC").Find(&res)
	return
}

func (mod *CateModel) GetSubIds(topCid int) []int {
	var cateIds = make([]int, 0)
	_, cateList := mod.GetList(topCid, 1, 0, 0)
	for _, item := range cateList {
		cateIds = append(cateIds, int(item.Id))
	}
	return cateIds
}
