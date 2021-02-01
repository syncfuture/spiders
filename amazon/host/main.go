package main

import (
	"encoding/json"
	"fmt"
	"strconv"
	"sync/atomic"

	"github.com/buaazp/fasthttprouter"
	"github.com/olivere/elastic/v7"
	"github.com/syncfuture/go/srand"
	"github.com/syncfuture/go/stask"
	"github.com/syncfuture/go/u"
	"github.com/syncfuture/spiders/amazon/dal"
	"github.com/syncfuture/spiders/amazon/dal/es"
	"github.com/syncfuture/spiders/amazon/model"
	"github.com/valyala/fasthttp"
)

const (
	_key = "be09d781115fe3491743fa205ea786852513f474"
)

var (
	_reviewDAL dal.IReviewDAL
	_itemDAL   dal.IItemDAL
)

func main() {
	// cp := sconfig.NewJsonConfigProvider()
	// log.Init(cp)
	// store := webshare.NewWebShareProxyStore(_key)
	// scraper, err := amazon.NewReviewsScraper(store, "B000SDKDM4")
	// if u.LogError(err) {
	// 	return
	// }

	// fromDate := time.Now().Add(-5 * 24 * time.Hour)
	// reviews, err := scraper.FetchPages(&fromDate)
	// if u.LogError(err) {
	// 	return
	// }
	// log.Infof("%d reviews fetched", len(reviews))

	var err error
	_itemDAL, err = es.NewESItemDAL(
		elastic.SetURL("http://192.168.188.166:9200"),
	)
	if u.LogError(err) {
		return
	}

	_reviewDAL, err = es.NewESReviewDAL(
		elastic.SetURL("http://192.168.188.166:9200"),
	)
	if u.LogError(err) {
		return
	}

	router := fasthttprouter.New()
	router.GET("/reviews", allowCORS(getReviews))
	router.GET("/items", allowCORS(getItems))
	router.POST("/scrape", allowCORS(scrape))

	fasthttp.ListenAndServe(":7000", router.Handler)
}

func allowCORS(next fasthttp.RequestHandler) fasthttp.RequestHandler {
	return func(ctx *fasthttp.RequestCtx) {

		//ç
		ctx.Response.Header.Add("Access-Control-Allow-Origin", "*")
		ctx.Response.Header.Add("Access-Control-Allow-Credentials", "true")
		ctx.Response.Header.Add("Access-Control-Allow-Methods", "POST,OPTIONS,GET,PUT,DELETE")
		ctx.Response.Header.Add("Access-Control-Allow-Headers", "Access-Control-Allow-Origin, Content-Type, x-requested-with")

		next(ctx)
	}
}

func getReviews(ctx *fasthttp.RequestCtx) {
	reviews, err := _reviewDAL.GetReviews()
	if u.LogError(err) {
		ctx.WriteString(err.Error())
		return
	}

	json, err := json.Marshal(reviews)
	if u.LogError(err) {
		ctx.WriteString(err.Error())
		return
	}

	ctx.Response.Header.SetContentType("application/json; charset=utf-8")
	ctx.Write(json)
}

func getItems(ctx *fasthttp.RequestCtx) {
	query := getItemQuery(ctx)
	items, err := _itemDAL.GetItems(query)
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
	query := getItemQuery(ctx)
	items, err := _itemDAL.GetItems(query)
	if u.LogError(err) {
		ctx.WriteString(err.Error())
		return
	}

	count := int32(0)

	f := stask.NewFlowScheduler(20)
	f.SliceRun(items, func(i int, v interface{}) {
		item := v.(*model.ItemDTO)

		atomic.AddInt32(&count, 1)
		item.Status = srand.IntRange(0, 2)
	})

	ctx.Response.Header.SetContentType("application/json; charset=utf-8")
	json := fmt.Sprintf(`{"count":%d}`, count)
	ctx.WriteString(json)
}

func getItemQuery(ctx *fasthttp.RequestCtx) *model.ItemQuery {
	asin := string(ctx.FormValue("asin"))
	itemNo := string(ctx.FormValue("itemNo"))
	statusStr := string(ctx.FormValue("status"))
	pageSizeStr := string(ctx.FormValue("pageSize"))
	searchAfter := string(ctx.FormValue("searchAfter"))
	status, err := strconv.Atoi(statusStr)
	if err != nil {
		status = -1
	}
	pageSize, err := strconv.Atoi(pageSizeStr)
	if err != nil {
		pageSize = 100
	}

	return &model.ItemQuery{
		Status:      status,
		PageSize:    pageSize,
		SearchAfter: searchAfter,
		ASIN:        asin,
		ItemNo:      itemNo,
	}
}
