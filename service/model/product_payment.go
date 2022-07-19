package model

import (
	"fmt"
	"github.com/uptrace/bun"
	"loante/global"
)

type ProductPayment struct {
	bun.BaseModel `bun:"table:product_payment,alias:pp"`
	Id            int    `json:"id" bun:",pk"`
	ProductId     int    `json:"product_id"`
	MchId         int    `json:"mch_id"`
	PaymentId     int    `json:"payment_id"`
	Configuration string `json:"configuration"`
	IsOpenOut     int    `json:"is_open_out"`
	IsOpenIn      int    `json:"is_open_in"`
	CreateTime    string `json:"create_time"`
}

func (a *ProductPayment) Insert() {
	_, err := global.C.DB.NewInsert().Model(a).Returning("*").Exec(global.C.Ctx)
	if err != nil {
		global.Log.Error("%v err=%v", a, err.Error())
	}
}

func (a *ProductPayment) Update(where string) {
	_, err := global.C.DB.NewUpdate().Model(a).Where(where).Exec(global.C.Ctx)
	if err != nil {
		global.Log.Error("%v err=%v", a, err.Error())
	}
}

func (a *ProductPayment) Del(where string) {
	_, err := global.C.DB.NewDelete().Model(a).Where(where).Exec(global.C.Ctx)
	if err != nil {
		global.Log.Error("%v err=%v", a, err.Error())
	}
}

func (a *ProductPayment) One(where string) {
	err := global.C.DB.NewSelect().Model(a).Where(where).Scan(global.C.Ctx)
	if err != nil {
		global.Log.Error("%v err=%v", a, err.Error())
	}
}

func (a *ProductPayment) Page(where string, page, limit int) ([]ProductPayment, int) {
	var datas []ProductPayment
	count, _ := global.C.DB.NewSelect().Model(&datas).Where(where).Order(fmt.Sprintf("id desc")).Offset((page - 1) * limit).Limit(limit).ScanAndCount(global.C.Ctx)
	return datas, count
}

func (a *ProductPayment) SetColumn(col string, value interface{}, where string) error {
	_, err := global.C.DB.NewUpdate().Model(a).SetColumn(col, "?", value).Where(where).Exec(global.C.Ctx)
	if err != nil {
		global.Log.Error("%v err=%v", a, err.Error())
		return err
	}
	return nil
}
