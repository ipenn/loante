package model

import (
	"fmt"
	"github.com/uptrace/bun"
	"loante/global"
)

type ProductPaymentDefault struct {
	bun.BaseModel `bun:"table:product_payment_default,alias:ppd"`
	Id            	int	`json:"id" bun:",pk"`
	ProductId     	int	`json:"product_id"`
	MchId          	int `json:"mch_id"`
	OutPaymentId 	int `json:"out_payment_id"`
	InPaymentId 	int `json:"in_payment_id"`

	Product		*ProductLittle `json:"product" bun:"rel:belongs-to,join:product_id=id"`
	Merchant	*MerchantLittle `json:"merchant" bun:"rel:belongs-to,join:mch_id=id"`
	OutPayment	*PaymentLittle `json:"out_payment" bun:"rel:belongs-to,join:out_payment_id=id"`
	InPayment	*PaymentLittle `json:"in_payment" bun:"rel:belongs-to,join:in_payment_id=id"`
}

func (a *ProductPaymentDefault) Insert() {
	_, err := global.C.DB.NewInsert().Model(a).Returning("*").Exec(global.C.Ctx)
	if err != nil {
		global.Log.Error("%v err=%v", a, err.Error())
	}
}

func (a *ProductPaymentDefault) Update(where string) {
	_, err := global.C.DB.NewUpdate().Model(a).Where(where).Exec(global.C.Ctx)
	if err != nil {
		global.Log.Error("%v err=%v", a, err.Error())
	}
}

func (a *ProductPaymentDefault) Del(where string) {
	_, err := global.C.DB.NewDelete().Model(a).Where(where).Exec(global.C.Ctx)
	if err != nil {
		global.Log.Error("%v err=%v", a, err.Error())
	}
}

func (a *ProductPaymentDefault) One(where string) {
	err := global.C.DB.NewSelect().Model(a).Where(where).Scan(global.C.Ctx)
	if err != nil {
		global.Log.Error("%v err=%v", a, err.Error())
	}
}

func (a *ProductPaymentDefault) Page(where string, page, limit int) ([]ProductPaymentDefault, int) {
	var datas []ProductPaymentDefault
	count, _ := global.C.DB.NewSelect().Model(&datas).Relation("Merchant").Relation("OutPayment").Relation("InPayment").Relation("Product").Where(where).Order(fmt.Sprintf("ppd.id desc")).Offset((page-1)*limit).Limit(limit).ScanAndCount(global.C.Ctx)
	return datas, count
}
