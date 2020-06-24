package es

import (
	"context"

	elastic "github.com/olivere/elastic/v7"
	"github.com/syncfuture/go/config"
	log "github.com/syncfuture/go/slog"
	"github.com/syncfuture/go/u"
	"github.com/syncfuture/spiders/amazon/protoc/product"
)

const (
	_skuInfoIndexName = "sku-info-amazon"
)

var (
	_esClient *elastic.Client
)

func init() {
	configProvider := config.NewJsonConfigProvider()
	esURLs := configProvider.GetStringSlice("ES.Addrs")

	cfg := []elastic.ClientOptionFunc{
		elastic.SetURL(esURLs...),
		elastic.SetSniff(false),
		elastic.SetInfoLog(&logger{Type: "info"}),
		elastic.SetErrorLog(&logger{Type: "error"}),
	}

	var err error
	_esClient, err = elastic.NewClient(cfg...)
	u.LogFaltal(err)
}

type ESAmazonProductDAL struct {
}

func (x *ESAmazonProductDAL) SaveProductSKUInfos(skus ...*product.ProductSKU) error {
	bulkService := _esClient.Bulk()

	for _, entry := range skus {
		request := elastic.NewBulkIndexRequest().Index(_skuInfoIndexName).Id(entry.ID).Doc(entry)
		bulkService.Add(request)
	}

	resp, err := bulkService.Do(context.Background())
	if u.LogError(err) {
		log.Error(resp.Items)
	} else {
		log.Debugf("[%d] items saved", len(resp.Succeeded()))
		if len(resp.Succeeded()) > 0 {
			log.Debugf("%v", resp.Succeeded()[0])
		}
	}
	return err
}
