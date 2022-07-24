package v1

import (
	"fmt"
	"github.com/gofiber/fiber/v2"
	"github.com/xuri/excelize/v2"
	"loante/global"
	"loante/service/model"
	"loante/service/payments"
	"loante/service/req"
	"loante/service/resp"
	"loante/tools"
	"math"
	"strings"
)

type payFlow struct {
}

func NewPayFlow() *payFlow {
	return new(payFlow)
}

type ordersQueryReq struct {
	req.PageReq
	Id               int    `json:"id" query:"id"`
	BorrowId         int    `json:"borrow_id" query:"borrow_id"`
	MchId            int    `json:"mch_id" query:"mch_id"`
	ProductId        int    `json:"product_id" query:"product_id"`
	UserName         string `json:"user_name" query:"user_name"`   //用户名
	IdCardNo         string `json:"id_card_no" query:"id_card_no"` //身份证号
	UrgeId           int    `json:"urge_id" query:"urge_id"`
	UrgeAdminId      int    `json:"urge_admin_id" query:"urge_admin_id"`
	Phone            string `json:"phone" query:"phone"`
	PaymentId        int    `json:"payment_id" query:"payment_id"`
	Status           int    `json:"status" query:"status"`           //还款状态
	PayChannel       int    `json:"pay_channel" query:"pay_channel"` //还款通道
	StartTime        string `json:"start_time" query:"start_time"`
	EndTime          string `json:"end_time" query:"end_time"`
	PaymentRequestNo string `json:"payment_request_no" query:"payment_request_no"` //发起还款请求编号
	ResStartTime     string `json:"res_start_time" query:"res_start_time"`
	ResEndTime       string `json:"res_end_time" query:"res_end_time"`
}

//Repayments 还款记录
func (a *payFlow) Repayments(c *fiber.Ctx) error {
	input := new(ordersQueryReq)
	if err := tools.ParseBody(c, input); err != nil {
		return resp.Err(c, 1, err.Error())
	}
	where := "o.type < 3"
	if input.MchId > 0 {
		where += fmt.Sprintf("o.mch_id = %d", input.MchId)
	}
	if input.BorrowId > 0 {
		where += fmt.Sprintf("o.bid = %d", input.BorrowId)
	}
	if input.ProductId > 0 {
		where += fmt.Sprintf("o.product_id = %d", input.ProductId)
	}
	if len(input.UserName) > 0 {
		where += fmt.Sprintf("u.aadhaar_name = '%s'", input.UserName)
	}
	if len(input.Phone) > 0 {
		where += fmt.Sprintf("u.phone = '%s'", input.Phone)
	}
	if input.PaymentId > 0 {
		where += fmt.Sprintf("o.payment = '%d'", input.PaymentId)
	}
	if input.Status > 0 {
		where += fmt.Sprintf("o.status = '%d'", input.Status)
	}
	if input.PayChannel > 0 {
		//where += fmt.Sprintf("o.status = '%d'", input.PayChannel)
	}
	if len(input.StartTime) > 0 {
		where += fmt.Sprintf("o.create_time >= '%s'", input.StartTime)
	}
	if len(input.EndTime) > 0 {
		where += fmt.Sprintf("o.create_time < '%s'", input.EndTime)
	}
	if len(input.PaymentRequestNo) > 0 {
		where += fmt.Sprintf("o.payment_request_no < '%s'", input.PaymentRequestNo)
	}

	lists, count := new(model.Orders).Page(where, input.Page, input.Size)
	return resp.OK(c, map[string]interface{}{
		"count": count,
		"list":  lists,
	})
}

func (a *payFlow) RepaymentsExport(c *fiber.Ctx) error {
	input := new(ordersQueryReq)
	if err := tools.ParseBody(c, input); err != nil {
		return resp.Err(c, 1, err.Error())
	}
	//导出文件
	t := tools.GetFormatTime()
	path := fmt.Sprintf("./static/xlsx/%s/%s.xlsx", t[0:7], tools.NewUUID())
	lists, _ := new(model.Orders).Gets("o.id > 0")
	f := excelize.NewFile()
	index := f.NewSheet("Sheet1")
	f.SetActiveSheet(index)
	for i, item := range lists {
		iStr := fmt.Sprintf("%d", i+1)
		f.SetCellValue("Sheet1", "A"+iStr, item.Id)
	}
	if err := f.SaveAs(path); err != nil {
		fmt.Println(err)
	}
	return resp.OK(c, map[string]interface{}{
		"path": strings.Trim(path, ","),
	})
}

type reconciliationQueryReq struct {
	req.PageReq
	BorrowId	int	`json:"borrow_id" query:"borrow_id"`  //借款的订单编号
	Name    string `json:"name" query:"name"` //用户名
	Phone       string `json:"phone" query:"phone"`			//手机号
	StartTime       string `json:"start_time" query:"start_time"`	//应还款时间
	EndTime       string `json:"end_time" query:"end_time"`		//应还款时间
	//Status      int `json:"status" query:"status"` //借款的订单状态
}
//Reconciliation 平账记录
func (a *payFlow)Reconciliation(c *fiber.Ctx) error {
	input := new(reconciliationQueryReq)
	if err := tools.ParseBody(c, input); err != nil {
		return resp.Err(c, 1, err.Error())
	}
	where := "o.type = 3"
	if input.BorrowId > 0{ //订单编号
		where += fmt.Sprintf("o.bid = %d", input.BorrowId)
	}
	if len(input.Name) > 0{
		where += fmt.Sprintf("user.aadhaar_name = '%s'", input.Name)
	}
	if len(input.Phone) > 0{
		where += fmt.Sprintf("user.phone = '%s'", input.Phone)
	}
	if len(input.StartTime) > 0{
		where += fmt.Sprintf("borrow.end_time > '%s'", input.StartTime)
	}
	if len(input.EndTime) > 0{
		where += fmt.Sprintf("borrow.end_time < '%s'", input.EndTime)
	}
	//if input.Status > 0{
	//	where += fmt.Sprintf("borrow.status = '%d'", input.Status)
	//}
	lists, count := new(model.Orders).Page(where, input.Page, input.Size)
	return resp.OK(c, map[string]interface{}{
		"count": count,
		"list":  lists,
	})
}

type depositQueryReq struct {
	req.PageReq
	BorrowId	int	`json:"borrow_id" query:"borrow_id"`  //借款的订单编号
	Name    string `json:"name" query:"name"` //用户名
	Phone       string `json:"phone" query:"phone"`			//手机号
	StartTime       string `json:"start_time" query:"start_time"`	//应还款时间
	EndTime       string `json:"end_time" query:"end_time"`		//应还款时间
}

//Deposits 入账
func (a *payFlow) Deposits(c *fiber.Ctx) error {
	input := new(depositQueryReq)
	if err := tools.ParseBody(c, input); err != nil {
		return resp.Err(c, 1, err.Error())
	}
	where := "o.type = 4"
	if input.BorrowId > 0{ //订单编号
		where += fmt.Sprintf("o.bid = %d", input.BorrowId)
	}
	if len(input.Name) > 0{
		where += fmt.Sprintf("user.aadhaar_name = '%s'", input.Name)
	}
	if len(input.Phone) > 0{
		where += fmt.Sprintf("user.phone = '%s'", input.Phone)
	}
	if len(input.StartTime) > 0{
		where += fmt.Sprintf("borrow.end_time > '%s'", input.StartTime)
	}
	if len(input.EndTime) > 0{
		where += fmt.Sprintf("borrow.end_time < '%s'", input.EndTime)
	}

	lists, count := new(model.Orders).Page(where, input.Page, input.Size)
	return resp.OK(c, map[string]interface{}{
		"count": count,
		"list":  lists,
	})
}

type loansQueryReq struct {
	req.PageReq
	BorrowId 	int `json:"borrow_id" query:"borrow_id"`  //借款id
	ProductId 	int `json:"product_id" query:"product_id"` //产品id
	PaymentId 	int `json:"payment_id" query:"payment_id"` //支付Id
	PaySuccess 	string `json:"pay_success" query:"pay_success"` //是否支付成功
	RequestNo 	string `json:"request_no" query:"request_no"` //放款请求编号
	StartTime 	string `json:"start_time" query:"start_time"` //创建时间
	EndTime   	string `json:"end_time" query:"end_time"`		//创建时间
	LoadStartTime 	string `json:"load_start_time" query:"load_start_time"` //放款时间
	LoadEndTime   	string `json:"load_end_time" query:"load_end_time"` //放款时间
}

//Loans 放款记录 从borrow表获取
func (a *payFlow) Loans(c *fiber.Ctx) error {
	input := new(loansQueryReq)
	if err := tools.ParseBody(c, input); err != nil {
		return resp.Err(c, 1, err.Error())
	}
	where := "b.id > 0"
	if input.BorrowId > 0{
		where += fmt.Sprintf(" and b.id = %d", input.BorrowId)
	}
	if input.ProductId > 0{
		where += fmt.Sprintf(" and b.product_id = %d", input.ProductId)
	}
	if input.PaymentId > 0{
		where += fmt.Sprintf(" and b.payment = %d", input.PaymentId)
	}
	if len(input.PaySuccess) > 0{
		where += fmt.Sprintf(" and loan_time > '2000-01-01'")
	}
	if len(input.RequestNo) > 0{
		where += fmt.Sprintf(" and payment_request_no = '%s'", input.RequestNo)
	}
	if len(input.StartTime) > 0{
		where += fmt.Sprintf(" and create_time >= '%s'", input.StartTime)
	}
	if len(input.EndTime) > 0{
		where += fmt.Sprintf(" and create_time < '%s'", input.EndTime)
	}
	if len(input.LoadStartTime) > 0{
		where += fmt.Sprintf(" and loan_time >= '%s'", input.LoadStartTime)
	}
	if len(input.LoadEndTime) > 0{
		where += fmt.Sprintf(" and loan_time < '%s'", input.LoadEndTime)
	}
	lists, count := new(model.Borrow).Page(where, input.Page, input.Size)
	return resp.OK(c, map[string]interface{}{
		"count": count,
		"list":  lists,
	})
}

//BatchLoans 批量重放款
func (a *payFlow) BatchLoans(c *fiber.Ctx) error {
	input := new(depositQueryReq)
	if err := tools.ParseBody(c, input); err != nil {
		return resp.Err(c, 1, err.Error())
	}
	lists, count := new(model.Orders).Page("o.id > 0", input.Page, input.Size)
	return resp.OK(c, map[string]interface{}{
		"count": count,
		"list":  lists,
	})
}

type utrPageReq struct {
	req.PageReq
	Name string	`query:"name" json:"name"`
	Phone   string	`query:"phone" json:"phone"`
	UtrCode    string	`query:"utr_code" json:"utr_code"`
	MchId     int	`query:"mch_id" json:"mch_id"`
	ProductId int	`query:"product_id" json:"product_id"`
	Status    int	`query:"status" json:"status"`
	StartTime string	`query:"start_time" json:"start_time"`
	EndTime   string	`query:"end_time" json:"end_time"`
}

//Utrs UTR对账单
func (a *payFlow)Utrs(c *fiber.Ctx) error {
	input := new(utrPageReq)
	if err := tools.ParseBody(c, input); err != nil {
		return resp.Err(c, 1, err.Error())
	}
	where := "bu.id > 0"
	if len(input.Name) > 0{
		where += fmt.Sprintf(" and user.aadhaar_name = '%s'", input.Name)
	}
	if len(input.Phone) > 0{
		where += fmt.Sprintf(" and user.phone = '%s'", input.Phone)
	}
	if len(input.UtrCode) > 0{
		where += fmt.Sprintf(" and bu.utr_code = '%s'", input.UtrCode)
	}
	if input.MchId > 0{
		where += fmt.Sprintf(" and bu.mch_id = '%d'", input.MchId)
	}
	if input.ProductId > 0{
		where += fmt.Sprintf(" and bu.product_id = '%d'", input.ProductId)
	}
	if input.Status > 0{
		where += fmt.Sprintf(" and bu.status = '%d'", input.Status)
	}
	if len(input.StartTime) > 0{
		where += fmt.Sprintf(" and bu.create >= '%s'", input.StartTime)
	}
	if len(input.EndTime) > 0{
		where += fmt.Sprintf("bu.create < '%s'", input.EndTime)
	}
	lists, count := new(model.BorrowUtr).Page(where, input.Page, input.Size)
	return resp.OK(c, map[string]interface{}{
		"count":count,
		"list":lists,
	})
}

type utrDismissedReq struct {
	req.PageReq
	Name string	`query:"name" json:"name"`
	Phone   string	`query:"phone" json:"phone"`
	MchId     int	`query:"mch_id" json:"mch_id"`
	ProductId int	`query:"product_id" json:"product_id"`
	StartTime string	`query:"start_time" json:"start_time"`
	EndTime   string	`query:"end_time" json:"end_time"`
}
//UtrsDismissed UTR驳回列表
func (a *payFlow)UtrsDismissed(c *fiber.Ctx) error {
	input := new(utrDismissedReq)
	if err := tools.ParseBody(c, input); err != nil {
		return resp.Err(c, 1, err.Error())
	}
	where := "bu.id > 0 and bu.status = 2"
	if len(input.Name) > 0{
		where += fmt.Sprintf(" and user.aadhaar_name = '%s'", input.Name)
	}
	if len(input.Phone) > 0{
		where += fmt.Sprintf(" and user.phone = '%s'", input.Phone)
	}
	if input.MchId > 0{
		where += fmt.Sprintf(" and bu.mch_id = '%d'", input.MchId)
	}
	if input.ProductId > 0{
		where += fmt.Sprintf(" and bu.product_id = '%d'", input.ProductId)
	}
	if len(input.StartTime) > 0{
		where += fmt.Sprintf(" and bu.create >= '%s'", input.StartTime)
	}
	if len(input.EndTime) > 0{
		where += fmt.Sprintf(" and bu.create < '%s'", input.EndTime)
	}
	lists, count := new(model.BorrowUtr).Page(where, input.Page, input.Size)
	return resp.OK(c, map[string]interface{}{
		"count": count,
		"list":  lists,
	})
}

type payPartialReq struct {
	BorrowId int	`query:"borrow_id" json:"borrow_id"`
	Type   int	`query:"type" json:"type"` //0 全额还款 1 部分还款 2 展期还款
	Amount   int	`query:"amount" json:"amount"`
	PaymentId   int	`query:"payment_id" json:"payment_id"`
}
//PayPartial 部分还款码<商户收款>
func (a *payFlow)PayPartial(c *fiber.Ctx) error {
	input := new(payPartialReq)
	if err := tools.ParseBody(c, input); err != nil {
		return resp.Err(c, 1, err.Error())
	}
	input.Amount = int(math.Abs(float64(input.Amount)))
	borrow := new(model.Borrow)
	borrow.One(fmt.Sprintf("id = %d", input.BorrowId))
	if borrow.Id == 0{
		return resp.Err(c, 1, "借款数据不存在！")
	}
	visit := new(model.BorrowVisit)
	visit.One(fmt.Sprintf("borrow_id = %d", borrow.Id))
	userData := new(model.User)
	userData.One(fmt.Sprintf("id = %d", borrow.Uid))
	//校对金额
	if input.Type != 1 && input.Type != 2{
		return resp.Err(c, 1, "还款方式不正确！")
	}
	if input.Type == 1 {
		if borrow.BeRepaidAmount + borrow.LatePaymentFee < input.Amount{
			input.Amount = borrow.BeRepaidAmount + borrow.LatePaymentFee
			input.Type = 0
		}
	}
	if input.Type == 2{ //展期
		////获取产品展期费用
		//productConfig := new(model.ProductDelayConfig)
		//productConfig.One(fmt.Sprintf("product_id = %d", borrow.ProductId))
		//if productConfig.Id == 0{
		//	return resp.Err(c, 1, "未找到产品展期配置！")
		//}
		//if productConfig.Status == 0{
		//	return resp.Err(c, 1, "产品展期未开启！")
		//}
		////borrow.PostponedPeriod += 1 展期还款成功以后才能 + 1
		////展期需要先还滞纳金 和 展期费用
		//input.Amount = int(float32(borrow.LoanAmount) * float32(productConfig.DelayRate)) + borrow.LatePaymentFee
		input.Amount = borrow.PostponeValuation + borrow.LatePaymentFee
	}
	//获取产品的默认支付通道
	if input.PaymentId == 0{
		payDefault := new(model.ProductPaymentDefault)
		payDefault.One(fmt.Sprintf("product_id = %d", borrow.ProductId))
		if payDefault.Id == 0{
			return resp.Err(c, 1, "产品的支付没有默认配置！")
		}
		input.PaymentId = payDefault.InPaymentId
	}
	//获取支付通道的配置
	pp := new(model.ProductPayment)
	pp.One(fmt.Sprintf("payment_id = %d and product_id = %d", input.PaymentId, borrow.ProductId))
	if pp.Id ==0{
		return resp.Err(c, 1, "未找到支付通道配置")
	}
	payment := new(model.Payment)
	payment.One(fmt.Sprintf("id = %d", 	input.PaymentId))
	if payment.Id == 0{
		return resp.Err(c, 1, "支付通道已经被删除！")
	}
	payModel :=  payments.SelectPay(payment.Name)
	paysData := payments.Pays{
		OrderId:fmt.Sprintf("%s-%d-%d-%d", tools.InviteCode(6), borrow.Id, borrow.ProductId, input.PaymentId),
		Amount:float64(input.Amount),
		CustomName:userData.AadhaarName,
		CustomMobile:userData.Phone,
		CustomEmail:userData.Email,
		BankAccount:userData.Bankcard,
		IfscCode:userData.Ifsc,
		Remark:"",
		NotifyUrl:global.C.Http.PayInNotify,
		CallbackUrl:global.C.Http.PayInNotify,
	}
	ret, data, err := (*payModel).PayIn(pp.Configuration, &paysData)
	if !ret {
		return  resp.Err(c, 1, err.Error())
	}
	//保存order
	orders := new(model.Orders)
	orders.MchId = borrow.MchId
	orders.ProductId = borrow.ProductId
	orders.Bid = borrow.Id
	orders.Uid = borrow.Uid
	orders.CreateTime = tools.GetFormatTime()
	orders.Type = input.Type
	orders.Payment = input.PaymentId
	orders.PostponePeriod = borrow.PostponedPeriod
	orders.ApplyAmount = int(input.Amount)
	orders.RepaidStatus = 2
	orders.UrgeId = visit.UrgeId
	orders.RepaidUrl = data["url"].(string)
	orders.RepaidType = 0
	orders.InvalidTime = tools.ToAddMinute(30)
	orders.LoanTime = borrow.LoanTime
	orders.EndTime = borrow.EndTime
	orders.PaymentRequestNo = paysData.OrderId
	orders.PaymentRespondNo = data["platId"].(string)
	orders.Insert()
	return resp.OK(c, map[string]string{
		"url":data["url"].(string),
	})
}

type payOutReq struct {
	BorrowId 			int		`json:"borrow_id"` //借贷
}
//PayOut 代付用于放款
func (a *payFlow)PayOut(c *fiber.Ctx) error {
	input := new(payOutReq)
	if err := tools.ParseBody(c, input); err != nil {
		return resp.Err(c, 1, err.Error())
	}
	//检测借贷是否存在
	borrow := new(model.Borrow)
	borrow.One(fmt.Sprintf("id = %d", input.BorrowId))
	if borrow.Id == 0{
		return resp.Err(c, 1, "借款数据不存在！")
	}
	if borrow.Status != 4{ //放款中
		return resp.Err(c, 1, "状态不正确！")
	}
	amount := float32(borrow.LoanAmount)
	//获取服务费
	productData := new(model.Product)
	productData.One(fmt.Sprintf("id = %d", borrow.ProductId))
	if productData.IsStopLending == 1{//已经停止
		return resp.Err(c, 1, "已停止放款！")
	}
	amount = amount - amount * productData.RateService
	userData := new(model.User)
	userData.One(fmt.Sprintf("id = %d", borrow.Uid))

	ppc, err := borrow.GetPaymentConfig("out")
	if err!= nil{
		return resp.Err(c, 1, err.Error())
	}
	payModel :=  payments.SelectPay(ppc.Name)
	paysData := payments.Pays{
		OrderId:fmt.Sprintf("%s-%d-%d", tools.InviteCode(6), borrow.Id, borrow.ProductId),
		Amount:float64(amount),
		CustomName:userData.AadhaarName,
		CustomMobile:userData.Phone,
		CustomEmail:userData.Email,
		BankAccount:userData.Bankcard,
		IfscCode:userData.Ifsc,
		Remark:"",
		NotifyUrl:"http://127.0.0.1/notify",
		CallbackUrl:"http://127.0.0.1/notify",
	}
	ret, err := (*payModel).PayOut(ppc.Config, &paysData)
	if !ret {
		return  resp.Err(c, 1, err.Error())
	}
	//开始
	borrow.Payment = ppc.Id
	borrow.PaymentRequestNo = paysData.OrderId
	borrow.PaymentRespond = paysData.PlatOrderId
	borrow.Update(fmt.Sprintf("id = %d", borrow.Id))
	return resp.OK(c, map[string]string{
		"url":"",
	})
}
func (a *payFlow)UtrsVerify(c *fiber.Ctx) error {
	input := new(req.IdReq)
	if err := tools.ParseBody(c, input); err != nil {
		return resp.Err(c, 1, err.Error())
	}
	//获取Borrow
	utr := new(model.BorrowUtr)
	utr.One(fmt.Sprintf("id = %d", input.Id))
	if utr.Id == 0{
		return resp.Err(c, 1, "UTR数据不存在！")
	}
	//获取产品的默认支付通道
	payDefault := new(model.ProductPaymentDefault)
	payDefault.One(fmt.Sprintf("product_id = %d", 	utr.ProductId))
	if payDefault.Id == 0{
		return resp.Err(c, 1, "产品的支付没有默认配置！")
	}
	payment := new(model.Payment)
	payment.One(fmt.Sprintf("id = %d", 	payDefault.InPaymentId))
	if payment.Id == 0{
		return resp.Err(c, 1, "支付通道已经被删除！")
	}
	//是否支持utr查单
	if payment.IsUtrQuery == 0{
		return resp.Err(c, 1, payment.Name+" 不支持UTR查单！")
	}
	//获取支付通道的配置
	pp := new(model.ProductPayment)
	pp.One(fmt.Sprintf("payment_id = %d and product_id = %d", payDefault.InPaymentId, utr.ProductId))
	if pp.Id ==0{
		return resp.Err(c, 1, "未找到支付通道配置")
	}

	return resp.OK(c, "")
}