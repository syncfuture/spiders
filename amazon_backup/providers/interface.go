package providers

import (
	"github.com/syncfuture/spiders/amazon/protoc/product"
	"github.com/syncfuture/spiders/amazon/providers/chrome"
)

type ISKUInfoProvider interface {
	Get(asin string) *product.ProductSKU
}

func NewSKUInfoProvider() ISKUInfoProvider {
	return new(chrome.ChromedpSKUInfoProvider)
}
