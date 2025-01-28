package controller

import (
	"fmt"
	"github.com/gofiber/fiber/v2"
	"net/http"
	"shantaram-cms/app/dao"
	"shantaram-cms/app/service"
	"shantaram-cms/pkg/database"
)

type Page struct {
	pageService *service.Page
}

func NewPage(
	pageService *service.Page,
) *Page {
	return &Page{
		pageService: pageService,
	}
}

func (c *Page) Insert(ctx *fiber.Ctx) error {
	var req database.Page

	if err := ctx.BodyParser(&req); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(dao.NoDataResponse{
			Error: true,
			Msg:   fmt.Sprintf("failed to parse body: %v", err),
		})
	}

	if err := c.pageService.Insert(&req); err != nil {
		return ctx.Status(http.StatusInternalServerError).JSON(dao.NoDataResponse{
			Error: true,
			Msg:   fmt.Sprintf("не удалось создать страницу: %v", err),
		})
	}

	return ctx.JSON(dao.NoDataResponse{
		Msg: "success",
	})
}

func (c *Page) Update(ctx *fiber.Ctx) error {
	var req database.Page

	if err := ctx.BodyParser(&req); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(dao.NoDataResponse{
			Error: true,
			Msg:   fmt.Sprintf("failed to parse body: %v", err),
		})
	}

	if err := c.pageService.Update(&req); err != nil {
		return ctx.Status(http.StatusInternalServerError).JSON(dao.NoDataResponse{
			Error: true,
			Msg:   fmt.Sprintf("не удалось обновить страницу: %v", err),
		})
	}

	return ctx.JSON(dao.NoDataResponse{
		Msg: "success",
	})
}

func (c *Page) Delete(ctx *fiber.Ctx) error {
	id := ctx.Params("id")

	if err := c.pageService.Delete(id); err != nil {
		return ctx.Status(http.StatusInternalServerError).JSON(dao.NoDataResponse{
			Error: true,
			Msg:   fmt.Sprintf("не удалось обновить страницу: %v", err),
		})
	}

	return ctx.JSON(dao.NoDataResponse{
		Msg: "success",
	})
}

func (c *Page) GetByID(ctx *fiber.Ctx) error {
	id := ctx.Params("id")

	page, err := c.pageService.GetByID(id)
	if err != nil {
		return ctx.Status(http.StatusInternalServerError).JSON(dao.NoDataResponse{
			Error: true,
			Msg:   fmt.Sprintf("не удалось получить страницу по id: %v", err),
		})
	}

	return ctx.JSON(dao.SuccessResponse[database.Page]{
		Data: *page,
	})
}

func (c *Page) GetAll(ctx *fiber.Ctx) error {
	pages, err := c.pageService.GetAll()
	if err != nil {
		return ctx.Status(http.StatusInternalServerError).JSON(dao.NoDataResponse{
			Error: true,
			Msg:   fmt.Sprintf("не удалось получить список страниц: %v", err),
		})
	}

	return ctx.JSON(dao.SuccessResponse[[]database.Page]{
		Data: pages,
	})
}
