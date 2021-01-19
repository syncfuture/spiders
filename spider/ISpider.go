package spider

import "github.com/PuerkitoBio/goquery"

type ISpider interface {
	GetDocument(url string) (*goquery.Document, error)
}
