package model

import (
	"fmt"
	"github.com/uptrace/bun"
	"loante/global"
)

type UrgeRules struct {
	bun.BaseModel `bun:"table:urge_rules,alias:ur"`
	Id            int            `json:"id"`
	CompanyId     int            `json:"company_id"`
	GroupId       int            `json:"group_id"`
	MaxDay        int            `json:"max_day"`
	MinDay        int            `json:"min_day"`
	Remark        string         `json:"remark"`
	IsAutoApply   int            `json:"is_auto_apply"`
	UrgeCompany *UrgeCompany `json:"urge_company" bun:"rel:belongs-to,join:company_id=id"`
	UrgeGroup   *UrgeGroup   `json:"urge_group" bun:"rel:belongs-to,join:group_id=id"`
}

func (a *UrgeRules) Insert() {
	_, err := global.C.DB.NewInsert().Model(a).Returning("*").Exec(global.C.Ctx)
	if err != nil {
		global.Log.Error("%v err=%v", a, err.Error())
	}
}

func (a *UrgeRules) Update(where string) {
	_, err := global.C.DB.NewUpdate().Model(a).Where(where).Exec(global.C.Ctx)
	if err != nil {
		global.Log.Error("%v err=%v", a, err.Error())
	}
}

func (a *UrgeRules) One(where string) {
	err := global.C.DB.NewSelect().Model(a).Where(where).Scan(global.C.Ctx)
	fmt.Println(err)
}

func (a *UrgeRules) Page(where string, page, limit int) ([]UrgeRules, int) {
	var datas []UrgeRules
	count, _ := global.C.DB.NewSelect().Model(&datas).Relation("UrgeGroup").Relation("UrgeCompany").Where(where).Order(fmt.Sprintf("ur.id desc")).Offset((page - 1) * limit).Limit(limit).ScanAndCount(global.C.Ctx)
	return datas, count
}
