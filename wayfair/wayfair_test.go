package wayfair

import (
	"testing"

	"github.com/olivere/elastic/v7"
	"github.com/stretchr/testify/assert"
	"github.com/syncfuture/go/sconfig"
	log "github.com/syncfuture/go/slog"
	"github.com/syncfuture/scraper/store/webshare"
	"github.com/syncfuture/spiders/wayfair/dal/es"
	"github.com/syncfuture/spiders/wayfair/model"
)

func TestGetReviews(t *testing.T) {
	cp := sconfig.NewJsonConfigProvider()
	log.Init(cp)

	store := webshare.NewDefaultWebShareProxyStore()
	t.Log(store)

	dal, err := es.NewESItemDAL(
		elastic.SetURL("http://192.168.188.200:9200"),
	)
	assert.NoError(t, err)

	query := &model.ItemQuery{}
	rs, err := dal.GetAllItems(query)
	assert.NoError(t, err)
	assert.NotEmpty(t, rs)

	// tryURL := "https://www.wayfair.com/bed-bath/pdp/product-zpcd6144.html"
	// urlScraper := NewURLScraper(store, "zpcd6144")

	// a, _, _ := urlScraper.FetchURL(tryURL)
	// log.Debug(a)

	// log.Debug(rs)
	// for _, item := range rs.Items {
	// 	scraper := NewReviewsScraper(cp, store, item)
	// 	reviews, err := scraper.FetchReviews()
	// 	if !u.LogError(err) {
	// 		item.Status = 1
	// 		dal.SaveReviews(reviews)
	// 	} else {
	// 		item.Status = -1
	// 	}
	// 	dal.SaveItems(item)
	// }
}
