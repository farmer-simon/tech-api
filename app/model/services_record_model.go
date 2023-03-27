package model

import (
	"bytes"
	"fmt"
	"github.com/farmer-simon/go-utils"
	"goskeleton/app/data_type"
	"math"
)

func CreateServicesRecordFactory(sqlType string) *ServicesRecordModel {
	return &ServicesRecordModel{BaseModel: BaseModel{DB: UseDbConn(sqlType)}}
}

type ServicesRecordModel struct {
	BaseModel
	data_type.ServicesRecordMod
}

// TableName 表名
func (mod *ServicesRecordModel) TableName() string {
	return "tech_services_record"
}

//InsertRecord 提交需求
func (mod *ServicesRecordModel) InsertRecord(membersId, servicesId int, content string) error {
	sql := "INSERT INTO `tech_services_record`(members_id, services_id, content, price, start_time, end_time, attach, state) VALUES(?,?,?,?,?,?,?,?)"
	res := mod.Exec(sql, membersId, servicesId, content, "", 0, 0, "", 1)
	return res.Error
}

//CloseRecord 关闭需求
func (mod *ServicesRecordModel) CloseRecord(id int) error {
	sql := "UPDATE `tech_services_record` SET state=? WHERE id=?"
	res := mod.Exec(sql, 2, id)
	return res.Error
}

//AcceptRecord 接受需求
func (mod *ServicesRecordModel) AcceptRecord(id int, price, startTime, endTime, attach string) error {
	sTime := utils.DatetimeToUnixTimestamp(startTime + " 00:00:00")
	eTime := utils.DatetimeToUnixTimestamp(endTime + " 23:59:59")
	sql := "UPDATE `tech_services_record` SET price=?, start_time=?, end_time=?, attach=?, state=? WHERE id=?"
	res := mod.Exec(sql, price, sTime, eTime, attach, 3, id)
	return res.Error
}

//GetRecordTotalByServicesIds 根据服务IDS获取各服务的需求数量
func (mod *ServicesRecordModel) GetRecordTotalByServicesIds(ids []int) map[int]int {
	var totals = make(map[int]int, 0)
	var results []data_type.RecordTotal

	sql := "SELECT services_id AS id, COUNT(0) AS total FROM `tech_services_record` WHERE services_id IN (?) GROUP BY services_id"
	mod.Raw(sql, ids).Scan(&results)
	for _, item := range results {
		totals[item.Id] = item.Total
	}
	return totals
}

//Index 根据服务IDS获取各服务的需求数量
func (mod *ServicesRecordModel) Index(id, page, limit int) (count int64, res []ServicesRecordModel) {
	offset := math.Max((float64(page)-1)*float64(limit), 0)
	mod.Raw("SELECT COUNT(0) FROM `tech_services_record` WHERE services_id=?", id).Count(&count)
	mod.Raw("SELECT r.*, m.nick_name, m.qq, m.phone, m.wechat  FROM tech_services_record r LEFT JOIN tech_members m ON r.members_id=m.id WHERE r.services_id=? LIMIT ?,? ", id, int(offset), limit).Find(&res)
	return
}

func (mod *ServicesRecordModel) MyList(membersId, page, limit int) (count int64, res []ServicesRecordModel) {
	offset := math.Max((float64(page)-1)*float64(limit), 0)
	mod.Raw("SELECT COUNT(0) FROM `tech_services_record` WHERE members_id=?", membersId).Count(&count)
	mod.Raw("SELECT r.*, s.title, m.nick_name as author  FROM tech_services_record r LEFT JOIN tech_services s ON r.services_id=r.id LEFT JOIN tech_members m ON s.members_id=m.id WHERE r.members_id=? LIMIT ?,? ", membersId, int(offset), limit).Find(&res)
	return
}

func (mod *ServicesRecordModel) GetById(id int) (res *ServicesRecordModel) {
	mod.Raw("SELECT r.*, m.nick_name, m.qq, m.phone, m.wechat  FROM tech_services_record r LEFT JOIN tech_members m ON r.members_id=m.id WHERE r.id=?", id).Find(&res)
	return
}

//AdminIndex 后台的查看
func (mod *ServicesRecordModel) AdminIndex(servicesId, state, page, limit int) (count int64, res []ServicesRecordModel) {
	offset := math.Max((float64(page)-1)*float64(limit), 0)

	var sql bytes.Buffer
	sql.WriteString(" FROM tech_services_record sr LEFT JOIN tech_members m ON sr.members_id=m.id LEFT JOIN tech_services s ON sr.services_id=s.id LEFT JOIN tech_members sm ON s.members_id=sm.id WHERE ")
	if state == 0 {
		sql.WriteString(" sr.state IN (1,2,3)")
	} else {
		sql.WriteString(fmt.Sprintf(" sr.state = %d", state))
	}
	if servicesId > 0 {
		sql.WriteString(fmt.Sprintf(" AND sr.services_id = %d", servicesId))
	}

	mod.Raw("SELECT COUNT(0) " + sql.String()).Count(&count)
	mod.Raw("SELECT sr.*, m.nick_name, m.qq, m.phone, m.wechat, s.title, sm.nick_name AS author"+""+sql.String()+" ORDER BY sr.id DESC LIMIT ?,? ", int(offset), limit).Find(&res)
	return
}
