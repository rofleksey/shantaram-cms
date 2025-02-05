package controller

import (
	"fmt"
	"github.com/gofiber/fiber/v2"
	"net/http"
	"shantaram-cms/app/dao"
	"shantaram-cms/app/service"
)

type Captcha struct {
	captchaService *service.Captcha
}

func NewCaptcha(
	captchaService *service.Captcha,
) *Captcha {
	return &Captcha{
		captchaService: captchaService,
	}
}

func (c *Captcha) Generate(ctx *fiber.Ctx) error {
	id, image, err := c.captchaService.Generate()
	if err != nil {
		return ctx.Status(http.StatusInternalServerError).JSON(dao.NoDataResponse{
			Error: true,
			Msg:   fmt.Sprintf("не удалось получить капчу: %v", err),
		})
	}

	return ctx.JSON(dao.SuccessResponse[dao.CaptchaGenerateResponse]{
		Data: dao.CaptchaGenerateResponse{
			ID:    id,
			Image: image,
		},
	})
}
