package dal

import "github.com/syncfuture/spiders/amazon"

type IItemDAL interface {
	GetItems(in *amazon.ItemQuery) (*amazon.ItemQueryResult, error)
	GetAllItems(in *amazon.ItemQuery) (r *amazon.ItemQueryResult, err error)
	SaveItems(...*amazon.ItemDTO) error
	ClearItems() error
	DeleteItems(...*amazon.ItemDTO) error
}
