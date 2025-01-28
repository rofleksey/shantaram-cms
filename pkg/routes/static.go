package routes

import (
	"github.com/gofiber/fiber/v2"
)

func StaticRoutes(app *fiber.App) {
	app.Static("/uploads", "./data/uploads", fiber.Static{
		ByteRange: true,
		Download:  true,
	})
}
