package httpClient

import (
	"compress/gzip"
	"crypto/tls"
	"io"
	"net/http"

	"github.com/PuerkitoBio/goquery"
	"github.com/syncfuture/spiders/spider/header"
	"github.com/syncfuture/spiders/spider/model"
)

type HttpClientSpider struct {
	client  *http.Client
	headers map[string]string
}

func NewHttpClientSpider(proxy *model.Proxy, headers map[string]string) (r *HttpClientSpider) {
	if headers == nil {
		hb := header.NewHeadersBuilder()
		headers = hb.Headers
	}

	r = new(HttpClientSpider)
	r.headers = headers
	if proxy != nil {
		r.client = &http.Client{
			Transport: &http.Transport{
				Proxy:           proxy.ToProxyURL(),
				TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
			},
		}
		return
	}

	r.client = http.DefaultClient
	return
}

func (x *HttpClientSpider) GetDocument(url string) (r *goquery.Document, err error) {
	msg, _ := http.NewRequest("GET", url, nil)

	for k, v := range x.headers {
		msg.Header.Add(k, v)
	}

	resp, err := x.client.Do(msg)
	if err != nil {
		return nil, err
	}

	var bodyReader io.ReadCloser

	defer func() {
		if bodyReader != nil {
			defer bodyReader.Close()
		}
	}()

	if resp.Header.Get("Content-Encoding") == "gzip" {
		bodyReader, _ = gzip.NewReader(resp.Body)
	} else {
		bodyReader = resp.Body
	}

	return goquery.NewDocumentFromReader(bodyReader)
}
