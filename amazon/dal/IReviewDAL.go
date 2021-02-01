package dal

import "github.com/syncfuture/scraper/amazon"

type IReviewDAL interface {
	SaveReviews(reviews []*amazon.ReviewDTO) error
	GetReviews() ([]*amazon.ReviewDTO, error)
}
