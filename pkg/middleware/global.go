package middleware

import (
	jwtMiddleware "github.com/gofiber/contrib/jwt"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/golang-jwt/jwt/v5"
	"shantaram-cms/pkg/config"
)

func FiberMiddleware(app *fiber.App, cfg *config.Config) {
	var args []any

	args = append(args, recover.New())

	args = append(args, logger.New(logger.Config{
		Next: func(c *fiber.Ctx) bool {
			return c.Path() == "/api/healthz"
		},
	}))

	args = append(args, cors.New())

	args = append(args, jwtMiddleware.New(jwtMiddleware.Config{
		SigningKey: jwtMiddleware.SigningKey{Key: []byte(cfg.JWTKey)},
		SuccessHandler: func(ctx *fiber.Ctx) error {
			tokenOpt := ctx.Locals("user")
			if tokenOpt == nil {
				ctx.Locals("username", "")
				return ctx.Next()
			}

			token := tokenOpt.(*jwt.Token)
			claims := token.Claims.(jwt.MapClaims)

			idOpt := claims["id"]
			if idOpt == nil {
				ctx.Locals("username", "")
				return ctx.Next()
			}

			ctx.Locals("username", idOpt.(string))
			return ctx.Next()
		},
		ErrorHandler: func(ctx *fiber.Ctx, err error) error {
			ctx.Locals("username", "")
			return ctx.Next()
		},
		TokenLookup: "query:token,header:Authorization",
		AuthScheme:  "Bearer",
	}))

	app.Use(args...)
}
