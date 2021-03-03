package dal

import "github.com/syncfuture/spiders/wayfair/model"

type IReviewDAL interface {
	GetReviews(in *model.ReviewQuery) (*model.ReviewQueryResult, error)
	GetAllReviews(in *model.ReviewQuery) (r *model.ReviewQueryResult, err error)
	SaveReviews(...*model.ReviewDTO) error
	DeleteReviews(...*model.ReviewDTO) error
}
