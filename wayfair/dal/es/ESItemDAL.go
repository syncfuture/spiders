package es

import (
	"context"
	"encoding/json"
	"io"

	"github.com/olivere/elastic/v7"
	log "github.com/syncfuture/go/slog"
	"github.com/syncfuture/go/u"
	"github.com/syncfuture/spiders/wayfair/dal"
	"github.com/syncfuture/spiders/wayfair/model"
)

const _itemIndex = "wayfair-items"

type ESItemDAL struct {
	esClient *elastic.Client
}

func NewESItemDAL(options ...elastic.ClientOptionFunc) (dal.IItemDAL, error) {
	var err error
	r := new(ESItemDAL)
	r.esClient, err = elastic.NewClient(options...)
	return r, err
}

func (x *ESItemDAL) GetItems(in *model.ItemQuery) (r *model.ItemQueryResult, err error) {
	r = new(model.ItemQueryResult)

	// searchService := x.esClient.Search(_itemIndex).
	searchService := x.esClient.Scroll(_itemIndex).
		Sort("SKU.keyword", false).
		Size(in.PageSize)

	// if in.Cursor != "" {
	// 	searchService.SearchAfter(in.Cursor)
	// }
	if in.Cursor != "" {
		searchService.ScrollId(in.Cursor)
	}

	filters := []elastic.Query{}
	if in.Status != "" {
		filters = append(filters, elastic.NewMatchQuery("Status", in.Status))
	}
	if in.SKU != "" {
		filters = append(filters, elastic.NewMatchQuery("SKU.keyword", in.SKU))
	}
	if in.ItemNo != "" {
		filters = append(filters, elastic.NewMatchQuery("ItemNo.keyword", in.ItemNo))
	}

	searchService.Query(elastic.NewBoolQuery().Filter(filters...))

	resp, err := searchService.Do(context.Background())
	if err == io.EOF || (resp != nil && len(resp.Hits.Hits) == 0) {
		return r, nil
	} else if u.LogErrorMsg(err, r) {
		return r, err
	}

	r.TotalCount = resp.TotalHits()

	for _, value := range resp.Hits.Hits {
		var doc *model.ItemDTO
		err = json.Unmarshal(value.Source, &doc)
		if !u.LogError(err) {
			r.Items = append(r.Items, doc)
		} else {
			break
		}
	}

	// r.Cursor = r.Items[len(r.Items)-1].SKU
	if r.TotalCount >= int64(in.PageSize) {
		r.Cursor = resp.ScrollId // 只有总条数大于分页数时才需要滚动查询，不做此判断ES总是会返回ScrollID
	}
	return
}

func (x *ESItemDAL) GetAllItems(in *model.ItemQuery) (*model.ItemQueryResult, error) {
	in.PageSize = 10000

	var r1, r2 *model.ItemQueryResult
	var err error
	r1, err = x.GetItems(in)
	if err != nil {
		return nil, err
	}
	in.Cursor = r1.Cursor
	for in.Cursor != "" {
		r2, err = x.GetItems(in)
		if err != nil {
			return nil, err
		}
		if len(r2.Items) > 0 {
			r1.Items = append(r1.Items, r2.Items...)
			in.Cursor = r2.Cursor
		} else {
			in.Cursor = ""
		}
	}

	r1.Cursor = "" // 获取所有不应该有Cusor返回
	return r1, err
}

func (x *ESItemDAL) SaveItems(items ...*model.ItemDTO) error {
	bulkService := x.esClient.Bulk().Index(_itemIndex)

	for _, item := range items {
		request := elastic.NewBulkIndexRequest().Id(item.SKU).Doc(item)
		bulkService.Add(request)
	}

	resp, err := bulkService.Do(context.Background())
	if err != nil {
		return err
	} else {
		log.Debugf("[%d] items saved", len(resp.Succeeded()))
	}
	return err
}

// func (x *ESItemDAL) SaveItem(item *model.ItemDTO) error {
// 	updateService := x.esClient.Update().Index(_itemIndex).
// 		Id(item.SKU).
// 		Type("items").
// 		Doc(map[string]interface{}{"Status": item.Status})
// 	resp, err := updateService.Do(context.Background())
// 	if err != nil {
// 		return err
// 	} else {
// 		log.Debug(resp.Result)
// 	}
// 	return err
// }

func (x *ESItemDAL) DeleteItems(items ...*model.ItemDTO) error {
	bulkService := x.esClient.Bulk().Index(_itemIndex)

	for _, item := range items {
		request := elastic.NewBulkDeleteRequest().Id(item.SKU)
		bulkService.Add(request)
	}

	resp, err := bulkService.Do(context.Background())
	if err != nil {
		return err
	} else {
		log.Debugf("[%d] items deleted", len(resp.Succeeded()))
	}
	return err
}
