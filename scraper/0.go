package scraper

import (
	"errors"
	"regexp"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/SyncSoftInc/proxy/protoc/proxy"
	"github.com/syncfuture/go/u"
)

const (
	_spaceRegexStr = `\s{2,}`
	_line          = "\n"
)

var (
	_spaceRegex = regexp.MustCompile(_spaceRegexStr)
	NotFound    = errors.New("404 Not Found")
)

func NewScrapeResult(statusCode int, doc *goquery.Document, p *proxy.Proxy, requestHeaders map[string]string, responseHeaders map[string][]string) *ScrapeResult {
	if p == nil {
		p = new(proxy.Proxy)
	}

	respHeaders := make(map[string]string, len(responseHeaders))
	if responseHeaders != nil {
		for k, v := range responseHeaders {
			respHeaders[k] = v[0]
		}
	}
	return &ScrapeResult{
		StatusCode:      statusCode,
		Proxy:           p,
		doc:             doc,
		RequestHeaders:  requestHeaders,
		ResponseHeaders: respHeaders,
	}
}

type IScraper interface {
	Get(url string) (*ScrapeResult, error)
}

type ScrapeResult struct {
	StatusCode      int
	Proxy           *proxy.Proxy
	doc             *goquery.Document
	RequestHeaders  map[string]string
	ResponseHeaders map[string]string
}

func (x *ScrapeResult) GetHtmlDocument() *goquery.Document {
	return x.doc
}

func (x *ScrapeResult) IsSucessStatusCode() bool {
	return x.StatusCode < 300
}

func (x *ScrapeResult) ToCompactHtml() string {
	html := x.ToHtml()
	if html == "" {
		return ""
	}
	html = _spaceRegex.ReplaceAllString(html, " ")
	return strings.ReplaceAll(html, _line, "")
}

func (x *ScrapeResult) ToHtml() string {
	if x.doc == nil {
		return ""
	}

	html, err := x.doc.Html()
	if u.LogError(err) {
		return ""
	}
	return html
}
