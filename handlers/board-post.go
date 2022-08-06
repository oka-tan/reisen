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

func BoardPost(pg *bun.DB, conf config.Config) func(echo.Context) error {
	return func(c echo.Context) error {
		board := c.Param("board")
		postNumberString := c.Param("post_number")

		postNumber, err := strconv.ParseInt(postNumberString, 10, 64)

		if err != nil {
			return c.Render(http.StatusBadRequest, "board-post-bad-request", map[string]interface{}{
				"boards": conf.Boards,
				"conf":   conf.TemplateConfig,
				"board":  board,
			})
		}

		post := db.Post{}

		err = pg.NewSelect().
			Model(&post).
			Where("board = ?", board).
			Where("post_number = ?", postNumber).
			Scan(context.Background())

		if err != nil {
			return c.Render(http.StatusNotFound, "board-post-not-found", map[string]interface{}{
				"boards":     conf.Boards,
				"conf":       conf.TemplateConfig,
				"board":      board,
				"postNumber": postNumber,
			})
		}

		return c.Redirect(http.StatusMovedPermanently, fmt.Sprintf("/%s/thread/%d#p%d", board, post.ThreadNumber, post.PostNumber))
	}
}
