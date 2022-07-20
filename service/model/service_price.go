package model

import (
	"fmt"
	"github.com/uptrace/bun"
	"loante/global"
)

type ServicePrice struct {
	bun.BaseModel `bun:"table:service_price,alias:sp"`
	Id            int     `json:"id" bun:",pk"`
	ServiceType   int     `json:"service_type"`
	DeductType    int     `json:"deduct_type"`
	Price         float64 `json:"price"`
}

func (a *ServicePrice) Insert() {
	_, err := global.C.DB.NewInsert().Model(a).Returning("*").Exec(global.C.Ctx)
	if err != nil {
		global.Log.Error("%v err=%v", a, err.Error())
	}
}

func (a *ServicePrice) Update(where string) {
	_, err := global.C.DB.NewUpdate().Model(a).Where(where).Exec(global.C.Ctx)
	if err != nil {
		global.Log.Error("%v err=%v", a, err.Error())
	}
}

func (a *ServicePrice) Del(where string) {
	_, err := global.C.DB.NewDelete().Model(a).Where(where).Exec(global.C.Ctx)
	if err != nil {
		global.Log.Error("%v err=%v", a, err.Error())
	}
}

func (a *ServicePrice) One(where string) {
	err := global.C.DB.NewSelect().Model(a).Where(where).Scan(global.C.Ctx)
	if err != nil {
		global.Log.Error("%v err=%v", a, err.Error())
	}
}

func (a *ServicePrice) Page(where string, page, limit int) ([]ServicePrice, int) {
	var datas []ServicePrice
	count, _ := global.C.DB.NewSelect().Model(&datas).Where(where).Order(fmt.Sprintf("id desc")).Offset((page - 1) * limit).Limit(limit).ScanAndCount(global.C.Ctx)
	return datas, count
}
