package model

import (
	"fmt"
	"github.com/uptrace/bun"
	"loante/global"
	"loante/tools"
)

type Merchant struct {
	bun.BaseModel `bun:"table:merchant,alias:m"`
	Id	int	`json:"id" bun:",pk"`
	Name	string	`json:"name"`
	Type	int	`json:"type"`
	CnyBalance	float64	`json:"cny_balance"`
	UsdBalance	float64	`json:"usd_balance"`
	CnyCredit	float64	`json:"cny_credit"`
	UsdCredit	float64	`json:"usd_credit"`
	ContactName	string	`json:"contact_name"`
	ContactMobile	string	`json:"contact_mobile"`
	ContactEmail	string	`json:"contact_email"`
	CreateTime	string	`json:"create_time"`
	UpdateTime	string	`json:"update_time"`
	Status	int	`json:"status"`
}

func (a *Merchant)Insert()  {
	a.CreateTime = tools.GetFormatTime()
	_, err := global.C.DB.NewInsert().Model(a).Returning("*").Exec(global.C.Ctx)
	if err != nil{
		global.Log.Error("%v err=%v", a, err.Error())
	}
}

func (a *Merchant)Update(where string)  {
	_, err := global.C.DB.NewUpdate().Model(a).Where(where).Exec(global.C.Ctx)
	if err != nil{
		global.Log.Error("%v err=%v", a, err.Error())
	}
}

func (a *Merchant)One(where string)  {
	err := global.C.DB.NewSelect().Model(a).Where(where).Scan(global.C.Ctx)
	if err != nil{
		global.Log.Error("%v err=%v", a, err.Error())
	}
}

func (a *Merchant)Page(where string, page, limit int) ([]Merchant, int) {
	var datas []Merchant
	count, _ := global.C.DB.NewSelect().Model(&datas).Where(where).Order(fmt.Sprintf("id desc")).Offset((page-1)*limit).Limit(limit).ScanAndCount(global.C.Ctx)
	return datas, count
}

func (a *Merchant)SetColumn(col string, value interface{}, where string) error  {
	_, err := global.C.DB.NewUpdate().Model(a).SetColumn(col, "?", value).Where(where).Exec(global.C.Ctx)
	if err != nil{
		global.Log.Error("%v err=%v", a, err.Error())
		return err
	}
	return nil
}

func (a *Merchant)Set(col string, where string) error  {
	_, err := global.C.DB.NewUpdate().Model(a).Set(col).Where(where).Exec(global.C.Ctx)
	if err != nil{
		global.Log.Error("%v err=%v", a, err.Error())
		return err
	}
	return nil
}
