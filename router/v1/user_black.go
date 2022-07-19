package v1

import (
	"fmt"
	"github.com/gofiber/fiber/v2"
	"loante/service/model"
	"loante/service/resp"
	"loante/tools"
)

type userBlack struct {
}

func NewUserBlack() *userBlack {
	return new(userBlack)
}

type userBlackCreate struct {
	Content     string `json:"content"`
	Description string `json:"description"`
	Type        string `json:"type"`
}

// UserBlackCreate 添加黑名单
func (a *userBlack) UserBlackCreate(c *fiber.Ctx) error {
	input := new(userBlackCreate)
	if err := tools.ParseBody(c, input); err != nil {
		return resp.Err(c, 1, err.Error())
	}
	black := new(model.UserBlack)
	black.Type = tools.ToInt(input.Type)
	black.Content = input.Content
	black.Description = input.Description
	black.CreateTime = tools.GetFormatTime()
	black.Insert()
	return resp.OK(c, "")
}

type userBlackList struct {
	Content string `json:"content"`
	Type    string `json:"type"`
	Page    int    `json:"page"`
	Size    int    `json:"size"`
}

// UserBlackList 黑名单列表
func (a *userBlack) UserBlackList(c *fiber.Ctx) error {
	input := new(userBlackList)
	if err := tools.ParseBody(c, input); err != nil {
		return resp.Err(c, 1, err.Error())
	}
	where := "ub.id>0"
	if input.Content != "" {
		where += " and ub.content like '%" + input.Content + "%' "
	}
	if input.Type != "" {
		where += fmt.Sprintf(" and ub.type=%d", tools.ToInt(input.Type))
	}
	black := new(model.UserBlack)
	lists, count := black.Page(where, input.Page, input.Size)
	return resp.OK(c, map[string]interface{}{
		"count": count,
		"list":  lists,
	})
}

type userBlackDel struct {
	Id string `json:"id"`
}

// UserBlackDel 删除用户黑名单
func (a *userBlack) UserBlackDel(c *fiber.Ctx) error {
	input := new(userBlackDel)
	if err := tools.ParseBody(c, input); err != nil {
		return resp.Err(c, 1, err.Error())
	}
	if input.Id == "" {
		return resp.Err(c, 1, "id不能为空")
	} else if tools.ToInt(input.Id) == 0 {
		return resp.Err(c, 1, "id不能为0")
	}
	black := new(model.UserBlack)
	black.Del(fmt.Sprintf("id=%d", tools.ToInt(input.Id)))
	return resp.OK(c, "")
}
