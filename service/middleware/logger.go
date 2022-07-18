package middleware

import (
	"github.com/gofiber/fiber/v2"
	"loante/service/model"
)

func Logger(c *fiber.Ctx) error {
	switch c.Method() {
	case "POST":
		//所以POST 写操作记录到数据库中
		log := new(model.AdminLog)
		log.Operate = string( c.Body())
		log.UserName =  c.Locals("userName").(string)
		log.RoleName =  c.Locals("roleName").(string)
		log.Path = c.Path()
		log.Ip = c.IP()
		log.Insert()
	}
	c.Next()
	return nil
}
