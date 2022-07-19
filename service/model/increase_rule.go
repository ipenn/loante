package model

import (
	"github.com/uptrace/bun"
	"loante/global"
	"loante/tools"
)

type IncreaseRule struct {
	bun.BaseModel    `bun:"table:increase_rule,alias:ir"`
	Id               int    `json:"id"`
	LoanProductCount int    `json:"loan_product_count"`
	MinCount         int    `json:"min_count"`
	CreateTime       string `json:"create_time"`
	UpdateTime       string `json:"update_time"`
}

func (a *IncreaseRule) One(where string) {
	err := global.C.DB.NewSelect().Model(a).Where(where).Scan(global.C.Ctx)
	if err != nil {
		global.Log.Error("%v err=%v", a, err.Error())
	}
}

func (a *IncreaseRule) Insert() {
	a.CreateTime = tools.GetFormatTime()
	_, err := global.C.DB.NewInsert().Model(a).Returning("*").Exec(global.C.Ctx)
	if err != nil {
		global.Log.Error("%v err=%v", a, err.Error())
	}
}

func (a *IncreaseRule) Update(where string) {
	_, err := global.C.DB.NewUpdate().Model(a).Where(where).Exec(global.C.Ctx)
	if err != nil {
		global.Log.Error("%v err=%v", a, err.Error())
	}
}

func (a *IncreaseRule) Del(where string) {
	_, err := global.C.DB.NewDelete().Model(a).Where(where).Exec(global.C.Ctx)
	if err != nil {
		global.Log.Error("%v err=%v", a, err.Error())
	}
}
func (a *IncreaseRule) Page (page, limit int) ([]IncreaseRule, int) {
	var d []IncreaseRule
	count, _ := global.C.DB.NewSelect().Model(&d).Order("ir.id desc").Offset((page - 1) * limit).Limit(limit).
		ScanAndCount(global.C.Ctx)
	return d, count
}
