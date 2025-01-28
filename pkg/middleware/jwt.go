package middleware

import (
	"github.com/gofiber/fiber/v2"
	"shantaram-cms/app/dao"
)

func AdminRestricted() func(*fiber.Ctx) error {
	return func(c *fiber.Ctx) error {
		if c.Locals("username").(string) != "admin" {
			return c.Status(fiber.StatusUnauthorized).JSON(dao.NoDataResponse{
				Error: true,
				Msg:   "unauthorized",
			})
		}

		return c.Next()
	}
}
