package model

import (
	"fmt"
	"github.com/uptrace/bun"
	"loante/global"
)

type ServiceFeeRule struct {
	bun.BaseModel `bun:"table:service_fee_rule,alias:sfr"`
	Id	int	`json:"id" bun:",pk"`
	Name	string	`json:"name"`
	Price	float64	`json:"price"`
	StartCount	int	`json:"start_count"`
	EndCount	int	`json:"end_count"`
}

func (a *ServiceFeeRule)Insert()  {
	_, err := global.C.DB.NewInsert().Model(a).Returning("*").Exec(global.C.Ctx)
	if err != nil{
		global.Log.Error("%v err=%v", a, err.Error())
	}
}

func (a *ServiceFeeRule)Update(where string)  {
	_, err := global.C.DB.NewUpdate().Model(a).Where(where).Exec(global.C.Ctx)
	if err != nil{
		global.Log.Error("%v err=%v", a, err.Error())
	}
}

func (a *ServiceFeeRule)Del(where string)  {
	_, err := global.C.DB.NewDelete().Model(a).Where(where).Exec(global.C.Ctx)
	if err != nil{
		global.Log.Error("%v err=%v", a, err.Error())
	}
}

func (a *ServiceFeeRule)One(where string)  {
	err := global.C.DB.NewSelect().Model(a).Where(where).Scan(global.C.Ctx)
	if err != nil{
		global.Log.Error("%v err=%v", a, err.Error())
	}
}

func (a *ServiceFeeRule)Page(where string, page, limit int) ([]ServiceFeeRule, int) {
	var datas []ServiceFeeRule
	count, _ := global.C.DB.NewSelect().Model(&datas).Where(where).Order(fmt.Sprintf("id desc")).Offset((page-1)*limit).Limit(limit).ScanAndCount(global.C.Ctx)
	return datas, count
}


