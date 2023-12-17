package model

import (
	"bytes"
	"fmt"
	"github.com/farmer-simon/go-utils"
	"goskeleton/app/data_type"
	"math"
)

func CreateNeedsRecordFactory(sqlType string) *NeedsRecordModel {
	return &NeedsRecordModel{BaseModel: BaseModel{DB: UseDbConn(sqlType)}}
}

type NeedsRecordModel struct {
	BaseModel
	data_type.NeedsRecordMod
}

// TableName 表名
func (mod *NeedsRecordModel) TableName() string {
	return "tech_needs_record"
}

// InsertRecord 提交需求
func (mod *NeedsRecordModel) InsertRecord(membersId, servicesId int, content, tenderPrice string) error {
	sql := "INSERT INTO `tech_needs_record`(members_id, needs_id, content, tender_price,price, start_time, end_time, attach, state) VALUES(?,?,?,?,?,?,?,?,?)"
	res := mod.Exec(sql, membersId, servicesId, content, tenderPrice, "", 0, 0, "", 1)
	return res.Error
}

// CloseRecord 关闭需求
func (mod *NeedsRecordModel) CloseRecord(id int) error {
	sql := "UPDATE `tech_needs_record` SET state=? WHERE id=?"
	res := mod.Exec(sql, 2, id)
	return res.Error
}

// AcceptRecord 接受需求
func (mod *NeedsRecordModel) AcceptRecord(id int, price, startTime, endTime, attach string) error {
	sTime := utils.DatetimeToUnixTimestamp(startTime + " 00:00:00")
	eTime := utils.DatetimeToUnixTimestamp(endTime + " 23:59:59")
	sql := "UPDATE `tech_needs_record` SET price=?, start_time=?, end_time=?, attach=?, state=? WHERE id=?"
	res := mod.Exec(sql, price, sTime, eTime, attach, 3, id)
	return res.Error
}

// GetRecordTotalByNeedsIds 根据服务IDS获取各服务的需求数量
func (mod *NeedsRecordModel) GetRecordTotalByNeedsIds(ids []int) map[int]int {
	var totals = make(map[int]int, 0)
	var results []data_type.RecordTotal

	sql := "SELECT needs_id AS id, COUNT(0) AS total FROM `tech_needs_record` WHERE needs_id IN (?) GROUP BY needs_id"
	mod.Raw(sql, ids).Scan(&results)
	for _, item := range results {
		totals[item.Id] = item.Total
	}
	return totals
}

// Index 根据服务IDS获取各服务的需求数量
func (mod *NeedsRecordModel) Index(id, page, limit int) (count int64, res []NeedsRecordModel) {
	offset := math.Max((float64(page)-1)*float64(limit), 0)
	mod.Raw("SELECT COUNT(0) FROM `tech_needs_record` WHERE needs_id=?", id).Count(&count)
	mod.Raw("SELECT r.*, m.nick_name, m.qq, m.phone, m.wechat  FROM tech_needs_record r LEFT JOIN tech_members m ON r.members_id=m.id WHERE r.needs_id=? LIMIT ?,? ", id, int(offset), limit).Find(&res)
	return
}

func (mod *NeedsRecordModel) MyList(membersId, page, limit int) (count int64, res []NeedsRecordModel) {
	offset := math.Max((float64(page)-1)*float64(limit), 0)
	mod.Raw("SELECT COUNT(0) FROM `tech_needs_record` WHERE members_id=?", membersId).Count(&count)
	mod.Raw("SELECT r.*, n.title, m.nick_name AS author  FROM tech_needs_record r LEFT JOIN tech_needs n ON r.needs_id=n.id LEFT JOIN tech_members m ON n.members_id=m.id WHERE r.members_id=? LIMIT ?,? ", membersId, int(offset), limit).Find(&res)
	return
}

func (mod *NeedsRecordModel) GetById(id int) (res *NeedsRecordModel) {
	mod.Raw("SELECT r.*, m.nick_name, m.qq, m.phone, m.wechat  FROM tech_needs_record r LEFT JOIN tech_members m ON r.members_id=m.id WHERE r.id=?", id).Find(&res)
	return
}

// AdminIndex 后台的查看
func (mod *NeedsRecordModel) AdminIndex(needsId, state, page, limit int) (count int64, res []NeedsRecordModel) {
	offset := math.Max((float64(page)-1)*float64(limit), 0)

	var sql bytes.Buffer
	sql.WriteString(" FROM tech_needs_record sr LEFT JOIN tech_members m ON sr.members_id=m.id LEFT JOIN tech_needs s ON sr.needs_id=s.id LEFT JOIN tech_members sm ON s.members_id=sm.id WHERE ")
	if state == 0 {
		sql.WriteString(" sr.state IN (1,2,3)")
	} else {
		sql.WriteString(fmt.Sprintf(" sr.state = %d", state))
	}
	if needsId > 0 {
		sql.WriteString(fmt.Sprintf(" AND sr.needs_id = %d", needsId))
	}

	mod.Raw("SELECT COUNT(0) " + sql.String()).Count(&count)
	mod.Raw("SELECT sr.*, m.nick_name, m.qq, m.phone, m.wechat, s.title, sm.nick_name AS author"+""+sql.String()+" ORDER BY sr.id DESC LIMIT ?,? ", int(offset), limit).Find(&res)
	return
}
