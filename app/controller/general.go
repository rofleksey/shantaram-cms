package controller

import (
	"fmt"
	"github.com/gofiber/fiber/v2"
	"net/http"
	"shantaram-cms/app/dao"
	"shantaram-cms/app/service"
)

type General struct {
	generalService *service.General
}

func NewGeneral(
	generalService *service.General,
) *General {
	return &General{
		generalService: generalService,
	}
}

func (c *General) Upsert(ctx *fiber.Ctx) error {
	id := ctx.Params("id")

	var req any

	if err := ctx.BodyParser(&req); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(dao.NoDataResponse{
			Error: true,
			Msg:   fmt.Sprintf("failed to parse body: %v", err),
		})
	}

	if err := c.generalService.Upsert(id, req); err != nil {
		return ctx.Status(http.StatusInternalServerError).JSON(dao.NoDataResponse{
			Error: true,
			Msg:   fmt.Sprintf("не удалось обновить настройки: %v", err),
		})
	}

	return ctx.JSON(dao.NoDataResponse{
		Msg: "success",
	})
}

func (c *General) GetByID(ctx *fiber.Ctx) error {
	id := ctx.Params("id")

	data, err := c.generalService.GetByID(id)
	if err != nil {
		return ctx.Status(http.StatusInternalServerError).JSON(dao.NoDataResponse{
			Error: true,
			Msg:   fmt.Sprintf("не удалось получить настройки по id: %v", err),
		})
	}

	return ctx.JSON(dao.SuccessResponse[any]{
		Data: data,
	})
}
