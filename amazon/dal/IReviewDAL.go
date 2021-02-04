package dal

import (
	"github.com/syncfuture/spiders/amazon"
)

type IReviewDAL interface {
	SaveReviews(reviews []*amazon.ReviewDTO) error
	GetReviews(*amazon.ReviewQuery) (*amazon.ReviewQueryResult, error)
	GetAllReviews(*amazon.ReviewQuery) (*amazon.ReviewQueryResult, error)
}
