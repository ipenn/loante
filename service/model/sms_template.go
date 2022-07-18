package model

import (
	"fmt"
	"github.com/uptrace/bun"
	"loante/global"
)

type SmsTemplate struct {
	bun.BaseModel     `bun:"table:sms_template,alias:st"`
	Id                int    `json:"id"`
	CompanyId         int    `json:"company_id"`
	SmsType           int    `json:"sms_type"`
	TemplateId        int    `json:"template_id"`
	Content           string `json:"content"`
	Description       string `json:"description"`
	SenderId          string `json:"sender_id"`
	PrincipalEntityId string `json:"principal_entity_id"`
}

func (a *SmsTemplate) Insert() {
	_, err := global.C.DB.NewInsert().Model(a).Returning("*").Exec(global.C.Ctx)
	if err != nil {
		global.Log.Error("%v err=%v", a, err.Error())
	}
}

func (a *SmsTemplate) Update(where string) {
	_, err := global.C.DB.NewUpdate().Model(a).Where(where).Exec(global.C.Ctx)
	if err != nil {
		global.Log.Error("%v err=%v", a, err.Error())
	}
}

func (a *SmsTemplate) Page(where string, page, limit int) ([]SmsTemplate, int) {
	var datas []SmsTemplate
	count, _ := global.C.DB.NewSelect().Model(&datas).Where(where).Order(fmt.Sprintf("id desc")).Offset((page - 1) * limit).Limit(limit).ScanAndCount(global.C.Ctx)
	return datas, count
}
