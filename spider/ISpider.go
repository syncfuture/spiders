package spider

import (
	"github.com/PuerkitoBio/goquery"
	"github.com/syncfuture/spiders/spider/model"
)

type ISpider interface {
	GetDocument(url string) (*goquery.Document, error)
	GetProxy() *model.Proxy
}
