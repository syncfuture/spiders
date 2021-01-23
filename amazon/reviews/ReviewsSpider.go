package revidews

import (
	"fmt"
	"regexp"
	"strconv"
	"time"

	"github.com/PuerkitoBio/goquery"
	"github.com/syncfuture/go/sid"
	log "github.com/syncfuture/go/slog"
	"github.com/syncfuture/scraper"
	"github.com/syncfuture/spiders/amazon/models"
)

const (
	_reviewsDomSelector = `div[data-hook="review"]`
	_dateSelector       = `.review-date`
	_ratingSelector     = `.review-rating`
	_nameSelector       = `.a-profile-name`
	_titleSelector      = `.review-title`
	_contentSelector    = `.review-text-content`
	_verifiedSelector   = `.review-format-strip`
	_ratingRegexStr     = `(\d+.\d+) out of 5 stars`
	_spaceRegexStr      = `\s{3,}|\n`
	_vpText             = "Verified Purchase"
	_dateRegexStr       = `Reviewed in (the )?(.+) on (\w+ \d+, \d+)` // Reviewed in the United States on January 20, 2021
	_urlTemplate        = "https://www.amazon.com/dp/product-reviews/%s/ref=cm_cr_arp_d_viewopt_srt?ie=UTF8&reviewerType=all_reviews&sortBy=recent&pageNumber=%d"
)

var (
	_ratingRegex = regexp.MustCompile(_ratingRegexStr)
	_dateRegex   = regexp.MustCompile(_dateRegexStr)
	_spaceRegex  = regexp.MustCompile(_spaceRegexStr)
)

type ReviewsSpider struct {
	asin        string
	scraperPool *scraper.ScraperPool
	PageInfo    *PageInfo
	idGenerator sid.IIDGenerator
}

func NewReviewsSpider(scraperPool *scraper.ScraperPool, asin string) (r *ReviewsSpider, err error) {
	r = new(ReviewsSpider)
	r.scraperPool = scraperPool
	r.asin = asin
	r.idGenerator = sid.NewSonyflakeIDGenerator()
	r.PageInfo = &PageInfo{
		PageSize:    10,
		scraperPool: r.scraperPool,
		url:         r.getURL(1),
	}
	err = r.PageInfo.FetchPageInfo()
	return
}
func (x *ReviewsSpider) FetchAll() []*models.ReviewDTO {
	for i := 1; i <= x.PageInfo.TotalPage; i++ {

	}
	return nil
}

func (x *ReviewsSpider) FetchIndex(pageIndex int) ([]*models.ReviewDTO, error) {
	if pageIndex <= 0 || pageIndex > x.PageInfo.TotalPage {
		log.Fatal("invalid pageindex")
	}
	httpScraper, err := x.scraperPool.GetHttpClientScraper()
	if err != nil {
		return nil, err
	}
	url := x.getURL(pageIndex)
	doc, err := httpScraper.GetDocument(url)

	reviewNodes := doc.Find(_reviewsDomSelector)
	if reviewNodes.Length() == 0 {
		return nil, fmt.Errorf("cannot find reviews dom '%s'", _reviewsDomSelector)
	}

	r := make([]*models.ReviewDTO, 0, reviewNodes.Length())

	reviewNodes.Each(func(i int, node *goquery.Selection) {
		dto := new(models.ReviewDTO)
		dto.ID = x.idGenerator.GenerateString()
		dto.AmazonID = x.getID(node)
		dto.Rating = x.getRating(node)
		dto.Location, dto.CreatedOn = x.getLocationDate(node)
		dto.CustomerName = x.getCustomerName(node)
		dto.Title = x.getTitle(node)
		dto.Content = x.getContent(node)
		dto.IsVerified = x.getVerified(node)
		r = append(r, dto)
	})

	return r, nil
}

func (x *ReviewsSpider) getURL(pageIndex int) string {
	return fmt.Sprintf(_urlTemplate, x.asin, pageIndex)
}
func (x *ReviewsSpider) getID(node *goquery.Selection) string {
	id, _ := node.Attr("id")
	return id
}
func (x *ReviewsSpider) getRating(node *goquery.Selection) float32 {
	nodeText := node.Find(_ratingSelector).Text()
	nodeArray := _ratingRegex.FindStringSubmatch(nodeText)
	if len(nodeArray) == 2 {
		a, err := strconv.ParseFloat(nodeArray[1], 32)
		if err == nil {
			return float32(a)
		}
	}
	log.Debugf("get rating failed, '%s'->'%s'", _ratingSelector, nodeText)
	return -1
}

func (x *ReviewsSpider) getLocationDate(node *goquery.Selection) (string, *time.Time) {
	nodeText := node.Find(_dateSelector).Text()
	nodeArray := _dateRegex.FindStringSubmatch(nodeText)
	var locStr, dateStr string
	if len(nodeArray) < 3 {
		log.Debugf("get rating failed, '%s'->'%s'", _ratingSelector, nodeText)
		return "", nil
	} else if len(nodeArray) == 3 {
		locStr = nodeArray[1]
		dateStr = nodeArray[2]
	} else if len(nodeArray) == 4 {
		locStr = nodeArray[2]
		dateStr = nodeArray[3]
	}

	date, err := time.Parse("January 2, 2006", dateStr)
	if err != nil {
		return locStr, nil
	}

	return locStr, &date
}

func (x *ReviewsSpider) getCustomerName(node *goquery.Selection) string {
	r := node.Find(_nameSelector).Text()
	return r
}

func (x *ReviewsSpider) getTitle(node *goquery.Selection) string {
	r := node.Find(_titleSelector).Text()
	r = _spaceRegex.ReplaceAllString(r, "")
	return r
}

func (x *ReviewsSpider) getContent(node *goquery.Selection) string {
	r := node.Find(_contentSelector).Text()
	r = _spaceRegex.ReplaceAllString(r, "")
	return r
}

func (x *ReviewsSpider) getVerified(node *goquery.Selection) bool {
	r := node.Find(_verifiedSelector).Text()
	return r == _vpText
}
