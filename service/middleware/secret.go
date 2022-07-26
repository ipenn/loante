package middleware

import (
	"github.com/gofiber/fiber/v2"
)

func Secret(c *fiber.Ctx) error {
	h := c.Get("sign")
	if h == "" {
		//return resp.Err(c, 1001, "sign error!")
	}
	c.Next()
	return nil
}
