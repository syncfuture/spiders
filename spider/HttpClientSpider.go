package spider

import (
	"compress/gzip"
	"io"
	"net/http"

	"github.com/syncfuture/go/spool"

	"github.com/syncfuture/spiders/spider/store"
)

var (
	_userAgents = []string{
		"Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/87.0.4280.141 Safari/537.36",
		"Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/87.0.4280.141 Safari/537.36 Edg/87.0.664.75",
		"Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/87.0.4280.67 Safari/537.36 Edg/87.0.664.55",
		"Mozilla/5.0 (Windows NT 10.0; WOW64; Trident/7.0; rv:11.0) like Gecko",
	}
	_headers = map[string]string{
		// "User-Agent":                "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/87.0.4280.67 Safari/537.36 Edg/87.0.664.55",
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
)

type HttpClientSpider struct {
	client     *http.Client
	proxyStore store.IProxyStore
	pool       spool.BufferPool
}

func (x *HttpClientSpider) Get(url string) (string, error) {
	msg, _ := http.NewRequest("GET", url, nil)

	for k, v := range _headers {
		msg.Header.Add(k, v)
	}
	userAgent := _userAgents[randInt(0, len(_userAgents)-1)]
	msg.Header.Add("User-Agent", userAgent)

	resp, err := x.client.Do(msg)
	if err != nil {
		return "", err
	}

	var bodyReader io.ReadCloser

	buffer := x.pool.GetBuffer()
	defer func() {
		x.pool.PutBuffer(buffer)
		if bodyReader != nil {
			defer bodyReader.Close()
		}
	}()

	if resp.Header.Get("Content-Encoding") == "gzip" {
		bodyReader, _ = gzip.NewReader(resp.Body)
	} else {
		bodyReader = resp.Body
	}

	buffer.ReadFrom(bodyReader)

	return buffer.String(), err
}
