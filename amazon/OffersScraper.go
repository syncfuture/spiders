package amazon

import (
	"fmt"
	"net/url"
	"strconv"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/syncfuture/scraper"
	"github.com/syncfuture/scraper/httpClient"
	"github.com/syncfuture/scraper/store"
)

type OffersScraper struct {
	skuBase
	pageSize         int
	newOnly          bool
	primeOnly        bool
	freeShippingOnly bool
}

func NewOffsersScraper(proxyStore store.IProxyStore, asin string, primeOnly, freeShippingOnly, newOnly bool) (r *OffersScraper) {
	r = new(OffersScraper)
	r.ProxyStore = proxyStore
	r.ASIN = asin
	r.newOnly = newOnly
	r.primeOnly = primeOnly
	r.freeShippingOnly = freeShippingOnly
	r.pageSize = 10
	return r
}

func (x *OffersScraper) FetchPagedOffers(cusor *OfferResultDTO) (*OfferResultDTO, error) {
	if cusor == nil {
		cusor = &OfferResultDTO{
			NextPageURL: x.buildURL(1),
		}
	} else if cusor.NextPageURL == "" {
		return nil, nil
	}
	httpScraper := httpClient.NewRandomScraper(x.ProxyStore, nil)
	result, err := httpScraper.Get(cusor.NextPageURL)
	if err != nil {
		return nil, err
	}

	doc := result.GetHtmlDocument()
	offerNodes := doc.Find(_offserDomSelector)
	if offerNodes.Length() == 0 {
		// log.Debug(doc.ToCompactHtml())newBuyBoxPrice price_inside_buybox
		return nil, fmt.Errorf("cannot find offer list dom '%s'", _offserDomSelector)
	}

	offers := make([]*OfferDTO, 0, offerNodes.Length())

	offerNodes.Each(func(i int, node *goquery.Selection) {
		dto := new(OfferDTO)
		dto.Price = x.getPrice(node)
		dto.Seller = x.getSeller(node)
		dto.Condition = x.getCondition(node)
		dto.FreeShipping = x.getShipping(node)
		offers = append(offers, dto)
	})

	return &OfferResultDTO{
		Offers:      offers,
		NextPageURL: x.getNextPageURL(result),
	}, nil
}

func (x *OffersScraper) FetchAllOffers() ([]*OfferDTO, error) {
	r := make([]*OfferDTO, 0, 10)

	var results *OfferResultDTO
	var err error
	for {
		results, err = x.FetchPagedOffers(results)
		if err != nil {
			return nil, err
		} else if results == nil {
			break
		}
		r = append(r, results.Offers...)
	}

	return r, nil
}

func (x *OffersScraper) buildURL(pageIndex int) (r string) {
	r = _offerURLURLBase + x.ASIN + _offerURLRef + strconv.Itoa(pageIndex)
	a, _ := url.Parse(r)

	query := a.Query()
	query.Set(_amazonIE, "UTF8")
	query.Set(_offerURLAll, "true")

	if pageIndex > 1 {
		offset := (pageIndex - 1) * x.pageSize
		query.Set(_offerURLStartIndex, strconv.Itoa(offset))
	}

	if x.primeOnly {
		query.Set(_offerURLPrime, "true")
	}

	if x.newOnly {
		query.Set(_offerURLNew, "true")
	}

	if x.freeShippingOnly {
		query.Set(_offerURLFreeShipping, "true")
	}

	r = r + "?" + query.Encode()

	return
}

func (x *OffersScraper) getPrice(node *goquery.Selection) float32 {
	str := node.Find(_offerPriceSelector).Text()
	return x.getPriceFromString(str)
}

func (x *OffersScraper) getSeller(node *goquery.Selection) string {
	str := node.Find(_offerSellerSelector).Text()
	return TrimSpaceAndLines(str)
}

func (x *OffersScraper) getCondition(node *goquery.Selection) string {
	str := node.Find(_offerConditionSelector).Text()
	return TrimSpaceAndLines(str)
}

func (x *OffersScraper) getShipping(node *goquery.Selection) bool {
	str := node.Find(_offerShippingSelector).Text()
	return strings.Contains(str, _offerFreeShippingtext)
}

func (x *OffersScraper) getNextPageURL(result *scraper.ScrapeResult) string {
	doc := result.GetHtmlDocument()
	nodes := doc.Find(_offerPageSelector)
	if nodes.Length() == 0 {
		return ""
	}

	href, ok := nodes.Attr("href")
	if !ok {
		return ""
	}
	return _amazonURLRoot + href
}
