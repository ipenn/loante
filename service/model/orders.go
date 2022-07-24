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
	Remark        string      `json:"remark"`
	PayNotifyTime string      `json:"pay_notify_time"`
	User          *UserLittle `json:"user" bun:"rel:belongs-to,join:uid=id"`
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
	_amount := int(amount)
	a.RepaidStatus = 1
	a.ActualAmount = _amount
	a.PayNotifyTime = tools.GetFormatTime()
	borrowData := new(Borrow)
	borrowData.One(fmt.Sprintf("id = %d", a.Bid))
	if a.Type == 0 { //全额还款
		if borrowData.Status == 6 {
			borrowData.Status = 8
		}
		if borrowData.Status == 7 {
			borrowData.Status = 9
		}
		a.LatePaymentFee = borrowData.LatePaymentFee
		borrowData.BeRepaidAmount = 0
		borrowData.LatePaymentFee = 0 //更新剩余滞纳金
		borrowData.CompleteTime = tools.GetFormatTime()
	} else if a.Type == 1 {
		if _amount >= borrowData.LatePaymentFee{
			a.LatePaymentFee = borrowData.LatePaymentFee
			borrowData.BeRepaidAmount -= _amount - borrowData.LatePaymentFee
			borrowData.LatePaymentFee = 0
		}else{
			a.LatePaymentFee = borrowData.LatePaymentFee - _amount
		}
	} else if a.Type == 2 { //展期还款
		//获取产品的展期天数
		productData := new(ProductDelayConfig)
		productData.One(fmt.Sprintf("id = %d", borrowData.ProductId))
		borrowData.PostponedPeriod += 1
		borrowData.Postponed = 1
		borrowData.Status = 5
		borrowData.LatePaymentFee = 0 //更新剩余滞纳金
		endTimeUnix := tools.StrToUnixTime(borrowData.EndTime) + int64(productData.DelayDay *24*3600)
		borrowData.EndTime = tools.UnixTimeToStr(endTimeUnix)
		borrowData.ExpireDay -= productData.DelayDay
		a.LatePaymentFee = borrowData.LatePaymentFee
	}
	a.Update(fmt.Sprintf("id = %d", a.Id))
	borrowData.Update(fmt.Sprintf("id = %d", borrowData.Id))
	//更新产品额度
	if borrowData.Status == 8 || borrowData.Status == 9 {
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

type ForStatistics struct {
	Count        int     `json:"count"`
	SuccessCount int     `json:"success_count"`
	SuccessRate  float64 `json:"success_rate"`
	Payment      int     `json:"payment"`
	CreateTime   string  `json:"create_time"`
	Name         string  `json:"name"`
	Type         int     `json:"type"`
}

func (a *Orders) ForStatistics(where string) []ForStatistics {
	//var datas []Borrow
	//count, _ := global.C.DB.NewSelect().Model(&Borrow{}).GroupExpr("HOUR(loan_time) desc").
	//	GroupExpr("payment DESC").Where(where).Offset((page - 1) * limit).Limit(limit).ScanAndCount(global.C.Ctx)
	//return datas, count
	var ordersForStatistics []OrdersForStatistics
	rows, _ := global.C.DB.QueryContext(global.C.Ctx, `SELECT
	COUNT( o.id ) AS count,
	o.payment,
	date_format(o.create_time,'%y-%m-%d %H') as create_time,
	p.name,
	0 AS type 
FROM
	orders o
	LEFT JOIN payment p ON o.payment = p.id 
WHERE
	 DATE_SUB(CURDATE(), INTERVAL 3 DAY) <= date(o.create_time) `+where+`
GROUP BY
	DAY(o.create_time),
	HOUR ( o.create_time ) DESC,
	o.payment DESC UNION ALL
SELECT
	COUNT( o.id ) AS count,
	o.payment,
	date_format(o.create_time,'%y-%m-%d %H') as create_time,
	p.name,
	1 AS type 
FROM
	orders o
	LEFT JOIN payment p ON o.payment = p.id 
WHERE
	o.repaid_status =1 and DATE_SUB(CURDATE(), INTERVAL 3 DAY) <= date(o.create_time) `+where+`
GROUP BY
	DAY(o.create_time),
	HOUR ( o.create_time ) DESC,
	o.payment DESC`)
	var result []ForStatistics
	global.C.DB.ScanRows(global.C.Ctx, rows, &ordersForStatistics)
	for _, statistic := range ordersForStatistics {
		flag := true
		if statistic.Type == 1 {
			break
		}
		for _, forStatistic := range ordersForStatistics {
			if statistic.Payment == forStatistic.Payment && statistic.Type == 0 && forStatistic.Type == 1 && statistic.CreateTime == forStatistic.CreateTime {
				var r ForStatistics
				r.Payment = statistic.Payment
				r.Count = statistic.Count
				r.SuccessCount = forStatistic.Count
				r.SuccessRate = tools.ToFloat64(forStatistic.Count) / tools.ToFloat64(statistic.Count)
				r.CreateTime = statistic.CreateTime
				r.Name = statistic.Name
				result = append(result, r)
				flag = false
			}
		}
		if flag {
			var r ForStatistics
			r.Payment = statistic.Payment
			r.Count = statistic.Count
			r.SuccessCount = 0
			r.SuccessRate = 0
			r.CreateTime = statistic.CreateTime
			r.Name = statistic.Name
			result = append(result, r)
		}
	}
	return result
}
