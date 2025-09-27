package middleware

import (
	"context"
	"log/slog"
	"net/http"
	"runtime/debug"
	"shantaram/pkg/config"
	"shantaram/pkg/util"
	"slices"
	"strings"

	"github.com/elliotchance/pie/v2"
	jwtware "github.com/gofiber/contrib/jwt"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/samber/do"
)

func FiberMiddleware(app *fiber.App, di *do.Injector) {
	cfg := do.MustInvoke[*config.Config](di)

	staticOrigins := []string{
		cfg.BaseApiURL,
		"capacitor://localhost", "http://localhost", "https://localhost", "http://localhost:4321", "http://localhost:5173",
		"http://localhost:1234", "http://localhost:3000", "http://localhost:9000", "http://localhost:8080",
	}

	// cors
	app.Use(cors.New(cors.Config{
		AllowHeaders:     "Origin, Content-Type, Accept, Authorization",
		AllowMethods:     "POST, GET, OPTIONS, DELETE, PUT, PATCH, HEAD",
		AllowCredentials: true,
		AllowOriginsFunc: func(origin string) bool {
			if pie.Contains(staticOrigins, origin) {
				return true
			}

			return false
		},
	}))

	// retrieve user ip
	app.Use(func(ctx *fiber.Ctx) error {
		ctx.SetUserContext(context.WithValue(ctx.UserContext(), util.IpContextKey, ctx.IP()))

		return ctx.Next()
	})

	ignorePaths := []string{"/api/healthz"}

	// log requests
	app.Use(NewLogWithConfig(slog.Default(), LogConfig{
		WithUserAgent:    true,
		WithRequestBody:  true,
		WithResponseBody: true,
		Filters: []func(*fiber.Ctx) bool{
			// ignore spam endpoints
			func(c *fiber.Ctx) bool {
				return !slices.Contains(ignorePaths, c.Path())
			},
			// ignore successful GET-s
			func(ctx *fiber.Ctx) bool {
				reqMethod := strings.ToLower(string(ctx.Context().Method()))
				return !(reqMethod == "get" && (ctx.Response().StatusCode() == http.StatusOK || ctx.Response().StatusCode() == http.StatusNotModified || ctx.Response().StatusCode() == http.StatusPartialContent)) //nolint:staticcheck
			},
		},
	}))

	app.Use(recover.New(recover.Config{
		EnableStackTrace: true,
		StackTraceHandler: func(ctx *fiber.Ctx, e any) {
			stackStr := util.TrimSuffixToNRunes(string(debug.Stack()), 2048)

			slog.Error("Panic",
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
