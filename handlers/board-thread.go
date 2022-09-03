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

//BoardThread prints out a thread
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

			return c.Render(http.StatusInternalServerError, "board-thread-error", model)
		}

		if len(thread) == 0 {
			model := map[string]interface{}{
				"board":        board,
				"boards":       conf.Boards,
				"conf":         conf.TemplateConfig,
				"threadNumber": threadNumber,
				"noIndex":      true,
			}

			return c.Render(http.StatusNotFound, "board-thread-not-found", model)
		}

		var title string
		var description string
		if thread[0].Subject != nil {
			subject := *thread[0].Subject

			title = fmt.Sprintf("%d - %s", threadNumber, subject)
			description = fmt.Sprintf("/%s/ %d - %s", board, threadNumber, subject)
		} else {
			title = fmt.Sprintf("%d - Untitled Thread", threadNumber)
			description = fmt.Sprintf("/%s/ %d - Untitled Thread", board, threadNumber)
		}

		model := map[string]interface{}{
			"board":              board,
			"boards":             conf.Boards,
			"conf":               conf.TemplateConfig,
			"threadNumber":       threadNumber,
			"op":                 thread[0],
			"thread":             thread[1:],
			"enableLatex":        conf.IsLatexEnabled(board),
			"enableTegaki":       conf.IsTegakiEnabled(board),
			"enableCountryFlags": conf.AreCountryFlagsEnabled(board),
			"enableBoardFlags":   conf.AreBoardFlagsEnabled(board),
			"title":              title,
			"description":        description,
		}

		return c.Render(http.StatusOK, "board-thread", model)
	}
}

func truncate(s string, max int) string {
	if len(s) < max {
		return s
	}
	return s[:max]
}
