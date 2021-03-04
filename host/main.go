package main

import (
	"github.com/kataras/iris/v12"
	"github.com/kataras/iris/v12/core/router"
	"github.com/kataras/iris/v12/middleware/logger"
	"github.com/kataras/iris/v12/middleware/recover"
	"github.com/syncfuture/go/sconfig"
	log "github.com/syncfuture/go/slog"
	"github.com/syncfuture/scraper/store"
	"github.com/syncfuture/scraper/store/webshare"
)

var (
	_proxyStore store.IProxyStore
	_cp         sconfig.IConfigProvider
	_listenAddr string
	_debug      bool
)

func init() {
	_cp = sconfig.NewJsonConfigProvider()
	_debug = _cp.GetBool("Debug")
	log.Init(_cp)
	_proxyStore = webshare.NewDefaultWebShareProxyStore()
	_listenAddr = _cp.GetStringDefault("ListenAddr", ":7000")
}

func main() {
	app := iris.New()
	logLevel := _cp.GetStringDefault("Log.Level", "info")
	app.Logger().SetLevel(logLevel)
	app.Use(recover.New())
	app.Use(logger.New())

	var api router.Party

	if _debug {
		// Debug mode
		app.HandleDir("/", "./dist")
		crs := func(ctx iris.Context) {
			ctx.Header("Access-Control-Allow-Origin", "*")
			ctx.Header("Access-Control-Allow-Credentials", "true")
			ctx.Header("Access-Control-Allow-Methods", "DELETE")
			ctx.Header("Access-Control-Allow-Headers", "Access-Control-Allow-Origin, Content-Type, x-requested-with")
			ctx.Next()
		}

		api = app.Party("/api", crs).AllowMethods(iris.MethodOptions)
	} else {
		// Production mode
		app.HandleDir("/", "./dist", iris.DirOptions{
			Asset:      Asset,
			AssetInfo:  AssetInfo,
			AssetNames: AssetNames,
		})
		api = app.Party("/api")
	}

	amazonAPI := api.Party("/amazon")
	amazonHttpHandler := NewAmazonHttpHandlers(_cp, _proxyStore)
	amazonAPI.Get("/reviews", amazonHttpHandler.GetReviews)
	amazonAPI.Post("/reviews/export", amazonHttpHandler.ExportReviews)
	amazonAPI.Get("/items", amazonHttpHandler.GetItems)
	amazonAPI.Post("/scrape", amazonHttpHandler.PostScrape)

	wayfairAPI := api.Party("/wayfair")
	wayfairHttpHandler := NewWayfairHttpHandlers(_cp, _proxyStore)
	wayfairAPI.Get("/reviews", wayfairHttpHandler.GetReviews)
	wayfairAPI.Post("/reviews/export", wayfairHttpHandler.ExportReviews)
	wayfairAPI.Get("/items", wayfairHttpHandler.GetItems)
	wayfairAPI.Post("/scrape", wayfairHttpHandler.PostScrape)

	app.Run(iris.Addr(_listenAddr))
}
