package controller

import (
	"encoding/json"
	"encoding/xml"
	"fmt"
	"github.com/gofiber/fiber/v2"
	"net/http"
	"shantaram-cms/app/dao"
	"shantaram-cms/app/service"
	"shantaram-cms/pkg/database"
	"time"
)

type Sitemap struct {
	db          *database.Database
	pageService *service.Page
}

func NewSitemap(
	db *database.Database,
	pageService *service.Page,
) *Sitemap {
	return &Sitemap{
		db:          db,
		pageService: pageService,
	}
}

func (c *Sitemap) GetSitemap(ctx *fiber.Ctx) error {
	ctx.Type("xml", "utf-8")

	sitemap := dao.Sitemap{
		Xmlns: "http://www.sitemaps.org/schemas/sitemap/0.9",
		URLs:  []dao.SitemapEntry{},
	}

	pages, err := c.pageService.GetAll()
	if err != nil {
		return ctx.Status(http.StatusInternalServerError).JSON(dao.NoDataResponse{
			Error: true,
			Msg:   fmt.Sprintf("не удалось получить список страниц: %v", err),
		})
	}

	for _, page := range pages {
		priority := float32(0.5)

		if page.ID == "index" {
			priority = 1.0
		}

		sitemap.URLs = append(sitemap.URLs, dao.SitemapEntry{
			Loc:        "https://shantaram-spb.ru/p/" + page.ID,
			LastMod:    page.Updated.Format(time.RFC3339),
			ChangeFreq: "daily",
			Priority:   priority,
		})
	}

	menuBytez, err := c.db.GetGeneralByIDData("menu")
	if err != nil {
		return ctx.Status(http.StatusInternalServerError).JSON(dao.NoDataResponse{
			Error: true,
			Msg:   fmt.Sprintf("не удалось получить меню: %v", err),
		})
	}

	var menuSettings database.MenuSettings

	if err := json.Unmarshal(menuBytez, &menuSettings); err != nil {
		return ctx.Status(http.StatusInternalServerError).JSON(dao.NoDataResponse{
			Error: true,
			Msg:   fmt.Sprintf("не удалось спарсить меню: %v", err),
		})
	}

	sitemap.URLs = append(sitemap.URLs, dao.SitemapEntry{
		Loc:        "https://shantaram-spb.ru/menu",
		LastMod:    menuSettings.Updated.Format(time.RFC3339),
		ChangeFreq: "daily",
		Priority:   1.0,
	})

	tagMap := make(map[string]struct{})
	for _, group := range menuSettings.Groups {
		for _, tag := range group.Tags {
			tagMap[tag] = struct{}{}
		}
	}

	for tag := range tagMap {
		sitemap.URLs = append(sitemap.URLs, dao.SitemapEntry{
			Loc:        "https://shantaram-spb.ru/menu/" + tag,
			LastMod:    menuSettings.Updated.Format(time.RFC3339),
			ChangeFreq: "daily",
			Priority:   0.5,
		})
	}

	xmlData, err := xml.MarshalIndent(sitemap, "", "  ")
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).SendString("Error generating sitemap")
	}

	return ctx.SendString(xml.Header + string(xmlData))
}
