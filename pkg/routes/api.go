package routes

import (
	"github.com/gofiber/fiber/v2"
	"shantaram-cms/app/controller"
	"shantaram-cms/pkg/middleware"
)

func PublicRoutes(
	healthController *controller.Health,
	authController *controller.Auth,
	pageController *controller.Page,
	fileController *controller.File,
	app *fiber.App,
) {

	route := app.Group("/v1")

	route.Get("/healthz", healthController.Health)

	route.Post("/login", authController.Login)

	route.Get("/page/:id", pageController.GetByID)
	route.Get("/pages", middleware.AdminRestricted(), pageController.GetAll)
	route.Delete("/page/:id", middleware.AdminRestricted(), pageController.Delete)
	route.Put("/page", middleware.AdminRestricted(), pageController.Update)
	route.Post("/page", middleware.AdminRestricted(), pageController.Insert)

	route.Get("/files", middleware.AdminRestricted(), fileController.GetAll)
	route.Delete("/file/:id", middleware.AdminRestricted(), fileController.Delete)
	route.Post("/file", middleware.AdminRestricted(), fileController.Upload)
	route.Get("/file/stats", middleware.AdminRestricted(), fileController.Stats)
}
