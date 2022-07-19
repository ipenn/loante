package model

import (
	"fmt"
	"github.com/uptrace/bun"
	"loante/global"
)

type Orders struct {
	bun.BaseModel    `bun:"table:orders,alias:o"`
	Id               int             `json:"id" bun:",pk"`
	MchId            int             `json:"-"`
	ProductId        int             `json:"-"`
	Bid              int             `json:"bid"`
	Uid              int             `json:"-"`
	CreateTime       string          `json:"create_time"`
	Type             int             `json:"type"`
	Payment          int             `json:"payment"`
	PostponePeriod   int             `json:"postpone_period"`
	ActualAmount     int             `json:"actual_amount"`
	ApplyAmount      int             `json:"apply_amount"`
	RepaidStatus     int             `json:"repaid_status"`
	UrgeId           int             `json:"urge_id"`
	RepaidUrl        string          `json:"repaid_url"`
	RepaidType       int             `json:"repaid_type"`
	InvalidTime      string          `json:"invalid_time"`
	LoanTime         string          `json:"loan_time"`
	EndTime          string          `json:"end_time"`
	PaymentRequestNo string          `json:"payment_request_no"`
	PaymentRespondNo string          `json:"payment_respond_no"`
	LatePaymentFee   int             `json:"late_payment_fee"`
	User             *UserLittle     `json:"user" bun:"rel:belongs-to,join:uid=id"`
	Borrow           *BorrowLittle   `json:"borrow" bun:"rel:belongs-to,join:bid=id"`
	Merchant         *MerchantLittle `json:"merchant" bun:"rel:belongs-to,join:mch_id=id"`
	Product          *ProductLittle  `json:"product" bun:"rel:belongs-to,join:product_id=id"`
	PaymentCom       *PaymentLittle  `json:"payment_com" bun:"rel:belongs-to,join:payment=id"`
}

func (a *Orders) Insert() {
	res, _ := global.C.DB.NewInsert().Model(a).Ignore().On("DUPLICATE KEY UPDATE").Exec(global.C.Ctx)
	id, _ := res.LastInsertId()
	a.Id = int(id)
}

func (a *Orders) One(where string) {
	err := global.C.DB.NewSelect().Model(a).Where(where).Scan(global.C.Ctx)
	fmt.Println(err)
}

func (a *Orders) Gets(where string) ([]Borrow, int) {
	var datas []Borrow
	count, _ := global.C.DB.NewSelect().Model(&datas).Where(where).ScanAndCount(global.C.Ctx)
	return datas, count
}

func (a *Orders) Count(where string) int {
	count, _ := global.C.DB.NewSelect().Model(a).Where(where).Count(global.C.Ctx)
	return count
}

func (a *Orders) Update(where string) {
	_, err := global.C.DB.NewUpdate().Model(a).Where(where).Exec(global.C.Ctx)
	if err != nil {
		global.Log.Error("%v err=%v", a, err.Error())
	}
}

func (a *Orders) Page(where string, page, limit int) ([]Orders, int) {
	var datas []Orders
	count, _ := global.C.DB.NewSelect().Model(&datas).Relation("Borrow").Relation("Merchant").Relation("Product").Relation("PaymentCom").Relation("User").Where(where).Order(fmt.Sprintf("o.id desc")).Offset((page - 1) * limit).Limit(limit).ScanAndCount(global.C.Ctx)
	return datas, count
}

func (a *Orders) Del(where string) {
	global.C.DB.NewDelete().Model(a).Where(where).Exec(global.C.Ctx)
}
