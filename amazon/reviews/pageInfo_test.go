package revidews

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/syncfuture/scraper"
	"github.com/syncfuture/scraper/store/webshare"
)

func TestPageInfo_FetchPageInfo(t *testing.T) {
	s := webshare.NewWebShareProxyStore("be09d781115fe3491743fa205ea786852513f474")
	sp := scraper.NewScraperPool(s)
	spider, err := NewReviewsSpider(sp, "B08164VTWH")
	assert.NoError(t, err)
	assert.Greater(t, spider.PageInfo.TotalPage, 0)

	reviews, err := spider.FetchIndex(2)
	assert.NotEmpty(t, reviews)
}
