package v1

import (
	"fmt"
	"github.com/gofiber/fiber/v2"
	"loante/service/model"
	"loante/service/req"
	"loante/service/resp"
	"loante/tools"
)

type utm_source struct {

}
func NewUtmSource() *utm_source {
	return new(utm_source)
}

type utmSourceReq struct {
	req.PageReq
	KeyWords string	`json:"key_words" query:"key_words"`
	Name string	`json:"name" query:"name"`
	NeedReview string	`json:"need_review" query:"need_review"`
	Status string	`json:"status" query:"status"`
}
//Lists 渠道
func (a *utm_source)Lists(c *fiber.Ctx) error {
	input := new(utmSourceReq)
	if err := tools.ParseBody(c, input); err != nil {
		return resp.Err(c, 1, err.Error())
	}
	where := fmt.Sprintf("id > 0")
	if len(input.Name) > 0{
		where += fmt.Sprintf(" and name='%s'", input.Name)
	}
	if len(input.KeyWords) > 0{
		where += fmt.Sprintf(" and keywords='%s'", input.KeyWords)
	}
	if len(input.Status) > 0{
		where += fmt.Sprintf(" and status='%s'", input.Status)
	}
	if len(input.NeedReview) > 0{
		where += fmt.Sprintf(" and is_need_review='%s'", input.NeedReview)
	}
	lists, count := new(model.ReferrerConfig).Page(where, input.Page, input.Size)
	return resp.OK(c, map[string]interface{}{
		"count":count,
		"list":lists,
	})
}

type utmSourceCreateReq struct {
	Id int	`json:"id"`
	KeyWords string	`json:"key_words"`
	Name string	`json:"name"`
	AppToken string	`json:"app_token"`
	Remark string	`json:"remark"`
}
func (a *utm_source)Create(c *fiber.Ctx) error {
	input := new(utmSourceCreateReq)
	if err := tools.ParseBody(c, input); err != nil {
		return resp.Err(c, 1, err.Error())
	}
	data := new(model.ReferrerConfig)
	if input.Id > 0{
		data.One(fmt.Sprintf("id = %d", input.Id))
		if data.Id == 0{
			return resp.Err(c, 1, "未找到数据")
		}
	}else{
		data.Status = 1
		data.IsRejectApply = 0
	}
	data.Name = input.Name
	data.Keyworks = input.KeyWords
	data.Remark = input.Remark
	data.AppToken = input.AppToken
	if data.Id > 0{
		data.Update(fmt.Sprintf("id = %d", input.Id))
	}else{
		data.Insert()
	}
	return resp.OK(c, "")
}

type utmSourceModifyReq struct {
	req.IdReq
	Key 	string	`json:"key"`
	Value 	string	`json:"value"`
}
func (a *utm_source)Modify(c *fiber.Ctx) error {
	input := new(utmSourceModifyReq)
	if err := tools.ParseBody(c, input); err != nil {
		return resp.Err(c, 1, err.Error())
	}
	data := new(model.ReferrerConfig)
	data.One(fmt.Sprintf("id = %d", input.Id))
	if data.Id == 0{
		return resp.Err(c, 1, "未找到数据")
	}
	err := data.Set(input.Key, input.Value, fmt.Sprintf("id = %d", input.Id))
	if err != nil{
		return resp.Err(c, 1, "保存失败")
	}
	return resp.OK(c, "")
}

type utmRiskReq struct {
	req.PageReq
	RiskModel string `json:"risk_model" query:"risk_model"`
}
func (a *utm_source)RiskConfig(c *fiber.Ctx) error {
	input := new(utmRiskReq)
	if err := tools.ParseBody(c, input); err != nil {
		return resp.Err(c, 1, err.Error())
	}
	where := fmt.Sprintf("id > 0")
	if len(input.RiskModel) > 0{
		where += fmt.Sprintf(" and risk_model='%s'", input.RiskModel)
	}
	lists, count := new(model.ReferrerRiskConfig).Page(where, input.Page, input.Size)
	return resp.OK(c, map[string]interface{}{
		"count":count,
		"list":lists,
	})
}

type utmRiskCreateReq struct {
	Id int	`json:"id"`
	ReferrerId int `json:"referrer_id"`
	StatCompay int `json:"stat_compay"`
	RiskModel int `json:"risk_model"`
	NewMinScore int `json:"new_min_score"`
	NewMaxScore	int	`json:"new_max_score"`
	OldJumpRisk	int	`json:"old_jump_risk"`
	OldMinScore	int	`json:"old_min_score"`
	OldMaxScore	int	`json:"old_max_score"`
	PlatformOldMinScore	int	`json:"platform_old_min_score"`
	PlatformOldMaxScore	int	`json:"platform_old_max_score"`
	Remark	string	`json:"remark"`
}
func (a *utm_source)RiskCreate(c *fiber.Ctx) error {
	input := new(utmRiskCreateReq)
	if err := tools.ParseBody(c, input); err != nil {
		return resp.Err(c, 1, err.Error())
	}
	data := new(model.ReferrerRiskConfig)
	if input.Id > 0{
		data.One(fmt.Sprintf("id = %d", input.Id))
		if data.Id == 0{
			return resp.Err(c, 1, "未找到数据")
		}
	}
	//检测 ReferrerId 是否存在
	refc := new(model.ReferrerConfig)
	refc.One(fmt.Sprintf("id = %d", input.ReferrerId))
	if refc.Id == 0{
		return resp.Err(c, 1, "未找到渠道")
	}
	data.ReferrerId = input.ReferrerId
	data.StatCompay = input.StatCompay
	data.RiskModel = input.RiskModel
	data.NewMinScore = input.NewMinScore
	data.NewMaxScore = input.NewMaxScore
	data.OldJumpRisk = input.OldJumpRisk
	data.OldMinScore = input.OldMinScore
	data.OldMaxScore = input.OldMaxScore
	data.PlatformOldMinScore = input.PlatformOldMinScore
	data.PlatformOldMaxScore = input.PlatformOldMaxScore
	data.Remark = input.Remark
	if data.Id > 0{
		data.Update(fmt.Sprintf("id = %d", input.Id))
	}else{
		data.Insert()
	}
	return resp.OK(c, "")
}