package main

import (
	"regexp"
	"strings"

	log "github.com/syncfuture/go/slog"
	"github.com/syncfuture/go/u"
	"github.com/syncfuture/spiders/amazon/dal"
	"github.com/syncfuture/spiders/amazon/providers"
)

var (
	cleanerRegex = regexp.MustCompile(`\([^\)]+\)|[\r\n]`)
)

func main() {
	// https://www.amazon.com/product/product-reviews/B004Z4PMF0
	// https://www.amazon.com/product/product-reviews/B004Z4PMF0/?pageNumber=1&sortBy=recent
	// https://www.amazon.com/ask/questions/asin/B004Z4PMF0/1?sort=SUBMIT_DATE

	skuInfoProvider := providers.NewSKUInfoProvider()
	skuInfo := skuInfoProvider.Get("B07CZYCDS5")

	store := dal.NewAmazonProductDAL()

	err := store.SaveProductSKUInfos(skuInfo)
	if u.LogError(err) {
		return
	}

	log.Info(skuInfo, " saved")
}

func format(str *string) string {
	if str == nil {
		return ""
	}
	return strings.TrimSpace(strings.Trim(cleanerRegex.ReplaceAllString(*str, ""), "."))
}
