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

type visitBaseReq struct{
	req.PageReq
	Name           	   string `json:"name" query:"name"`            //用户名
	Phone              string `json:"phone" query:"phone"`                //手机号
	BorrowId           int    `json:"borrow_id" query:"borrow_id"`            //订单编号
	Tag                string `json:"tag" query:"tag"`                  //标签
	ExpireDay        int    `json:"expire_day" query:"expire_day"`        //逾期天数
	Wish               int `json:"wish" query:"wish"`                 //还款意愿
	Status             string `json:"status" query:"status"`               //订单状态
	MchId              int    `json:"mch_id" query:"mch_id"`               //商户名称
	ProductId          int    `json:"product_id" query:"product_id"`           //产品名称
	PayStartTime       string `json:"pay_start_time" query:"pay_start_time"`       //还款开始时间
	PayEndTime         string `json:"pay_end_time" query:"pay_end_time"`         //还款结束时间
	VisitStartTime     string `json:"visit_start_time" query:"visit_start_time"`     //最近一次跟进时间
	VisitEndTime       string `json:"visit_end_time" query:"visit_end_time"`       //最近一次跟进时间
	RepaymentStartTime string `json:"repayment_start_time" query:"repayment_start_time"` //应还款时间
	RepaymentEndTime   string `json:"repayment_end_time" query:"repayment_end_time"`   //应还款时间
	RemindCompanyId    int    `json:"remind_company_id" query:"remind_company_id"`    //预提醒公司
	RemindGroupId      int    `json:"remind_group_id" query:"remind_group_id"`      //预提醒组
	RemindId           int    `json:"remind_id" query:"remind_id"`            //预提醒员
	UrgeCompanyId    int    `json:"urge_company_id" query:"urge_company_id"`    //催收公司
	UrgeGroupId      int    `json:"urge_group_id" query:"urge_group_id"`      //催收组
	UrgeId           int    `json:"urge_id" query:"urge_id"`            //催收员
}

func (a *visit) visitBorrowQuery(input *visitBaseReq) string {
	where := " "
	if len(input.Name) > 0 {
		where += fmt.Sprintf(" and user.aadhaar_name='%s'", input.Name)
	}
	if len(input.Phone) > 0 {
		where += fmt.Sprintf(" and user.phone='%s'", input.Phone)
	}
	if len(input.Status) > 0 {
		where += fmt.Sprintf(" and borrow.status='%s'", input.Status)
	}
	if len(input.Tag) > 0 {
		where += fmt.Sprintf(" and bv.tag='%s'", input.Tag)
	}
	if input.Wish > 0 {
		where += fmt.Sprintf(" and bv.wish='%s'", input.Wish)
	}
	if input.ExpireDay > 0 {
		where += fmt.Sprintf(" and borrow.expire_day='%d'", input.ExpireDay)
	}
	if input.BorrowId > 0 {
		where += fmt.Sprintf(" and bv.borrow_id='%d'", input.BorrowId)
	}
	if input.ProductId > 0 {
		where += fmt.Sprintf(" and bv.product_id='%d'", input.ProductId)
	}
	if input.MchId > 0 {
		where += fmt.Sprintf(" and bv.mch_id='%d'", input.MchId)
	}
	if input.RemindCompanyId > 0 {
		where += fmt.Sprintf(" and bv.remind_company_id='%d'", input.RemindCompanyId)
	}
	if input.RemindGroupId > 0 {
		where += fmt.Sprintf(" and bv.remind_group_id='%d'", input.RemindGroupId)
	}
	if input.RemindId > 0 {
		where += fmt.Sprintf(" and bv.remind_id='%d'", input.RemindId)
	}
	if input.UrgeCompanyId > 0 {
		where += fmt.Sprintf(" and bv.urge_company_id='%d'", input.UrgeCompanyId)
	}
	if input.UrgeGroupId > 0 {
		where += fmt.Sprintf(" and bv.urge_group_id='%d'", input.UrgeGroupId)
	}
	if input.UrgeId > 0 {
		where += fmt.Sprintf(" and bv.urge_id='%d'", input.UrgeId)
	}
	if len(input.PayStartTime) > 0 {
		where += fmt.Sprintf(" and borrow.complete_time >='%s'", input.PayStartTime)
	}
	if len(input.PayEndTime) > 0 {
		where += fmt.Sprintf(" and borrow.complete_time <'%s'", input.PayEndTime)
	}
	if len(input.VisitStartTime) > 0 {
		where += fmt.Sprintf(" and borrow.remind_last_time >='%s'", input.VisitStartTime)
	}
	if len(input.VisitEndTime) > 0 {
		where += fmt.Sprintf(" and borrow.remind_last_time <'%s'", input.VisitEndTime)
	}
	if len(input.RepaymentStartTime) > 0 {
		where += fmt.Sprintf(" and borrow.end_time >='%s'", input.RepaymentStartTime)
	}
	if len(input.RepaymentEndTime) > 0 {
		where += fmt.Sprintf(" and borrow.end_time <'%s'", input.RepaymentEndTime)
	}
	return where
}

//RemindBorrowAll 预提醒订单  预提醒中的订单 + 预提醒阶段就完成的订单
func (a *visit) RemindBorrowAll(c *fiber.Ctx) error {
	input := new(visitBaseReq)
	if err := tools.ParseBody(c, input); err != nil {
		return resp.Err(c, 1, err.Error())
	}
	where := "bv.remind_company_id > 0 and bv.urge_company_id = 0" + a.visitBorrowQuery(input) //在预提醒阶段完成的订单
	lists, count := new(model.BorrowVisit).RemindPage(where, input.Page, input.Size)
	return resp.OK(c, map[string]interface{}{
		"count": count,
		"list":  lists,
	})
}

//RemindBorrowing 预提醒中订单  预提醒中的订单
func (a *visit) RemindBorrowing(c *fiber.Ctx) error {
	input := new(visitBaseReq)
	if err := tools.ParseBody(c, input); err != nil {
		return resp.Err(c, 1, err.Error())
	}
	where := "bv.remind_company_id > 0 and bv.urge_company_id = 0 and borrow.status != 8" + a.visitBorrowQuery(input) //在预提醒阶段完成的订单
	lists, count := new(model.BorrowVisit).RemindPage(where, input.Page, input.Size)
	return resp.OK(c, map[string]interface{}{
		"count": count,
		"list":  lists,
	})
}

//RemindBorrowed 预提醒完成的订单  预提醒阶段就完成的订单
func (a *visit) RemindBorrowed(c *fiber.Ctx) error {
	input := new(visitBaseReq)
	if err := tools.ParseBody(c, input); err != nil {
		return resp.Err(c, 1, err.Error())
	}
	where := "bv.remind_company_id > 0 and bv.urge_company_id = 0 and borrow.status = 8" + a.visitBorrowQuery(input) //在预提醒阶段完成的订单
	lists, count := new(model.BorrowVisit).RemindPage(where, input.Page, input.Size)
	return resp.OK(c, map[string]interface{}{
		"count": count,
		"list":  lists,
	})
}

//RemindDetail 预提醒记录
func (a *visit) RemindDetail(c *fiber.Ctx) error {
	input := new(detailReq)
	if err := tools.ParseBody(c, input); err != nil {
		return resp.Err(c, 1, err.Error())
	}
	where := "bvd.type = 0 " + input.Query()
	lists, count := new(model.BorrowVisitDetail).RemindPage(where, input.Page, input.Size)
	return resp.OK(c, map[string]interface{}{
		"count": count,
		"list":  lists,
	})
}
//Details 催记记录
func (a *visit) Details(c *fiber.Ctx) error {
	input := new(detailReq)
	if err := tools.ParseBody(c, input); err != nil {
		return resp.Err(c, 1, err.Error())
	}
	where := "bvd.id > 0 " + input.Query()
	lists, count := new(model.BorrowVisitDetail).RemindPage(where, input.Page, input.Size)
	return resp.OK(c, map[string]interface{}{
		"count": count,
		"list":  lists,
	})
}

//UrgeBorrowAll 催收订单  催收中的订单 + 催收阶段完成的订单
func (a *visit) UrgeBorrowAll(c *fiber.Ctx) error {
	input := new(visitBaseReq)
	if err := tools.ParseBody(c, input); err != nil {
		return resp.Err(c, 1, err.Error())
	}
	where := "(bv.urge_company_id > 0)" + a.visitBorrowQuery(input)
	lists, count := new(model.BorrowVisit).UrgePage(where, input.Page, input.Size)
	return resp.OK(c, map[string]interface{}{
		"count": count,
		"list":  lists,
	})
}

//UrgeBorrowing 催收中的订单
func (a *visit) UrgeBorrowing(c *fiber.Ctx) error {
	input := new(visitBaseReq)
	if err := tools.ParseBody(c, input); err != nil {
		return resp.Err(c, 1, err.Error())
	}
	where := "(bv.urge_company_id > 0 and borrow.status != 9)" + a.visitBorrowQuery(input)
	lists, count := new(model.BorrowVisit).UrgePage(where, input.Page, input.Size)
	return resp.OK(c, map[string]interface{}{
		"count": count,
		"list":  lists,
	})
}

//UrgeBorrowed 催收阶段完成的订单
func (a *visit) UrgeBorrowed(c *fiber.Ctx) error {
	input := new(visitBaseReq)
	if err := tools.ParseBody(c, input); err != nil {
		return resp.Err(c, 1, err.Error())
	}
	where := "(bv.urge_company_id > 0 and borrow.status = 9)" + a.visitBorrowQuery(input)
	lists, count := new(model.BorrowVisit).UrgePage(where, input.Page, input.Size)
	return resp.OK(c, map[string]interface{}{
		"count": count,
		"list":  lists,
	})
}

type detailReq struct {
	req.PageReq
	Name          string `json:"name" query:"name"`                       //用户名
	Phone         string `json:"phone" query:"phone"`                     //债务人手机号
	BorrowId      int    `json:"borrow_id" query:"borrow_id"`             //订单编号
	UserId      	int    `json:"user_id" query:"user_id"`             //订单编号
	MchId      		int    `json:"mch_id" query:"mch_id"`             //商户id
	ProductId      int    `json:"product_id" query:"product_id"`             //产品id
	OverDueDays   int    `json:"over_due_days" query:"over_due_days"`     //到期天数 //待定
	ContactPhone  string `json:"contact_phone" query:"contact_phone"`     //联系人手机号
	Wish          string `json:"wish" query:"wish"`                       //还款意愿
	Relationship  string `json:"relationship" query:"relationship"`       //联系人与债主的关系
	Tag           string `json:"tag" query:"tag"`                         //标签
	RemindId        int    `json:"remind_id" query:"remind_id"`                 //催收人员
	RemindCompanyId int    `json:"remind_company_id" query:"remind_company_id"` //催收公司
	RemindGroupId   int    `json:"remind_group_id" query:"remind_group_id"`     //催收组
	ExpireDay   	string    `json:"expire_day" query:"expire_day"`     //逾期天数

	UrgeId        int    `json:"urge_id" query:"urge_id"`                 //催收人员
	UrgeCompanyId int    `json:"urge_company_id" query:"urge_company_id"` //催收公司
	UrgeGroupId   int    `json:"urge_group_id" query:"urge_group_id"`     //催收组

	StartTime     string `json:"start_time" query:"start_time"`           //开始时间
	EndTime       string `json:"end_time" query:"end_time"`               //结束时间
}

func (input *detailReq)Query() string {
	where := ""
	if len(input.Name) > 0 {
		where += fmt.Sprintf(" and user.name = '%s'", input.Name)
	}
	if len(input.Phone) > 0 {
		where += fmt.Sprintf(" and user.phone = '%s'", input.Phone)
	}
	if input.BorrowId > 0 {
		where += fmt.Sprintf(" and bvd.borrow_id = '%d'", input.BorrowId)
	}
	if input.UserId > 0 {
		where += fmt.Sprintf(" and bvd.user_id = '%d'", input.UserId)
	}
	if len(input.ContactPhone) > 0 {
		where += fmt.Sprintf(" and bvd.contact_phone = '%s'", input.ContactPhone)
	}
	if len(input.Wish) > 0 {
		where += fmt.Sprintf(" and bvd.wish = '%s'", input.Wish)
	}
	if len(input.Relationship) > 0 {
		where += fmt.Sprintf(" and bvd.relationship = '%s'", input.Relationship)
	}
	if len(input.Tag) > 0 {
		where += fmt.Sprintf(" and bvd.tag = '%s'", input.Tag)
	}
	if input.RemindId > 0 {
		where += fmt.Sprintf(" and bvd.urge_id = '%d'", input.RemindId)
	}
	if input.RemindCompanyId > 0 {
		where += fmt.Sprintf(" and bvd.urge_company_id = '%d'", input.RemindCompanyId)
	}
	if input.RemindGroupId > 0 {
		where += fmt.Sprintf(" and bvd.urge_group_id = '%d'", input.RemindGroupId)
	}
	if input.UrgeId > 0 {
		where += fmt.Sprintf(" and bvd.urge_id = '%d'", input.UrgeId)
	}
	if input.UrgeCompanyId > 0 {
		where += fmt.Sprintf(" and bvd.urge_company_id = '%d'", input.UrgeCompanyId)
	}
	if input.UrgeGroupId > 0 {
		where += fmt.Sprintf(" and bvd.urge_group_id = '%d'", input.UrgeGroupId)
	}
	if len(input.StartTime) > 0 {
		where += fmt.Sprintf(" and bvd.create_time >= '%s'", input.StartTime)
	}
	if len(input.EndTime) > 0 {
		where += fmt.Sprintf(" and bvd.create_time <'%s'", input.EndTime)
	}
	if len(input.ExpireDay) > 0 {
		where += fmt.Sprintf(" and borrow.expire_day ='%s'", input.ExpireDay)
	}
	return where
}

//UrgeDetail 催收记录
func (a *visit) UrgeDetail(c *fiber.Ctx) error {
	input := new(detailReq)
	if err := tools.ParseBody(c, input); err != nil {
		return resp.Err(c, 1, err.Error())
	}
	where := "bvd.type = 1" + input.Query()
	lists, count := new(model.BorrowVisitDetail).UrgePage(where, input.Page, input.Size)
	return resp.OK(c, map[string]interface{}{
		"count": count,
		"list":  lists,
	})
}

type urgeReportReq struct {
	req.PageReq
	StartTime     string
	EndTime       string
	MchId         int
	Type          int
	UrgeCompanyId int
	UrgeGroupId   int
	UrgeId        int
	Group         []string
}

func (a *visit) report(input *urgeReportReq) string {
	where := ""
	if input.MchId > 0 {
		where += fmt.Sprintf(" and mch_id = %d", input.MchId)
	}
	if input.UrgeCompanyId > 0 {
		where += fmt.Sprintf(" and urge_company_id = %d", input.UrgeCompanyId)
	}
	if len(input.StartTime) > 0 {
		where += fmt.Sprintf(" and urge_company_id = %d", input.UrgeCompanyId)
	}
	return where
}

//UrgeReport 催收业绩
func (a *visit) UrgeReport(c *fiber.Ctx) error {
	input := new(urgeReportReq)
	if err := tools.ParseBody(c, input); err != nil {
		return resp.Err(c, 1, err.Error())
	}
	where := "1" + a.report(input)
	group := ""
	lists, count := new(model.BorrowVisit).GroupPage(where, group, input.Page, input.Size)
	return resp.OK(c, map[string]interface{}{
		"count": count,
		"list":  lists,
	})
}

type visitActionReq struct {
	BorrowId                int     `json:"borrow_id"`
	Relationship            int  `json:"relationship"`
	ContactName             string  `json:"contact_name"`
	ContactPhone            string  `json:"contact_phone"`
	PromisedRepaymentAmount string  `json:"promised_repayment_amount"`
	PromisedRepaymentTime   string  `json:"promised_repayment_time"`
	NextVisitTime           string  `json:"next_visit_time"`
	Tag                     int  `json:"tag"`
	Wish                    int  `json:"wish"`
	Remark                  string  `json:"remark"`
}

//UrgeAction 新增催记
func (a *visit) UrgeAction(c *fiber.Ctx) error {
	input := new(visitActionReq)
	if err := tools.ParseBody(c, input); err != nil {
		return resp.Err(c, 1, err.Error())
	}
	borrowData := new(model.Borrow)
	borrowData.One(fmt.Sprintf("id = %d", input.BorrowId))
	if borrowData.Id == 0 {
		return resp.Err(c, 1, "借款不存在")
	}
	borrowVisitData := new(model.BorrowVisit)
	borrowVisitData.One(fmt.Sprintf("borrow_id = %d", input.BorrowId))
	if borrowVisitData.BorrowId == 0 {
		return resp.Err(c, 1, "借款还没有分配")
	}
	if borrowVisitData.UrgeId == 0 {
		return resp.Err(c, 1, "借款还没有分配给催收员")
	}
	//开始新增
	t := tools.GetFormatTime()
	visitData := new(model.BorrowVisitDetail)
	visitData.Type = 1
	visitData.MchId = borrowData.MchId
	visitData.ProductId = borrowData.ProductId
	visitData.BorrowId = borrowData.Id
	visitData.UserId = borrowData.Uid
	visitData.UrgeCompanyId = borrowVisitData.RemindCompanyId
	visitData.UrgeId = borrowVisitData.RemindId
	visitData.UrgeGroupId = borrowVisitData.RemindGroupId
	visitData.NextVisitTime = input.NextVisitTime
	visitData.Wish = input.Wish
	visitData.Tag = input.Tag
	visitData.Remark = input.Remark
	visitData.Relationship = input.Relationship
	visitData.ContactName = input.ContactName
	visitData.ContactPhone = input.ContactPhone
	visitData.CreateTime = t
	visitData.PromisedRepaymentAmount = input.PromisedRepaymentAmount
	visitData.PromisedRepaymentTime = input.PromisedRepaymentTime
	visitData.Insert()
	borrowVisitData.Wish = input.Wish
	borrowVisitData.Tag = input.Tag
	borrowVisitData.UrgeLastTime = t
	borrowVisitData.Update(fmt.Sprintf("borrow_id = %d", input.BorrowId))
	return resp.OK(c, "")
}

//RemindAction 新增预提醒
func (a *visit) RemindAction(c *fiber.Ctx) error {
	input := new(visitActionReq)
	if err := tools.ParseBody(c, input); err != nil {
		return resp.Err(c, 1, err.Error())
	}
	borrowData := new(model.Borrow)
	borrowData.One(fmt.Sprintf("id = %d", input.BorrowId))
	if borrowData.Id == 0 {
		return resp.Err(c, 1, "借款不存在")
	}
	borrowVisitData := new(model.BorrowVisit)
	borrowVisitData.One(fmt.Sprintf("borrow_id = %d", input.BorrowId))
	if borrowVisitData.BorrowId == 0 {
		return resp.Err(c, 1, "借款还没有分配")
	}
	if borrowVisitData.RemindId == 0 {
		return resp.Err(c, 1, "借款还没有分配给提醒员")
	}
	//开始新增
	t := tools.GetFormatTime()
	visitData := new(model.BorrowVisitDetail)
	visitData.Type = 0
	visitData.MchId = borrowData.MchId
	visitData.ProductId = borrowData.ProductId
	visitData.BorrowId = borrowData.Id
	visitData.UserId = borrowData.Uid
	visitData.UrgeCompanyId = borrowVisitData.UrgeCompanyId
	visitData.UrgeId = borrowVisitData.UrgeId
	visitData.UrgeGroupId = borrowVisitData.UrgeGroupId
	visitData.NextVisitTime = input.NextVisitTime
	visitData.Wish = input.Wish
	visitData.Tag = input.Tag
	visitData.Remark = input.Remark
	visitData.Relationship = input.Relationship
	visitData.ContactName = input.ContactName
	visitData.ContactPhone = input.ContactPhone
	visitData.CreateTime = t
	visitData.PromisedRepaymentAmount = input.PromisedRepaymentAmount
	visitData.PromisedRepaymentTime = input.PromisedRepaymentTime
	visitData.Insert()
	borrowVisitData.Wish = input.Wish
	borrowVisitData.Tag = input.Tag
	borrowVisitData.RemindLastTime = t
	borrowVisitData.Update(fmt.Sprintf("borrow_id = %d", input.BorrowId))
	return resp.OK(c, "")
}

type assignReq struct {
	BorrowId	int		`json:"borrow_id" query:"borrow_id""`
	AdminId	int		`json:"admin_id" query:"admin_id""`
}
func (a *visit) RemindAssign(c *fiber.Ctx) error {
	input := new(assignReq)
	if err := tools.ParseBody(c, input); err != nil {
		return resp.Err(c, 1, err.Error())
	}
	borrowData := new(model.Borrow)
	borrowData.One(fmt.Sprintf("id = %d", input.BorrowId))
	if borrowData.Id == 0 {
		return resp.Err(c, 1, "借款不存在")
	}
	//判断admin_id 跟商户匹配
	adminData := new(model.Admin)
	adminData.One(fmt.Sprintf("id = %d", input.AdminId))
	if adminData.RoleId != 7{
		return resp.Err(c, 1, "请分配给提醒专员")
	}
	if adminData.MchId != borrowData.MchId {
		return resp.Err(c, 1, "提醒专员跟商户不匹配")
	}
	visitData := new(model.BorrowVisit)
	visitData.One(fmt.Sprintf("borrow_id = %d", input.BorrowId))
	if visitData.BorrowId == 0{
		return resp.Err(c, 1, "没有到分配的时间")
	}
	visitData.RemindId = input.AdminId
	visitData.RemindAssignTime = tools.GetFormatTime()
	visitData.Update(fmt.Sprintf("borrow_id = %d", input.BorrowId))

	return resp.OK(c, "")
}

func (a *visit) UrgeAssign(c *fiber.Ctx) error {
	input := new(assignReq)
	if err := tools.ParseBody(c, input); err != nil {
		return resp.Err(c, 1, err.Error())
	}
	borrowData := new(model.Borrow)
	borrowData.One(fmt.Sprintf("id = %d", input.BorrowId))
	if borrowData.Id == 0 {
		return resp.Err(c, 1, "借款不存在")
	}
	//判断admin_id 跟商户匹配
	adminData := new(model.Admin)
	adminData.One(fmt.Sprintf("id = %d", input.AdminId))
	if adminData.RoleId != 3{
		return resp.Err(c, 1, "请分配给催收专员")
	}
	if adminData.MchId != borrowData.MchId {
		return resp.Err(c, 1, "催收专员跟商户不匹配")
	}
	visitData := new(model.BorrowVisit)
	visitData.One(fmt.Sprintf("borrow_id = %d", input.BorrowId))
	if visitData.BorrowId == 0{
		return resp.Err(c, 1, "没有到分配的时间")
	}
	visitData.UrgeId = input.AdminId
	visitData.UrgeAssignTime = tools.GetFormatTime()
	visitData.Update(fmt.Sprintf("borrow_id = %d", input.BorrowId))
	return resp.OK(c, "")
}

type utrCreateReq struct {
	BorrowId int	`json:"borrow_id"`
	Path    string	`json:"path"`
	UtrCode string	`json:"utr_code"`
	Remark  string	`json:"remark"`
}
//UtrCreate Utr提交 创建
func (a *visit) UtrCreate(c *fiber.Ctx) error  {
	input := new(utrCreateReq)
	if err := tools.ParseBody(c, input); err != nil {
		return resp.Err(c, 1, err.Error())
	}
	//验证borrow
	borrowData := new(model.Borrow)
	borrowData.One(fmt.Sprintf("id = %d", input.BorrowId))
	if borrowData.Id == 0{
		return resp.Err(c, 1, "没有找到记录")
	}

	utrData := new(model.BorrowUtr)
	utrData.BorrowId = borrowData.Id
	utrData.UrgeId = c.Locals("adminId").(int)
	utrData.MchId = borrowData.MchId
	utrData.ProductId = borrowData.ProductId
	utrData.UserId = borrowData.Uid
	utrData.CreateTime = tools.GetFormatTime()
	utrData.UtrPath = input.Path
	utrData.Status = 1 //等待处理
	utrData.Type = 1 //后台催收员提交
	utrData.UtrCode = input.UtrCode
	utrData.Remark = input.Remark
	utrData.Remark = input.Remark
	utrData.Insert()
	return resp.OK(c, "")
}

type utrExamineReq struct {
	req.IdReq
	Status    int	`json:"status" query:"status"`
	RejectReason    string	`json:"reject_reason" query:"reject_reason"`
}
//UtrExamine Utr提交 审核
func (a *visit) UtrExamine(c *fiber.Ctx) error  {
	input := new(utrExamineReq)
	if err := tools.ParseBody(c, input); err != nil {
		return resp.Err(c, 1, err.Error())
	}
	//验证borrow
	utrData := new(model.BorrowUtr)
	utrData.One(fmt.Sprintf("id = %d", input.Id))
	if utrData.Id == 0{
		return resp.Err(c, 1, "没有找到记录")
	}
	//status 状态
	if input.Status > 4 || input.Status < 1{
		return resp.Err(c, 1, "状态不正确")
	}
	utrData.Status = input.Status
	utrData.RejectReason = input.RejectReason
	utrData.Update(fmt.Sprintf("id = %d", input.Id))
	if input.Status == 3{
		//已到账
		//borrowData := new(model.Borrow)
		//borrowData.One(fmt.Sprintf("id = %d", utrData.BorrowId))
	}
	return resp.OK(c, "")
}

type statUrgeReq struct {
	req.PageReq
	StartTime string	`json:"start_time"`
	EndTime string	`json:"end_time"`
	MchId int	`json:"mch_id"`
	Type string	`json:"type"`
	UrgeCompanyId int	`json:"urge_company_id"`
	UrgeGroupId int	`json:"urge_group_id"`
	UrgeId int	`json:"urge_id"`
	Group string	`json:"group"`
}
//StatUrge 催收业绩
func (a *visit)StatUrge(c *fiber.Ctx) error {
	input := new(statUrgeReq)
	if err := tools.ParseBody(c, input); err != nil {
		return resp.Err(c, 1, err.Error())
	}
	where := "su.id > 0"
	if len(input.StartTime) > 0{
		where += fmt.Sprintf(" and st.stat_time >= '%s'", input.StartTime)
	}
	if len(input.EndTime) > 0{
		where += fmt.Sprintf(" and st.end_time <= '%s'", input.EndTime)
	}
	if input.MchId > 0{
		where += fmt.Sprintf(" and st.mch_id = '%d'", input.MchId)
	}
	if len(input.Type)  > 0{
		where += fmt.Sprintf(" and st.type = '%s'", input.Type)
	}
	if input.UrgeCompanyId > 0{
		where += fmt.Sprintf(" and st.urge_company_id = '%d'", input.UrgeCompanyId)
	}
	if input.UrgeGroupId > 0{
		where += fmt.Sprintf(" and st.urge_group_id = '%d'", input.UrgeGroupId)
	}
	if input.UrgeId > 0{
		where += fmt.Sprintf(" and st.urge_id = '%d'", input.UrgeId)
	}
	statData := new(model.StatUrge)
	lists, count := statData.GroupPage(where, input.Group, input.Page, input.Size)
	return resp.OK(c, map[string]interface{}{
		"count": count,
		"list":  lists,
	})
}