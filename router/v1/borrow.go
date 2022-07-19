package v1

import (
	"fmt"
	"github.com/gofiber/fiber/v2"
	"loante/service/model"
	"loante/service/req"
	"loante/service/resp"
	"loante/tools"
)

type borrow struct{}

func NewBorrow() *borrow {
	return new(borrow)
}

type borrowQueryReq struct {
	req.PageReq
	ProductId       int    `query:"product_id" json:"product_id"`
	UserId          int    `query:"user_id" json:"user_id"`
	Name            string `query:"name" json:"name"`
	NoApplying      string `query:"no_applying" json:"no_applying"`             //去除申请中的
	ProcessingInPay string `query:"processing_in_pay" json:"processing_in_pay"` //支付公司放款处理中
}

func (a *borrow) Query(c *fiber.Ctx) error {
	input := new(borrowQueryReq)
	if err := tools.ParseBody(c, input); err != nil {
		return resp.Err(c, 1, err.Error())
	}
	where := "b.id > 0"
	if len(input.NoApplying) > 0 {
		where += " and b.status > 4"
	}
	if len(input.ProcessingInPay) > 0 {
		where += " and b.status = 4"
	}
	if input.ProductId > 0 {
		where += fmt.Sprintf(" and b.product_id =%d", input.ProductId)
	}
	if input.UserId > 0 {
		where += fmt.Sprintf(" and b.uid =%d", input.UserId)
	}
	if input.UserId > 0 {
		where += fmt.Sprintf(" and u.aadhaar_name ='%s'", input.Name)
	}
	lists, count := new(model.Borrow).Page("b.id > 0", input.Page, input.Size)
	return resp.OK(c, map[string]interface{}{
		"count": count,
		"list":  lists,
	})
}
