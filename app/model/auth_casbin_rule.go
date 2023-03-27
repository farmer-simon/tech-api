package model

func CreateCasbinRuleFactory(sqlType string) *CasbinRule {
	return &CasbinRule{BaseModel: BaseModel{DB: UseDbConn(sqlType)}}
}

type CasbinRule struct {
	BaseModel `json:"-"`
	Ptype     string ` json:"ptype"`
	V0        string `json:"v0"`
	V1        string `json:"v1"`
	V2        string `json:"v2"`
	V3        string `json:"v3"`
	V4        string `json:"v4"`
	V5        string `json:"v5"`
}

// 表名
func (c *CasbinRule) TableName() string {
	return "tech_auth_casbin_rule"
}

// 查询

// 新增
func (c *CasbinRule) InsertData(pType, v0, v1, v2, v3, v4, v5 string) bool {
	sql := `
	INSERT INTO  tech_auth_casbin_rule(ptype,v0,v1,v2,v3,v4,v5)  
	SELECT  ?,?,?,?,?,?,?  FROM   DUAL   WHERE   NOT  EXISTS(SELECT 1 FROM tech_auth_casbin_rule WHERE ptype=?  AND  v0=? AND v1=? AND v2=? AND v3=? AND v4=? AND v5=? )
	`
	if res := c.Exec(sql, pType, v0, v1, v2, v3, v4, v5, pType, v0, v1, v2, v3, v4, v5); res.Error == nil {
		return true
	}
	return false
}

// 修改
func (c *CasbinRule) UpdateData(id int, pType, v0, v1, v2, v3, v4, v5 string) bool {
	sql := "update tech_auth_casbin_rule  set ptype=?, v0=?,v1=?,v2=?,v3=?,v4=?,v5=? where  id=? "
	if res := c.Exec(sql, pType, v0, v1, v2, v3, v4, v5, id); res.Error == nil {
		return true
	}
	return false
}

//删除
func (c *CasbinRule) DeleteData(id int) bool {
	if res := c.Delete(c, id); res.Error == nil {
		return true
	} else {
		return false
	}
}
