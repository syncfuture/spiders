package main

import (
	"github.com/syncfuture/go/sconfig"
	log "github.com/syncfuture/go/slog"
	"github.com/syncfuture/host/sfasthttp"
)

func main() {
	cp := sconfig.NewJsonConfigProvider()
	host := sfasthttp.NewFHWebHost(cp)

	amazonHttpHandler := NewAmazonHttpHandlers(cp)

	host.GET("/api/amazon/reviews", amazonHttpHandler.GetReviews)
	host.POST("/api/amazon/reviews/export", amazonHttpHandler.ExportReviews)
	host.GET("/api/amazon/items", amazonHttpHandler.GetItems)
	host.POST("/api/amazon/scrape", amazonHttpHandler.PostScrape)

	wayfairHttpHandler := NewWayfairHttpHandlers(cp)
	host.GET("/api/wayfair/reviews", wayfairHttpHandler.GetReviews)
	host.POST("/api/wayfair/reviews/export", wayfairHttpHandler.ExportReviews)
	host.GET("/api/wayfair/items", wayfairHttpHandler.GetItems)
	host.POST("/api/wayfair/scrape", wayfairHttpHandler.PostScrape)
	host.POST("/api/wayfair/scrape/cancel", wayfairHttpHandler.PostCancel)
	host.GET("/api/wayfair/scrape/status", wayfairHttpHandler.GetStatus)

	log.Fatal(host.Run())
}
