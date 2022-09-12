package handlers

import (
	"context"
	"fmt"
	"net/http"
	"reisen/config"
	"reisen/db"
	"strconv"
	"time"

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
			Where("NOT hidden").
			Order("post_number ASC").
			Scan(context.Background())

		if err != nil {
			c.Response().Header().Set("Cache-Control", "no-store")

			model := map[string]interface{}{
				"board":        board,
				"boards":       conf.Boards,
				"conf":         conf.TemplateConfig,
				"threadNumber": threadNumber,
				"noIndex":      true,
			}

			return c.Render(http.StatusInternalServerError, "board-thread-error", model)
		}

		/*
		 * The second condition translates to either "the
		 * OP isn't available in the db" or "the OP is hidden"
		 */
		if len(thread) == 0 || !thread[0].Op {
			c.Response().Header().Set("Cache-Control", "no-store")

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

		//360 days
		if time.Now().Sub(thread[0].LastModified) > 24*8640*time.Hour {
			//2 weeks
			c.Response().Header().Set("Cache-Control", "max-age=604800, public, immutable")
		} else {
			//10 minutes
			c.Response().Header().Set("Cache-Control", "max-age=600, public, immutable")
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
