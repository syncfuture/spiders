package amazon

import (
	"regexp"
	"strconv"
	"time"

	log "github.com/syncfuture/go/slog"
	"github.com/syncfuture/go/u"
	"github.com/syncfuture/spiders/scraper/store"
)

const (
	_amazonURLRoot          = `https://www.amazon.com`
	_amazonRef              = "ref"
	_amazonRefValue1        = "ppx_yo_dt_b_asin_title_o03_s00"
	_amazonIE               = "ie"
	_amazonIEValue1         = "UTF8"
	_amazonPSC              = "psc"
	_amazonPSCValue1        = "1"
	_amazonPriceRegexStr    = `\$(\d+(\.\d+)?)`
	_reviewDomSelector      = `div[data-hook="review"]`
	_reviewStripSelector    = `.review-format-strip a[data-hook="format-strip"]`
	_reviewASINRegexStr     = `/product-reviews/(\w+)/`
	_reviewDateSelector     = `.review-date`
	_reviewRatingSelector   = `.review-rating`
	_reviewUserNameSelector = `.a-profile-name`
	_reviewTitleSelector    = `.review-title`
	_reviewContentSelector  = `.review-text-content`
	_reviewVerifiedSelector = `.review-format-strip`
	_reviewRatingRegexStr   = `(\d+.\d+) out of 5 stars`
	_spaceLineRegexStr      = `\s{3,}|\n`
	_reviewVPText           = "Verified Purchase"
	_reviewDateRegexStr     = `Reviewed in (the )?(.+) on (\w+ \d+, \d+)` // Reviewed in the United States on January 20, 2021
	// _reviewURLBase          = "https://www.amazon.com/dp/product-reviews/%s/ref=cm_cr_arp_d_viewopt_srt?ie=UTF8&reviewerType=all_reviews&sortBy=recent&pageNumber=%d"
	_reviewURLBase          = "https://www.amazon.com/dp/product-reviews/%s/ref=cm_cr_getr_d_rvw_fmt?ie=UTF8&formatType=current_format&sortBy=recent&pageNumber=%d"
	_offerURLURLBase        = "https://www.amazon.com/gp/offer-listing/"
	_offerURLRef            = "/ref=olp_page_"
	_offerURLPrime          = "f_primeEligible"
	_offerURLNew            = "f_new"
	_offerURLFreeShipping   = "f_freeShipping"
	_offerURLStartIndex     = "startIndex"
	_offerURLAll            = "f_all"
	_offserDomSelector      = ".olpOffer"
	_offerPriceSelector     = ".olpOfferPrice"
	_offerSellerSelector    = ".olpSellerName"
	_offerConditionSelector = ".olpCondition"
	_offerShippingSelector  = ".olpShippingInfo"
	_offerPageSelector      = ".a-pagination li.a-last a"
	_offerFreeShippingtext  = "FREE Shipping"
	_skuURLBase             = "https://www.amazon.com/gp/product/"
	_skuPriceSelector1      = "#price_inside_buybox"
	_skuPriceSelector2      = "#newBuyBoxPrice"
	// _skuBuyBoxSelector      = "#buybox"
)

var (
	_amazonPriceRegex  = regexp.MustCompile(_amazonPriceRegexStr)
	_reviewASINRegex   = regexp.MustCompile(_reviewASINRegexStr)
	_reviewRatingRegex = regexp.MustCompile(_reviewRatingRegexStr)
	_reviewDateRegex   = regexp.MustCompile(_reviewDateRegexStr)
	_pst               *time.Location
)

func init() {
	var err error
	_pst, err = time.LoadLocation("America/Los_Angeles")
	u.LogFaltal(err)
}

type skuBase struct {
	ASIN       string
	ProxyStore store.IProxyStore
}

func (x *skuBase) getPriceFromString(priceStr string) float32 {
	array := _amazonPriceRegex.FindStringSubmatch(priceStr)
	if len(array) < 2 {
		log.Warn(priceStr)
		return -1
	}
	r, err := strconv.ParseFloat(array[1], 32)
	if u.LogError(err) {
		return -1
	}
	return float32(r)
}
