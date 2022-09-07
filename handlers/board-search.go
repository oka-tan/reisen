package handlers

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"reisen/config"
	"reisen/db"
	"reisen/lnx"
	"strconv"
	"strings"

	"github.com/labstack/echo/v4"
	"github.com/uptrace/bun"
)

//BoardSearch is the regular search endpoint
func BoardSearch(pg *bun.DB, lnxService lnx.Service, conf config.Config) func(echo.Context) error {
	return func(c echo.Context) error {
		board := c.Param("board")
		search := strings.TrimSpace(c.QueryParam("search"))

		postNumber, err := strconv.ParseInt(search, 10, 64)

		if err == nil {
			post := db.Post{}

			err := pg.NewSelect().
				Model(&post).
				Where("board = ?", board).
				Where("post_number = ?", postNumber).
				Where("NOT hidden").
				Scan(context.Background())

			if err != nil {
				return c.Render(http.StatusNotFound, "board-search-thread-not-found", map[string]interface{}{
					"boards":     conf.Boards,
					"conf":       conf.TemplateConfig,
					"board":      board,
					"postNumber": postNumber,
				})
			}

			return c.Redirect(http.StatusMovedPermanently, fmt.Sprintf("/%s/thread/%d", board, post.ThreadNumber))
		}

		offset64, _ := strconv.ParseInt(c.QueryParam("offset"), 10, 32)
		offset := int(offset64)

		searchResult, err := lnxService.Search(board, search, offset)

		if err != nil {
			log.Printf("Error searching Lnx: %v\n", err)

			return c.Render(http.StatusBadRequest, "board-search-bad-request", map[string]interface{}{
				"boards": conf.Boards,
				"conf":   conf.TemplateConfig,
				"board":  board,
				"search": search,
			})
		}

		if len(searchResult.Hits) == 0 {
			return c.Render(http.StatusBadRequest, "board-search-no-results", map[string]interface{}{
				"boards": conf.Boards,
				"conf":   conf.TemplateConfig,
				"board":  board,
				"search": search,
			})
		}

		posts := make([]*db.Post, 0, 20)

		err = pg.NewSelect().
			Model(&posts).
			Where("board = ?", board).
			Where("post_number IN (?)", bun.In(searchResult.Hits)).
			Order("post_number DESC").
			Limit(20).
			Scan(context.Background())

		if err != nil {
			return c.Render(http.StatusInternalServerError, "board-search-server-error", map[string]interface{}{
				"boards": conf.Boards,
				"conf":   conf.TemplateConfig,
				"board":  board,
				"search": search,
			})
		}

		nextPageOffset := 0
		if searchResult.Count > offset+20 {
			nextPageOffset = offset + 20
		}

		previousPageOffset := 0
		if offset > 0 {
			previousPageOffset = offset - 20
		}

		return c.Render(http.StatusOK, "board-search", map[string]interface{}{
			"boards":             conf.Boards,
			"conf":               conf.TemplateConfig,
			"board":              board,
			"search":             search,
			"posts":              posts,
			"count":              searchResult.Count,
			"countMoreThanOne":   searchResult.Count > 1,
			"offset":             offset,
			"nextPageOffset":     nextPageOffset,
			"previousPageOffset": previousPageOffset,
			"enableLatex":        conf.IsLatexEnabled(board),
			"enableTegaki":       conf.IsTegakiEnabled(board),
			"enableCountryFlags": conf.AreCountryFlagsEnabled(board),
			"enableBoardFlags":   conf.AreBoardFlagsEnabled(board),
		})
	}
}
