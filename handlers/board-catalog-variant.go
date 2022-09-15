//Package handlers contains all request handlers in reisen.
//
//Since golang dislikes dependency injection
//they're all functions taking a couple of dependencies (
//the db connection, etc) and returning the handler function proper.
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

//BoardCatalogVariant is the catalog-like variant
//of the index page.
//The "default" variant hasn't been made yet
func BoardCatalogVariant(pg *bun.DB, conf config.Config) func(echo.Context) error {
	return func(c echo.Context) error {
		board := c.Param("board")

		var threads []*db.Post

		q := pg.NewSelect().
			Model(&threads).
			Where("op").
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
			c.Response().Header().Set(echo.HeaderCacheControl, "no-store")
			return c.Render(http.StatusInternalServerError, "board-error", map[string]interface{}{
				"boards":  conf.Boards,
				"conf":    conf.TemplateConfig,
				"board":   board,
				"noIndex": true,
			})
		}

		if len(threads) == 0 {
			c.Response().Header().Set(echo.HeaderCacheControl, "no-store")

			return c.Render(http.StatusNoContent, "board-empty", map[string]interface{}{
				"boards":  conf.Boards,
				"conf":    conf.TemplateConfig,
				"board":   board,
				"noIndex": true,
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
			"boards":             conf.Boards,
			"conf":               conf.TemplateConfig,
			"board":              board,
			"threads":            threads,
			"keyset":             keyset,
			"rkeyset":            rkeyset,
			"enableLatex":        conf.IsLatexEnabled(board),
			"enableTegaki":       conf.IsTegakiEnabled(board),
			"enableCountryFlags": conf.AreCountryFlagsEnabled(board),
			"enablePolFlags":     conf.ArePolFlagsEnabled(board),
			"enableMlpFlags":     conf.AreMlpFlagsEnabled(board),
			"noIndex":            true,
		}

		//Should take around this long for responses to be refreshed on the db in the first place
		c.Response().Header().Set(echo.HeaderCacheControl, "public, max-age=600, immutable")

		return c.Render(http.StatusOK, "board-catalog-variant", model)
	}
}
