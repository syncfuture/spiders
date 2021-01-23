package revidews

import (
	"errors"
	"math"
	"regexp"
	"strconv"
	"strings"

	log "github.com/syncfuture/go/slog"
	"github.com/syncfuture/scraper"
)

const (
	_pageInfoDomID    = "#filter-info-section"
	_pageInfoRegexStr = `(\d+) global reviews`
)

var (
	_pageInfoRegex = regexp.MustCompile(_pageInfoRegexStr)
)

type PageInfo struct {
	PageSize    float64
	TotalCount  float64
	TotalPage   int
	url         string
	scraperPool *scraper.ScraperPool
}

func (x *PageInfo) FetchPageInfo() error {
	if x.PageSize <= 0 {
		log.Fatal("invalid pagesize")
	}

	httpScraper, err := x.scraperPool.GetHttpClientScraper()
	if err != nil {
		return err
	}
	doc, err := httpScraper.GetDocument(x.url)
	if err != nil {
		return err
	}

	infoSecsion := doc.Find(_pageInfoDomID)
	if infoSecsion.Length() == 0 {
		html, _ := doc.Html()
		log.Debug(httpScraper.GetProxy().Host, strings.Replace(html, "\n", "", -1))
		// 有可能是代理被封了，移除掉
		x.scraperPool.ProxyStore.Remove(httpScraper.GetProxy())
		return errors.New("cannot find total reviews dom")
	}

	nodeText := infoSecsion.Text()
	matches := _pageInfoRegex.FindStringSubmatch(nodeText)
	if len(matches) != 2 {
		log.Debug(httpScraper.GetProxy().Host, nodeText)
		return errors.New("cannot find total reviews text")
	}

	x.TotalCount, err = strconv.ParseFloat(matches[1], 32)
	if err != nil {
		log.Debug(matches)
		return err
	}
	x.TotalPage = int(math.Ceil(x.TotalCount / x.PageSize))

	return nil
}
