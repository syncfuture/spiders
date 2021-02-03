package main

import (
	"encoding/json"
	"fmt"
	"strconv"
	"sync"
	"sync/atomic"
	"time"

	"github.com/buaazp/fasthttprouter"
	"github.com/olivere/elastic/v7"
	"github.com/syncfuture/go/sconfig"
	log "github.com/syncfuture/go/slog"
	"github.com/syncfuture/go/stask"
	"github.com/syncfuture/go/u"
	"github.com/syncfuture/scraper/amazon"
	"github.com/syncfuture/scraper/store"
	"github.com/syncfuture/scraper/store/webshare"
	"github.com/syncfuture/spiders/amazon/dal"
	"github.com/syncfuture/spiders/amazon/dal/es"
	"github.com/syncfuture/spiders/amazon/model"
	"github.com/valyala/fasthttp"
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
)

func init() {
	_cp = sconfig.NewJsonConfigProvider()
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

	router := fasthttprouter.New()
	router.GET("/reviews", allowCORS(getReviews))
	router.POST("/reviews/export", allowCORS(exportReviews))
	router.OPTIONS("/reviews/export", allowCORS(options))

	router.GET("/items", allowCORS(getItems))

	router.POST("/scrape", allowCORS(scrape))
	router.OPTIONS("/scrape", allowCORS(options))

	log.Infof("Listen on %s", _listenAddr)
	fasthttp.ListenAndServe(_listenAddr, router.Handler)
}

func allowCORS(next fasthttp.RequestHandler) fasthttp.RequestHandler {
	return func(ctx *fasthttp.RequestCtx) {

		ctx.Response.Header.Add("Access-Control-Allow-Origin", "*")
		ctx.Response.Header.Add("Access-Control-Allow-Credentials", "true")
		ctx.Response.Header.Add("Access-Control-Allow-Methods", "POST,OPTIONS,GET,PUT,DELETE")
		ctx.Response.Header.Add("Access-Control-Allow-Headers", "Access-Control-Allow-Origin, Content-Type, x-requested-with")

		next(ctx)
	}
}
func options(ctx *fasthttp.RequestCtx) {}

func getReviews(ctx *fasthttp.RequestCtx) {
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

	ctx.Response.Header.SetContentType("application/json; charset=utf-8")
	ctx.Write(json)
}

func getItems(ctx *fasthttp.RequestCtx) {
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

	ctx.Response.Header.SetContentType("application/json; charset=utf-8")
	ctx.Write(json)
}

func scrape(ctx *fasthttp.RequestCtx) {
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
		item := v.(*model.ItemDTO)

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

	ctx.Response.Header.SetContentType("application/json; charset=utf-8")
	json := fmt.Sprintf(`{"count":%d}`, count)
	ctx.WriteString(json)
}

func getItemQuery(ctx *fasthttp.RequestCtx) *model.ItemQuery {
	asin := string(ctx.FormValue("asin"))
	itemNo := string(ctx.FormValue("itemNo"))
	statusStr := string(ctx.FormValue("status"))
	pageSizeStr := string(ctx.FormValue("pageSize"))
	cursor := string(ctx.FormValue("cursor"))
	// status, err := strconv.Atoi(statusStr)
	// if err != nil {
	// 	status = -1
	// }
	pageSize, err := strconv.Atoi(pageSizeStr)
	if err != nil {
		pageSize = 10000
	}

	return &model.ItemQuery{
		Status:   statusStr,
		PageSize: pageSize,
		Cursor:   cursor,
		ASIN:     asin,
		ItemNo:   itemNo,
	}
}

func getReviewQuery(ctx *fasthttp.RequestCtx) *model.ReviewQuery {
	asin := string(ctx.FormValue("asin"))
	itemNo := string(ctx.FormValue("itemNo"))
	pageSizeStr := string(ctx.FormValue("pageSize"))
	cursor := string(ctx.FormValue("cursor"))
	pageSize, err := strconv.Atoi(pageSizeStr)
	if err != nil {
		pageSize = 10000
	}

	return &model.ReviewQuery{
		PageSize: pageSize,
		Cursor:   cursor,
		ASIN:     asin,
		ItemNo:   itemNo,
	}
}

func exportReviews(ctx *fasthttp.RequestCtx) {
	// query := getReviewQuery(ctx)
	// result, err := _reviewDAL.GetAllReviews(query)
	// if u.LogError(err) {
	// 	ctx.WriteString(err.Error())
	// 	return
	// }

}