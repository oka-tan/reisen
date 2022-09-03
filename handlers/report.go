package handlers

import (
	"context"
	"net/http"
	"reisen/config"
	"reisen/db"
	"strconv"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/uptrace/bun"
)

func ReportGET(db *bun.DB, conf config.Config) func(echo.Context) error {
	return func(c echo.Context) error {
		board := c.Param("board")
		postNumberString := c.Param("post_number")
		postNumber, _ := strconv.ParseInt(postNumberString, 10, 64)

		model := map[string]interface{}{
			"board":      board,
			"postNumber": postNumber,
			"conf":       conf.TemplateConfig,
		}

		return c.Render(http.StatusOK, "report", model)
	}
}

func ReportPOST(pg *bun.DB, conf config.Config) func(echo.Context) error {
	return func(c echo.Context) error {
		count, err := pg.NewSelect().
			Model(&db.Report{}).
			Where("user_ip = ?", c.RealIP()).
			Where("created_at > ?", time.Now().Add(-time.Hour)).
			Where("NOT handled").
			Count(context.Background())

		if count >= 5 {
			model := map[string]interface{}{
				"conf":    conf.TemplateConfig,
				"message": "You can't make more than five reports per hour",
			}

			return c.Render(http.StatusBadRequest, "report-result", model)
		}

		board := c.Param("board")

		postNumberString := c.Param("post_number")
		postNumber, _ := strconv.ParseInt(postNumberString, 10, 64)

		userIP := c.RealIP()
		reportType := c.FormValue("reportType")
		comment := c.FormValue("comment")
		createdAt := time.Now()

		report := db.Report{
			Board:      board,
			PostNumber: postNumber,
			UserIP:     userIP,
			ReportType: reportType,
			Comment:    comment,
			CreatedAt:  createdAt,
			Handled:    false,
		}

		_, err = pg.NewInsert().
			Model(&report).
			On("CONFLICT (board, post_number, user_ip) DO UPDATE SET comment = EXCLUDED.comment, report_type = EXCLUDED.report_type, created_at = EXCLUDED.created_at").
			Exec(context.Background())

		if err != nil {
			model := map[string]interface{}{
				"conf":    conf.TemplateConfig,
				"message": "Internal server error inserting report",
			}

			return c.Render(http.StatusInternalServerError, "report-result", model)
		}

		model := map[string]interface{}{
			"conf":    conf.TemplateConfig,
			"message": "Report successfully created",
		}

		return c.Render(http.StatusOK, "report-result", model)
	}
}
