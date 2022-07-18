package model

import (
	"fmt"
	"github.com/uptrace/bun"
	"loante/global"
)

type RemindRules struct {
	bun.BaseModel `bun:"table:remind_rules,alias:rr"`
	Id            int            `json:"id"`
	CompanyId     int            `json:"company_id"`
	GroupId       int            `json:"group_id"`
	MaxDay        int            `json:"max_day"`
	MinDay        int            `json:"min_day"`
	Remark        string         `json:"remark"`
	IsAutoApply   int            `json:"is_auto_apply"`
	RemindCompany *RemindCompany `json:"remind_company" bun:"rel:belongs-to,join:company_id=id"`
	RemindGroup   *RemindGroup   `json:"remind_group" bun:"rel:belongs-to,join:group_id=id"`
}

func (a *RemindRules) Insert() {
	_, err := global.C.DB.NewInsert().Model(a).Returning("*").Exec(global.C.Ctx)
	if err != nil {
		global.Log.Error("%v err=%v", a, err.Error())
	}
}

func (a *RemindRules) Update(where string) {
	_, err := global.C.DB.NewUpdate().Model(a).Where(where).Exec(global.C.Ctx)
	if err != nil {
		global.Log.Error("%v err=%v", a, err.Error())
	}
}

func (a *RemindRules) One(where string) {
	err := global.C.DB.NewSelect().Model(a).Where(where).Scan(global.C.Ctx)
	fmt.Println(err)
}

func (a *RemindRules) Page(where string, page, limit int) ([]RemindRules, int) {
	var datas []RemindRules
	count, _ := global.C.DB.NewSelect().Model(&datas).Relation("RemindGroup").Relation("RemindCompany").Where(where).Order(fmt.Sprintf("rr.id desc")).Offset((page - 1) * limit).Limit(limit).ScanAndCount(global.C.Ctx)
	return datas, count
}
