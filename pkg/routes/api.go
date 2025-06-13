package routes

import (
	"github.com/gofiber/contrib/websocket"
	"github.com/gofiber/fiber/v2"
	"shantaram-cms/app/controller"
	"shantaram-cms/pkg/middleware"
)

func PublicRoutes(
	healthController *controller.Health,
	sitemapController *controller.Sitemap,
	authController *controller.Auth,
	pageController *controller.Page,
	fileController *controller.File,
	orderController *controller.Order,
	captchaController *controller.Captcha,
	generalController *controller.General,
	wsController *controller.WebSocket,
	statsController *controller.Stats,
	app *fiber.App,
) {

	route := app.Group("/v1")

	route.Get("/healthz", healthController.Health)
	route.Get("/sitemap", sitemapController.GetSitemap)

	route.Post("/login", authController.Login)

	route.Get("/ws", middleware.WebSocketUpgrade(), websocket.New(wsController.GlobalHandler))

	route.Get("/page/:id", pageController.GetByID)
	route.Get("/pages", middleware.AdminRestricted(), pageController.GetAll)
	route.Delete("/page/:id", middleware.AdminRestricted(), pageController.Delete)
	route.Put("/page", middleware.AdminRestricted(), pageController.Update)
	route.Post("/page", middleware.AdminRestricted(), pageController.Insert)

	route.Get("/files", middleware.AdminRestricted(), fileController.GetAll)
	route.Delete("/file/:id", middleware.AdminRestricted(), fileController.Delete)
	route.Post("/file", middleware.AdminRestricted(), fileController.Upload)
	route.Get("/file/stats", middleware.AdminRestricted(), fileController.Stats)

	route.Get("/order/:id", middleware.AdminRestricted(), orderController.GetByID)
	route.Get("/orders", middleware.AdminRestricted(), orderController.GetPaginated)
	route.Delete("/order/:id", middleware.AdminRestricted(), orderController.Delete)
	route.Put("/order", middleware.AdminRestricted(), orderController.Update)
	route.Post("/order", orderController.Insert)

	route.Get("/general/:id", generalController.GetByID)
	route.Post("/general/:id", middleware.AdminRestricted(), generalController.Upsert)

	route.Get("/captcha", captchaController.Generate)

	route.Get("/stats", middleware.AdminRestricted(), statsController.Get)
}
