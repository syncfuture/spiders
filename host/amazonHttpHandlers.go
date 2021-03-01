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
	"github.com/syncfuture/scraper/store"
	"github.com/syncfuture/spiders/amazon"
	"github.com/syncfuture/spiders/amazon/dal"
	"github.com/syncfuture/spiders/amazon/dal/es"
	"github.com/tealeg/xlsx"
)

const (
	_notFoundError = "404 Not Found"
)

type amazonHttpHandlers struct {
	configProvier sconfig.IConfigProvider
	reviewDAL     dal.IReviewDAL
	itemDAL       dal.IItemDAL
	scrapeLocker  *sync.Mutex
	proxyStore    store.IProxyStore
	maxConcurrent int
	bufferPool    spool.BufferPool
}

func NewAmazonHttpHandlers(cp sconfig.IConfigProvider, proxyStore store.IProxyStore) *amazonHttpHandlers {
	addrs := cp.GetStringSlice("ES.Addrs")

	itemDAL, err := es.NewESItemDAL(
		elastic.SetURL(addrs...),
	)
	u.LogFaltal(err)

	reviewDAL, err := es.NewESReviewDAL(
		elastic.SetURL(addrs...),
	)
	u.LogFaltal(err)

	return &amazonHttpHandlers{
		configProvier: cp,
		itemDAL:       itemDAL,
		reviewDAL:     reviewDAL,
		scrapeLocker:  new(sync.Mutex),
		proxyStore:    proxyStore,
		maxConcurrent: cp.GetIntDefault("MaxConcurrent", 15),
		bufferPool:    spool.NewSyncBufferPool(4096),
	}
}

func (x *amazonHttpHandlers) GetReviews(ctx iris.Context) {
	query := x.getReviewQuery(ctx)
	result, err := x.reviewDAL.GetAllReviews(query)
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

func (x *amazonHttpHandlers) GetItems(ctx iris.Context) {
	query := x.getItemQuery(ctx)
	items, err := x.itemDAL.GetAllItems(query)
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

func (x *amazonHttpHandlers) PostScrape(ctx iris.Context) {
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
		item := v.(*amazon.ItemDTO)

		atomic.AddInt32(&count, 1)

		s := amazon.NewReviewsScraper(x.proxyStore, item.ASIN)
		reviews, err := s.FetchPages(&fromDate)
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
			// 关联E&E ItemNo
			for _, review := range reviews {
				review.CustomerNo = item.ItemNo
			}

			err = x.reviewDAL.SaveReviews(reviews)
			if u.LogError(err) {
				item.Status = -1
				x.itemDAL.SaveItems(item)
				return
			}
		}

		item.Status = 1
		x.itemDAL.SaveItems(item)
	})

	// err = _itemDAL.SaveItems(result.Items...)
	// if u.LogError(err) {
	// 	return
	// }

	ctx.ContentType("application/json; charset=utf-8")
	json := fmt.Sprintf(`{"count":%d}`, count)
	ctx.WriteString(json)
}

func (x *amazonHttpHandlers) ExportReviews(ctx iris.Context) {
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

	buffer := x.bufferPool.GetBuffer()
	defer x.bufferPool.PutBuffer(buffer)

	wb.Write(buffer)

	ctx.ContentType("application/octet-stream")
	ctx.Header("Content-Disposition", fmt.Sprintf("attachment;filename=%s.xlsx", time.Now().Format("20060102-150405")))
	ctx.Write(buffer.Bytes())
}

func (x *amazonHttpHandlers) getItemQuery(ctx iris.Context) *amazon.ItemQuery {
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

func (x *amazonHttpHandlers) getReviewQuery(ctx iris.Context) *amazon.ReviewQuery {
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
