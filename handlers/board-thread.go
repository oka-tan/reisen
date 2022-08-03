package handlers

import (
	"context"
	"net/http"
	"reisen/config"
	"reisen/db"
	"strconv"

	"github.com/labstack/echo/v4"
	"github.com/uptrace/bun"
)

func BoardThread(pg *bun.DB, conf config.Config) func(echo.Context) error {
	return func(c echo.Context) error {
		board := c.Param("board")
		threadNumber, _ := strconv.ParseInt(c.Param("thread_number"), 10, 64)

		var thread []*db.Post

		err := pg.NewSelect().
			Model(&thread).
			Where("board = ?", board).
			Where("thread_number = ?", threadNumber).
			Order("post_number ASC").
			Scan(context.Background())

		if err != nil {
			model := map[string]interface{}{
				"board":        board,
				"boards":       conf.Boards,
				"conf":         conf.TemplateConfig,
				"threadNumber": threadNumber,
			}

			return c.Render(http.StatusOK, "board-thread-error", model)
		}

		if len(thread) == 0 {
			model := map[string]interface{}{
				"board":        board,
				"boards":       conf.Boards,
				"conf":         conf.TemplateConfig,
				"threadNumber": threadNumber,
			}

			return c.Render(http.StatusOK, "board-thread-not-found", model)
		}

		model := map[string]interface{}{
			"board":        board,
			"boards":       conf.Boards,
			"conf":         conf.TemplateConfig,
			"threadNumber": threadNumber,
			"op":           thread[0],
			"thread":       thread[1:],
		}

		return c.Render(http.StatusOK, "board-thread", model)
	}
}
