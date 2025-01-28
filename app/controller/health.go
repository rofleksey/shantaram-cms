package controller

import (
	"github.com/gofiber/fiber/v2"
)

type Health struct{}

func NewHealth() *Health {
	return &Health{}
}

func (c *Health) Health(ctx *fiber.Ctx) error {
	return ctx.SendStatus(fiber.StatusOK)
}
