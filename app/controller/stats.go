package controller

import (
	"fmt"
	"github.com/gofiber/fiber/v2"
	"net/http"
	"shantaram-cms/app/dao"
	"shantaram-cms/app/service"
)

type Stats struct {
	statsService *service.Stats
}

func NewStats(
	statsService *service.Stats,
) *Stats {
	return &Stats{
		statsService: statsService,
	}
}

func (c *Stats) Get(ctx *fiber.Ctx) error {
	stats, err := c.statsService.GetStats()

	if err != nil {
		return ctx.Status(http.StatusInternalServerError).JSON(dao.NoDataResponse{
			Error: true,
			Msg:   fmt.Sprintf("не удалось получить статистику: %v", err),
		})
	}

	return ctx.JSON(dao.SuccessResponse[dao.Stats]{
		Data: stats,
	})
}
