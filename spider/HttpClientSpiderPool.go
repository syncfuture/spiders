package spider

import (
	"net/http"

	"github.com/syncfuture/go/u"

	"github.com/syncfuture/spiders/spider/model"
	"github.com/syncfuture/spiders/spider/store"
)

type HttpClientSpiderPool struct {
	proxies []*model.Proxy
}

func NewHttpClientSpiderPool(proxyStore store.IProxyStore) *HttpClientSpiderPool {
	proxies, err := proxyStore.GetProxies()
	u.LogFaltal(err)

	return &HttpClientSpiderPool{
		proxies: proxies,
	}
}

func (x HttpClientSpiderPool) GetSpider() ISpider {
	i := randInt(0, len(x.proxies)-1)
	proxy := x.proxies[i]

	return &HttpClientSpider{
		client: &http.Client{
			Transport: &http.Transport{
				Proxy: proxy.ToProxyURL(),
			},
		},
	}
}
