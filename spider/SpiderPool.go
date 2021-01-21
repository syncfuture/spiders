package spider

import (
	"github.com/syncfuture/spiders/spider/httpClient"
	"github.com/syncfuture/spiders/spider/store"
)

type SpiderPool struct {
	ProxyStore store.IProxyStore
}

func NewSpiderPool(proxyStore store.IProxyStore) *SpiderPool {
	return &SpiderPool{
		ProxyStore: proxyStore,
	}
}

func (x SpiderPool) GetHttpClientSpider() (ISpider, error) {
	proxy, err := x.ProxyStore.GetRandomProxy()
	if err != nil {
		return nil, err
	}
	return httpClient.NewHttpClientSpider(proxy, nil), err
}

func (x SpiderPool) GetHttpClientSpiderWithHeaders(headers map[string]string) (ISpider, error) {
	proxy, err := x.ProxyStore.GetRandomProxy()
	if err != nil {
		return nil, err
	}
	return httpClient.NewHttpClientSpider(proxy, headers), err
}
