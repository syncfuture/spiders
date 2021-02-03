package dal

import (
	"github.com/syncfuture/scraper/amazon"
	"github.com/syncfuture/spiders/amazon/model"
)

type IReviewDAL interface {
	SaveReviews(reviews []*amazon.ReviewDTO) error
	GetReviews(*model.ReviewQuery) (*model.ReviewQueryResult, error)
	GetAllReviews(*model.ReviewQuery) (*model.ReviewQueryResult, error)
}
