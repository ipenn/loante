package model

import (
	"fmt"
	"github.com/uptrace/bun"
	"loante/global"
	"loante/service/sms"
	"strings"
)

type SmsTemplate struct {
	bun.BaseModel     `bun:"table:sms_template,alias:st"`
	Id                int    `json:"id"`
	CompanyId         int    `json:"company_id"`
	SmsType           int    `json:"sms_type"`
	TemplateId        string    `json:"template_id"`
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

func (a *SmsTemplate)Send(Type int, phone string, ps []string) bool {
	lists, _ := a.Page(fmt.Sprintf("sms_type = %d", Type),1,1000)
	var args []interface{}
	for _, item := range ps{
		args = append(args, item)
	}
	for _, item := range lists{
		smsModel := sms.SelectSms(item.CompanyId)
		content := item.Content
		//if len(ps)> 0{
		//	if ps[0] != interface{}(nil){
		//		content = tools.Sprintf(item.Content, ps)
		//	}
		//}
		content = fmt.Sprintf(content,args...)
		content = strings.ReplaceAll(content, "%!(EXTRA []interface {}=[])","")
		ret, err := (*smsModel).Send(phone, content, Type)
		if ret{
			return true
		}
		global.Log.Error("发送短信失败%v", err.Error())
	}
	return false
}