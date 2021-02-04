package amazon

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/syncfuture/go/sconfig"
	log "github.com/syncfuture/go/slog"
	"github.com/syncfuture/scraper/store/webshare"
)

const (
	_testWebshareKey = "be09d781115fe3491743fa205ea786852513f474"
	_testASIN        = "B00D4KN4Q0"
)

var (
	_testStore = webshare.NewWebShareProxyStore(_testWebshareKey)
)

func init() {
	cp := sconfig.NewJsonConfigProvider()
	log.Init(cp)
}

func TestReviewsScraper_FetchPage(t *testing.T) {
	sc := NewReviewsScraper(_testStore, _testASIN)
	reviews, _, err := sc.FetchPage("")
	assert.NoError(t, err)
	assert.NotEmpty(t, reviews)
}

func TestReviewsScraper_FetchAllPages(t *testing.T) {
	sc := NewReviewsScraper(_testStore, _testASIN)

	fromDate := time.Now().AddDate(0, -1, 0)
	reviews, err := sc.FetchPages(&fromDate)
	assert.NoError(t, err)
	assert.NotEmpty(t, reviews)
}

func TestReviewsScraper_FetchToDatePages(t *testing.T) {
	sc := NewReviewsScraper(_testStore, _testASIN)

	toDate := time.Now().Add(-5 * 24 * time.Hour)
	reviews, err := sc.FetchPages(&toDate)
	assert.NoError(t, err)
	assert.NotEmpty(t, reviews)
}

func TestOffersScraper_FetchPage(t *testing.T) {
	sc := NewOffsersScraper(_testStore, _testASIN, false, false, false)

	var results *OfferResultDTO
	var err error
	results, err = sc.FetchPagedOffers(results)
	assert.NoError(t, err)
	assert.NotEmpty(t, results.Offers)
}

func TestOffersScraper_FetchAllOffers(t *testing.T) {
	sc := NewOffsersScraper(_testStore, _testASIN, false, false, false)

	results, err := sc.FetchAllOffers()
	assert.NoError(t, err)
	assert.NotEmpty(t, results)
}

func TestSKUInfoScraper_Fetch(t *testing.T) {
	s := NewSKUInfoScraper(_testStore, _testASIN)
	sku, err := s.Fetch()
	assert.NoError(t, err)
	assert.NotNil(t, sku)
	t.Log(sku)
}
