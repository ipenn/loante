package resp

import "github.com/gofiber/fiber/v2"

type FiberRes struct {
	Code int  `json:"code"`
	Msg string	`json:"msg"`
	Data interface{}	`json:"data"`
}
func OK(c *fiber.Ctx, data interface{}) error{
	return c.JSON(FiberRes{
		Code: 0,
		Msg: "",
		Data: data,
	})
}

func Err(c *fiber.Ctx,code int, data interface{}) error {
	return c.JSON(FiberRes{
		Code: code,
		Msg: data.(string),
		Data: "",
	})
}
