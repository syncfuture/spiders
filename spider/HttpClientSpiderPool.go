package spider

import (
	"crypto/tls"
	"net/http"

	"github.com/syncfuture/go/u"

	"github.com/syncfuture/spiders/spider/model"
	"github.com/syncfuture/spiders/spider/store"
)

var (
	_userAgents = []string{
		"Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/87.0.4280.141 Safari/537.36",
		"Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/87.0.4280.141 Safari/537.36 Edg/87.0.664.75",
		"Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/87.0.4280.67 Safari/537.36 Edg/87.0.664.55",
		"Mozilla/5.0 (Windows NT 10.0; WOW64; Trident/7.0; rv:11.0) like Gecko",
	}
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
				Proxy:           proxy.ToProxyURL(),
				TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
			},
		},
		headers: x.createHeaders(),
	}
}

func (x HttpClientSpiderPool) GetSpiderWithHeaders(headers map[string]string) ISpider {
	i := randInt(0, len(x.proxies)-1)
	proxy := x.proxies[i]

	return &HttpClientSpider{
		client: &http.Client{
			Transport: &http.Transport{
				Proxy:           proxy.ToProxyURL(),
				TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
			},
		},
		headers: headers,
	}
}

func (x HttpClientSpiderPool) createHeaders() map[string]string {
	userAgent := _userAgents[randInt(0, len(_userAgents)-1)]
	r := map[string]string{
		"User-Agent":                userAgent,
		"Connection":                "keep-alive",
		"DNT":                       "1",
		"Upgrade-Insecure-Requests": "1",
		"Accept":                    "text/html,application/xhtml+xml,application/xml;q=0.9,image/webp,image/apng,*/*;q=0.8,application/signed-exchange;v=b3;q=0.9",
		"Sec-Fetch-Site":            "none",
		"Sec-Fetch-Mode":            "navigate",
		"Sec-Fetch-Dest":            "document",
		"Accept-Encoding":           "gzip, deflate, br",
		"Accept-Language":           "en-US,en;q=0.9",
	}
	return r
}
