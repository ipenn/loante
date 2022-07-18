package v1

import (
	"github.com/gofiber/fiber/v2"
	"loante/service/model"
	"loante/service/req"
	"loante/service/resp"
	"loante/tools"
)

type borrow struct {
	
}

func NewBorrow() *borrow {
	return new(borrow)
}

type borrowQueryReq struct {
	req.PageReq
	ProductId int `query:"product_id" json:"product_id"`
}
func (a *borrow)Query(c *fiber.Ctx) error {
	input := new(borrowQueryReq)
	if err := tools.ParseBody(c, input); err != nil {
		return resp.Err(c, 1, err.Error())
	}
	lists, count := new(model.Borrow).Page("b.id > 0", input.Page, input.Size)
	return resp.OK(c, map[string]interface{}{
		"count":count,
		"list":lists,
	})
}