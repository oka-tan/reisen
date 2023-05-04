package handlers

import (
	"net/http"
	"reisen/config"

	"github.com/labstack/echo/v4"
	"github.com/uptrace/bun"
)

// About is the endpoint at /about containing
// some data on contacting the archive owner and stuff
func About(pg *bun.DB, conf config.Config) func(c echo.Context) error {
	return func(c echo.Context) error {
		model := map[string]interface{}{
			"boards":      conf.Boards,
			"conf":        conf.TemplateConfig,
			"title":       "About",
			"description": "About",
		}

		return c.Render(http.StatusOK, "about", model)
	}
}
