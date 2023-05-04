package handlers

import (
	"context"
	"encoding/base64"
	"net/http"
	"reisen/config"
	"reisen/db"
	"strconv"

	"github.com/labstack/echo/v4"
	"github.com/uptrace/bun"
)

// BoardViewSame handles requests at the /:board/view-same/:hash endpoint
// and lists out paginated posts with the given hash.
func BoardViewSame(pg *bun.DB, conf config.Config) func(echo.Context) error {
	return func(c echo.Context) error {
		board := c.Param("board")
		media4chanHashString := c.Param("media_4chan_hash")
		media4chanHash, _ := base64.URLEncoding.DecodeString(media4chanHashString)

		var posts []*db.Post

		q := pg.NewSelect().
			Model(&posts).
			Where("media_4chan_hash = ?", media4chanHash).
			Where("board = ?", board).
			Where("NOT hidden")

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
			return c.Render(http.StatusInternalServerError, "board-view-same-error", map[string]interface{}{
				"boards": conf.Boards,
				"conf":   conf.TemplateConfig,
				"board":  board,
				"title":  "Error",
			})
		}

		if len(posts) == 0 {
			return c.Render(http.StatusNotFound, "board-view-same-empty", map[string]interface{}{
				"boards": conf.Boards,
				"conf":   conf.TemplateConfig,
				"board":  board,
				"title":  "Error",
			})
		}

		if keyset > 0 {
			rkeyset = posts[0].PostNumber

			if len(posts) == 25 {
				keyset = posts[23].PostNumber
				posts = posts[:24]
			} else {
				keyset = 0
			}
		} else if rkeyset > 0 {
			keyset = posts[len(posts)-1].PostNumber
			if len(posts) == 25 {
				rkeyset = posts[1].PostNumber
				posts = posts[1:]
			} else {
				rkeyset = 0
			}
		} else {
			rkeyset = 0
			if len(posts) == 25 {
				keyset = posts[23].PostNumber
				posts = posts[:24]
			} else {
				keyset = 0
			}
		}

		model := map[string]interface{}{
			"boards":             conf.Boards,
			"conf":               conf.TemplateConfig,
			"board":              board,
			"posts":              posts,
			"keyset":             keyset,
			"rkeyset":            rkeyset,
			"title":              "View Same",
			"enableLatex":        conf.IsLatexEnabled(board),
			"enableTegaki":       conf.IsTegakiEnabled(board),
			"enableCountryFlags": conf.AreCountryFlagsEnabled(board),
			"enablePolFlags":     conf.ArePolFlagsEnabled(board),
			"enableMlpFlags":     conf.AreMlpFlagsEnabled(board),
			"media4chanHash":     media4chanHashString,
		}

		return c.Render(http.StatusOK, "board-view-same", model)
	}
}
