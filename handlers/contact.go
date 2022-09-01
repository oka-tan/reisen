package handlers

import (
	"net/http"
	"reisen/config"

	"github.com/labstack/echo/v4"
	"github.com/uptrace/bun"
)

//Contact is the endpoint at /contact containing
//some data on contacting the archive owner
func Contact(pg *bun.DB, conf config.Config) func(c echo.Context) error {
	return func(c echo.Context) error {
		model := map[string]interface{}{
			"boards":      conf.Boards,
			"conf":        conf.TemplateConfig,
			"title":       "Contact",
			"description": "Contact",
		}

		return c.Render(http.StatusOK, "contact", model)
	}
}
