package v1

import (
	"fmt"
	"github.com/gofiber/fiber/v2"
	"loante/service/model"
	"loante/service/req"
	"loante/service/resp"
	"loante/tools"
)

type productDelayConfig struct{}

func NewProductDelayConfig() *productDelayConfig {
	return new(productDelayConfig)
}

type productDelayConfigList struct {
	req.PageReq
	MchId     string `query:"mchId" json:"mch_id "`
	ProductId string `query:"productId" json:"product_id"`
}

func (a *productDelayConfig) ProductDelayConfig(c *fiber.Ctx) error {
	input := new(productDelayConfigList)
	if err := tools.ParseBody(c, input); err != nil {
		return resp.Err(c, 1, err.Error())
	}
	where := "pdc.id>0"
	if input.MchId != "" {
		where += " and pdc.mch_id=" + input.MchId
	}
	if input.ProductId != "" {
		where += " and pdc.product_id =" + input.ProductId
	}
	lists, count := new(model.ProductDelayConfig).Page(where, input.Page, input.Size)
	return resp.OK(c, map[string]interface{}{
		"count": count,
		"list":  lists,
	})
}

type productDelayConfigCreateOrUpdate struct {
	Id          string `json:"id"`
	MchId       string `json:"mchId"`
	ProductId   string `json:"productId"`
	DelayDay    string `json:"delayDay"`
	DelayRate   string `json:"delayRate"`
	Status      string `json:"status"`
	IsShowDelay string `json:"isShowDelay"`
	CreateTime  string `json:"createTime"`
}

func (a *productDelayConfig) ProductDelayConfigCreateOrUpdate(c *fiber.Ctx) error {
	input := new(productDelayConfigCreateOrUpdate)
	if err := tools.ParseBody(c, input); err != nil {
		return resp.Err(c, 1, err.Error())
	}
	p := new(model.ProductDelayConfig)
	p.One(fmt.Sprintf("product_id=%d", tools.ToInt(input.ProductId)))
	p.MchId = tools.ToInt(input.MchId)
	p.ProductId = tools.ToInt(input.ProductId)
	p.DelayDay = tools.ToInt(input.DelayDay)
	p.DelayRate = tools.ToFloat64(input.DelayRate)
	p.Status = tools.ToInt(input.Status)
	p.IsShowDelay = tools.ToInt(input.IsShowDelay)
	p.CreateTime = input.CreateTime
	if tools.ToInt(input.Id) == 0 {
		if p.Id != 0 {
			return resp.Err(c, 1, "该产品展期规则已存在")
		}
		p.CreateTime = tools.GetFormatTime()
		p.Insert()
	} else {
		p.Id = tools.ToInt(input.Id)
		p.Update(fmt.Sprintf("id=%d", tools.ToInt(input.Id)))
	}
	return resp.OK(c, "")
}
