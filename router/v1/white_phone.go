package v1

import (
	"fmt"
	"github.com/gofiber/fiber/v2"
	"loante/service/model"
	"loante/service/resp"
	"loante/tools"
)

type whitePhone struct {
}

func NewWhitePhone() *whitePhone {
	return new(whitePhone)
}

type whitePhoneCreate struct {
	Phone       string `json:"phone"`
	Description string `json:"description"`
}

// WhitePhoneCreate 添加白名单
func (a *whitePhone) WhitePhoneCreate(c *fiber.Ctx) error {
	input := new(whitePhoneCreate)
	if err := tools.ParseBody(c, input); err != nil {
		return resp.Err(c, 1, err.Error())
	}
	white := new(model.WhitePhone)
	white.Phone = input.Phone
	white.Description = input.Description
	white.CreateTime = tools.GetFormatTime()
	white.Insert()
	return resp.OK(c, "")
}

type whitePhoneList struct {
	Phone string `json:"phone"`
	Page  int    `json:"page"`
	Size  int    `json:"size"`
}

// WhitePhoneList 白名单列表
func (a *whitePhone) WhitePhoneList(c *fiber.Ctx) error {
	input := new(whitePhoneList)
	if err := tools.ParseBody(c, input); err != nil {
		return resp.Err(c, 1, err.Error())
	}
	where := "wb.id>0"
	if input.Phone != "" {
		where += " and wb.phone like '%" + input.Phone + "%' "
	}
	white := new(model.WhitePhone)
	lists, count := white.Page(where, input.Page, input.Size)
	return resp.OK(c, map[string]interface{}{
		"count": count,
		"list":  lists,
	})
}

type whitePhoneDel struct {
	Id string `json:"id"`
}

// WhitePhoneDel 删除用户白名单
func (a *whitePhone) WhitePhoneDel(c *fiber.Ctx) error {
	input := new(whitePhoneDel)
	if err := tools.ParseBody(c, input); err != nil {
		return resp.Err(c, 1, err.Error())
	}
	if input.Id == "" {
		return resp.Err(c, 1, "id不能为空")
	} else if tools.ToInt(input.Id) == 0 {
		return resp.Err(c, 1, "id不能为0")
	}
	black := new(model.WhitePhone)
	black.Del(fmt.Sprintf("id=%d", tools.ToInt(input.Id)))
	return resp.OK(c, "")
}
