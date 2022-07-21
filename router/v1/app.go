package v1

import (
	"fmt"
	"github.com/gofiber/fiber/v2"
	"loante/global"
	"loante/service/model"
	"loante/service/resp"
	"loante/service/sms"
	"loante/tools"
)

type app struct {
	
}

func NewApp() *app {
	return new(app)
}

//type 是哪类型的短信
type smsSendReq struct {
	MchId   int	`query:"mch_id" json:"mch_id"`
	Type   int	`query:"type" json:"type"`
	Phone   string	`query:"phone" json:"phone"`
	Params   []string	`query:"params" json:"params"`
}

func mySprintf(tpl string, arg []interface{}) string {
	return  fmt.Sprintf(tpl, arg...)
}
func (a *app)SmsSend(c *fiber.Ctx) error {
	input := new(smsSendReq)
	if err := tools.ParseBody(c, input); err != nil {
		return resp.Err(c, 1, err.Error())
	}
	var ps []interface{}
	for _, item := range input.Params{
		ps = append(ps, item)
	}
	//判断商户是否存在
	merchant := new(model.Merchant)
	merchant.One(fmt.Sprintf("id = %d", input.MchId))
	if merchant.Id == 0{
		return  resp.Err(c,1,"没有找到商户")
	}
	//根据type查找短信模板
	tpl := new(model.SmsTemplate)
	lists, _ := tpl.Page(fmt.Sprintf("sms_type = %d", input.Type),1,1000)
	for _, item := range lists{
		smsModel := sms.SelectSms(item.CompanyId)
		content := mySprintf(item.Content, ps)
		ret, err := (*smsModel).Send(input.Phone, content);
		if ret{
			break
		}
		global.Log.Error("发送短信失败%v", err.Error())
	}
	return resp.OK(c,"")
}