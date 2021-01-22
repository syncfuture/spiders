package dal

import (
	"github.com/syncfuture/spiders/amazon/dal/es"
	"github.com/syncfuture/spiders/amazon/protoc/product"
)

type IAmazonProductDAL interface {
	SaveProductSKUInfos(...*product.ProductSKU) error
}

func NewAmazonProductDAL() IAmazonProductDAL {
	return new(es.ESAmazonProductDAL)
}
