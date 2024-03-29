package es

import (
	"context"
	"encoding/json"
	"io"

	"github.com/olivere/elastic/v7"
	log "github.com/syncfuture/go/slog"
	"github.com/syncfuture/go/u"
	"github.com/syncfuture/spiders/amazon"
	"github.com/syncfuture/spiders/amazon/dal"
)

const _itemIndex = "amazon-items"

type ESItemDAL struct {
	esClient *elastic.Client
}

func NewESItemDAL(options ...elastic.ClientOptionFunc) (dal.IItemDAL, error) {
	var err error
	r := new(ESItemDAL)
	r.esClient, err = elastic.NewClient(options...)
	return r, err
}

func (x *ESItemDAL) GetItems(in *amazon.ItemQuery) (r *amazon.ItemQueryResult, err error) {
	r = new(amazon.ItemQueryResult)

	// searchService := x.esClient.Search(_itemIndex).
	searchService := x.esClient.Scroll(_itemIndex).
		Sort("ASIN.keyword", false).
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
	if in.ASIN != "" {
		filters = append(filters, elastic.NewMatchQuery("ASIN.keyword", in.ASIN))
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
		var doc *amazon.ItemDTO
		err = json.Unmarshal(value.Source, &doc)
		if !u.LogError(err) {
			r.Items = append(r.Items, doc)
		} else {
			break
		}
	}

	// r.Cursor = r.Items[len(r.Items)-1].ASIN
	if r.TotalCount >= int64(in.PageSize) {
		r.Cursor = resp.ScrollId // 只有总条数大于分页数时才需要滚动查询，不做此判断ES总是会返回ScrollID
	}
	return
}

func (x *ESItemDAL) GetAllItems(in *amazon.ItemQuery) (*amazon.ItemQueryResult, error) {
	in.PageSize = 10000

	var r1, r2 *amazon.ItemQueryResult
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

func (x *ESItemDAL) SaveItems(items ...*amazon.ItemDTO) error {
	if len(items) == 0 {
		return nil
	}
	bulkService := x.esClient.Bulk().Index(_itemIndex)

	for _, item := range items {
		request := elastic.NewBulkIndexRequest().Id(item.ASIN).Doc(item)
		bulkService.Add(request)
	}

	resp, err := bulkService.Do(context.Background())
	if err != nil {
		return err
	} else {
		log.Infof("[%d] items saved", len(resp.Succeeded()))
	}
	return err
}

// func (x *ESItemDAL) SaveItem(item *amazon.ItemDTO) error {
// 	updateService := x.esClient.Update().Index(_itemIndex).
// 		Id(item.ASIN).
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

func (x *ESItemDAL) DeleteItems(items ...*amazon.ItemDTO) error {
	if len(items) == 0 {
		return nil
	}
	bulkService := x.esClient.Bulk().Index(_itemIndex)

	for _, item := range items {
		request := elastic.NewBulkDeleteRequest().Id(item.ASIN)
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

func (x *ESItemDAL) ClearItems() error {
	deleteIndexService := x.esClient.DeleteIndex(_itemIndex)

	resp, err := deleteIndexService.Do(context.Background())
	if err != nil {
		return err
	} else {
		log.Debugf("items clear acknowledged: %t", resp.Acknowledged)
	}
	return err
}
