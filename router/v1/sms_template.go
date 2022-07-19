package v1

import (
	"fmt"
	"github.com/gofiber/fiber/v2"
	"loante/service/model"
	"loante/service/req"
	"loante/service/resp"
	"loante/tools"
)

type smsTemplate struct{}

func NewSmsTemplate() *smsTemplate {
	return new(smsTemplate)
}

type smsTemplateList struct {
	req.PageReq
	CompanyId string `query:"companyId" json:"company_id"`
	SmsType   string `query:"smsType" json:"sms_type"`
}

func (a *smsTemplate) SmsTemplate(c *fiber.Ctx) error {
	input := new(smsTemplateList)
	if err := tools.ParseBody(c, input); err != nil {
		return resp.Err(c, 1, err.Error())
	}
	where := "id>0"
	if input.CompanyId != "" {
		where += " and company_id=" + input.CompanyId
	}
	if input.SmsType != "" {
		where += " and sms_type=" + input.SmsType
	}
	lists, count := new(model.SmsTemplate).Page(where, input.Page, input.Size)
	return resp.OK(c, map[string]interface{}{
		"count": count,
		"list":  lists,
	})
}

type smsTemplateCreateOrUpdate struct {
	ID                string `json:"id"`
	CompanyId         string `json:"companyId"`
	SmsType           string `json:"smsType"`
	TemplateId        string `json:"templateId"`
	Content           string `json:"content"`
	Description       string `json:"description"`
	SenderId          string `json:"senderId"`
	PrincipalEntityId string `json:"principalEntityId"`
}

func (a *smsTemplate) SmsTemplateCreateOrUpdate(c *fiber.Ctx) error {
	input := new(smsTemplateCreateOrUpdate)
	if err := tools.ParseBody(c, input); err != nil {
		return resp.Err(c, 1, err.Error())
	}
	s := new(model.SmsTemplate)
	s.CompanyId = tools.ToInt(input.CompanyId)
	s.SmsType = tools.ToInt(input.SmsType)
	s.TemplateId = tools.ToInt(input.TemplateId)
	s.Content = input.Content
	s.Description = input.Description
	s.SenderId = input.SenderId
	s.PrincipalEntityId = input.PrincipalEntityId
	if tools.ToInt(input.ID) == 0 {
		s.Insert()
	} else {
		s.Id = tools.ToInt(input.ID)
		s.Update(fmt.Sprintf("id=%d", input.ID))
	}
	return resp.OK(c, "")
}