package v1

import (
	"fmt"
	"github.com/gofiber/fiber/v2"
	"loante/service/model"
	"loante/service/req"
	"loante/service/resp"
	"loante/tools"
)

type product struct{}

func NewProduct() *product {
	return new(product)
}

type productList struct {
	req.PageReq
	MchId         string `query:"mchId" json:"mch_id"`
	ProductName   string `query:"productName" json:"product_name"`
	IsAutoLending string `query:"isAutoLending" json:"is_auto_lending"`
	IsStopLending string `query:"isStopLending" json:"is_stop_lending"`
	Status        string `query:"status" json:"status"`
}

func (a *product) Product(c *fiber.Ctx) error {
	input := new(productList)
	if err := tools.ParseBody(c, input); err != nil {
		return resp.Err(c, 1, err.Error())
	}
	where := "p.id>0"
	if input.MchId != "" {
		where += " and m.id=" + input.MchId
	}
	if input.ProductName != "" {
		where += " and p.product_name like '%" + input.ProductName + "%'"
	}
	if input.IsAutoLending != "" {
		where += " and p.is_auto_lending =" + input.IsAutoLending
	}
	if input.IsStopLending != "" {
		where += " and p.is_stop_lending =" + input.IsStopLending
	}
	if input.Status != "" {
		where += " and p.status =" + input.Status
	}
	lists, count := new(model.Product).Page(where, input.Page, input.Size)
	return resp.OK(c, map[string]interface{}{
		"count": count,
		"list":  lists,
	})
}

type productCreteOrUpdate struct {
	Id                  string `json:"id"`
	ProductName         string `json:"productName"`
	IconPath            string `json:"iconPath"`
	MchId               string `json:"mchId"`
	DayMaxApply         string `json:"dayMaxApply"`
	MaxApply            string `json:"maxApply"`
	DayApplyPass        string `json:"dayApplyPass"`
	StartAmount         string `json:"startAmount"`
	TodayApplyCount     string `json:"todayApplyCount"`
	TotalApplyCount     string `json:"totalApplyCount"`
	ApplyStartTime      string `json:"applyStartTime"`
	ApplyEndTime        string `json:"applyEndTime"`
	TotalMaxApplyCount  string `json:"totalMaxApplyCount"`
	UpTime              string `json:"upTime"`
	DownTime            string `json:"downTime"`
	IsAutoLending       string `json:"isAutoLending"`
	IsRejectNew         string `json:"isRejectNew"`
	IsRejectOld         string `json:"isRejectOld"`
	IsStopLending       string `json:"isStopLending"`
	Status              string `json:"status"`
	CreateTime          string `json:"createTime"`
	Description         string `json:"description"`
	RateNormalInterest  string `json:"rateNormalInterest"`
	RateOverdueInterest string `json:"rateOverdueInterest"`
	RateService         string `json:"rateService"`
	RateTax             string `json:"rateTax"`
}

func (a *product) ProductCreateOrUpdate(c *fiber.Ctx) error {
	input := new(productCreteOrUpdate)
	if err := tools.ParseBody(c, input); err != nil {
		return resp.Err(c, 1, err.Error())
	}
	p := new(model.Product)
	p.ProductName = input.ProductName
	p.IconPath = input.IconPath
	p.MchId = tools.ToInt(input.MchId)
	p.DayMaxApply = tools.ToInt(input.DayMaxApply)
	p.MaxApply = tools.ToInt(input.MaxApply)
	p.DayApplyPass = tools.ToInt(input.DayApplyPass)
	p.StartAmount = tools.ToInt(input.StartAmount)
	p.TodayApplyCount = tools.ToInt(input.TodayApplyCount)
	//p.TotalApplyCount = tools.ToInt(input.TotalApplyCount)
	p.ApplyStartTime = input.ApplyStartTime
	p.ApplyEndTime = input.ApplyEndTime
	p.TotalMaxApplyCount = tools.ToInt(input.TotalMaxApplyCount)
	p.UpTime = input.UpTime
	p.DownTime = input.DownTime
	p.IsAutoLending = tools.ToInt(input.IsAutoLending)
	p.IsRejectNew = tools.ToInt(input.IsRejectNew)
	p.IsRejectOld = tools.ToInt(input.IsRejectOld)
	p.IsStopLending = tools.ToInt(input.IsStopLending)
	p.Status = tools.ToInt(input.Status)
	p.Description = input.Description
	p.RateNormalInterest = tools.ToFloat32(input.RateNormalInterest)
	p.RateOverdueInterest = tools.ToFloat32(input.RateOverdueInterest)
	p.RateService = tools.ToFloat32(input.RateService)
	p.RateTax = tools.ToFloat32(input.RateTax)
	if tools.ToInt(input.Id) == 0 {
		p.CreateTime = tools.GetFormatTime()
		p.Insert()
	} else {
		p2 := model.Product{}
		p2.One(fmt.Sprintf("id=%d", tools.ToInt(input.Id)))
		if p2.Id == 0 {
			return resp.Err(c, 1, "没有产品")
		}
		p2.Id = tools.ToInt(input.Id)
		p2.CreateTime = p.CreateTime
		p2.Update(fmt.Sprintf("id=%d", tools.ToInt(input.Id)))
	}
	return resp.OK(c, "")
}

<<<<<<< HEAD
<<<<<<< HEAD
<<<<<<< HEAD
type productPreceptReq struct {
	req.PageReq
	ProductId int `json:"product_id" query:"product_id"`
	Status int `json:"status" query:"status"`
}
//ProductPrecept 提额列表
func (a *product) ProductPrecept(c *fiber.Ctx) error  {
	input := new(productPreceptReq)
	if err := tools.ParseBody(c, input); err != nil {
		return resp.Err(c, 1, err.Error())
	}
	where := "ppt.id>0"
	if input.ProductId > 0 {
		where += fmt.Sprintf(" and ppt.product_id= %d", input.ProductId)
	}
	if input.Status > 0 {
		where += fmt.Sprintf(" and ppt.status = %d", input.Status)
	}
	lists, count := new(model.ProductPrecept).Page(where, input.Page, input.Size)
	return resp.OK(c, map[string]interface{}{
		"count": count,
		"list":  lists,
	})
}

type productPreceptCreateReq struct {
	req.IdReq
	ProductId int `json:"product_id"`
	Amount float64 `json:"amount"`
	MinCount int `json:"min_count"`
	Status int `json:"status"`
}
//ProductPreceptCreate 提额规则添加或修改
func (a *product) ProductPreceptCreate(c *fiber.Ctx) error  {
	input := new(productPreceptCreateReq)
	if err := tools.ParseBody(c, input); err != nil {
		return resp.Err(c, 1, err.Error())
	}
	ppt := new(model.ProductPrecept)
	ppt.One(fmt.Sprintf("id = %d", input.Id))
	ppt.ProductId = input.ProductId
	ppt.Amount = input.Amount
	ppt.MinCount = input.MinCount
	ppt.Status = input.Status
	if input.Id == 0{
		ppt.Insert()
	}else{
		ppt.Update(fmt.Sprintf("id = %d", input.Id))
	}
	return resp.OK(c, "")
}
//ProductPreceptDel 提额规则删除
func (a *product) ProductPreceptDel(c *fiber.Ctx) error  {
	input := new(req.IdReq)
	if err := tools.ParseBody(c, input); err != nil {
		return resp.Err(c, 1, err.Error())
	}
	ppt := new(model.ProductPrecept)
	ppt.Del(fmt.Sprintf("id = %d", input.Id))
	return resp.OK(c, "")
}
=======
=======
>>>>>>> 5fa5f02c1373b226cd4ab46bcdfa3326f6ae89d0
=======
>>>>>>> 5fa5f02c1373b226cd4ab46bcdfa3326f6ae89d0
type productUpdateForMch struct {
	Id                  string `json:"id"`
	ProductName         string `json:"productName"`
	IconPath            string `json:"iconPath"`
	IsAutoLending       string `json:"isAutoLending"`
	IsStopLending       string `json:"isStopLending"`
	Status              string `json:"status"`
	RateNormalInterest  string `json:"rateNormalInterest"`
	RateOverdueInterest string `json:"rateOverdueInterest"`
	RateService         string `json:"rateService"`
	RateTax             string `json:"rateTax"`
}

func (a *product) ProductUpdateForMch(c *fiber.Ctx) error {
	input := new(productUpdateForMch)
	if err := tools.ParseBody(c, input); err != nil {
		return resp.Err(c, 1, err.Error())
	}
	p := new(model.Product)
	p.One(fmt.Sprintf("id=%d", tools.ToInt(input.Id)))
	if p.Id == 0 {
		return resp.Err(c, 1, "id不能为0")
	}
	p.ProductName = input.ProductName
	p.IconPath = input.IconPath
	p.IsAutoLending = tools.ToInt(input.IsAutoLending)
	p.IsStopLending = tools.ToInt(input.IsStopLending)
	p.Status = tools.ToInt(input.Status)
	p.RateNormalInterest = tools.ToFloat32(input.RateNormalInterest)
	p.RateOverdueInterest = tools.ToFloat32(input.RateOverdueInterest)
	p.RateService = tools.ToFloat32(input.RateService)
	p.RateTax = tools.ToFloat32(input.RateTax)

	p.Update(fmt.Sprintf("id=%d", p.Id))
	return resp.OK(c, "")
}
<<<<<<< HEAD
<<<<<<< HEAD
>>>>>>> 5fa5f02c1373b226cd4ab46bcdfa3326f6ae89d0
=======
>>>>>>> 5fa5f02c1373b226cd4ab46bcdfa3326f6ae89d0
=======
>>>>>>> 5fa5f02c1373b226cd4ab46bcdfa3326f6ae89d0
