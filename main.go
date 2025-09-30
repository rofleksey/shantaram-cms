package main

import (
	"context"
	"log/slog"
	"os"
	"os/signal"
	"shantaram/app/api"
	"shantaram/app/controller"
	"shantaram/app/service/auth"
	"shantaram/app/service/limits"
	"shantaram/app/service/menu"
	"shantaram/app/service/order"
	"shantaram/app/service/params"
	"shantaram/app/service/pubsub"
	"shantaram/app/service/telegram"
	"shantaram/pkg/config"
	"shantaram/pkg/database"
	"shantaram/pkg/middleware"
	"shantaram/pkg/migration"
	"shantaram/pkg/routes"
	"shantaram/pkg/telemetry"
	"shantaram/pkg/tlog"
	"time"

	"github.com/exaring/otelpgx"
	"github.com/getsentry/sentry-go"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/log"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/samber/do"
)

func main() {
	appCtx, cancel := context.WithCancel(context.Background())
	defer cancel()

	di := do.New()
	do.ProvideValue(di, appCtx)

	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("config load failed: %v", err)
	}
	do.ProvideValue(di, cfg)

	if err = telemetry.InitSentry(cfg); err != nil {
		log.Fatalf("sentry init failed: %v", err)
	}
	defer sentry.Flush(3 * time.Second)

	tel, err := telemetry.Init(cfg)
	if err != nil {
		log.Fatalf("telemetry init failed: %v", err)
	}
	defer tel.Shutdown(appCtx)
	do.ProvideValue(di, tel)

	if err = tlog.Init(cfg, tel); err != nil {
		log.Fatalf("logging init failed: %v", err)
	}

	metrics, err := telemetry.NewMetrics(cfg, tel.Meter)
	if err != nil {
		log.Fatalf("metrics init failed: %v", err)
	}
	do.ProvideValue(di, metrics)

	tracing := telemetry.NewTracing(cfg, tel.Tracer)
	do.ProvideValue(di, tracing)

	slog.ErrorContext(appCtx, "Service restarted")

	dbConnStr := "postgres://" + cfg.DB.User + ":" + cfg.DB.Pass + "@" + cfg.DB.Host + "/" + cfg.DB.Database + "?sslmode=disable&pool_max_conns=30&pool_min_conns=5&pool_max_conn_lifetime=1h&pool_max_conn_idle_time=30m&pool_health_check_period=1m&connect_timeout=10"

	dbConf, err := pgxpool.ParseConfig(dbConnStr)
	if err != nil {
		log.Fatalf("pgxpool.ParseConfig() failed: %v", err)
	}

	dbConf.ConnConfig.RuntimeParams = map[string]string{
		"statement_timeout":                   "30000",
		"idle_in_transaction_session_timeout": "60000",
	}
	dbConf.ConnConfig.Tracer = otelpgx.NewTracer(
		otelpgx.WithMeterProvider(tel.MeterProvider),
		otelpgx.WithTracerProvider(tel.TracerProvider),
	)

	dbConn, err := pgxpool.NewWithConfig(appCtx, dbConf)
	if err != nil {
		log.Fatalf("failed to connect to database: %v", err)
	}
	defer dbConn.Close()

	if err = otelpgx.RecordStats(dbConn); err != nil {
		log.Fatalf("unable to record database stats: %v", err)
	}

	if err = database.InitSchema(appCtx, dbConn); err != nil {
		log.Fatalf("failed to init schema: %v", err)
	}

	do.ProvideValue(di, dbConn)

	queries := database.New(dbConn)
	do.ProvideValue(di, queries)

	if err = migration.Migrate(appCtx, di); err != nil {
		log.Fatalf("failed to migrate: %v", err)
	}

	do.Provide(di, pubsub.New)
	do.Provide(di, auth.New)
	do.Provide(di, limits.New)
	do.Provide(di, telegram.New)
	do.Provide(di, menu.New)
	do.Provide(di, order.New)
	do.Provide(di, params.New)

	go do.MustInvoke[*params.Service](di).RunHeaderDeadline(appCtx)

	wsController := controller.NewWS(di)

	server := controller.NewStrictServer(di)
	handler := api.NewStrictHandler(server, nil)

	app := fiber.New(fiber.Config{
		AppName:          "Shantaram API",
		ErrorHandler:     middleware.ErrorHandler,
		ProxyHeader:      "X-Forwarded-For",
		ReadTimeout:      time.Second * 60,
		WriteTimeout:     time.Second * 60,
		DisableKeepalive: false,
	})

	middleware.FiberMiddleware(app, di)
	routes.StaticRoutes(app)
	routes.WSRoutes(app, wsController)

	apiGroup := app.Group("/v1")
	api.RegisterHandlersWithOptions(apiGroup, handler, api.FiberServerOptions{
		BaseURL: "",
		Middlewares: []api.MiddlewareFunc{
			middleware.NewOpenAPIValidator(),
		},
	})

	routes.NotFoundRoute(app)

	go func() {
		sigint := make(chan os.Signal, 1)
		signal.Notify(sigint, os.Interrupt)
		<-sigint

		log.Info("Shutting down server...")

		_ = app.Shutdown()
		cancel()
	}()

	log.Info("Server started on port 8080")
	if err := app.Listen(":8080"); err != nil {
		log.Warnf("Server stopped! Reason: %v", err)
	}

	log.Info("Waiting for services to finish...")
	_ = di.Shutdown()
}
