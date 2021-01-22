package revidews

import (
	"github.com/syncfuture/spiders/spider"
	"github.com/syncfuture/spiders/spider/store/webshare"
)

type ReviewsSpider struct {
	pool *spider.SpiderPool
}

const (
	_webshareKey = "be09d781115fe3491743fa205ea786852513f474"
)

func NewReviewsSpider() (r *ReviewsSpider) {
	r = new(ReviewsSpider)
	ps := webshare.NewWebShareProxyStore("_webshareKey")
	r.pool = spider.NewSpiderPool(ps)
	return
}

func (x *ReviewsSpider) GetPageInfo() {
	x.pool.GetHttpClientSpider()
}
