package handlers

import (
	"context"
	"encoding/base64"
	"fmt"
	"net/http"
	"reisen/config"
	"reisen/db"
	"strconv"

	"github.com/labstack/echo/v4"
	"github.com/uptrace/bun"
)

//MoochImage is an endpoint for image mooching by other archives
func MoochImage(pg *bun.DB, conf config.Config) func(echo.Context) error {
	return func(c echo.Context) error {
		board := c.Param("board")
		postNumberString := c.Param("post_number")

		postNumber, err := strconv.ParseInt(postNumberString, 10, 64)

		if err != nil {
			return c.String(http.StatusBadRequest, "Bad Request")
		}

		post := db.Post{}

		err = pg.NewSelect().
			Model(&post).
			Where("board = ?", board).
			Where("post_number = ?", postNumber).
			Scan(context.Background())

		if err != nil {
			return c.String(http.StatusNotFound, "Post not found")
		}

		if post.MediaInternalHash == nil || post.Hidden {
			return c.String(http.StatusNotFound, "Image not available")
		}

		return c.Redirect(
			http.StatusMovedPermanently,
			fmt.Sprintf("%s/%s", conf.TemplateConfig.ThumbnailsUrl, base64.URLEncoding.EncodeToString(*(post.MediaInternalHash))),
		)
	}
}
