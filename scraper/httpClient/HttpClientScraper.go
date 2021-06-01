package httpClient

import (
	"compress/gzip"
	"crypto/tls"
	"io"
	"net/http"
	"net/url"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/SyncSoftInc/proxy/protoc/proxy"
	log "github.com/syncfuture/go/slog"
	"github.com/syncfuture/go/spool"
	"github.com/syncfuture/spiders/scraper"
	"github.com/syncfuture/spiders/scraper/header"
	"github.com/syncfuture/spiders/scraper/store"
)

var (
	_blockedErrs = []string{"automated access"}
	_expireErrs  = []string{"Proxy Authentication Required"}
)

type HttpClientScraper struct {
	proxyStore    store.IProxyStore
	headers       map[string]string
	bufferPool    spool.IBufferPool
	ExpireChecker func(*proxy.Proxy, string)
	BlockChecker  func(*proxy.Proxy, string)
}

func NewRandomScraper(proxyStore store.IProxyStore, headers map[string]string) scraper.IScraper {
	if headers == nil {
		hb := header.NewHeadersBuilder()
		headers = hb.Headers
	}

	r := new(HttpClientScraper)
	r.proxyStore = proxyStore
	r.headers = headers
	r.bufferPool = spool.NewSyncBufferPool(2048)
	r.ExpireChecker = func(p *proxy.Proxy, s string) {
		for _, err := range _expireErrs {
			if strings.Contains(s, err) {
				p.Score = -1
			}
		}
	}
	r.BlockChecker = func(p *proxy.Proxy, s string) {
		for _, err := range _blockedErrs {
			if strings.Contains(s, err) {
				p.Score = 0
			}
		}
	}

	return r
}

func (x *HttpClientScraper) Get(targetURL string) (r *scraper.ScrapeResult, err error) {
	msg, _ := http.NewRequest("GET", targetURL, nil)

	for k, v := range x.headers {
		msg.Header.Add(k, v)
	}

	p := x.proxyStore.Rent()     // 租借一个代理
	defer x.proxyStore.Return(p) // 用完归还
	c := &http.Client{
		Transport: &http.Transport{
			Proxy: func(*http.Request) (*url.URL, error) {
				return url.Parse(p.URI)
			},
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		},
	}
	log.Debugf("[%s] GET %s", p.URI, targetURL)
	resp, err := c.Do(msg)
	if err != nil {
		// 验证代理是否可用
		x.ExpireChecker(p, err.Error())
		return scraper.NewScrapeResult(http.StatusBadRequest, nil, p, x.headers, nil), err
	}
	statusCode := resp.StatusCode

	if statusCode == http.StatusNotFound {
		return scraper.NewScrapeResult(statusCode, nil, p, x.headers, resp.Header), scraper.NotFound
	}

	var bodyReader io.ReadCloser
	buffer := x.bufferPool.GetBuffer()
	defer func() {
		x.bufferPool.PutBuffer(buffer)
		if bodyReader != nil {
			bodyReader.Close()
		}
	}()

	if resp.Header.Get("Content-Encoding") == "gzip" {
		bodyReader, _ = gzip.NewReader(resp.Body)
	} else {
		bodyReader = resp.Body
	}

	buffer.ReadFrom(bodyReader)

	// 验证代理是否可用
	x.ExpireChecker(p, buffer.String())
	if p.Score <= 0 {
		statusCode = http.StatusBadRequest
	}

	doc, err := goquery.NewDocumentFromReader(buffer)
	if err != nil {
		return scraper.NewScrapeResult(statusCode, doc, p, x.headers, resp.Header), err
	}

	return scraper.NewScrapeResult(statusCode, doc, p, x.headers, resp.Header), err
}

// func (x HttpClientScraper) ExpireCheck(err error, proxy *proxy.Proxy) bool {
// 	if err != nil {
// 		if webshare.IsExiredError(err.Error()) {
// 			// 代理过期
// 			proxy.Expired = true
// 		}
// 		return true
// 	}

// 	return false
// }

// const _automatedAccess = "automated access"

// func (x HttpClientScraper) BlockCheck(html string, proxy *proxy.Proxy) bool {
// 	if strings.Contains(html, _automatedAccess) {
// 		// 代理被封
// 		proxy.Blocked = true
// 		return true
// 	}
// 	return false
// }
