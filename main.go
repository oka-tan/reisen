package main

import (
	"database/sql"
	"fmt"
	"net/http"
	"reisen/config"
	"reisen/handlers"
	"reisen/lnx"
	"strings"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	gommon_log "github.com/labstack/gommon/log"
	"github.com/uptrace/bun"
	"github.com/uptrace/bun/dialect/pgdialect"
	"github.com/uptrace/bun/driver/pgdriver"
)

func main() {
	conf := config.LoadConfig()

	//Generate robots.txt
	robots(conf)

	//db connection
	sqlpg := sql.OpenDB(pgdriver.NewConnector(pgdriver.WithDSN(conf.PostgresConfig.ConnectionString)))
	pg := bun.NewDB(sqlpg, pgdialect.New())

	//search engine interface
	lnxService := lnx.NewService(conf.LnxConfig.Host, conf.LnxConfig.Port)

	e := echo.New()

	//inject the templating engine
	e.Renderer = newTemplater()

	e.Use(middleware.SecureWithConfig(middleware.SecureConfig{
		XSSProtection:         "",
		XFrameOptions:         "",
		ContentSecurityPolicy: conf.CspConfig,
	}))

	//panic recovery
	e.Use(middleware.RecoverWithConfig(middleware.RecoverConfig{
		StackSize: 1 << 10,
		LogLevel:  gommon_log.ERROR,
	}))

	//GZIP compression
	if conf.ForceGzip {
		e.Use(forceGzip)
	}

	e.Use(middleware.Gzip())

	//If we're forcing gzip we can trim the "Vary: Accept-Encoding" header
	if conf.ForceGzip {
		e.Use(trimVaryHeader)
	}

	e.File("/favicon.ico", "static/favicon.ico", cacheControlHeader("public, immutable, max-age=604800"))
	e.File("/robots.txt", "static/robots.txt")
	e.Group("/static", cacheControlHeader("public, immutable, max-age=604800")).Static("/", "static")

	e.GET("/", handlers.Index(pg, conf), cacheControlHeader("public, immutable, max-age=60"))
	e.GET("/contact", handlers.Contact(pg, conf), cacheControlHeader("public, immutable, max-age=60"))
	e.GET("/search-reference", handlers.SearchReference(pg, conf), cacheControlHeader("public, immutable, max-age=60"))
	e.GET("/:board", handlers.BoardCatalogVariant(pg, conf))
	e.GET("/:board/thread/:thread_number", handlers.BoardThread(pg, conf))
	e.GET("/:board/search", handlers.BoardSearch(pg, lnxService, conf), cacheControlHeader("no-store"))
	e.GET("/:board/post/:post_number", handlers.BoardPost(pg, conf))
	e.GET("/:board/view-same/:media_4chan_hash", handlers.BoardViewSame(pg, conf), cacheControlHeader("no-store"))
	e.GET("/:board/mooch-image/:post_number", handlers.MoochImage(pg, conf))

	e.GET("/:board/report/:post_number", handlers.ReportGET(pg, conf), cacheControlHeader("no-store"))
	e.POST("/:board/report/:post_number", handlers.ReportPOST(pg, conf), cacheControlHeader("no-store"))

	e.Logger.Fatal(e.Start(fmt.Sprintf(":%d", conf.Port)))
}

func cacheControlHeader(s string) echo.MiddlewareFunc {
	return func(cacheHeader string) echo.MiddlewareFunc {
		return func(next echo.HandlerFunc) echo.HandlerFunc {
			return func(c echo.Context) error {
				c.Response().Header().Set(echo.HeaderCacheControl, s)
				return next(c)
			}
		}
	}(s)
}

//Rejects requests without the Accept-Encoding header
//explicitly saying it accepts GZIP
func forceGzip(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		if !strings.Contains(c.Request().Header.Get(echo.HeaderAcceptEncoding), "gzip") {
			return c.String(http.StatusBadRequest, "Requests must accept GZIP encoding")
		}

		return next(c)
	}
}

//Removes the 'Vary' header from responses
func trimVaryHeader(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		c.Response().Header().Del(echo.HeaderVary)
		return next(c)
	}
}
