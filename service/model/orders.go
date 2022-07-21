package model

import (
	"fmt"
	"github.com/uptrace/bun"
	"loante/global"
	"loante/tools"
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
	Remark   		 string           `json:"remark"`
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

//PayAfter 还款成功以后处理订单逻辑
func (a *Orders) PayAfter(amount float64) {
	//更新 orders 支付状态
	a.RepaidStatus = 1
	a.ActualAmount = int(amount)
	borrowData := new(Borrow)
	borrowData.One(fmt.Sprintf("id = %d", a.Bid))
	if a.Type == 0 { //全额还款
		if borrowData.Status == 6 {
			borrowData.Status = 8
		}
		if borrowData.Status == 7 {
			borrowData.Status = 9
		}
		borrowData.BeRepaidAmount = 0
	} else if a.Type == 1 {
		borrowData.BeRepaidAmount -= a.ActualAmount
	} else if a.Type == 2 { //展期还款
		//获取产品的展期天数
		productData := new(ProductDelayConfig)
		productData.One(fmt.Sprintf("id = %d", borrowData.ProductId))
		borrowData.PostponedPeriod += 1
		borrowData.Postponed = 1
		borrowData.Status = 5
		endTimeUnix := tools.StrToUnixTime(borrowData.EndTime) + int64(productData.DelayDay *24*3600)
		borrowData.EndTime = tools.UnixTimeToStr(endTimeUnix)
	}
	a.Update(fmt.Sprintf("id = %d", a.Id))
	borrowData.Update(fmt.Sprintf("id = %d", borrowData.Id))
	//更新产品额度
	if borrowData.Status == 8 || borrowData.Status == 9{
		new(UserQuota).Increase(borrowData.ProductId, borrowData.Uid, borrowData.Status)
	}
}


type OrdersForStatistics struct {
	Count      int    `json:"count"`
	Payment    int    `json:"payment"`
	CreateTime string `json:"create_time"`
	Name       string `json:"name"`
	Type       int    `json:"type"`
}

func (a *Orders) ForStatistics(where string) []OrdersForStatistics {
	//var datas []Borrow
	//count, _ := global.C.DB.NewSelect().Model(&Borrow{}).GroupExpr("HOUR(loan_time) desc").
	//	GroupExpr("payment DESC").Where(where).Offset((page - 1) * limit).Limit(limit).ScanAndCount(global.C.Ctx)
	//return datas, count
	var ordersForStatistics []OrdersForStatistics
	rows, _ := global.C.DB.QueryContext(global.C.Ctx, `SELECT
	COUNT( o.id ) AS count,
	o.payment,
	o.create_time,
	p.name,
	0 AS type 
FROM
	orders o
	LEFT JOIN payment p ON o.payment = p.id 
WHERE
	o.repaid_status IN ( 0,1 ) and DATE_SUB(CURDATE(), INTERVAL 7 DAY) <= date(o.create_time) `+where+`
GROUP BY
	HOUR ( o.create_time ) DESC,
	o.payment DESC UNION ALL
SELECT
	COUNT( o.id ) AS count,
	o.payment,
	o.create_time,
	p.name,
	0 AS type 
FROM
	orders o
	LEFT JOIN payment p ON o.payment = p.id 
WHERE
	o.repaid_status =2 and DATE_SUB(CURDATE(), INTERVAL 7 DAY) <= date(o.create_time) `+where+`
GROUP BY
	HOUR ( o.create_time ) DESC,
	o.payment DESC`)

	global.C.DB.ScanRows(global.C.Ctx, rows, &ordersForStatistics)
	return ordersForStatistics
}
