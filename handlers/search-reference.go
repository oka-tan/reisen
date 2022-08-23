package handlers

import (
	"net/http"
	"reisen/config"

	"github.com/labstack/echo/v4"
	"github.com/uptrace/bun"
)

func SearchReference(pg *bun.DB, conf config.Config) func(c echo.Context) error {
	return func(c echo.Context) error {
		model := map[string]interface{}{
			"boards":      conf.Boards,
			"conf":        conf.TemplateConfig,
			"title":       "Search Reference",
			"description": "Search Reference",
		}

		return c.Render(http.StatusOK, "search-reference", model)
	}
}
