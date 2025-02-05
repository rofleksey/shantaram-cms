package controller

import (
	"fmt"
	"github.com/gofiber/fiber/v2"
	"net/http"
	"shantaram-cms/app/dao"
	"shantaram-cms/app/service"
	"shantaram-cms/pkg/database"
	"strconv"
)

type Order struct {
	orderService *service.Order
}

func NewOrder(
	orderService *service.Order,
) *Order {
	return &Order{
		orderService: orderService,
	}
}

func (c *Order) Insert(ctx *fiber.Ctx) error {
	var req dao.NewOrderRequest

	if err := ctx.BodyParser(&req); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(dao.NoDataResponse{
			Error: true,
			Msg:   fmt.Sprintf("failed to parse body: %v", err),
		})
	}

	if len(req.Items) == 0 {
		return ctx.Status(fiber.StatusBadRequest).JSON(dao.NoDataResponse{
			Error: true,
			Msg:   "empty order",
		})
	}

	if err := c.orderService.Insert(req); err != nil {
		return ctx.Status(http.StatusInternalServerError).JSON(dao.NoDataResponse{
			Error: true,
			Msg:   fmt.Sprintf("не удалось создать заказ: %v", err),
		})
	}

	return ctx.JSON(dao.NoDataResponse{
		Msg: "success",
	})
}

func (c *Order) Update(ctx *fiber.Ctx) error {
	var req database.Order

	if err := ctx.BodyParser(&req); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(dao.NoDataResponse{
			Error: true,
			Msg:   fmt.Sprintf("failed to parse body: %v", err),
		})
	}

	if err := c.orderService.Update(&req); err != nil {
		return ctx.Status(http.StatusInternalServerError).JSON(dao.NoDataResponse{
			Error: true,
			Msg:   fmt.Sprintf("не удалось обновить заказ: %v", err),
		})
	}

	return ctx.JSON(dao.NoDataResponse{
		Msg: "success",
	})
}

func (c *Order) Delete(ctx *fiber.Ctx) error {
	idStr := ctx.Params("id")

	id, err := strconv.Atoi(idStr)
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(dao.NoDataResponse{
			Error: true,
			Msg:   fmt.Sprintf("invalid id: %v", err),
		})
	}

	if err := c.orderService.Delete(uint64(id)); err != nil {
		return ctx.Status(http.StatusInternalServerError).JSON(dao.NoDataResponse{
			Error: true,
			Msg:   fmt.Sprintf("не удалось обновить заказ: %v", err),
		})
	}

	return ctx.JSON(dao.NoDataResponse{
		Msg: "success",
	})
}

func (c *Order) GetByID(ctx *fiber.Ctx) error {
	idStr := ctx.Params("id")

	id, err := strconv.Atoi(idStr)
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(dao.NoDataResponse{
			Error: true,
			Msg:   fmt.Sprintf("invalid id: %v", err),
		})
	}

	page, err := c.orderService.GetByID(uint64(id))
	if err != nil {
		return ctx.Status(http.StatusInternalServerError).JSON(dao.NoDataResponse{
			Error: true,
			Msg:   fmt.Sprintf("не удалось получить заказ по id: %v", err),
		})
	}

	return ctx.JSON(dao.SuccessResponse[database.Order]{
		Data: *page,
	})
}

func (c *Order) GetPaginated(ctx *fiber.Ctx) error {
	offsetStr := ctx.Query("offset")
	limitStr := ctx.Query("limit")

	offset, _ := strconv.Atoi(offsetStr)
	limit, _ := strconv.Atoi(limitStr)

	if offset < 0 {
		offset = 0
	}

	if limit < 0 {
		limit = 0
	}

	if limit > 100 {
		limit = 100
	}

	data, err := c.orderService.GetPaginated(offset, limit)
	if err != nil {
		return ctx.Status(http.StatusInternalServerError).JSON(dao.NoDataResponse{
			Error: true,
			Msg:   fmt.Sprintf("не удалось получить список страниц: %v", err),
		})
	}

	return ctx.JSON(dao.SuccessResponse[database.DataPage[database.Order]]{
		Data: *data,
	})
}
