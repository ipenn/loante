package model

import (
	"fmt"
	"github.com/uptrace/bun"
	"loante/global"
)

type PaymentMini struct {
	bun.BaseModel `bun:"table:payment,alias:p"`
	Id          int	`json:"id" bun:",pk"`
	Name        string	`json:"name"`
}
type PaymentLittle struct {
	bun.BaseModel `bun:"table:payment,alias:p"`
	PaymentMini
	IsOpenOut   int	`json:"is_open_out"`
	IsOpenIn    int	`json:"is_open_in"`
	LendingStartTime string	`json:"lending_start_time"`
	LendingEndTime string	`json:"lending_end_time"`

}
type Payment struct {
	bun.BaseModel `bun:"table:payment,alias:p"`
	PaymentLittle
	IsUtrQuery int    `json:"is_utr_query"`
	IsUtrFill  int    `json:"is_utr_fill"`
	Fields     string `json:"fields"`
}

func (a *Payment) Insert() {
	_, err := global.C.DB.NewInsert().Model(a).Returning("*").Exec(global.C.Ctx)
	if err != nil {
		global.Log.Error("%v err=%v", a, err.Error())
	}
}

func (a *Payment) Update(where string) {
	_, err := global.C.DB.NewUpdate().Model(a).Where(where).Exec(global.C.Ctx)
	if err != nil {
		global.Log.Error("%v err=%v", a, err.Error())
	}
}

func (a *Payment) Del(where string) {
	_, err := global.C.DB.NewDelete().Model(a).Where(where).Exec(global.C.Ctx)
	if err != nil {
		global.Log.Error("%v err=%v", a, err.Error())
	}
}

func (a *Payment) One(where string) {
	err := global.C.DB.NewSelect().Model(a).Where(where).Scan(global.C.Ctx)
	if err != nil {
		global.Log.Error("%v err=%v", a, err.Error())
	}
}

func (a *PaymentLittle) Gets() ([]PaymentLittle, int) {
	var datas []PaymentLittle
	count, _ := global.C.DB.NewSelect().Model(&datas).ScanAndCount(global.C.Ctx)
	return datas, count
}

func (a *Payment) Page(where string, page, limit int) ([]Payment, int) {
	var datas []Payment
	count, _ := global.C.DB.NewSelect().Model(&datas).Where(where).Order(fmt.Sprintf("id desc")).Offset((page - 1) * limit).Limit(limit).ScanAndCount(global.C.Ctx)
	return datas, count
}

func (a *Payment) SetColumn(col string, value interface{}, where string) error {
	_, err := global.C.DB.NewUpdate().Model(a).SetColumn(col, "?", value).Where(where).Exec(global.C.Ctx)
	if err != nil {
		global.Log.Error("%v err=%v", a, err.Error())
		return err
	}
	return nil
}
