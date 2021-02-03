package es

import (
	"context"
	"encoding/json"
	"io"

	"github.com/olivere/elastic/v7"
	log "github.com/syncfuture/go/slog"
	"github.com/syncfuture/go/u"
	"github.com/syncfuture/scraper/amazon"
	"github.com/syncfuture/spiders/amazon/dal"
	"github.com/syncfuture/spiders/amazon/model"
)

const _reviewIndex = "amazon-reviews"

type ESReviewDAL struct {
	esClient *elastic.Client
}

func NewESReviewDAL(options ...elastic.ClientOptionFunc) (dal.IReviewDAL, error) {
	var err error
	r := new(ESReviewDAL)
	r.esClient, err = elastic.NewClient(options...)
	return r, err
}

func (x *ESReviewDAL) SaveReviews(reviews []*amazon.ReviewDTO) error {
	bulkService := x.esClient.Bulk().Index(_reviewIndex)

	for _, review := range reviews {
		request := elastic.NewBulkIndexRequest().Id(review.AmazonID).Doc(review)
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

func (x *ESReviewDAL) GetReviews(in *model.ReviewQuery) (r *model.ReviewQueryResult, err error) {
	r = new(model.ReviewQueryResult)
	// default
	// searchService := x.esClient.Search().Index(_reviewIndex).
	searchService := x.esClient.Scroll(_reviewIndex).
		Sort("ID.keyword", false).
		Size(in.PageSize)

		// if in.Cursor != "" {
		// 	searchService.SearchAfter(in.Cursor)
		// }
	if in.Cursor != "" {
		searchService.ScrollId(in.Cursor)
	}

	filters := []elastic.Query{}
	// if in.Marketplace != "" {
	// 	filters = append(filters, elastic.NewMatchQuery("marketplace-id.keyword", strings.ToUpper(in.Marketplace)))
	// }
	if in.ASIN != "" {
		filters = append(filters, elastic.NewMatchQuery("ASIN.keyword", in.ASIN))
	}
	if in.ItemNo != "" {
		filters = append(filters, elastic.NewMatchQuery("ItemNo.keyword", in.ItemNo))
	}
	// if in.From != "" {
	// 	from, err := time.Parse(time.RFC3339, in.From)
	// 	if err == nil {
	// 		filters = append(filters, elastic.NewMatchQuery("ItemNo.keyword", in.ItemNo))
	// 	}
	// }
	// filters = append(filters, elastic.NewBoolQuery().MustNot(elastic.NewExistsQuery("CustomerNo")))

	searchService.Query(elastic.NewBoolQuery().Filter(filters...))

	resp, err := searchService.Do(context.Background())
	if err == io.EOF || (resp != nil && len(resp.Hits.Hits) == 0) {
		return r, nil
	} else if u.LogErrorMsg(err, r) {
		return r, err
	}

	r.TotalCount = resp.TotalHits()

	for _, value := range resp.Hits.Hits {
		var doc *amazon.ReviewDTO
		err = json.Unmarshal(value.Source, &doc)
		if !u.LogError(err) {
			r.Reviews = append(r.Reviews, doc)
		} else {
			break
		}
	}

	// r.Cursor = r.Reviews[len(r.Reviews)-1].ASIN
	if r.TotalCount >= int64(in.PageSize) {
		r.Cursor = resp.ScrollId // 只有总条数大于分页数时才需要滚动查询，不做此判断ES总是会返回ScrollID
	}
	return
}

func (x *ESReviewDAL) GetAllReviews(in *model.ReviewQuery) (*model.ReviewQueryResult, error) {
	in.PageSize = 10000

	var r1, r2 *model.ReviewQueryResult
	var err error
	r1, err = x.GetReviews(in)
	if err != nil {
		return nil, err
	}
	in.Cursor = r1.Cursor
	for in.Cursor != "" {
		r2, err = x.GetReviews(in)
		if err != nil {
			return nil, err
		}
		if len(r2.Reviews) > 0 {
			r1.Reviews = append(r1.Reviews, r2.Reviews...)
			in.Cursor = r2.Cursor
		} else {
			in.Cursor = ""
		}
	}

	r1.Cursor = "" // 获取所有不应该有Cusor返回
	return r1, err
}
