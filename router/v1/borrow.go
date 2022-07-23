package v1

import (
	"fmt"
	"github.com/gofiber/fiber/v2"
	"loante/service/model"
	"loante/service/req"
	"loante/service/resp"
	"loante/tools"
	"math"
)

type borrow struct{}

func NewBorrow() *borrow {
	return new(borrow)
}

type borrowQueryReq struct {
	req.PageReq
	ProductId       int    `query:"product_id" json:"product_id"` //产品编号
	UserId          int    `query:"user_id" json:"user_id"` //客户编号
	Name            string `query:"name" json:"name"` //客户名称
	NoApplying      string `query:"no_applying" json:"no_applying"`             //去除申请中的
	ProcessingInPay string `query:"processing_in_pay" json:"processing_in_pay"` //支付公司放款处理中
	StartTime       string		`json:"start_time" query:"start_time"`	//创建开始时间
	EndTime         string		`json:"end_time" query:"end_time"`	//创建结束时间
	Phone           string		`json:"phone" query:"phone"` //手机号码
	IdNo            string		`json:"id_no" query:"id_no"` //身份证号
	Id              int		`json:"id" query:"id"`    //订单编号
	Status          int		`json:"status" query:"status"`    //订单状态
	Postponed       int		`json:"postponed" query:"postponed"`    //订单类型 是否展期
	LoanType        int		`json:"loan_type" query:"loan_type"`    //贷款类型 平台首贷 平台复贷产品首贷 平台复贷产品复贷
	RiskModel int	`json:"risk_model" query:"risk_model"` //模型类型
	Payment          int	`json:"payment" query:"payment"`    //支付公司
	BeRepaidAmount   int	`json:"be_repaid_amount" query:"be_repaid_amount"`    //待还金额(小于)
	PaymentRequestNo string	`json:"payment_request_no" query:"payment_request_no"` //放款请求单号
	LoanStartTime    string	`json:"loan_start_time" query:"loan_start_time"` //放款时间
	LoanEndTime      string	`json:"loan_end_time" query:"loan_end_time"` //放款时间

	//string //批次编号
	// INT //APP包名
	// INT //渠道来源
}

func (a *borrow) Query(c *fiber.Ctx) error {
	input := new(borrowQueryReq)
	if err := tools.ParseBody(c, input); err != nil {
		return resp.Err(c, 1, err.Error())
	}
	where := "b.id > 0"
	if len(input.NoApplying) > 0 {
		where += " and b.status > 4"
	}
	if len(input.ProcessingInPay) > 0 {
		where += " and b.status = 4"
	}
	if input.ProductId > 0 {
		where += fmt.Sprintf(" and b.product_id =%d", input.ProductId)
	}
	if input.UserId > 0 {
		where += fmt.Sprintf(" and b.uid =%d", input.UserId)
	}
	if len(input.Name) > 0 {
		where += fmt.Sprintf(" and user.aadhaar_name ='%s'", input.Name)
	}
	if len(input.StartTime) > 0 {
		where += fmt.Sprintf(" and b.create_time >= '%s'", input.StartTime)
	}
	if len(input.EndTime) > 0 {
		where += fmt.Sprintf(" and b.create_time < '%s'", input.EndTime)
	}
	if len(input.Phone) > 0 {
		where += fmt.Sprintf(" and user.phone = '%s'", input.Phone)
	}
	if len(input.Phone) > 0 { //身份证号
		where += fmt.Sprintf(" and user.phone = '%s'", input.Phone)
	}
	if input.Id > 0 {
		where += fmt.Sprintf(" and b.id = '%d'", input.Id)
	}
	if input.Status > -1 {
		where += fmt.Sprintf(" and b.status = '%d'", input.Status)
	}
	if input.Postponed > -1 {
		where += fmt.Sprintf(" and b.postponed = '%d'", input.Postponed)
	}
	if input.LoanType > -1 {
		where += fmt.Sprintf(" and b.loan_type = '%d'", input.LoanType)
	}
	if input.RiskModel > 0 {
		where += fmt.Sprintf(" and b.risk_model = '%d'", input.RiskModel)
	}
	if input.Payment > 0 {
		where += fmt.Sprintf(" and b.payment = '%d'", input.Payment)
	}
	if input.BeRepaidAmount > 0 {
		where += fmt.Sprintf(" and b.be_repaid_amount < '%d'", input.BeRepaidAmount)
	}
	if len(input.PaymentRequestNo) > 0 {
		where += fmt.Sprintf(" and b.payment_request_no = '%s'", input.PaymentRequestNo)
	}
	if len(input.LoanStartTime) > 0 {
		where += fmt.Sprintf(" and b.loan_time >= '%s'", input.LoanStartTime)
	}
	if len(input.LoanEndTime) > 0 {
		where += fmt.Sprintf(" and b.loan_time < '%s'", input.LoanEndTime)
	}
	lists, count := new(model.Borrow).Page(where, input.Page, input.Size)
	return resp.OK(c, map[string]interface{}{
		"count": count,
		"list":  lists,
	})
}

//QueryExport 订单导出
func (a *borrow)QueryExport(c *fiber.Ctx) error {
	input := new(borrowQueryReq)
	if err := tools.ParseBody(c, input); err != nil {
		return resp.Err(c, 1, err.Error())
	}
	where := "b.id > 0"
	if len(input.NoApplying) > 0 {
		where += " and b.status > 4"
	}
	if len(input.ProcessingInPay) > 0 {
		where += " and b.status = 4"
	}
	if input.ProductId > 0 {
		where += fmt.Sprintf(" and b.product_id =%d", input.ProductId)
	}
	if input.UserId > 0 {
		where += fmt.Sprintf(" and b.uid =%d", input.UserId)
	}
	if input.UserId > 0 {
		where += fmt.Sprintf(" and user.aadhaar_name ='%s'", input.Name)
	}
	lists, count := new(model.Borrow).Gets(where)
	return resp.OK(c, map[string]interface{}{
		"count": count,
		"list":  lists,
	})
}

type reconciliationReq struct {
	req.IdReq
	Amount float64	`json:"amount"`
	Remark  string `json:"remark"`
}

//Reconciliation 平账操作  不管还多少金额 最终需要更新订单的待还金额为0 标记已还款
func (a *borrow)Reconciliation(c *fiber.Ctx) error {
	input := new(reconciliationReq)
	if err := tools.ParseBody(c, input); err != nil {
		return resp.Err(c, 1, err.Error())
	}
	borrowData := new(model.Borrow)
	borrowData.One(fmt.Sprintf("id = %d", input.Id))
	input.Amount = math.Abs(input.Amount)
	if borrowData.Id == 0{
		return resp.Err(c, 1, "未找到记录")
	}
	if borrowData.Status == 8 || borrowData.Status == 9{
		return resp.Err(c, 1, "借贷已经完成")
	}
	t := tools.GetFormatTime()
	fundData := new(model.Orders)
	fundData.MchId = borrowData.MchId
	fundData.ProductId = borrowData.ProductId
	fundData.Remark = input.Remark
	fundData.Bid = borrowData.Id
	fundData.Uid = borrowData.Uid
	fundData.RepaidType = 1 //人工平账
	fundData.RepaidStatus = 1
	fundData.Type = 3
	fundData.ApplyAmount = int(input.Amount)
	fundData.ActualAmount = int(input.Amount)
	fundData.CreateTime = t
	fundData.Insert()
	if fundData.Id > 0{
		//borrowData.BeRepaidAmount -= input.Amount
		borrowData.BeRepaidAmount = 0	//待还金额
		borrowData.LatePaymentFee = 0	//滞纳金
		borrowData.CompleteTime = t		//还款完成时间
		borrowData.Status = 8
		if borrowData.EndTime < t{
			borrowData.Status = 9
		}
		if borrowData.BeRepaidAmount <= 0{ //已经还完
		}
		borrowData.Update(fmt.Sprintf("id = %d", input.Id))
		//提额
		new(model.UserQuota).Increase(borrowData.ProductId, borrowData.Uid, borrowData.Status)
	}
	return resp.OK(c, "")
}

type depositReq struct {
	req.IdReq
	PaymentId int	`json:"payment_id"`
	Amount float64	`json:"amount"`
	OrderNo string `json:"order_no"`
	Remark  string `json:"remark"`
}
//Deposit 入账操作 不管还多少金额 最终需要更新订单的待还金额为0 标记已还款
func (a *borrow)Deposit(c *fiber.Ctx) error {
	input := new(depositReq)
	if err := tools.ParseBody(c, input); err != nil {
		return resp.Err(c, 1, err.Error())
	}
	borrowData := new(model.Borrow)
	borrowData.One(fmt.Sprintf("id = %d", input.Id))
	input.Amount = math.Abs(input.Amount)
	if borrowData.Id == 0{
		return resp.Err(c, 1, "未找到记录")
	}
	if borrowData.Status == 8 || borrowData.Status == 9{
		return resp.Err(c, 1, "借贷已经完成")
	}
	t := tools.GetFormatTime()
	fundData := new(model.Orders)
	fundData.MchId = borrowData.MchId
	fundData.Remark = input.Remark
	fundData.Payment = input.PaymentId  // 0 表示私卡账户
	fundData.Bid = borrowData.Id
	fundData.Uid = borrowData.Uid
	fundData.RepaidType = 1
	fundData.RepaidStatus = 1
	fundData.Type = 3
	fundData.ApplyAmount = int(input.Amount)
	fundData.ActualAmount = int(input.Amount)
	fundData.CreateTime = t
	fundData.Insert()
	if fundData.Id > 0{
		//borrowData.BeRepaidAmount -= input.Amount
		borrowData.BeRepaidAmount = 0	//待还金额
		borrowData.LatePaymentFee = 0	//滞纳金
		borrowData.Status = 8
		borrowData.CompleteTime = t		//还款完成时间
		if borrowData.EndTime < t{
			borrowData.Status = 9
		}
		if borrowData.BeRepaidAmount <= 0{ //已经还完
		}
		borrowData.Update(fmt.Sprintf("id = %d", input.Id))
		//提额
		new(model.UserQuota).Increase(borrowData.ProductId, borrowData.Uid, borrowData.Status)
	}
	return resp.OK(c, "")
}

type borrowFundsReq struct {
	req.IdReq
	req.PageReq
}
func (a *borrow)Funds(c *fiber.Ctx) error {
	input := new(borrowFundsReq)
	if err := tools.ParseBody(c, input); err != nil {
		return resp.Err(c, 1, err.Error())
	}
	borrowData := new(model.Borrow)
	borrowData.One(fmt.Sprintf("id = %d", input.Id))
	if borrowData.Id == 0{
		return resp.Err(c, 1, "未找到记录")
	}
	order := new(model.Orders)
	lists, count := order.Page(fmt.Sprintf("o.bid = %d and o.repaid_status = 1", input.Id), input.Page, input.Size)
	return resp.OK(c, map[string]interface{}{
		"count": count,
		"list":  lists,
	})
}

func (a *borrow)SetLoanFail(c *fiber.Ctx) error {
	input := new(req.IdReq)
	if err := tools.ParseBody(c, input); err != nil {
		return resp.Err(c, 1, err.Error())
	}
	borrowData := new(model.Borrow)
	borrowData.One(fmt.Sprintf("id = %d", input.Id))
	if borrowData.Id == 0{
		return resp.Err(c, 1, "未找到记录")
	}
	borrowData.Status = 0
	borrowData.Update(fmt.Sprintf("id = %d", input.Id))
	return resp.OK(c, "")
}
