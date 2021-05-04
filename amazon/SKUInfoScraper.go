package amazon

import (
	"net/url"

	"github.com/PuerkitoBio/goquery"
	"github.com/syncfuture/go/srand"
	"github.com/syncfuture/spiders/scraper/httpClient"
	"github.com/syncfuture/spiders/scraper/store"
)

type SKUInfoScraper struct {
	skuBase
}

func NewSKUInfoScraper(proxyStore store.IProxyStore, asin string) *SKUInfoScraper {
	r := new(SKUInfoScraper)
	r.ASIN = asin
	r.ProxyStore = proxyStore
	return r
}

func (x *SKUInfoScraper) Fetch() (*AmazonItemDTO, error) {
	httpScraper := httpClient.NewRandomScraper(x.ProxyStore, nil)
	url := x.buildURL()
	result, err := httpScraper.Get(url)
	if err != nil {
		return nil, err
	}

	r := new(AmazonItemDTO)
	r.ASIN = x.ASIN
	doc := result.GetHtmlDocument()
	r.Price = x.getPrice(doc.Selection)

	return r, err
}

func (x *SKUInfoScraper) buildURL() string {
	r := _skuURLBase + x.ASIN
	seed := srand.IntRange(0, 1)
	if seed == 1 {
		r = r + "/" + _amazonRef + "=" + _amazonRefValue1
	}

	uri, _ := url.Parse(r)

	q := uri.Query()
	seed = srand.IntRange(0, 1)
	if seed == 1 {
		q.Set(_amazonIE, _amazonIEValue1)
	}
	seed = srand.IntRange(0, 1)
	if seed == 1 {
		q.Set(_amazonPSC, _amazonPSCValue1)
	}

	r = r + "?" + q.Encode()
	return r
}

func (x *SKUInfoScraper) getPrice(selection *goquery.Selection) float32 {
	priceNodes := selection.Find(_skuPriceSelector1)
	if priceNodes.Length() == 0 {
		priceNodes = selection.Find(_skuPriceSelector2)
		if priceNodes.Length() == 0 {
			return -1
		}
	}
	priceStr := priceNodes.Text()

	return x.getPriceFromString(priceStr)
}
