package model

import (
	"bytes"
	"fmt"
	"github.com/farmer-simon/go-utils"
	"goskeleton/app/data_type"
	"math"
	"time"
)

func CreateServicesFactory(sqlType string) *ServicesModel {
	return &ServicesModel{BaseModel: BaseModel{DB: UseDbConn(sqlType)}}
}

type ServicesModel struct {
	BaseModel
	data_type.ServicesMod
}

// TableName 表名
func (mod *ServicesModel) TableName() string {
	return "tech_services"
}

//Index 列表
func (mod *ServicesModel) Index(q *data_type.ServicesQuery, offset, limit int) (count int64, res []NeedsModel) {
	var sql bytes.Buffer
	sql.WriteString(" FROM tech_services n LEFT JOIN tech_members m ON n.members_id=m.id LEFT JOIN tech_category c ON n.category_id=c.id WHERE n.is_del=0 ")
	sort := " ORDER BY n.id DESC"
	if q.QueryType == "my" {
		sql.WriteString(fmt.Sprintf(" AND n.members_id = %d", q.MembersId))
		if q.State == 1 {
			sql.WriteString(" AND n.state IN (1,2)")
		} else if q.State == 3 {
			sql.WriteString(" AND n.state = 3")
		}
	} else if q.QueryType == "home" {
		sql.WriteString(" AND n.state = 3")
		if q.Cid != 0 {
			sql.WriteString(fmt.Sprintf(" AND n.category_id = %d", int(q.Cid)))
		}
		if q.Title != "" {
			sql.WriteString(" AND n.title LIKE '%" + q.Title + "%'")
		}
		if q.Recommend == 1 {
			sql.WriteString(" AND n.recommend = 1")
		}
		if q.Sort == "sale" {
			sort = "ORDER BY n.hits DESC"
		}
	} else if q.QueryType == "admin" {
		if q.Cid != 0 {
			sql.WriteString(fmt.Sprintf(" AND n.category_id = %d", int(q.Cid)))
		}
		if q.StartTime != "" {
			sql.WriteString(fmt.Sprintf(" AND n.created_at>=%s", q.StartTime))
		}
		if q.EndTime != "" {
			sql.WriteString(fmt.Sprintf(" AND n.created_at<%s", q.EndTime))
		}
		if q.State != 0 {
			sql.WriteString(fmt.Sprintf(" AND n.state =%d", int(q.State)))
		}
		if q.Phone != "" {
			sql.WriteString(fmt.Sprintf(" AND n.phone = %s", q.Phone))
		}
		if q.Name != "" {
			sql.WriteString(fmt.Sprintf(" AND (m.nick_name = '%s' OR m.real_name='%s')", q.Name, q.Name))
		}
		if q.Title != "" {
			sql.WriteString(" AND n.title LIKE '%" + q.Title + "%'")
		}
	}

	mod.Raw("SELECT COUNT(0) " + sql.String()).Count(&count)
	mod.Debug().Raw("SELECT n.*, m.nick_name, m.avatar,c.name AS cate_name "+sql.String()+sort+" LIMIT ?,? ", offset, limit).Find(&res)
	//处理过期时间(返回剩余秒数)
	currentTime := time.Now().Unix()
	for i, item := range res {
		item.ExpireTime = int64(math.Max(0, float64(item.ExpireTime-currentTime)))
		res[i] = item
	}

	return
}

//Insert 添加
func (mod *ServicesModel) Insert(membersId, categoryId int, price float64, title, content, phone, wechat, qq string) (*ServicesModel, error) {
	mod.MembersId = membersId
	mod.CategoryId = categoryId
	mod.Price = price
	mod.Title = title
	mod.Content = content
	mod.Phone = phone
	mod.Wechat = wechat
	mod.Qq = qq
	mod.Recommend = 0
	mod.State = 2
	settings := CreateSettingsFactory("").GetSettings()
	if utils.String2Int(settings["content_publish_state"].(string)) == 2 {
		mod.State = 3
	}
	if res := mod.Omit("NickName", "Avatar", "CateName", "CreatedAt", "UpdatedAt").Create(&mod); res.Error != nil {
		return nil, res.Error
	}
	return mod, nil
}

// Edit 更新扩展资料
func (mod *ServicesModel) Edit(id, categoryId int, price float64, title, content, phone, wechat, qq string) error {
	state := 2
	settings := CreateSettingsFactory("").GetSettings()
	if utils.String2Int(settings["content_publish_state"].(string)) == 2 {
		state = 3
	}
	sql := "UPDATE tech_services SET category_id=?, price=?, title=?, content=?,phone=?, wechat=?, qq=?, state=? WHERE id=?"
	if res := mod.Exec(sql, categoryId, price, title, content, phone, wechat, qq, state, id); res.Error != nil {
		return res.Error
	}
	return nil
}

//GetById 根据用户ID查询一条信息
func (mod *ServicesModel) GetById(id int) *ServicesModel {
	sql := "SELECT  n.*, c.name AS cate_name FROM  `tech_services` n LEFT JOIN `tech_category` c ON n.category_id=c.id  WHERE n.`is_del`=0 and n.id=? LIMIT 1"
	result := mod.Raw(sql, id).First(&mod)
	if result.Error == nil && mod.Id > 0 {
		return mod
	}
	return nil
}

//Verify 更新状态
func (mod *ServicesModel) Verify(id, state int, rejectReason string) bool {
	sql := "UPDATE tech_services SET state=? ,reject_reason=?  WHERE id=? LIMIT 1"
	if mod.Debug().Exec(sql, state, rejectReason, id).RowsAffected > 0 {
		return true
	}
	return false
}

//Delete 更新状态
func (mod *ServicesModel) Delete(id int) bool {
	sql := "UPDATE tech_services SET is_del=1  WHERE id=? LIMIT 1"
	if mod.Exec(sql, id).RowsAffected > 0 {
		return true
	}
	return false
}
