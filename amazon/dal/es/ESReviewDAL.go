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
		request := elastic.NewBulkIndexRequest().Id(review.ID).Doc(review)
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

func (x *ESReviewDAL) GetReviews() (r []*amazon.ReviewDTO, err error) {
	r = make([]*amazon.ReviewDTO, 0)
	// default
	searchService := x.esClient.Search().Index(_reviewIndex)

	filters := []elastic.Query{}
	// if in.Marketplace != "" {
	// 	filters = append(filters, elastic.NewMatchQuery("marketplace-id.keyword", strings.ToUpper(in.Marketplace)))
	// }

	searchService.Query(elastic.NewBoolQuery().Filter(filters...))
	// searchService.Aggregation("min-date", elastic.NewMinAggregation().Field("posted-date").Format("MM/dd/YYYY"))

	resp, err := searchService.Do(context.Background())
	if err == io.EOF || (resp != nil && len(resp.Hits.Hits) == 0) {
		return r, nil
	} else if u.LogErrorMsg(err, r) {
		return r, err
	}

	for _, value := range resp.Hits.Hits {
		var doc *amazon.ReviewDTO
		err = json.Unmarshal(value.Source, &doc)
		if !u.LogError(err) {
			r = append(r, doc)
		} else {
			break
		}
	}

	return
}
