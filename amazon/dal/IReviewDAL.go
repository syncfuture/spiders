package dal

import (
	"github.com/syncfuture/scraper/amazon"
	"github.com/syncfuture/spiders/amazon/model"
)

type IReviewDAL interface {
	SaveReviews(reviews []*amazon.ReviewDTO) error
	GetReviews() (*model.ReviewQueryResult, error)
}
