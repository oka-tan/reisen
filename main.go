package main

import (
	"database/sql"
	"reisen/config"
	"reisen/handlers"
	"reisen/lnx"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	gommon_log "github.com/labstack/gommon/log"
	"github.com/uptrace/bun"
	"github.com/uptrace/bun/dialect/pgdialect"
	"github.com/uptrace/bun/driver/pgdriver"
	"golang.org/x/crypto/acme/autocert"
)

func main() {
	conf := config.LoadConfig()

	sqlpg := sql.OpenDB(pgdriver.NewConnector(pgdriver.WithDSN(conf.PostgresConfig.ConnectionString)))
	pg := bun.NewDB(sqlpg, pgdialect.New())

	lnxService := lnx.NewService(conf.LnxConfig.Host, conf.LnxConfig.Port)

	e := echo.New()
	e.Renderer = NewTemplater()

	e.Use(middleware.SecureWithConfig(middleware.SecureConfig{
		XSSProtection:         "",
		XFrameOptions:         "",
		ContentSecurityPolicy: conf.CspConfig,
	}))
	e.Use(middleware.RecoverWithConfig(middleware.RecoverConfig{
		StackSize: 1 << 10,
		LogLevel:  gommon_log.ERROR,
	}))
	e.Use(middleware.Gzip())

	e.File("/favicon.ico", "static/favicon.ico", CacheControlHeader)
	e.Group("/static", CacheControlHeader).Static("/", "static")

	e.GET("/", handlers.Index(pg, conf))
	e.GET("/contact", handlers.Contact(pg, conf))
	e.GET("/search-reference", handlers.SearchReference(pg, conf))
	e.GET("/:board", handlers.BoardCatalogVariant(pg, conf))
	e.GET("/:board/thread/:thread_number", handlers.BoardThread(pg, conf))
	e.GET("/:board/search", handlers.BoardSearch(pg, lnxService, conf))
	e.GET("/:board/post/:post_number", handlers.BoardPost(pg, conf))
	e.GET("/:board/view-same/:media_4chan_hash", handlers.BoardViewSame(pg, conf))

	if conf.Production {
		e.AutoTLSManager.HostPolicy = autocert.HostWhitelist(conf.Hosts...)
		e.AutoTLSManager.Cache = autocert.DirCache("/var/www/.cache")
		e.Pre(middleware.HTTPSNonWWWRedirect())

		go func() {
			e2 := echo.New()
			e2.Pre(middleware.HTTPSNonWWWRedirect())
			e2.Pre(middleware.HTTPSRedirect())

			e2.Logger.Fatal(e2.Start(":80"))
		}()

		e.Logger.Fatal(e.StartAutoTLS(":443"))
	} else {
		e.Logger.Fatal(e.Start(":1323"))
	}
}

func CacheControlHeader(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		c.Response().Header().Set(echo.HeaderCacheControl, "max-age=604800")
		return next(c)
	}
}
