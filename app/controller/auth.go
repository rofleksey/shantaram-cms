package controller

import (
	"errors"
	"fmt"
	"github.com/gofiber/fiber/v2"
	"net/http"
	"shantaram-cms/app/dao"
	"shantaram-cms/app/service"
	"shantaram-cms/pkg/util"
)

type Auth struct {
	authService *service.Auth
}

func NewAuth(
	authService *service.Auth,
) *Auth {
	return &Auth{
		authService: authService,
	}
}

func (c *Auth) Login(ctx *fiber.Ctx) error {
	var req dao.LoginRequest

	if err := ctx.BodyParser(&req); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(dao.NoDataResponse{
			Error: true,
			Msg:   fmt.Sprintf("failed to parse body: %v", err),
		})
	}

	token, err := c.authService.AuthAdmin(req.Pass)
	if err != nil {
		if errors.Is(err, util.ErrInvalidCredentials) {
			return ctx.Status(http.StatusUnauthorized).JSON(dao.NoDataResponse{
				Error: true,
				Msg:   fmt.Sprintf("не удалось войти: %v", err),
			})
		}

		return ctx.Status(http.StatusInternalServerError).JSON(dao.NoDataResponse{
			Error: true,
			Msg:   fmt.Sprintf("не удалось войти: %v", err),
		})
	}

	return ctx.JSON(dao.SuccessResponse[dao.LoginResponse]{
		Error: false,
		Data: dao.LoginResponse{
			ID:       1,
			Username: "admin",
			Roles:    []string{"admin"},
			Token:    token,
		},
	})
}
