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

type merchant struct{}

func NewMerchant() *merchant {
	return new(merchant)
}

type merchantPageReq struct {
	req.PageReq
	Name      string `json:"name" query:"name"`
	EndTime   string `json:"end_time" query:"end_time"`
	StartTime string `json:"start_time" query:"start_time"`
}

func (a *merchant) Lists(c *fiber.Ctx) error {
	input := new(merchantPageReq)
	if err := tools.ParseBody(c, input); err != nil {
		return resp.Err(c, 1, err.Error())
	}
	where := "id > 0"
	if len(input.Name) > 0 {
		where += fmt.Sprintf(" and name='%s'", input.Name)
	}
	if len(input.StartTime) > 0 {
		where += fmt.Sprintf(" and create_time > '%s'", input.StartTime)
	}
	if len(input.EndTime) > 0 {
		where += fmt.Sprintf(" and create_time < '%s'", input.EndTime)
	}
	lists, count := new(model.Merchant).Page(where, input.Page, input.Size)
	return resp.OK(c, map[string]interface{}{
		"count": count,
		"list":  lists,
	})
}

type merchantCreateReq struct {
	Id            int     `json:"id"`
	Name          string  `json:"name"`
	Type          int     `json:"type"`
	Status        int     `json:"status"`
	Password      string  `json:"password"`
	CnyCredit     float64 `json:"cny_credit"`
	UsdCredit     float64 `json:"usd_credit"`
	ContactName   string  `json:"contact_name"`
	ContactMobile string  `json:"contact_mobile"`
	ContactEmail  string  `json:"contact_email"`
}

func (a *merchant) Create(c *fiber.Ctx) error {
	input := new(merchantCreateReq)
	if err := tools.ParseBody(c, input); err != nil {
		return resp.Err(c, 1, err.Error())
	}
	//检测商户名称
	mch := new(model.Merchant)
	mch.One(fmt.Sprintf("name = '%s'", input.Name))
	if mch.Id > 0 && mch.Id != input.Id {
		return resp.Err(c, 1, "商户名称已经存在")
	}
	if input.Id > 0 {
		mch.One(fmt.Sprintf("id = '%d'", input.Id))
	}
	mch.Type = input.Type
	mch.Name = input.Name
	mch.Status = input.Status
	mch.CnyCredit = input.CnyCredit
	mch.UsdCredit = input.UsdCredit
	mch.ContactName = input.ContactName
	mch.ContactMobile = input.ContactMobile
	mch.ContactEmail = input.ContactEmail
	mch.UpdateTime = tools.GetFormatTime()
	if input.Id == 0 {
		mch.Insert()
		//需要同步生成一个商户账号
		admin := new(model.Admin)
		admin.AdminName = input.Name
		admin.RoleId = 8
		admin.Status = 1
		admin.Password = input.Password
		admin.MchId = mch.Id
		admin.Email = mch.ContactEmail
		admin.Mobile = mch.ContactMobile
		admin.Insert()
	} else {
		mch.Update(fmt.Sprintf("id = %d", input.Id))
	}
	return resp.OK(c, "")
}

func (a *merchant) Modify(c *fiber.Ctx) error {
	input := new(req.ModifyReq)
	if err := tools.ParseBody(c, input); err != nil {
		return resp.Err(c, 1, err.Error())
	}
	mch := new(model.Merchant)
	mch.One(fmt.Sprintf("id = '%d'", input.Id))
	if mch.Id == 0 {
		return resp.Err(c, 1, "商户不存在")
	}
	mch.SetColumn(input.Key, input.Value, fmt.Sprintf("id = %d", input.Id))
	return resp.OK(c, "")
}

type merchantFundReq struct {
	req.PageReq
	MchId     int    `json:"mch_id" query:"mch_id"`
	Type      int `json:"type" query:"type"`
	StartTime string `json:"start_time" query:"start_time"`
	EndTime   string `json:"end_time" query:"end_time"`
}

func (a *merchant) Funds(c *fiber.Ctx) error {
	input := new(merchantFundReq)
	if err := tools.ParseBody(c, input); err != nil {
		return resp.Err(c, 1, err.Error())
	}
	where := "mf.id > 0"
	if input.MchId > 0 {
		where += fmt.Sprintf(" and mf.mch_id='%d'", input.MchId)
	}
	if input.Type > 0 {
		where += fmt.Sprintf(" and mf.type='%d'", input.Type)
	}
	if len(input.StartTime) > 0 {
		where += fmt.Sprintf(" and mf.create_time > '%s'", input.StartTime)
	}
	if len(input.EndTime) > 0 {
		where += fmt.Sprintf(" and mf.create_time < '%s'", input.EndTime)
	}
	lists, count := new(model.MerchantFund).Page(where, input.Page, input.Size)
	return resp.OK(c, map[string]interface{}{
		"count": count,
		"list":  lists,
	})
}

type fundCreateReq struct {
	MchId       int     `json:"mch_id"`
	InAccountNo string  `json:"in_account_no"`
	Currency    int     `json:"currency"`
	Rate        float64 `json:"rate"`
	Amount      float64 `json:"amount"`
	FundNo      string  `json:"fund_no"`
	Remark      string  `json:"remark"`
	Path        string  `json:"path"`
	Type        int  	`json:"type"` //充值类型   1="现金充值",2="现金退款",3="服务扣款",4="短信扣款",5="风控服务费扣款"
}

func (a *merchant) FundCreate(c *fiber.Ctx) error {
	input := new(fundCreateReq)
	if err := tools.ParseBody(c, input); err != nil {
		return resp.Err(c, 1, err.Error())
	}
	//查找商户
	mch := new(model.Merchant)
	mch.One(fmt.Sprintf("id = '%d'", input.MchId))
	if mch.Id == 0 {
		return resp.Err(c, 1, "商户不存在")
	}
	if input.Type == 1{
		input.Amount = math.Abs(input.Amount)
	}
	if input.Type == 2{
		input.Amount = math.Abs(input.Amount) * -1
	}
	fund := new(model.MerchantFund)
	fund.MchId = input.MchId
	fund.InAccountNo = input.InAccountNo
	fund.Currency = input.Currency
	fund.Rate = input.Rate
	fund.Amount = input.Amount
	fund.FundNo = input.FundNo
	fund.Remark = input.Remark
	fund.Path = input.Path
	fund.Type = input.Type
	fund.Insert()
	//充值到商户Balance
	if input.Currency == 1{ //充值人民币
		mch.CnyBalance += input.Amount
	}else{
		mch.CnyBalance += input.Amount
	}
	mch.Update(fmt.Sprintf("id = %d", input.MchId))
	return resp.OK(c, "")
}

//ServiceRule 进件计价规则
func (a *merchant) ServiceRule(c *fiber.Ctx) error {
	input := new(req.PageReq)
	if err := tools.ParseBody(c, input); err != nil {
		return resp.Err(c, 1, err.Error())
	}
	lists, count := new(model.ServiceFeeRule).Page("id > 0", input.Page, input.Size)
	return resp.OK(c, map[string]interface{}{
		"count": count,
		"list":  lists,
	})
}

type serviceRuleCreateReq struct {
	Id         int     `json:"id"`
	Name       string  `json:"name"`
	Price      float64 `json:"price"`
	StartCount int     `json:"start_count"`
	EndCount   int     `json:"end_count"`
}

//ServiceRuleCreate 进件计价规则创建
func (a *merchant) ServiceRuleCreate(c *fiber.Ctx) error {
	input := new(serviceRuleCreateReq)
	if err := tools.ParseBody(c, input); err != nil {
		return resp.Err(c, 1, err.Error())
	}
	rule := new(model.ServiceFeeRule)
	if input.Id > 0 {
		rule.One(fmt.Sprintf("id = %d", input.Id))
		if rule.Id == 0 {
			return resp.Err(c, 1, "规则不存在")
		}
	}
	rule.Name = input.Name
	rule.Price = input.Price
	rule.StartCount = input.StartCount
	rule.EndCount = input.EndCount
	if input.Id > 0 {
		rule.Update(fmt.Sprintf("id = %d", input.Id))
	} else {
		rule.Insert()
	}
	return resp.OK(c, "")
}

//ServiceRuleDel 进件计价规则删除
func (a *merchant) ServiceRuleDel(c *fiber.Ctx) error {
	input := new(req.IdReq)
	if err := tools.ParseBody(c, input); err != nil {
		return resp.Err(c, 1, err.Error())
	}
	rule := new(model.ServiceFeeRule)
	rule.One(fmt.Sprintf("id = %d", input.Id))
	if rule.Id == 0 {
		return resp.Err(c, 1, "规则不存在")
	}
	rule.Del(fmt.Sprintf("id = %d", input.Id))
	return resp.OK(c, "")
}

//ServicePrice 服务定价
func (a *merchant) ServicePrice(c *fiber.Ctx) error {
	input := new(req.PageReq)
	if err := tools.ParseBody(c, input); err != nil {
		return resp.Err(c, 1, err.Error())
	}
	lists, count := new(model.ServicePrice).Page("id > 0", input.Page, input.Size)
	return resp.OK(c, map[string]interface{}{
		"count": count,
		"list":  lists,
	})
}

type servicePriceCreateReq struct {
	Id          int     `json:"id"`
	ServiceType int     `json:"service_type"`
	DeductType  int     `json:"deduct_type"`
	Price       float64 `json:"price"`
}

//ServicePriceCreate 服务定价创建
func (a *merchant) ServicePriceCreate(c *fiber.Ctx) error {
	input := new(servicePriceCreateReq)
	if err := tools.ParseBody(c, input); err != nil {
		return resp.Err(c, 1, err.Error())
	}
	row := new(model.ServicePrice)
	if input.Id > 0 {
		row.One(fmt.Sprintf("id = %d", input.Id))
		if row.Id == 0 {
			return resp.Err(c, 1, "服务定价不存在")
		}
	}
	row.ServiceType = input.ServiceType
	row.Price = input.Price
	row.DeductType = input.DeductType
	if input.Id > 0 {
		row.Update(fmt.Sprintf("id = %d", input.Id))
	} else {
		row.Insert()
	}
	return resp.OK(c, "")
}

//ServicePriceDel 服务定价删除
func (a *merchant) ServicePriceDel(c *fiber.Ctx) error {
	input := new(req.IdReq)
	if err := tools.ParseBody(c, input); err != nil {
		return resp.Err(c, 1, err.Error())
	}
	rule := new(model.ServicePrice)
	rule.One(fmt.Sprintf("id = %d", input.Id))
	if rule.Id == 0 {
		return resp.Err(c, 1, "定价不存在")
	}
	rule.Del(fmt.Sprintf("id = %d", input.Id))
	return resp.OK(c, "")
}
