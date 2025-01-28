package routes

import (
	"github.com/gofiber/fiber/v2"
	"shantaram-cms/app/dao"
)

func NotFoundRoute(a *fiber.App) {
	a.Use(
		func(c *fiber.Ctx) error {
			return c.Status(fiber.StatusNotFound).JSON(dao.NoDataResponse{
				Error: true,
				Msg:   "route not found",
			})
		},
	)
}
