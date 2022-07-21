package v1

import (
	"fmt"
	"github.com/gofiber/fiber/v2"
	"loante/service/model"
	"loante/service/req"
	"loante/service/resp"
	"loante/tools"
)

type visit struct {

}

func NewVisit() *visit {
	return new(visit)
}


type remindPageReq struct {
	req.PageReq
	UserName           string		`json:"user_name"` //用户名
	Phone              string		`json:"phone"` //手机号
	BorrowId           int		`json:"borrow_id"`    //订单编号
	BorrowComment      string		`json:"borrow_comment"`    //标签
	OverDueDays        int		`json:"over_due_days"`    //到期天数
	Wish               string		`json:"wish"`    //还款意愿
	Status             int		`json:"status"`    //订单状态
	MchId              int		`json:"mch_id"`    //商户名称
	ProductId          int		`json:"product_id"`    //产品名称
	RemindCompanyId    int		`json:"remind_company_id"`    //催收公司
	RemindId           int		`json:"remind_id"`    //催收员
	PayStartTime       string		`json:"pay_start_time"`    //还款开始时间
	PayEndTime         string		`json:"pay_end_time"`    //还款结束时间
	VisitStartTime     string		`json:"visit_start_time"`    //最近一次跟进时间
	VisitEndTime       string		`json:"visit_end_time"`    //最近一次跟进时间
	RepaymentStartTime string		`json:"repayment_start_time"`    //应还款时间
	RepaymentEndTime   string		`json:"repayment_end_time"`    //应还款时间
}
func (a *visit)remindBorrowQuery(input *remindPageReq) string{
	where := " "
	return where
}
//RemindBorrowAll 预提醒订单  预提醒中的订单 + 预提醒阶段就完成的订单
func (a *visit)RemindBorrowAll(c *fiber.Ctx) error {
	input := new(remindPageReq)
	if err := tools.ParseBody(c, input); err != nil {
		return resp.Err(c, 1, err.Error())
	}
	where := "(b.status = 8 or bv.remind_id = 0)" + a.remindBorrowQuery(input) //在预提醒阶段完成的订单
	lists, count := new(model.BorrowVisit).RemindPage(where, input.Page, input.Size)
	return resp.OK(c, map[string]interface{}{
		"count":count,
		"list":lists,
	})
}
//RemindBorrowing 预提醒中订单  预提醒中的订单
func (a *visit)RemindBorrowing(c *fiber.Ctx) error {
	input := new(remindPageReq)
	if err := tools.ParseBody(c, input); err != nil {
		return resp.Err(c, 1, err.Error())
	}
	where := "bv.remind_id = 0" + a.remindBorrowQuery(input) //在预提醒阶段完成的订单
	lists, count := new(model.BorrowVisit).RemindPage(where, input.Page, input.Size)
	return resp.OK(c, map[string]interface{}{
		"count":count,
		"list":lists,
	})
}
//RemindBorrowed 预提醒完成的订单  预提醒阶段就完成的订单
func (a *visit)RemindBorrowed(c *fiber.Ctx) error {
	input := new(remindPageReq)
	if err := tools.ParseBody(c, input); err != nil {
		return resp.Err(c, 1, err.Error())
	}
	where := "b.status = 8" + a.remindBorrowQuery(input) //在预提醒阶段完成的订单
	lists, count := new(model.BorrowVisit).RemindPage(where, input.Page, input.Size)
	return resp.OK(c, map[string]interface{}{
		"count":count,
		"list":lists,
	})
}

type remindDetailReq struct {
	req.PageReq
	UserName           string		`json:"user_name"` //用户名
	Phone              string		`json:"phone"` //债务人手机号
	BorrowId           int			`json:"borrow_id"`    //订单编号
	OverDueDays        int		`json:"over_due_days"` //到期天数
	ContactPhone string		`json:"contact_phone"`	//联系人手机号
	Wish string		`json:"wish"`	//还款意愿
	Relationship int		`json:"relationship"`	//联系人与债主的关系
	Tag int		`json:"tag"`	//标签
	UrgeId        int		`json:"urge_id"`			//催收人员
	UrgeCompanyId int		`json:"urge_company_id"`	//催收公司
	UrgeGroupId   int		`json:"urge_group_id"`   	//催收组
	StartTime   	string		`json:"start_time"` 	//开始时间
	EndTime   	string		`json:"end_time"`			//结束时间
}
//RemindDetail 预提醒记录
func (a *visit)RemindDetail(c *fiber.Ctx) error {
	input := new(remindDetailReq)
	if err := tools.ParseBody(c, input); err != nil {
		return resp.Err(c, 1, err.Error())
	}
	where := "bvd.type = 0 "
	lists, count := new(model.BorrowVisitDetail).Page(where, input.Page, input.Size)
	return resp.OK(c, map[string]interface{}{
		"count":count,
		"list":lists,
	})
}

//UrgeBorrowAll 催收订单  催收中的订单 + 催收阶段完成的订单
func (a *visit)UrgeBorrowAll(c *fiber.Ctx) error {
	input := new(remindPageReq)
	if err := tools.ParseBody(c, input); err != nil {
		return resp.Err(c, 1, err.Error())
	}
	where := "(bv.remind_id > 0)" + a.remindBorrowQuery(input)
	lists, count := new(model.BorrowVisit).UrgePage(where, input.Page, input.Size)
	return resp.OK(c, map[string]interface{}{
		"count":count,
		"list":lists,
	})
}


//UrgeBorrowing 催收中的订单
func (a *visit)UrgeBorrowing(c *fiber.Ctx) error {
	input := new(remindPageReq)
	if err := tools.ParseBody(c, input); err != nil {
		return resp.Err(c, 1, err.Error())
	}
	where := "(bv.remind_id > 0 and b.status != 9)" + a.remindBorrowQuery(input)
	lists, count := new(model.BorrowVisit).UrgePage(where, input.Page, input.Size)
	return resp.OK(c, map[string]interface{}{
		"count":count,
		"list":lists,
	})
}


//UrgeBorrowed 催收阶段完成的订单
func (a *visit)UrgeBorrowed(c *fiber.Ctx) error {
	input := new(remindPageReq)
	if err := tools.ParseBody(c, input); err != nil {
		return resp.Err(c, 1, err.Error())
	}
	where := "(bv.remind_id > 0 and b.status = 9)" + a.remindBorrowQuery(input)
	lists, count := new(model.BorrowVisit).UrgePage(where, input.Page, input.Size)
	return resp.OK(c, map[string]interface{}{
		"count":count,
		"list":lists,
	})
}


//UrgeDetail 催收记录
func (a *visit)UrgeDetail(c *fiber.Ctx) error {
	input := new(remindDetailReq)
	if err := tools.ParseBody(c, input); err != nil {
		return resp.Err(c, 1, err.Error())
	}
	where := "bvd.type = 1"
	lists, count := new(model.BorrowVisitDetail).Page(where, input.Page, input.Size)
	return resp.OK(c, map[string]interface{}{
		"count":count,
		"list":lists,
	})
}

type urgeReportReq struct {
	req.PageReq
	StartTime string
	EndTime string
	MchId	int
	Type int
	UrgeCompanyId	int
	UrgeGroupId	int
	UrgeId	int
	Group []string
}
func (a *visit)report(input *urgeReportReq) string{
	where := ""
	if input.MchId > 0{
		where += fmt.Sprintf(" and mch_id = %d", input.MchId)
	}
	if input.UrgeCompanyId > 0{
		where += fmt.Sprintf(" and urge_company_id = %d", input.UrgeCompanyId)
	}
	if len(input.StartTime) > 0{
		where += fmt.Sprintf(" and urge_company_id = %d", input.UrgeCompanyId)
	}
	return where
}
//UrgeReport 催收业绩
func (a *visit)UrgeReport(c *fiber.Ctx) error {
	input := new(urgeReportReq)
	if err := tools.ParseBody(c, input); err != nil {
		return resp.Err(c, 1, err.Error())
	}
	where := "1" + a.report(input)
	group := ""
	lists, count := new(model.BorrowVisit).GroupPage(where,group, input.Page, input.Size)
	return resp.OK(c, map[string]interface{}{
		"count":count,
		"list":lists,
	})
}

type visitActionReq struct {
	BorrowId            	int			`json:"borrow_id"`
	Relationship            string			`json:"relationship"`
	ContactName             string		`json:"contact_name"`
	ContactPhone            string		`json:"contact_phone"`
	PromisedRepaymentAmount float64		`json:"promised_repayment_amount"`
	PromisedRepaymentTime   string		`json:"promised_repayment_time"`
	NextVisitTime           string		`json:"next_visit_time"`
	Tag                     string		`json:"tag"`
	Wish                    string		`json:"wish"`
	Remark                  string		`json:"remark"`
}
//UrgeAction 新增催记
func (a *visit)UrgeAction(c *fiber.Ctx) error {
	input := new(visitActionReq)
	if err := tools.ParseBody(c, input); err != nil {
		return resp.Err(c, 1, err.Error())
	}
	borrowData := new(model.Borrow)
	borrowData.One(fmt.Sprintf("id = %d", input.BorrowId))
	if borrowData.Id == 0{
		return resp.Err(c, 1, "借款不存在")
	}
	borrowVisitData := new(model.BorrowVisit)
	borrowVisitData.One(fmt.Sprintf("id = %d", input.BorrowId))
	if borrowVisitData.BorrowId == 0{
		return resp.Err(c, 1, "借款还没有分配")
	}
	if borrowVisitData.UrgeId == 0{
		return resp.Err(c, 1, "借款还没有分配给催收员")
	}
	//开始新增
	visitData := new(model.BorrowVisitDetail)
	visitData.UrgeCompanyId = borrowVisitData.UrgeCompanyId
	visitData.UrgeId = borrowVisitData.UrgeId
	//visitData.UrgeGroupId = borrowVisitData.Urge
	visitData.NextVisitTime = input.NextVisitTime
	visitData.Wish = input.Wish
	visitData.Tag = input.Tag
	visitData.UserId = borrowData.Uid
	visitData.Remark = input.Remark
	visitData.Relationship = input.Relationship
	visitData.ContactName = input.ContactName
	visitData.ContactPhone = input.ContactPhone
	visitData.PromisedRepaymentAmount = input.PromisedRepaymentAmount
	visitData.PromisedRepaymentTime = input.PromisedRepaymentTime
	visitData.Insert()
	borrowVisitData.Wish =  input.Wish
	borrowVisitData.Tag =  input.Tag
	borrowVisitData.UrgeLastTime =  tools.GetFormatTime()
	borrowVisitData.Update(fmt.Sprintf("borrow_id = %d", input.BorrowId))
	return resp.OK(c, "")
}

