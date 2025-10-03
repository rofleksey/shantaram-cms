package middleware

import (
	"context"
	"log/slog"
	"net/http"
	"runtime/debug"
	"shantaram/pkg/config"
	"shantaram/pkg/telemetry"
	"shantaram/pkg/util"
	"slices"
	"strings"

	"github.com/elliotchance/pie/v2"
	sentryotel "github.com/getsentry/sentry-go/otel"
	jwtware "github.com/gofiber/contrib/jwt"
	"github.com/gofiber/contrib/otelfiber/v2"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/rofleksey/meg"
	"github.com/samber/do"
	slogfiber "github.com/samber/slog-fiber"
)

func FiberMiddleware(app *fiber.App, di *do.Injector) {
	cfg := do.MustInvoke[*config.Config](di)
	tel := do.MustInvoke[*telemetry.Telemetry](di)

	staticOrigins := []string{
		cfg.BaseApiURL, cfg.BaseFrontURL, cfg.BaseWWWFrontURL, cfg.BaseAdminURL,
		"capacitor://localhost", "http://localhost", "https://localhost", "http://localhost:4321", "http://localhost:5173",
		"http://localhost:1234", "http://localhost:3000", "http://localhost:9000", "http://localhost:8080",
	}

	// cors
	app.Use(cors.New(cors.Config{
		AllowHeaders:     "Origin, Content-Type, Accept, Authorization, Sentry-Trace, Baggage",
		AllowMethods:     "POST, GET, OPTIONS, DELETE, PUT, PATCH, HEAD",
		AllowCredentials: true,
		AllowOriginsFunc: func(origin string) bool {
			if pie.Contains(staticOrigins, origin) {
				return true
			}

			return false
		},
	}))

	app.Use(otelfiber.Middleware(
		otelfiber.WithMeterProvider(tel.MeterProvider),
		otelfiber.WithTracerProvider(tel.TracerProvider),
		otelfiber.WithPropagators(sentryotel.NewSentryPropagator()),
		otelfiber.WithCollectClientIP(true),
	))

	// retrieve user ip
	app.Use(func(ctx *fiber.Ctx) error {
		ctx.SetUserContext(context.WithValue(ctx.UserContext(), util.IpContextKey, ctx.IP()))

		return ctx.Next()
	})

	ignorePaths := []string{"/api/healthz"}

	// log requests
	app.Use(slogfiber.NewWithConfig(slog.Default(), slogfiber.Config{
		Filters: []slogfiber.Filter{
			func(c *fiber.Ctx) bool {
				return !slices.Contains(ignorePaths, c.Path())
			},
			func(ctx *fiber.Ctx) bool {
				reqMethod := strings.ToLower(string(ctx.Context().Method()))
				return !(reqMethod == "get" && (ctx.Response().StatusCode() == http.StatusOK || ctx.Response().StatusCode() == http.StatusNotModified || ctx.Response().StatusCode() == http.StatusPartialContent)) //nolint:staticcheck
			},
		},
		WithTraceID: true,
	}))

	app.Use(recover.New(recover.Config{
		EnableStackTrace: true,
		StackTraceHandler: func(ctx *fiber.Ctx, e any) {
			stackStr := meg.TrimSuffixToNRunes(string(debug.Stack()), 2048)

			slog.ErrorContext(ctx.Context(), "Panic",
				slog.Any("error", e),
				slog.String("stack", stackStr),
			)
		},
	}))

	app.Use(jwtware.New(jwtware.Config{
		SigningKey: jwtware.SigningKey{Key: []byte(cfg.JWT.Secret)},
		SuccessHandler: func(ctx *fiber.Ctx) error {
			tokenOpt := ctx.Locals("user")
			if tokenOpt == nil {
				return ctx.Next()
			}

			ctx.Locals("admin", true)
			newUserCtx := context.WithValue(ctx.UserContext(), "admin", true)
			newUserCtx = context.WithValue(newUserCtx, util.UsernameContextKey, "admin")
			ctx.SetUserContext(newUserCtx)

			return ctx.Next()
		},
		ErrorHandler: func(ctx *fiber.Ctx, _ error) error {
			return ctx.Next()
		},
		TokenLookup: "query:token,header:Authorization",
		AuthScheme:  "Bearer",
	}))
}
