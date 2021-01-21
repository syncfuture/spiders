package spider

import (
	"regexp"
	"sort"
	"strconv"
	"strings"
	"testing"
	"time"

	log "github.com/syncfuture/go/slog"
	"github.com/syncfuture/go/u"
	"github.com/syncfuture/spiders/spider/store/webshare"
)

func TestHttpClientSpider_Get(t *testing.T) {
	// store := redis.NewRedisProxyStore("ams:Proxies", &sredis.RedisConfig{
	// 	Addrs:    []string{"localhost:6379"},
	// 	Password: "Famous901",
	// })
	store := webshare.NewWebShareProxyStore("be09d781115fe3491743fa205ea786852513f474")
	pool := NewSpiderPool(store)

	for i := 0; i < 100; i++ {
		check(pool)
		time.Sleep(1000 * time.Millisecond)
	}
}

func check(pool *SpiderPool) {
	spider, err := pool.GetHttpClientSpider()
	u.LogFaltal(err)
	// url := "https://api.ip.sb/ip"
	// url := "http://ipv4.webshare.io/"
	// url := "https://www.amazon.com/AMD-Ryzen-5900X-24-Thread-Processor/dp/B08164VTWH"
	url := "https://www.amazon.com/gp/offer-listing/B08164VTWH//ref=olp_f_new?f_primeEligible=true&f_new=true"
	doc, err := spider.GetDocument(url)
	var html string
	if doc != nil {
		html, _ = doc.Html()
		html = strings.Replace(html, "\n", "", -1)
	}
	if u.LogError(err) {
		log.Warn(html)
		// 可能代理有问题, 移除
		invalidProxy := spider.GetProxy()
		invalidProxy.Blocked = true
		pool.ProxyStore.SaveProxy(invalidProxy)
		return
	}

	nodes := doc.Find(".olpOfferPrice")

	regex := regexp.MustCompile(`\$(\d+[\.\d]+)`)

	matches := regex.FindAllStringSubmatch(nodes.Text(), -1)
	if len(matches) == 0 {
		log.Warn(html)
		// 可能代理有问题, 移除
		invalidProxy := spider.GetProxy()
		invalidProxy.Blocked = true
		pool.ProxyStore.SaveProxy(invalidProxy)
		return
	}
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

	log.Info(prices)
}
