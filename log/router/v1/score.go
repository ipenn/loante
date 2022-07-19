package v1

import (
	"fmt"
	"github.com/gofiber/fiber/v2"
	"loante/service/model"
	"loante/service/resp"
	"loante/tools"
	"math"
)

type score struct {
}

func NewScore() *score {
	return new(score)
}
func (a *score) Stat(c *fiber.Ctx) error {
	uId := c.Locals("userId").(int)
	return resp.OK(c, map[string]interface{}{
		"total":  new(model.Score).Balance(fmt.Sprintf("user_id = %d and status = 1 and amount > 0", uId)),
		"equity": new(model.Score).Balance(fmt.Sprintf("user_id = %d and status >= 0", uId)),
	})
}

type pageScoreReq struct {
	Page  int `json:"page"`
	Limit int `json:"limit"`
}

func (a *score) List(c *fiber.Ctx) error {
	input := new(pageScoreReq)
	if err := tools.ParseBody(c, input); err != nil {
		return resp.Err(c, 1, err.Error())
	}
	uId := c.Locals("userId").(int)
	where := fmt.Sprintf("user_id = %d and amount > 0", uId)
	list, count := new(model.Score).Page(where, input.Page, input.Limit)
	return resp.OK(c, map[string]interface{}{
		"list":  list,
		"count": count,
	})
}

func (a *score) Exchange(c *fiber.Ctx) error {
	input := new(pageScoreReq)
	if err := tools.ParseBody(c, input); err != nil {
		return resp.Err(c, 1, err.Error())
	}
	uId := c.Locals("userId").(int)
	where := fmt.Sprintf("user_id = %d and amount < 0", uId)
	list, count := new(model.Score).Page(where, input.Page, input.Limit)
	return resp.OK(c, map[string]interface{}{
		"list":  list,
		"count": count,
	})
}

type scoreExReq struct {
	Amount float64 `json:"amount"`
}

func (a *score) ExchangeAction(c *fiber.Ctx) error {
	input := new(scoreExReq)
	if err := tools.ParseBody(c, input); err != nil {
		return resp.Err(c, 1, err.Error())
	}
	input.Amount = math.Abs(input.Amount)
	uId := c.Locals("userId").(int)
	balance := new(model.Score).Balance(fmt.Sprintf("user_id = %d and status >= 0", uId))
	if input.Amount > balance {
		return resp.Err(c, 1, "可兑换金额不足")
	}
	scoreModal := new(model.Score)
	scoreModal.UserId = fmt.Sprintf("%d", uId)
	scoreModal.Amount = input.Amount * -1
	scoreModal.CreateTime = tools.GetFormatTime()
	scoreModal.Status = 1
	scoreModal.FromUserId = uId
	scoreModal.Comment = fmt.Sprintf("兑换")
	scoreModal.Insert()
	return resp.OK(c, map[string]interface{}{})
}
