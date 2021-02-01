package dal

import "github.com/syncfuture/spiders/amazon/model"

type IItemDAL interface {
	GetItems(in *model.ItemQuery) (*model.ItemQueryResult, error)
	SaveItems(...*model.ItemDTO) error
	DeleteItems(...*model.ItemDTO) error
}
