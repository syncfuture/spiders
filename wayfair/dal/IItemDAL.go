package dal

import "github.com/syncfuture/spiders/wayfair"

type IItemDAL interface {
	GetItems(in *wayfair.ItemQuery) (*wayfair.ItemQueryResult, error)
	GetAllItems(in *wayfair.ItemQuery) (r *wayfair.ItemQueryResult, err error)
	SaveItems(...*wayfair.ItemDTO) error
	// SaveItem(*wayfair.ItemDTO) error
	DeleteItems(...*wayfair.ItemDTO) error
}
