package v1

import (
	"fmt"
	"github.com/gofiber/fiber/v2"
	"github.com/xuri/excelize/v2"
	"loante/service/model"
	"loante/service/req"
	"loante/service/resp"
	"loante/tools"
	"strings"
)

type payFlow struct {

}

func NewPayFlow() *payFlow {
	return new(payFlow)
}

type ordersQueryReq struct {
	req.PageReq
	Id			int	`json:"id" query:"id"`
	BorrowId         int    `json:"borrow_id" query:"borrow_id"`
	MchId            int    `json:"mch_id" query:"mch_id"`
	ProductId        int    `json:"product_id" query:"product_id"`
	UserName         string `json:"user_name" query:"user_name"`	  //用户名
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
func (a *payFlow)Repayments(c *fiber.Ctx) error {
	input := new(ordersQueryReq)
	if err := tools.ParseBody(c, input); err != nil {
		return resp.Err(c, 1, err.Error())
	}
	where := "o.id > 0"
	if input.MchId >0 {
		where += fmt.Sprintf("o.mch_id = %d", input.MchId)
	}
	if input.BorrowId >0 {
		where += fmt.Sprintf("o.bid = %d", input.BorrowId)
	}
	if input.ProductId >0 {
		where += fmt.Sprintf("o.product_id = %d", input.ProductId)
	}
	if len(input.UserName) >0 {
		where += fmt.Sprintf("u.aadhaar_name = '%s'", input.UserName)
	}
	if len(input.Phone) >0 {
		where += fmt.Sprintf("u.phone = '%s'", input.Phone)
	}
	if input.PaymentId >0 {
		where += fmt.Sprintf("o.payment = '%d'", input.PaymentId)
	}
	if input.Status >0 {
		where += fmt.Sprintf("o.status = '%d'", input.Status)
	}
	if input.PayChannel >0 {
		//where += fmt.Sprintf("o.status = '%d'", input.PayChannel)
	}
	if len(input.StartTime)> 0{
		where += fmt.Sprintf("o.create_time >= '%s'", input.StartTime)
	}
	if len(input.EndTime)> 0{
		where += fmt.Sprintf("o.create_time < '%s'", input.EndTime)
	}
	if len(input.PaymentRequestNo)> 0{
		where += fmt.Sprintf("o.payment_request_no < '%s'", input.PaymentRequestNo)
	}

	lists, count := new(model.Orders).Page(where, input.Page, input.Size)
	return resp.OK(c, map[string]interface{}{
		"count":count,
		"list":lists,
	})
}

func (a *payFlow)RepaymentsExport(c *fiber.Ctx) error {
	input := new(ordersQueryReq)
	if err := tools.ParseBody(c, input); err != nil {
		return resp.Err(c, 1, err.Error())
	}
	//导出文件
	t := tools.GetFormatTime()
	path := fmt.Sprintf("./static/xlsx/%s/%s.xlsx", t[0:7],tools.NewUUID())
	lists, _ := new(model.Orders).Gets("o.id > 0")
	f := excelize.NewFile()
	index := f.NewSheet("Sheet1")
	f.SetActiveSheet(index)
	for i, item := range lists{
		iStr := fmt.Sprintf("%d", i+1)
		f.SetCellValue("Sheet1", "A" + iStr, item.Id)
	}
	if err := f.SaveAs(path); err != nil {
		fmt.Println(err)
	}
	return resp.OK(c, map[string]interface{}{
		"path":strings.Trim(path,","),
	})
}

type reconciliationQueryReq struct {
	req.PageReq
	Id	int	`json:"id" query:"id"`
	BorrowId	string	`json:"borrow_id" query:"borrow_id"`
	UserName    string `json:"user_name" query:"user_name"`
	Phone       string `json:"phone" query:"phone"`
	Status      int `json:"status" query:"status"`
}
//Reconciliation 平账
func (a *payFlow)Reconciliation(c *fiber.Ctx) error {
	input := new(reconciliationQueryReq)
	if err := tools.ParseBody(c, input); err != nil {
		return resp.Err(c, 1, err.Error())
	}
	lists, count := new(model.Orders).Page("o.id > 0", input.Page, input.Size)
	return resp.OK(c, map[string]interface{}{
		"count":count,
		"list":lists,
	})
}


type depositQueryReq struct {
	req.PageReq
	Id	int	`json:"id" query:"id"`
	BorrowId	string	`json:"borrow_id" query:"borrow_id,omitempty"`
	UserName    string `json:"user_name" query:"user_name"`
	Phone       string `json:"phone" query:"phone"`
	Status      int `json:"status" query:"status"`
}
//Deposits 入账
func (a *payFlow)Deposits(c *fiber.Ctx) error {
	input := new(depositQueryReq)
	if err := tools.ParseBody(c, input); err != nil {
		return resp.Err(c, 1, err.Error())
	}
	lists, count := new(model.Orders).Page("o.id > 0", input.Page, input.Size)
	return resp.OK(c, map[string]interface{}{
		"count":count,
		"list":lists,
	})
}

type loansQueryReq struct {
	req.PageReq
	Id	int	`json:"id" query:"id"`
	BorrowId	string	`json:"borrow_id" query:"borrow_id"`
	UserName    string `json:"user_name" query:"user_name"`
	Phone       string `json:"phone" query:"phone"`
	Status      int `json:"status" query:"status"`
}
//Loans 放款
func (a *payFlow)Loans(c *fiber.Ctx) error {
	input := new(depositQueryReq)
	if err := tools.ParseBody(c, input); err != nil {
		return resp.Err(c, 1, err.Error())
	}
	lists, count := new(model.Orders).Page("o.id > 0", input.Page, input.Size)
	return resp.OK(c, map[string]interface{}{
		"count":count,
		"list":lists,
	})
}

//BatchLoans 批量重放款
func (a *payFlow)BatchLoans(c *fiber.Ctx) error {
	input := new(depositQueryReq)
	if err := tools.ParseBody(c, input); err != nil {
		return resp.Err(c, 1, err.Error())
	}
	lists, count := new(model.Orders).Page("o.id > 0", input.Page, input.Size)
	return resp.OK(c, map[string]interface{}{
		"count":count,
		"list":lists,
	})
}
