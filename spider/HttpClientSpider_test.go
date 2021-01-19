package spider

import (
	"regexp"
	"sort"
	"strconv"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/syncfuture/go/sredis"
	"github.com/syncfuture/spiders/spider/store/redis"
)

func TestHttpClientSpider_Get(t *testing.T) {
	store := redis.NewRedisProxyStore("ams:Proxies", &sredis.RedisConfig{
		Addrs:    []string{"localhost:6379"},
		Password: "Famous901",
	})
	pool := NewHttpClientSpiderPool(store)
	spider := pool.GetSpider()
	// url := "https://api.ip.sb/ip"
	// url := "http://ipv4.webshare.io/"
	// url := "https://www.amazon.com/AMD-Ryzen-5900X-24-Thread-Processor/dp/B08164VTWH"
	url := "https://www.amazon.com/gp/offer-listing/B08164VTWH//ref=olp_f_new?f_primeEligible=true&f_new=true"
	doc, err := spider.GetDocument(url)

	assert.NoError(t, err)
	assert.NotEmpty(t, doc)
	// t.Log(doc.Find(".olpOfferPrice").Length())
	nodes := doc.Find(".olpOfferPrice").Text()
	// t.Log(doc.Find(".olpOfferPrice"))
	// t.Log(doc.Text())

	regex := regexp.MustCompile(`\$(\d+[\.\d]+)`)

	matches := regex.FindAllStringSubmatch(nodes, -1)
	prices := make([]float64, 0, len(matches))

	for _, m := range matches {
		if len(m) == 2 {
			price, _ := strconv.ParseFloat(m[1], 32)
			prices = append(prices, price)
		}
	}

	sort.Slice(prices, func(i, j int) bool {
		return prices[i] < prices[j]
	})

	if prices[0] <= 550 {
		t.Log("Found")
	} else {
		t.Log("Not found")
	}
}
