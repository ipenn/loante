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
	ProductId       int    `query:"product_id" json:"product_id"`
	UserId          int    `query:"user_id" json:"user_id"`
	Name            string `query:"name" json:"name"`
	NoApplying      string `query:"no_applying" json:"no_applying"`             //去除申请中的
	ProcessingInPay string `query:"processing_in_pay" json:"processing_in_pay"` //支付公司放款处理中
	StartTime       string		`json:"start_time" query:"start_time"`
	EndTime         string		`json:"end_time" query:"end_time"`
	Phone           string		`json:"phone" query:"phone"`
	IdNo            string		`json:"id_no" query:"id_no"` //身份证号
	Id              int		`json:"id" query:"id"`    //订单编号
	Status          int		`json:"status" query:"status"`    //订单状态
	Postponed       int		`json:"postponed" query:"postponed"`    //订单类型 是否展期
	LoanType        int		`json:"loan_type" query:"loan_type"`    //贷款类型 平台首贷 平台复贷产品首贷 平台复贷产品复贷
							//string //批次编号
							// INT //APP包名
							// INT //渠道来源
	RiskModel int	`json:"risk_model" query:"risk_model"` //模型类型
	Payment          int	`json:"payment" query:"payment"`    //支付公司
	BeRepaidAmount   int	`json:"be_repaid_amount" query:"be_repaid_amount"`    //待还金额(小于)
	PaymentRequestNo string	`json:"payment_request_no" query:"payment_request_no"` //放款请求单号
	LoanStartTime    string	`json:"loan_start_time" query:"loan_start_time"` //放款时间
	LoanEndTime      string	`json:"loan_end_time" query:"loan_end_time"` //放款时间
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
	if input.UserId > 0 {
		where += fmt.Sprintf(" and u.aadhaar_name ='%s'", input.Name)
	}
	lists, count := new(model.Borrow).Page("b.id > 0", input.Page, input.Size)
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
	if len(input.NoApplying) > 0{
		where += " and b.status > 4"
	}
	if len(input.ProcessingInPay) > 0{
		where += " and b.status = 4"
	}
	if input.ProductId > 0{
		where += fmt.Sprintf(" and b.product_id =%d", input.ProductId)
	}
	if input.UserId > 0{
		where += fmt.Sprintf(" and b.uid =%d", input.UserId)
	}
	if input.UserId > 0{
		where += fmt.Sprintf(" and u.aadhaar_name ='%s'", input.Name)
	}
	lists, count := new(model.Borrow).Gets("b.id > 0")
	return resp.OK(c, map[string]interface{}{
		"count":count,
		"list":lists,
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
	if borrowData.Id > 0{
		return resp.Err(c, 1, "未找到记录")
	}
	if borrowData.Status == 8 || borrowData.Status == 9{
		return resp.Err(c, 1, "借贷已经完成")
	}
	fundData := new(model.BorrowFund)
	fundData.Amount = input.Amount
	fundData.Remark = input.Remark
	fundData.BorrowId = borrowData.Id
	fundData.UserId = borrowData.Uid
	fundData.BeRepaidAmount = float64(borrowData.BeRepaidAmount)
	fundData.RemainingAmount = float64(borrowData.BeRepaidAmount) - input.Amount
	fundData.Insert()
	if fundData.Id > 0{
		//borrowData.BeRepaidAmount -= input.Amount
		borrowData.BeRepaidAmount = 0
		borrowData.Status = 8
		if borrowData.EndTime < tools.GetFormatTime(){
			borrowData.Status = 9
		}
		if borrowData.BeRepaidAmount <= 0{ //已经还完
		}
		borrowData.Update(fmt.Sprintf("id = %d", input.Id))
	}
	return resp.OK(c, "")
}

type depositReq struct {
	req.IdReq
	PaymentId int	`json:"payment_id"`
	Amount float64	`json:"amount"`
	OrderNo string `json:"OrderNo"`
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
	if borrowData.Id > 0{
		return resp.Err(c, 1, "未找到记录")
	}
	if borrowData.Status == 8 || borrowData.Status == 9{
		return resp.Err(c, 1, "借贷已经完成")
	}
	fundData := new(model.BorrowFund)
	fundData.OrderNo = input.OrderNo
	fundData.Amount = input.Amount
	fundData.Remark = input.Remark
	fundData.BorrowId = borrowData.Id
	fundData.UserId = borrowData.Uid
	fundData.PaymentId = input.PaymentId
	fundData.BeRepaidAmount = float64(borrowData.BeRepaidAmount)
	fundData.RemainingAmount = float64(borrowData.BeRepaidAmount) - input.Amount
	fundData.Insert()
	if fundData.Id > 0{
		//borrowData.BeRepaidAmount -= input.Amount
		borrowData.BeRepaidAmount = 0
		borrowData.Status = 8
		if borrowData.EndTime < tools.GetFormatTime(){
			borrowData.Status = 9
		}
		if borrowData.BeRepaidAmount <= 0{ //已经还完
		}
		borrowData.Update(fmt.Sprintf("id = %d", input.Id))
	}
	return resp.OK(c, "")
}

