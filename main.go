package main

import (
	"context"
	"github.com/gofiber/fiber/v2"
	_ "go.uber.org/automaxprocs"
	"log"
	"os"
	"os/signal"
	"path/filepath"
	"shantaram-cms/app/controller"
	"shantaram-cms/app/service"
	"shantaram-cms/pkg/config"
	"shantaram-cms/pkg/database"
	"shantaram-cms/pkg/middleware"
	"shantaram-cms/pkg/routes"
)

func main() {
	appCtx, cancel := context.WithCancel(context.Background())
	defer cancel()

	os.MkdirAll(filepath.Join("data", "uploads"), os.ModePerm)
	os.MkdirAll(filepath.Join("data", "temp"), os.ModePerm)

	exitChan := make(chan struct{})

	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("failed to load config: %s", err.Error())
	}

	db, err := database.New()
	if err != nil {
		log.Fatalf("failed to connect database: %s", err.Error())
	}

	if err := db.Init(); err != nil {
		log.Fatalf("failed to init database: %s", err.Error())
	}

	authService := service.NewAuth(cfg)
	pageService := service.NewPage(db)
	uploadsService := service.NewUploads(appCtx)

	healthController := controller.NewHealth()
	authController := controller.NewAuth(authService)
	pageController := controller.NewPage(pageService)

	app := fiber.New(fiber.Config{
		BodyLimit: 1024 * 1024 * 100, // 100 mb
	})

	middleware.FiberMiddleware(app, cfg)
	routes.StaticRoutes(app)
	routes.PublicRoutes(healthController, authController, pageController, app)
	routes.NotFoundRoute(app)

	go func() {
		sigint := make(chan os.Signal, 1)
		signal.Notify(sigint, os.Interrupt)
		<-sigint

		log.Println("Shutting down server...")

		if err := app.Shutdown(); err != nil {
			log.Printf("Server is shutting down! Reason: %v", err)
		}

		close(exitChan)
	}()

	if err := app.Listen(":8080"); err != nil {
		log.Printf("Server stopped! Reason: %v", err)
	}

	<-exitChan
	cancel()

	log.Println("Waiting for services to finish...")

	uploadsService.CancelAndJoin()
}
