package main

import (
	"encoding/json"
	"fmt"
	"strconv"
	"sync"
	"sync/atomic"
	"time"

	"github.com/kataras/iris/v12"
	"github.com/olivere/elastic/v7"
	"github.com/syncfuture/go/sconfig"
	"github.com/syncfuture/go/spool"
	"github.com/syncfuture/go/stask"
	"github.com/syncfuture/go/u"
	"github.com/syncfuture/scraper/scdp"
	"github.com/syncfuture/scraper/store"
	"github.com/syncfuture/spiders/wayfair"
	"github.com/syncfuture/spiders/wayfair/dal"
	"github.com/syncfuture/spiders/wayfair/dal/es"
	"github.com/syncfuture/spiders/wayfair/model"
	"github.com/tealeg/xlsx"
)

type wayfairHttpHandlers struct {
	configProvier sconfig.IConfigProvider
	reviewDAL     dal.IReviewDAL
	itemDAL       dal.IItemDAL
	scrapeLocker  *sync.Mutex
	// proxyStore    store.IProxyStore
	CDPClient     *scdp.ChromeDPClient
	maxConcurrent int
	bufferPool    spool.BufferPool
}

func NewWayfairHttpHandlers(cp sconfig.IConfigProvider, proxyStore store.IProxyStore) *wayfairHttpHandlers {
	addrs := cp.GetStringSlice("ES.Addrs")

	itemDAL, err := es.NewESItemDAL(
		elastic.SetURL(addrs...),
		elastic.SetSniff(false),
	)
	u.LogFaltal(err)

	reviewDAL, err := es.NewESReviewDAL(
		elastic.SetURL(addrs...),
		elastic.SetSniff(false),
	)
	u.LogFaltal(err)

	return &wayfairHttpHandlers{
		configProvier: cp,
		itemDAL:       itemDAL,
		reviewDAL:     reviewDAL,
		scrapeLocker:  new(sync.Mutex),
		// proxyStore:    proxyStore,
		CDPClient:     scdp.NewChromeDPClient(cp, proxyStore),
		maxConcurrent: cp.GetIntDefault("WayfairMaxConcurrent", 1),
		bufferPool:    spool.NewSyncBufferPool(4096),
	}
}

func (x *wayfairHttpHandlers) GetReviews(ctx iris.Context) {
	ctx.ContentType("application/json; charset=utf-8")
	query := x.getReviewQuery(ctx)
	result, err := x.reviewDAL.GetAllReviews(query)

	if u.LogError(err) {
		// ctx.WriteString(err.Error())
		ctx.WriteString("{}")
		return
	}

	json, err := json.Marshal(result)
	if u.LogError(err) {
		// ctx.WriteString(err.Error())
		ctx.WriteString("{}")
		return
	}

	ctx.Write(json)
}

func (x *wayfairHttpHandlers) GetItems(ctx iris.Context) {
	ctx.ContentType("application/json; charset=utf-8")
	query := x.getItemQuery(ctx)
	result, err := x.itemDAL.GetAllItems(query)

	if u.LogError(err) {
		// ctx.WriteString(err.Error())
		ctx.WriteString("{}")
		return
	}

	json, err := json.Marshal(result)
	if u.LogError(err) {
		// ctx.WriteString(err.Error())
		ctx.WriteString("{}")
		return
	}

	ctx.Write(json)
}

func (x *wayfairHttpHandlers) PostScrape(ctx iris.Context) {
	x.scrapeLocker.Lock()
	defer x.scrapeLocker.Unlock()

	query := x.getItemQuery(ctx)
	result, err := x.itemDAL.GetAllItems(query)
	if u.LogError(err) {
		ctx.WriteString(err.Error())
		return
	}

	count := int32(0)
	fromDate := time.Now().AddDate(0, -1, 0) // 一个月内的评论

	f := stask.NewFlowScheduler(x.maxConcurrent)
	f.SliceRun(&result.Items, func(i int, v interface{}) {
		item := v.(*model.ItemDTO)

		atomic.AddInt32(&count, 1)

		s := wayfair.NewReviewsScraper(x.CDPClient)
		reviews, err := s.FetchReviews(item, fromDate)
		if u.LogError(err) {
			if err.Error() == _notFoundError {
				item.Status = 404
			} else {
				item.Status = -1
			}
			x.itemDAL.SaveItems(item)
			return
		}

		if len(reviews) > 0 { // 有评论才存储
			err = x.reviewDAL.SaveReviews(reviews...)
			if u.LogError(err) {
				item.Status = -1
				x.itemDAL.SaveItems(item)
				return
			}
		}

		item.Status = 1
		x.itemDAL.SaveItems(item)
	})

	ctx.ContentType("application/json; charset=utf-8")
	json := fmt.Sprintf(`{"count":%d}`, count)
	ctx.WriteString(json)
}

func (x *wayfairHttpHandlers) ExportReviews(ctx iris.Context) {
	query := x.getReviewQuery(ctx)
	result, err := x.reviewDAL.GetAllReviews(query)
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
	header.AddCell().Value = "ReviewID"
	header.AddCell().Value = "SKU"
	header.AddCell().Value = "ItemNOs"
	header.AddCell().Value = "ReviewerName"
	header.AddCell().Value = "HasVerifiedBuyerStatus"
	header.AddCell().Value = "IsUSReviewer"
	header.AddCell().Value = "ReviewerBadgeText"
	header.AddCell().Value = "ReviewerBadgeID"
	header.AddCell().Value = "RatingStars"
	header.AddCell().Value = "Date"
	header.AddCell().Value = "Headline"
	header.AddCell().Value = "ProductComments"
	header.AddCell().Value = "HeadlineTranslation"
	header.AddCell().Value = "ProductCommentsTranslation"
	header.AddCell().Value = "LanguageCode"
	header.AddCell().Value = "ReviewHelpful"
	header.AddCell().Value = "IsReviewHelpfulUpvoted"
	header.AddCell().Value = "ProductName"
	header.AddCell().Value = "ProductUrl"
	// header.AddCell().Value = "CreatedOn"
	header.AddCell().Value = "CustomerPhotos"

	for _, review := range result.Reviews {
		row := sheet.AddRow()
		row.AddCell().Value = strconv.Itoa(review.ReviewID)
		row.AddCell().Value = review.SKU
		row.AddCell().Value = review.ItemNOs
		row.AddCell().Value = review.ReviewerName
		row.AddCell().Value = strconv.FormatBool(review.HasVerifiedBuyerStatus)
		row.AddCell().Value = strconv.FormatBool(review.IsUSReviewer)
		row.AddCell().Value = review.ReviewerBadgeText
		row.AddCell().Value = strconv.Itoa(review.ReviewerBadgeID)
		row.AddCell().Value = strconv.Itoa(review.RatingStars)
		row.AddCell().Value = review.Date
		row.AddCell().Value = review.Headline
		row.AddCell().Value = review.ProductComments
		row.AddCell().Value = review.HeadlineTranslation
		row.AddCell().Value = review.ProductCommentsTranslation
		row.AddCell().Value = review.LanguageCode
		row.AddCell().Value = strconv.Itoa(review.ReviewHelpful)
		row.AddCell().Value = strconv.FormatBool(review.IsReviewHelpfulUpvoted)
		row.AddCell().Value = review.ProductName
		row.AddCell().Value = review.ProductUrl
		// row.AddCell().Value = review.CreatedOn

		var photoStr string
		if len(review.CustomerPhotos) > 0 {
			data, err := json.Marshal(review.CustomerPhotos)
			if err != nil {
				photoStr = err.Error()
			} else {
				photoStr = string(data)
			}
		}
		row.AddCell().Value = photoStr
	}

	buffer := x.bufferPool.GetBuffer()
	defer x.bufferPool.PutBuffer(buffer)

	wb.Write(buffer)

	ctx.ContentType("application/octet-stream")
	ctx.Header("Content-Disposition", fmt.Sprintf("attachment;filename=%s.xlsx", time.Now().Format("20060102-150405")))
	ctx.Write(buffer.Bytes())
}

func (x *wayfairHttpHandlers) getItemQuery(ctx iris.Context) *model.ItemQuery {
	sku := string(ctx.FormValue("sku"))
	itemNo := string(ctx.FormValue("itemNo"))
	statusStr := string(ctx.FormValue("status"))
	pageSizeStr := string(ctx.FormValue("pageSize"))
	cursor := string(ctx.FormValue("cursor"))
	pageSize, err := strconv.Atoi(pageSizeStr)
	if err != nil {
		pageSize = 10000
	}

	return &model.ItemQuery{
		Status:   statusStr,
		PageSize: pageSize,
		Cursor:   cursor,
		SKU:      sku,
		ItemNo:   itemNo,
	}
}

func (x *wayfairHttpHandlers) getReviewQuery(ctx iris.Context) *model.ReviewQuery {
	sku := string(ctx.FormValue("sku"))
	itemNo := string(ctx.FormValue("itemNo"))
	pageSizeStr := string(ctx.FormValue("pageSize"))
	cursor := string(ctx.FormValue("cursor"))
	fromDate := string(ctx.FormValue("fromDate"))
	pageSize, err := strconv.Atoi(pageSizeStr)
	if err != nil {
		pageSize = 10000
	}

	return &model.ReviewQuery{
		PageSize: pageSize,
		Cursor:   cursor,
		SKU:      sku,
		ItemNo:   itemNo,
		FromDate: fromDate,
	}
}
