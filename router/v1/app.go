package v1

import (
	"fmt"
	"github.com/gofiber/fiber/v2"
	"loante/service/model"
	"loante/service/resp"
	"loante/tools"
)

type app struct {
}

func NewApp() *app {
	return new(app)
}

//type 是哪类型的短信
type smsSendReq struct {
	MchId  int      `query:"mch_id" json:"mch_id"`
	Type   int      `query:"type" json:"type"`
	Phone  string   `query:"phone" json:"phone"`
	Params []string `query:"params" json:"params"`
}

func (a *app) SmsSend(c *fiber.Ctx) error {
	input := new(smsSendReq)
	//fmt.Println(string(c.Body()))
	if err := tools.ParseBody(c, input); err != nil {
		return resp.Err(c, 1, err.Error())
	}
	var ps []string
	for _, item := range input.Params {
		ps = append(ps, item)
	}
	//根据type查找短信模板
	if new(model.SmsTemplate).Send(input.Type, input.Phone, ps) {
		//判断商户是否存在
		if input.MchId > 0{
			merchantData := new(model.Merchant)
			merchantData.One(fmt.Sprintf("id = %d", input.MchId))
			if merchantData.Id > 0 {
				merchantData.AddService(1, 1) //扣费
			}
		}
		return resp.OK(c, "")
	}
	return resp.Err(c, 1, "发送失败")
}
