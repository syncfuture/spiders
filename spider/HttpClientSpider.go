package spider

import (
	"compress/gzip"
	"io"
	"net/http"

	"github.com/PuerkitoBio/goquery"

	"github.com/syncfuture/spiders/spider/store"
)

type HttpClientSpider struct {
	client     *http.Client
	headers    map[string]string
	proxyStore store.IProxyStore
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
