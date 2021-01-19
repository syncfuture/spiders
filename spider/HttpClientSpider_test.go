package spider

import (
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
	url := "https://www.amazon.com/AMD-Ryzen-5900X-24-Thread-Processor/dp/B08164VTWH"
	html, err := spider.Get(url)

	assert.NoError(t, err)
	assert.NotEmpty(t, html)
	t.Log(html)
}
