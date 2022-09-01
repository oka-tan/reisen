package main

import (
	"database/sql"
	"fmt"
	"reisen/config"
	"reisen/handlers"
	"reisen/lnx"

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

	e.Use(middleware.Gzip())

	e.File("/favicon.ico", "static/favicon.ico", cacheControlHeader)
	e.File("/robots.txt", "static/robots.txt")
	e.Group("/static", cacheControlHeader).Static("/", "static")

	e.GET("/", handlers.Index(pg, conf))
	e.GET("/contact", handlers.Contact(pg, conf))
	e.GET("/search-reference", handlers.SearchReference(pg, conf))
	e.GET("/:board", handlers.BoardCatalogVariant(pg, conf))
	e.GET("/:board/thread/:thread_number", handlers.BoardThread(pg, conf))
	e.GET("/:board/search", handlers.BoardSearch(pg, lnxService, conf))
	e.GET("/:board/post/:post_number", handlers.BoardPost(pg, conf))
	e.GET("/:board/view-same/:media_4chan_hash", handlers.BoardViewSame(pg, conf))
	e.GET("/:board/mooch-image/:post_number", handlers.MoochImage(pg, conf))

	e.Logger.Fatal(e.Start(fmt.Sprintf(":%d", conf.Port)))
}

func cacheControlHeader(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		c.Response().Header().Set(echo.HeaderCacheControl, "max-age=604800")
		return next(c)
	}
}
