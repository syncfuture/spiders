package spider

import (
	"github.com/syncfuture/go/u"

	"github.com/syncfuture/spiders/spider/httpClient"
	"github.com/syncfuture/spiders/spider/model"
	"github.com/syncfuture/spiders/spider/store"
)

type SpiderPool struct {
	proxies []*model.Proxy
}

func NewSpiderPool(proxyStore store.IProxyStore) *SpiderPool {
	proxies, err := proxyStore.GetProxies()
	u.LogFaltal(err)

	return &SpiderPool{
		proxies: proxies,
	}
}

func (x SpiderPool) GetHttpClientSpider() ISpider {
	i := randInt(0, len(x.proxies)-1)
	proxy := x.proxies[i]

	return httpClient.NewHttpClientSpider(proxy, nil)
}

func (x SpiderPool) GetHttpClientSpiderWithHeaders(headers map[string]string) ISpider {
	i := randInt(0, len(x.proxies)-1)
	proxy := x.proxies[i]
	return httpClient.NewHttpClientSpider(proxy, headers)
}
