package es

import (
	"context"
	"encoding/json"
	"io"

	"github.com/olivere/elastic/v7"
	log "github.com/syncfuture/go/slog"
	"github.com/syncfuture/go/u"
	"github.com/syncfuture/spiders/amazon/dal"
	"github.com/syncfuture/spiders/amazon/model"
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

func (x *ESItemDAL) GetItems(in *model.ItemQuery) (r *model.ItemQueryResult, err error) {
	r = new(model.ItemQueryResult)

	searchService := x.esClient.Search(_itemIndex).
		// searchService := x.esClient.Scroll(_itemIndex).
		Sort("ASIN.keyword", false).
		Size(in.PageSize)

	if in.SearchAfter != "" {
		searchService.SearchAfter(in.SearchAfter)
	}

	filters := []elastic.Query{}
	if in.Status >= 0 {
		filters = append(filters, elastic.NewMatchQuery("Status", in.Status))
	}
	if in.ASIN != "" {
		filters = append(filters, elastic.NewMatchQuery("ASIN.keyword", in.ASIN))
	}
	if in.ItemNo != "" {
		filters = append(filters, elastic.NewMatchQuery("ItemNo.keyword", in.ItemNo))
	}

	// if in.ScrollID != "" {
	// 	searchService.ScrollId(in.ScrollID)
	// }
	searchService.Query(elastic.NewBoolQuery().Filter(filters...))

	resp, err := searchService.Do(context.Background())
	if err == io.EOF || (resp != nil && len(resp.Hits.Hits) == 0) {
		return r, nil
	} else if u.LogErrorMsg(err, r) {
		return r, err
	}

	r.TotalCount = resp.TotalHits()
	// if r.TotalCount >= int64(in.PageSize) {
	// 	r.ScrollID = resp.ScrollId // 只有总条数大于分页数时才需要滚动查询，不做此判断ES总是会返回ScrollID
	// }
	for _, value := range resp.Hits.Hits {
		var doc *model.ItemDTO
		err = json.Unmarshal(value.Source, &doc)
		if !u.LogError(err) {
			doc.SearchAfter = value.Sort[0].(string)
			r.Items = append(r.Items, doc)
		} else {
			break
		}
	}

	return
}

func (x *ESItemDAL) SaveItems(items ...*model.ItemDTO) error {
	bulkService := x.esClient.Bulk().Index(_itemIndex)

	for _, item := range items {
		request := elastic.NewBulkIndexRequest().Id(item.ASIN).Doc(item)
		bulkService.Add(request)
	}

	resp, err := bulkService.Do(context.Background())
	if err != nil {
		return err
	} else {
		log.Debugf("[%d] reviews saved", len(resp.Succeeded()))
	}
	return err
}

// func (x *ESItemDAL) SaveItem(item *model.ItemDTO) error {
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

func (x *ESItemDAL) DeleteItems(items ...*model.ItemDTO) error {
	bulkService := x.esClient.Bulk().Index(_itemIndex)

	for _, item := range items {
		request := elastic.NewBulkDeleteRequest().Id(item.ASIN)
		bulkService.Add(request)
	}

	resp, err := bulkService.Do(context.Background())
	if err != nil {
		return err
	} else {
		log.Debugf("[%d] reviews saved", len(resp.Succeeded()))
	}
	return err
}
