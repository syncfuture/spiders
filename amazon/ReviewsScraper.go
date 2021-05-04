package amazon

import (
	"fmt"
	"strconv"
	"time"

	"github.com/PuerkitoBio/goquery"
	"github.com/SyncSoftInc/proxy/protoc/proxy"
	log "github.com/syncfuture/go/slog"
	"github.com/syncfuture/spiders/scraper/httpClient"
	"github.com/syncfuture/spiders/scraper/store"
)

type ReviewsScraper struct {
	skuBase
	// PageInfo *PageInfo
}

func NewReviewsScraper(proxyStore store.IProxyStore, asin string) (r *ReviewsScraper) {
	r = new(ReviewsScraper)
	r.ProxyStore = proxyStore
	r.ASIN = asin
	// r.PageInfo = &PageInfo{
	// 	// PageSize:   10,
	// 	// proxyStore: r.ProxyStore,
	// 	url: r.buildURL(1),
	// }
	// err = r.PageInfo.FetchPageInfo()
	return
}

func (x *ReviewsScraper) FetchPages(fromDate *time.Time) (r []*ReviewDTO, err error) {
	r = make([]*ReviewDTO, 0, 50)
	var p *proxy.Proxy
	rs := &ReviewResult{
		NextPageURL: x.buildURL(1),
	}

	for rs.NextPageURL != "" {
		rs, p, err = x.FetchPage(rs.NextPageURL)
		if err != nil {
			if p.Score <= 0 {
				// 代理有问题，重试一次
				rs, p, err = x.FetchPage(rs.NextPageURL)
				if err != nil {
					// 再次失败就终止
					return nil, err
				}
			} else {
				return nil, err
			}
		}

		for _, review := range rs.Reviews {
			if fromDate != nil && review.CreatedOn.Before(*fromDate) {
				// 超出给定时间，终止获取
				return r, nil
			} else {
				r = append(r, review)
			}
		}
	}
	return
}

func (x *ReviewsScraper) FetchPage(url string) (*ReviewResult, *proxy.Proxy, error) {
	// if pageIndex <= 0 || pageIndex > x.PageInfo.TotalPage {
	// 	log.Fatal("invalid pageindex")
	// }
	httpScraper := httpClient.NewRandomScraper(x.ProxyStore, nil)
	if url == "" {
		url = x.buildURL(1)
	}
	result, err := httpScraper.Get(url)
	if err != nil {
		return &ReviewResult{NextPageURL: url}, result.Proxy, err
	}

	doc := result.GetHtmlDocument()
	reviewNodes := doc.Find(_reviewDomSelector)
	log.Debugf("[%s] %d review nodes found", result.Proxy.URI, reviewNodes.Length())
	if reviewNodes.Length() == 0 {
		return &ReviewResult{
			NextPageURL: "", // 没有评论，自然没有下一页
			Reviews:     make([]*ReviewDTO, 0),
		}, result.Proxy, nil
	}

	r := new(ReviewResult)
	r.Reviews = make([]*ReviewDTO, 0, reviewNodes.Length())

	reviewNodes.Each(func(i int, node *goquery.Selection) {
		dto := new(ReviewDTO)
		dto.AmazonID = x.getID(node)
		if dto.AmazonID == "" {
			log.Warnf("[%s]fetched empty amazon review id ", x.ASIN)
			return
		}
		dto.StripInfo, dto.ASIN = x.getStripInfo(node)
		if dto.ASIN == "" {
			dto.ASIN = x.ASIN
		}
		dto.Rating = x.getRating(node)
		dto.Location, dto.CreatedOn = x.getLocationDate(node)
		dto.CustomerName = x.getCustomerName(node)
		dto.Title = x.getTitle(node)
		dto.Content = x.getContent(node)
		dto.IsVerified = x.getVerified(node)
		r.Reviews = append(r.Reviews, dto)
	})

	nextLink := doc.Find("ul.a-pagination li.a-last a")
	if nextLink.Length() > 0 {
		r.NextPageURL, _ = nextLink.Attr("href")
		r.NextPageURL = _amazonURLRoot + r.NextPageURL
	}
	return r, result.Proxy, nil
}

func (x *ReviewsScraper) buildURL(pageIndex int) string {
	return fmt.Sprintf(_reviewURLBase, x.ASIN, pageIndex)
}
func (x *ReviewsScraper) getID(node *goquery.Selection) string {
	id, _ := node.Attr("id")
	return id
}
func (x *ReviewsScraper) getRating(node *goquery.Selection) float32 {
	nodeText := node.Find(_reviewRatingSelector).Text()
	nodeArray := _reviewRatingRegex.FindStringSubmatch(nodeText)
	if len(nodeArray) == 2 {
		a, err := strconv.ParseFloat(nodeArray[1], 32)
		if err == nil {
			return float32(a)
		}
	}
	log.Debugf("get rating failed, '%s'->'%s'", _reviewRatingSelector, nodeText)
	return -1
}
func (x *ReviewsScraper) getLocationDate(node *goquery.Selection) (string, *time.Time) {
	nodeText := node.Find(_reviewDateSelector).Text()
	nodeArray := _reviewDateRegex.FindStringSubmatch(nodeText)
	var country, dateStr string
	if len(nodeArray) < 3 {
		log.Debugf("get rating failed, '%s'->'%s'", _reviewRatingSelector, nodeText)
		return "", nil
	} else if len(nodeArray) == 3 {
		country = nodeArray[1]
		dateStr = nodeArray[2]
	} else if len(nodeArray) == 4 {
		country = nodeArray[2]
		dateStr = nodeArray[3]
	}

	date, err := time.ParseInLocation("January 2, 2006", dateStr, _pst)
	if err != nil {
		return country, nil
	}

	date = date.UTC()

	return country, &date
}
func (x *ReviewsScraper) getCustomerName(node *goquery.Selection) string {
	r := node.Find(_reviewUserNameSelector).Text()
	return r
}
func (x *ReviewsScraper) getTitle(node *goquery.Selection) string {
	r := node.Find(_reviewTitleSelector).Text()
	r = TrimSpaceAndLines(r)
	return r
}
func (x *ReviewsScraper) getContent(node *goquery.Selection) string {
	r := node.Find(_reviewContentSelector).Text()
	r = TrimSpaceAndLines(r)
	return r
}
func (x *ReviewsScraper) getVerified(node *goquery.Selection) bool {
	r := node.Find(_reviewVerifiedSelector).Text()
	return r == _reviewVPText
}
func (x ReviewsScraper) getStripInfo(node *goquery.Selection) (sizeColor string, asin string) {
	stripNode := node.Find(_reviewStripSelector)
	if stripNode.Length() == 0 {
		return
	}
	sizeColor = stripNode.Text()
	href, ok := stripNode.Attr("href")
	if !ok {
		return
	}

	matches := _reviewASINRegex.FindStringSubmatch(href)
	if len(matches) < 2 {
		return
	}

	asin = matches[1]
	return
}
