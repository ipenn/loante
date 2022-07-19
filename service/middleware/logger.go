package middleware

import (
	"github.com/gofiber/fiber/v2"
	"loante/service/model"
)

func Logger(c *fiber.Ctx) error {
	switch c.Method() {
	case "POST":
		//POST 写操作记录到数据库中
		log := new(model.AdminLog)
		log.ReqBody = string( c.Body())
		log.AdminName =  c.Locals("userName").(string)
		log.AdminId =  c.Locals("userId").(int)
		log.Path = c.Path()
		log.Insert()
	}
	c.Next()
	return nil
}
