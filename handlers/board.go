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

func Board(pg *bun.DB, conf config.Config) func(echo.Context) error {
	return func(c echo.Context) error {
		board := c.Param("board")

		var threads []*db.Post

		q := pg.NewSelect().Model(&threads).Where("op").Where("board = ?", board)
		keyset, kerr := strconv.ParseInt(c.QueryParam("keyset"), 10, 64)
		rkeyset, rkerr := strconv.ParseInt(c.QueryParam("rkeyset"), 10, 64)

		if kerr == nil {
			q.Where("post_number < ?", keyset)
		} else if rkerr == nil {
			q.Where("post_number > ?", rkeyset)
		}

		err := q.Order("post_number DESC").
			Limit(25).
			Scan(context.Background())

		if err != nil {
			return c.Render(http.StatusOK, "board-error", map[string]interface{}{
				"boards": conf.Boards,
				"conf":   conf.TemplateConfig,
				"board":  board,
			})
		}

		if len(threads) == 0 {
			return c.Render(http.StatusOK, "board-empty", map[string]interface{}{
				"boards": conf.Boards,
				"conf":   conf.TemplateConfig,
				"board":  board,
			})
		}

		if keyset > 0 {
			rkeyset = threads[0].PostNumber

			if len(threads) == 25 {
				keyset = threads[23].PostNumber
				threads = threads[:24]
			} else {
				keyset = 0
			}
		} else if rkeyset > 0 {
			keyset = threads[len(threads)-1].PostNumber
			if len(threads) == 25 {
				rkeyset = threads[1].PostNumber
				threads = threads[1:]
			} else {
				rkeyset = 0
			}
		} else {
			rkeyset = 0
			if len(threads) == 25 {
				keyset = threads[23].PostNumber
				threads = threads[:24]
			} else {
				keyset = 0
			}
		}

		model := map[string]interface{}{
			"boards":  conf.Boards,
			"conf":    conf.TemplateConfig,
			"board":   board,
			"threads": threads,
			"keyset":  keyset,
			"rkeyset": rkeyset,
		}

		return c.Render(http.StatusOK, "board", model)
	}
}
