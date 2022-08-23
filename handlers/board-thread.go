package handlers

import (
	"context"
	"fmt"
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
				"noIndex":      true,
			}

			return c.Render(http.StatusOK, "board-thread-error", model)
		}

		if len(thread) == 0 {
			model := map[string]interface{}{
				"board":        board,
				"boards":       conf.Boards,
				"conf":         conf.TemplateConfig,
				"threadNumber": threadNumber,
				"noIndex":      true,
			}

			return c.Render(http.StatusOK, "board-thread-not-found", model)
		}

		var title string
		var description string
		if thread[0].Subject != nil {
			title = fmt.Sprintf("%d - %s", threadNumber, *thread[0].Subject)
			description = fmt.Sprintf("/%s/ %d - %s", board, threadNumber, *thread[0].Subject)
		} else if thread[0].Comment != nil {
			title = fmt.Sprintf("%d - %s", threadNumber, truncate(*thread[0].Comment, 20))
			description = fmt.Sprintf("/%s/ %d - %s", board, threadNumber, *thread[0].Comment)
		} else {
			title = fmt.Sprintf("%d - Untitled Thread", threadNumber)
			description = fmt.Sprintf("/%s/ %d - Untitled Thread", board, threadNumber)
		}

		model := map[string]interface{}{
			"board":        board,
			"boards":       conf.Boards,
			"conf":         conf.TemplateConfig,
			"threadNumber": threadNumber,
			"op":           thread[0],
			"thread":       thread[1:],
			"enableLatex":  conf.IsLatexEnabled(board),
			"enableTegaki": conf.IsTegakiEnabled(board),
			"title":        title,
			"description":  description,
		}

		return c.Render(http.StatusOK, "board-thread", model)
	}
}

func truncate(s string, max int) string {
	if len(s) < max {
		return s
	} else {
		return s[:max]
	}
}
