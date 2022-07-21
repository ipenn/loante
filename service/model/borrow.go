package model

import (
	"errors"
	"fmt"
	"github.com/uptrace/bun"
	"loante/global"
	"loante/tools"
)

type BorrowLittle struct {
	bun.BaseModel  `bun:"table:borrow,alias:b"`
	Id             int    `json:"id" bun:",pk"`
	Uid            int    `json:"uid"`
	MchId          int    `json:"mch_id"`
	ProductId      int    `json:"product_id"`
	Status         int    `json:"status"`
	Payment        int    `json:"payment"`
	BeRepaidAmount int    `json:"be_repaid_amount"`
	EndTime        string `json:"end_time"`
	LoanAmount     int    `json:"loan_amount"`
}
type Borrow struct {
	bun.BaseModel `bun:"table:borrow,alias:b"`
	BorrowLittle
<<<<<<< HEAD
	Postponed          	int				`json:"postponed"`
	PostponedPeriod    	int				`json:"postponed_period"`
	PostponeValuation 	int				`json:"postpone_valuation"`
	LoanType         	int				`json:"loan_type"`
	Score            	int				`json:"score"`
	RiskModel          int             `json:"risk_model"`
	ScoreTime          string          `json:"score_time"`
	CreateTime         string          `json:"create_time"`
	LoanTime           string          `json:"loan_time"`
	PaymentRequestNo   string          `json:"payment_request_no"`
	PaymentRespond     string          `json:"payment_respond"`
	Remark             string          `json:"remark"`
	CompleteTime       string          `json:"complete_time"`
	ActualAmount       int             `json:"actual_amount"`
	LatePaymentFee     int             `json:"late_payment_fee"`      //滞纳金
	LatePaymentFeeRate float64             `json:"late_payment_fee_rate"` //滞纳金
	User               *UserLittle     `json:"user" bun:"rel:belongs-to,join:uid=id"`
	Merchant           *MerchantLittle `json:"merchant" bun:"rel:belongs-to,join:mch_id=id"`
	Product            *Product        `json:"product" bun:"rel:belongs-to,join:product_id=id"`
=======
	Postponed         int             `json:"postponed"`
	PostponedPeriod   int             `json:"postponed_period"`
	PostponeValuation int             `json:"postpone_valuation"`
	LoanType          int             `json:"loan_type"`
	Score             int             `json:"score"`
	RiskModel         int             `json:"risk_model"`
	ScoreTime         string          `json:"score_time"`
	CreateTime        string          `json:"create_time"`
	LoanTime          string          `json:"loan_time"`
	PaymentRequestNo  string          `json:"payment_request_no"`
	PaymentRespond    string          `json:"payment_respond"`
	Remark            string          `json:"remark"`
	CompleteTime      string          `json:"complete_time"`
	ActualAmount      int             `json:"actual_amount"`
	LatePaymentFee    int             `json:"late_payment_fee"` //滞纳金
	User              *UserLittle     `json:"user" bun:"rel:belongs-to,join:uid=id"`
	Merchant          *MerchantLittle `json:"merchant" bun:"rel:belongs-to,join:mch_id=id"`
	Product           *Product        `json:"product" bun:"rel:belongs-to,join:product_id=id"`
>>>>>>> 5fa5f02c1373b226cd4ab46bcdfa3326f6ae89d0
}

func (a *Borrow) Insert() {
	res, _ := global.C.DB.NewInsert().Model(a).Ignore().On("DUPLICATE KEY UPDATE").Exec(global.C.Ctx)
	id, _ := res.LastInsertId()
	a.Id = int(id)
}

func (a *Borrow) One(where string) {
	err := global.C.DB.NewSelect().Model(a).Where(where).Scan(global.C.Ctx)
	fmt.Println(err)
}

func (a *Borrow) Gets(where string) ([]Borrow, int) {
	var datas []Borrow
	count, _ := global.C.DB.NewSelect().Model(&datas).Where(where).ScanAndCount(global.C.Ctx)
	return datas, count
}

func (a *Borrow) Count(where string) int {
	count, _ := global.C.DB.NewSelect().Model(a).Where(where).Count(global.C.Ctx)
	return count
}

func (a *Borrow) Update(where string) {
	_, err := global.C.DB.NewUpdate().Model(a).Where(where).Exec(global.C.Ctx)
	if err != nil {
		global.Log.Error("%v err=%v", a, err.Error())
	}
}

func (a *Borrow) Page(where string, page, limit int) ([]Borrow, int) {
	var datas []Borrow
	count, _ := global.C.DB.NewSelect().Model(&datas).Where(where).Order(fmt.Sprintf("b.id desc")).Offset((page - 1) * limit).Limit(limit).ScanAndCount(global.C.Ctx)
	return datas, count
}

func (a *Borrow) Del(where string) {
	global.C.DB.NewDelete().Model(a).Where(where).Exec(global.C.Ctx)
}

<<<<<<< HEAD
type PPConfig struct {
	Config string
	Name string
	Id int
}
func (a *Borrow) GetPaymentConfig(t string) (PPConfig, error) {
	payDefault := new(ProductPaymentDefault)
	payDefault.One(fmt.Sprintf("product_id = %d", a.ProductId))
	ppc := PPConfig{}
	if payDefault.Id == 0{
		return ppc, errors.New("产品的支付没有默认配置！")
	}
	//获取支付通道的配置
	pp := new(ProductPayment)
	if t == "in"{
		pp.One(fmt.Sprintf("payment_id = %d and product_id = %d", payDefault.InPaymentId, payDefault.ProductId))
	}else{
		pp.One(fmt.Sprintf("payment_id = %d and product_id = %d", payDefault.OutPaymentId, payDefault.ProductId))
	}
	if pp.Id == 0{
		return ppc, errors.New("产品没有配置支付信息！")
	}
	if t == "in" && pp.IsOpenIn == 0{
		return ppc, errors.New("产品配置支付未开放代收！")
	}
	if t == "out" && pp.IsOpenOut == 0{
		return ppc, errors.New("产品配置支付未开放代付！")
	}
	//全局支付通道是否开放
	p := new(Payment)
	p.One(fmt.Sprintf("id = %d", pp.PaymentId))
	if t == "in" && p.IsOpenIn == 0{
		return ppc, errors.New("支付通道代收全局未开放！")
	}
	if t == "out" && p.IsOpenOut == 0{
		return ppc, errors.New("支付通道代付全局未开放！")
	}
	ppc.Config = pp.Configuration
	ppc.Name = p.Name
	ppc.Id = p.Id
	return ppc, nil
}

//PayAfter 放款成功以后 需要发短信 和 扣费
func (a *Borrow)PayAfter()  {
	//更新borrow
	a.LoanTime = tools.GetFormatTime()
	a.Status = 5
	a.Update(fmt.Sprintf("id = %d", a.Id))
	//发短信？
	userData := new(User)
	userData.One(fmt.Sprintf("id = %d", a.Uid))
	tpl := new(SmsTemplate)
	if tpl.Send(2, userData.Phone, []interface{}{}){
		mchData := new(Merchant)
		mchData.Id = a.MchId
		mchData.AddService(1, 1) //扣费
	}
=======
type BorrowForStatistics struct {
	Count      int    `json:"count"`
	Payment    int    `json:"payment"`
	CreateTime string `json:"create_time"`
	Name       string `json:"name"`
	Type       int    `json:"type"`
}

func (a *Borrow) ForStatistics(where string) []BorrowForStatistics {
	//var datas []Borrow
	//count, _ := global.C.DB.NewSelect().Model(&Borrow{}).GroupExpr("HOUR(loan_time) desc").
	//	GroupExpr("payment DESC").Where(where).Offset((page - 1) * limit).Limit(limit).ScanAndCount(global.C.Ctx)
	//return datas, count
	var borrowForStatistics []BorrowForStatistics
	rows, _ := global.C.DB.QueryContext(global.C.Ctx, `SELECT
	COUNT( b.id ) AS count,
	b.payment,
	b.create_time,
	p.name,
	0 AS type 
FROM
	borrow b
	LEFT JOIN payment p ON b.payment = p.id 
WHERE
	b.STATUS IN ( 5, 0 ) and DATE_SUB(CURDATE(), INTERVAL 7 DAY) <= date(b.create_time) `+where+`
GROUP BY
	HOUR ( b.create_time ) DESC,
	b.payment DESC UNION ALL
SELECT
	COUNT( b.id ) AS count,
	b.payment,
	b.create_time,
	p.name,
	0 AS type 
FROM
	borrow b
	LEFT JOIN payment p ON b.payment = p.id 
WHERE
	b.STATUS =5 and DATE_SUB(CURDATE(), INTERVAL 7 DAY) <= date(b.create_time) `+where+`
GROUP BY
	HOUR ( b.create_time ) DESC,
	b.payment DESC`)

	global.C.DB.ScanRows(global.C.Ctx, rows, &borrowForStatistics)
	return borrowForStatistics
>>>>>>> 5fa5f02c1373b226cd4ab46bcdfa3326f6ae89d0
}
