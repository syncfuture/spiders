package main

import (
	"encoding/json"
	"fmt"
	"strconv"
	"sync"
	"sync/atomic"
	"time"

	"github.com/kataras/iris/v12"
	"github.com/kataras/iris/v12/core/router"
	"github.com/kataras/iris/v12/middleware/logger"
	"github.com/kataras/iris/v12/middleware/recover"
	"github.com/olivere/elastic/v7"
	"github.com/syncfuture/go/sconfig"
	log "github.com/syncfuture/go/slog"
	"github.com/syncfuture/go/spool"
	"github.com/syncfuture/go/stask"
	"github.com/syncfuture/go/u"
	"github.com/syncfuture/scraper/store"
	"github.com/syncfuture/scraper/store/webshare"
	"github.com/syncfuture/spiders/amazon"
	"github.com/syncfuture/spiders/amazon/dal"
	"github.com/syncfuture/spiders/amazon/dal/es"
	"github.com/tealeg/xlsx"
)

const (
	_key = "be09d781115fe3491743fa205ea786852513f474"
)

var (
	_reviewDAL     dal.IReviewDAL
	_itemDAL       dal.IItemDAL
	_scrapeLocker  = new(sync.Mutex)
	_store         store.IProxyStore
	_cp            sconfig.IConfigProvider
	_listenAddr    string
	_maxConcurrent int
	_debug         bool
	_bufferPool    = spool.NewSyncBufferPool(4096)
)

func init() {
	_cp = sconfig.NewJsonConfigProvider()
	_debug = _cp.GetBool("Debug")
	log.Init(_cp)
	_store = webshare.NewWebShareProxyStore(_key)
	_listenAddr = _cp.GetStringDefault("ListenAddr", ":7000")
	_maxConcurrent = _cp.GetIntDefault("MaxConcurrent", 15)
}

func main() {
	addrs := _cp.GetStringSlice("ES.Addrs")

	var err error
	_itemDAL, err = es.NewESItemDAL(
		elastic.SetURL(addrs...),
	)
	if u.LogError(err) {
		return
	}

	_reviewDAL, err = es.NewESReviewDAL(
		elastic.SetURL(addrs...),
	)
	if u.LogError(err) {
		return
	}

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

	api.Get("/reviews", getReviews)
	api.Post("/reviews/export", exportReviews)
	api.Get("/items", getItems)
	api.Post("/scrape", scrape)

	app.Run(iris.Addr(_listenAddr))
}

func getReviews(ctx iris.Context) {
	query := getReviewQuery(ctx)
	result, err := _reviewDAL.GetAllReviews(query)
	if u.LogError(err) {
		ctx.WriteString(err.Error())
		return
	}

	json, err := json.Marshal(result.Reviews)
	if u.LogError(err) {
		ctx.WriteString(err.Error())
		return
	}

	ctx.ContentType("application/json; charset=utf-8")
	ctx.Write(json)
}

func getItems(ctx iris.Context) {
	query := getItemQuery(ctx)
	items, err := _itemDAL.GetAllItems(query)
	if u.LogError(err) {
		ctx.WriteString(err.Error())
		return
	}

	json, err := json.Marshal(items)
	if u.LogError(err) {
		ctx.WriteString(err.Error())
		return
	}

	ctx.ContentType("application/json; charset=utf-8")
	ctx.Write(json)
}

func scrape(ctx iris.Context) {
	_scrapeLocker.Lock()
	defer _scrapeLocker.Unlock()

	query := getItemQuery(ctx)
	result, err := _itemDAL.GetAllItems(query)
	if u.LogError(err) {
		ctx.WriteString(err.Error())
		return
	}

	count := int32(0)
	fromDate := time.Now().AddDate(0, -1, 0) // 一个月内的评论

	f := stask.NewFlowScheduler(_maxConcurrent)
	f.SliceRun(&result.Items, func(i int, v interface{}) {
		item := v.(*amazon.ItemDTO)

		atomic.AddInt32(&count, 1)

		s := amazon.NewReviewsScraper(_store, item.ASIN)
		reviews, err := s.FetchPages(&fromDate)
		if u.LogError(err) {
			item.Status = -1
			_itemDAL.SaveItems(item)
			return
		}

		if len(reviews) > 0 { // 有评论才存储
			// 关联E&E ItemNo
			for _, review := range reviews {
				review.CustomerNo = item.ItemNo
			}

			err = _reviewDAL.SaveReviews(reviews)
			if u.LogError(err) {
				item.Status = -1
				_itemDAL.SaveItems(item)
				return
			}
		}

		item.Status = 1
		_itemDAL.SaveItems(item)
	})

	// err = _itemDAL.SaveItems(result.Items...)
	// if u.LogError(err) {
	// 	return
	// }

	ctx.ContentType("application/json; charset=utf-8")
	json := fmt.Sprintf(`{"count":%d}`, count)
	ctx.WriteString(json)
}

func getItemQuery(ctx iris.Context) *amazon.ItemQuery {
	asin := string(ctx.FormValue("asin"))
	itemNo := string(ctx.FormValue("itemNo"))
	statusStr := string(ctx.FormValue("status"))
	pageSizeStr := string(ctx.FormValue("pageSize"))
	cursor := string(ctx.FormValue("cursor"))
	pageSize, err := strconv.Atoi(pageSizeStr)
	if err != nil {
		pageSize = 10000
	}

	return &amazon.ItemQuery{
		Status:   statusStr,
		PageSize: pageSize,
		Cursor:   cursor,
		ASIN:     asin,
		ItemNo:   itemNo,
	}
}

func getReviewQuery(ctx iris.Context) *amazon.ReviewQuery {
	asin := string(ctx.FormValue("asin"))
	itemNo := string(ctx.FormValue("itemNo"))
	pageSizeStr := string(ctx.FormValue("pageSize"))
	cursor := string(ctx.FormValue("cursor"))
	fromDate := string(ctx.FormValue("fromDate"))
	pageSize, err := strconv.Atoi(pageSizeStr)
	if err != nil {
		pageSize = 10000
	}

	return &amazon.ReviewQuery{
		PageSize: pageSize,
		Cursor:   cursor,
		ASIN:     asin,
		ItemNo:   itemNo,
		FromDate: fromDate,
	}
}

func exportReviews(ctx iris.Context) {
	query := getReviewQuery(ctx)
	result, err := _reviewDAL.GetAllReviews(query)
	if u.LogError(err) {
		ctx.WriteString(err.Error())
		return
	}

	wb := xlsx.NewFile()
	sheet, err := wb.AddSheet("Reviews")
	if u.LogError(err) {
		ctx.WriteString(err.Error())
		return
	}
	header := sheet.AddRow()
	header.AddCell().Value = "AmazonID"
	header.AddCell().Value = "ASIN"
	header.AddCell().Value = "ItemNo"
	header.AddCell().Value = "StripInfo"
	header.AddCell().Value = "Rating"
	header.AddCell().Value = "IsVerified"
	header.AddCell().Value = "Location"
	header.AddCell().Value = "CustomerName"
	header.AddCell().Value = "CreatedOn"
	header.AddCell().Value = "Title"
	header.AddCell().Value = "Content"

	for _, review := range result.Reviews {
		row := sheet.AddRow()
		row.AddCell().Value = review.AmazonID
		row.AddCell().Value = review.ASIN
		row.AddCell().Value = review.CustomerNo
		row.AddCell().Value = review.StripInfo
		row.AddCell().Value = strconv.Itoa(int(review.Rating))
		row.AddCell().Value = strconv.FormatBool(review.IsVerified)
		row.AddCell().Value = review.Location
		row.AddCell().Value = review.CustomerName
		row.AddCell().Value = review.CreatedOn.Format("01/02/2006")
		row.AddCell().Value = review.Title
		row.AddCell().Value = review.Content
	}

	buffer := _bufferPool.GetBuffer()
	defer _bufferPool.PutBuffer(buffer)

	wb.Write(buffer)

	ctx.ContentType("application/octet-stream")
	ctx.Header("Content-Disposition", fmt.Sprintf("attachment;filename=%s.xlsx", time.Now().Format("20060102-150405")))
	ctx.Write(buffer.Bytes())
}
