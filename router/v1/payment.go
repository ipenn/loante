package v1

import (
	"fmt"
	"github.com/gofiber/fiber/v2"
	"loante/service/model"
	"loante/service/req"
	"loante/service/resp"
	"loante/tools"
)

type payment struct {
}

func NewPayment() *payment {
	return new(payment)
}

func (a *payment)Lists(c *fiber.Ctx) error {
	input := new(req.PageReq)
	if err := tools.ParseBody(c, input); err != nil {
		return resp.Err(c, 1, err.Error())
	}
	lists, count := new(model.Payment).Page("id > 0", input.Page, input.Size)
	return resp.OK(c, map[string]interface{}{
		"count":count,
		"list":lists,
	})
}

type paymentCreateReq struct {
	Id               int    `json:"id"`
	LendingStartTime string `json:"lending_start_time"`
	LendingEndTime   string `json:"lending_end_time"`
	IsOpenOut        int    `json:"is_open_out"`
	IsOpenIn         int    `json:"is_open_in"`
	IsUtrQuery       int    `json:"is_utr_query"`
	IsUtrFill        int    `json:"is_utr_fill"`
}
func (a *payment)Modify(c *fiber.Ctx) error {
	input := new(paymentCreateReq)
	if err := tools.ParseBody(c, input); err != nil {
		return resp.Err(c, 1, err.Error())
	}
	row := new(model.Payment)
	row.One(fmt.Sprintf("id = %d", input.Id))
	if row.Id == 0{
		return resp.Err(c, 1, "支付通道不存在")
	}
	row.LendingStartTime = input.LendingStartTime
	row.LendingEndTime = input.LendingEndTime
	row.IsOpenOut = input.IsOpenOut
	row.IsOpenIn = input.IsOpenIn
	row.IsUtrQuery = input.IsUtrQuery
	row.IsUtrFill = input.IsUtrFill
	row.Update(fmt.Sprintf("id = %d", input.Id))

	return resp.OK(c, "")
}

func (a *payment)Set(c *fiber.Ctx) error {
	input := new(req.ModifyReq)
	if err := tools.ParseBody(c, input); err != nil {
		return resp.Err(c, 1, err.Error())
	}
	row := new(model.Payment)
	row.One(fmt.Sprintf("id = %d", input.Id))
	if row.Id == 0{
		return resp.Err(c, 1, "支付通道不存在")
	}
	row.SetColumn(input.Key, input.Value, fmt.Sprintf("id = %d", input.Id))
	return resp.OK(c, "")
}

type paymentConfigReq struct {
	req.PageReq
	MchId       int 	`query:"mch_id" json:"mch_id"`
	ProductId 	int 	`query:"product_id" json:"product_id"`
	PaymentId   int 	`query:"payment_id" json:"payment_id"`
}

func (a *payment)ConfigLists(c *fiber.Ctx) error {
	input := new(paymentConfigReq)
	if err := tools.ParseBody(c, input); err != nil {
		return resp.Err(c, 1, err.Error())
	}
	lists, count := new(model.ProductPayment).Page("id > 0", input.Page, input.Size)
	return resp.OK(c, map[string]interface{}{
		"count":count,
		"list":lists,
	})
}

type productPaymentCreateReq struct {
	Id int `json:"id"`
	MchId int `json:"mch_id"`
	ProductId int `json:"product_id"`
	PaymentId int `json:"payment_id"`
	Configuration string `json:"configuration"` //{"config":[{"merchantId":"121","merchantKey":"12132","desc":"xsaxsa"}],"use":1}
}
func (a *payment)ConfigCreate(c *fiber.Ctx) error {
	input := new(productPaymentCreateReq)
	if err := tools.ParseBody(c, input); err != nil {
		return resp.Err(c, 1, err.Error())
	}
	pp := new(model.ProductPayment)
	//判断是否已经存在
	if input.Id == 0{
		pp.One(fmt.Sprintf("mch_id = %d and product_id = %d and payment_id = %d", input.MchId, input.ProductId, input.PaymentId))
		if pp.Id > 0{
			return resp.Err(c, 1, "支付通道配置已经存在")
		}
	}else{
		pp.One(fmt.Sprintf("id = %d", input.Id))
		if pp.Id == 0{
			return resp.Err(c, 1, "支付通道配置不存在")
		}
	}
	pp.ProductId = input.ProductId
	pp.PaymentId = input.PaymentId
	pp.MchId = input.MchId
	pp.Configuration = input.Configuration
	if input.Id == 0 {
		pp.IsOpenIn = 1
		pp.IsOpenOut = 1
		pp.Insert()
	}else{
		pp.Update(fmt.Sprintf("id = %d", input.Id))
	}
	return resp.OK(c, "")
}

func (a *payment)ConfigDel(c *fiber.Ctx) error {
	input := new(req.IdReq)
	if err := tools.ParseBody(c, input); err != nil {
		return resp.Err(c, 1, err.Error())
	}
	pp := new(model.ProductPayment)
	pp.One(fmt.Sprintf("id = %d", input.Id))
	if pp.Id == 0{
		return resp.Err(c, 1, "没有找到数据")
	}
	pp.Del(fmt.Sprintf("id = %d", pp.Id))
	return resp.OK(c, "")
}

func (a *payment)ConfigSet(c *fiber.Ctx) error {
	input := new(req.ModifyReq)
	if err := tools.ParseBody(c, input); err != nil {
		return resp.Err(c, 1, err.Error())
	}
	row := new(model.ProductPayment)
	row.One(fmt.Sprintf("id = %d", input.Id))
	if row.Id == 0{
		return resp.Err(c, 1, "支付通道不存在")
	}
	row.SetColumn(input.Key, input.Value, fmt.Sprintf("id = %d", input.Id))
	return resp.OK(c, "")
}

type paymentDefaultReq struct {
	req.PageReq
	MchId       int 	`query:"mch_id" json:"mch_id"`
	ProductId 	int 	`query:"product_id" json:"product_id"`
}
func (a *payment)DefaultLists(c *fiber.Ctx) error {
	input := new(paymentDefaultReq)
	if err := tools.ParseBody(c, input); err != nil {
		return resp.Err(c, 1, err.Error())
	}
	lists, count := new(model.ProductPaymentDefault).Page("ppd.id > 0", input.Page, input.Size)
	return resp.OK(c, map[string]interface{}{
		"count":count,
		"list":lists,
	})
}

type paymentDefaultCreateReq struct {
	Id       int 	`json:"id"`
	MchId       int 	`json:"mch_id"`
	ProductId 	int 	`json:"product_id"`
	OutPaymentId 	int 	`json:"out_payment_id"`
	InPaymentId 	int 	`json:"in_payment_id"`
}

//DefaultCreate 创建选择默认的支付通道
func (a *payment)DefaultCreate(c *fiber.Ctx) error {
	input := new(paymentDefaultCreateReq)
	if err := tools.ParseBody(c, input); err != nil {
		return resp.Err(c, 1, err.Error())
	}
	defaultPP := new(model.ProductPaymentDefault)
	//判断是否已经存在
	if input.Id == 0{
		defaultPP.One(fmt.Sprintf("mch_id = %d and product_id = %d", input.MchId, input.ProductId))
		if defaultPP.Id > 0{
			return resp.Err(c, 1, "规则已经存在")
		}
	}else{
		defaultPP.One(fmt.Sprintf("id = %d", input.Id))
		if defaultPP.Id == 0{
			return resp.Err(c, 1, "规则不存在")
		}
	}
	defaultPP.MchId = input.MchId
	defaultPP.ProductId = input.ProductId
	defaultPP.OutPaymentId = input.OutPaymentId
	defaultPP.InPaymentId = input.InPaymentId
	if input.Id == 0{
		defaultPP.Insert()
	}else{
		defaultPP.Update(fmt.Sprintf("id = %d", input.Id))
	}
	return resp.OK(c, "")
}
//DefaultDel 默认支付通道选择 -> 删除
func (a *payment)DefaultDel(c *fiber.Ctx) error {
	input := new(req.IdReq)
	if err := tools.ParseBody(c, input); err != nil {
		return resp.Err(c, 1, err.Error())
	}
	defaultPP := new(model.ProductPaymentDefault)
	defaultPP.One(fmt.Sprintf("id = %d", input.Id))
	if defaultPP.Id == 0{
		return resp.Err(c, 1, "规则不存在")
	}
	defaultPP.Del(fmt.Sprintf("id = %d", input.Id))
	return resp.OK(c, "")
}