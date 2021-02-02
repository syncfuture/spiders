package dal

import "github.com/syncfuture/spiders/amazon/model"

type IItemDAL interface {
	GetItems(in *model.ItemQuery) (*model.ItemQueryResult, error)
	GetAllItems(in *model.ItemQuery) (r *model.ItemQueryResult, err error)
	SaveItems(...*model.ItemDTO) error
	// SaveItem(*model.ItemDTO) error
	DeleteItems(...*model.ItemDTO) error
}
