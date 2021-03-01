package wayfair

import (
	"context"
	"fmt"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/chromedp/chromedp"

	"github.com/syncfuture/go/sconfig"
	log "github.com/syncfuture/go/slog"
	"github.com/syncfuture/scraper/scdp"
	"github.com/syncfuture/scraper/store"
	"github.com/syncfuture/spiders/wayfair/model"
)

const (
	_expError  = "ERR_INVALID_AUTH_CREDENTIALS"
	_urlFormat = "https://www.wayfair.com/bed-bath/pdp/product-%s.html"
)

var (
	_helpfulRegex = regexp.MustCompile(`\d+`)
)

type skuBase struct {
	Item       *model.ItemDTO
	ProxyStore store.IProxyStore
	CDPClient  *scdp.ChromeDPClient
}

type ReviewsScraper struct {
	skuBase
	// PageInfo *PageInfo
}

func NewReviewsScraper(cp sconfig.IConfigProvider, proxyStore store.IProxyStore, item *model.ItemDTO) (r *ReviewsScraper) {
	r = new(ReviewsScraper)
	r.ProxyStore = proxyStore
	r.Item = item
	r.CDPClient = scdp.NewChromeDPClient(cp)
	return
}

func (x *ReviewsScraper) FetchReviews() (reviews []*model.ReviewDTO, err error) {
	url := fmt.Sprintf(_urlFormat, x.Item.SKU)
	proxy := x.ProxyStore.Lease()
	defer x.ProxyStore.Return(proxy)

	proxyURL := proxy.ToURL(false).String()
	x.CDPClient.FetchWithProxy(proxyURL, func(ctx context.Context) {
		// client.FetchWithHead(func(ctx context.Context) {
		// 跳转去产品页
		err = scdp.Navigate(ctx, url, 1920, 936)
		if err != nil {
			if strings.Contains(err.Error(), _expError) {
				proxy.Expired = true
			}
			return
		}

		// 检查是否被人机验证阻止
		captchaNodes := scdp.GetNodes(ctx, "h1.Captcha-title")
		if len(captchaNodes) > 0 {
			log.Warnf("%s blocked by captcha", proxy.Host)
			proxy.Blocked = true
			return
		}

		// 选择评论按时间排序，等待评论加载完毕
		err = chromedp.Run(ctx,
			chromedp.WaitVisible(`.ReviewsSearchSortFilter-dropdown`, chromedp.ByQuery),
			chromedp.Click(`.ReviewsSearchSortFilter-dropdown`, chromedp.ByQuery),
			chromedp.WaitVisible(`#productReviewsSubheaderDropdown-item-2`, chromedp.ByQuery),
			chromedp.Click(`#productReviewsSubheaderDropdown-item-2`, chromedp.ByQuery),
			chromedp.WaitVisible(`div.ProductReviewList-links`, chromedp.ByQuery),
		)
		if err != nil {
			log.Debug(err)
			return
		}

		for loadReviews(ctx, time.Now().Add(time.Hour*24*-30)) { // 加载30天内的评论
		}

		// 获取页面上所有评论
		timer := time.Now()
		reviews = x.getReviews(ctx)
		elapsed := time.Since(timer)
		log.Infof("[%d] found %d reviews", elapsed.Milliseconds(), len(reviews))
	})

	if proxy.Blocked || proxy.Expired {
		return x.FetchReviews() // 代理出问题，重试
	}
	return
}

func (x *ReviewsScraper) getReviews(ctx context.Context) (r []*model.ReviewDTO) {
	reviewNodes := scdp.GetNodes(ctx, "article.ProductReview")
	r = make([]*model.ReviewDTO, 0, len(reviewNodes))

	for _, reviewNode := range reviewNodes {
		dto := new(model.ReviewDTO)
		dto.Items = x.Item.Items
		dto.SKU = x.Item.SKU
		// var photoNodes []*cdp.Node
		// var date, helpful string
		// chromedp.Run(ctx,
		// 	chromedp.Text(`.ProductReview-reviewDetails .ProductReview-comments`, &dto.Comments, chromedp.ByQuery, chromedp.AtLeast(0), chromedp.FromNode(reviewNode)),
		// 	chromedp.Text(`.ProductReview-reviewDetails .pl-ReviewStars`, &dto.Rating, chromedp.ByQuery, chromedp.AtLeast(0), chromedp.FromNode(reviewNode)),
		// 	chromedp.Text(`.ProductReview-reviewDetails p[data-hb-id="pl-text"]`, &date, chromedp.ByQuery, chromedp.AtLeast(0), chromedp.FromNode(reviewNode)),
		// 	chromedp.Text(`.ProductReview-customerInfo .ProductReviewCustomerInfo-name`, &dto.Name, chromedp.ByQuery, chromedp.AtLeast(0), chromedp.FromNode(reviewNode)),
		// 	chromedp.Text(`.ProductReview-customerInfo .ProductReviewerComplianceBadge-tooltip`, &dto.Badge, chromedp.ByQuery, chromedp.AtLeast(0), chromedp.FromNode(reviewNode)),
		// 	chromedp.Text(`.ProductReview-helpfulButton`, &helpful, chromedp.ByQuery, chromedp.AtLeast(0), chromedp.FromNode(reviewNode)),
		// 	chromedp.Nodes(`.ProductReview-reviewPhotos .ProductReviewPhotos-item img`, &photoNodes, chromedp.ByQueryAll, chromedp.AtLeast(0), chromedp.FromNode(reviewNode)),
		// )
		// dto.Photos = getPhotos(photoNodes)
		// dto.Helpful = getHelpful(helpful)
		dto.Comments = scdp.GetText(ctx, ".ProductReview-reviewDetails .ProductReview-comments", chromedp.FromNode(reviewNode))
		dto.Rating = scdp.GetText(ctx, ".ProductReview-reviewDetails .pl-ReviewStars", chromedp.FromNode(reviewNode))
		date := scdp.GetText(ctx, `.ProductReview-reviewDetails p[data-hb-id="pl-text"]`, chromedp.FromNode(reviewNode))
		dto.Date, _ = time.Parse("01/02/2006", date)
		dto.Photos = scdp.GetAttributes(ctx, `.ProductReview-reviewPhotos .ProductReviewPhotos-item img`, "src", chromedp.FromNode(reviewNode))
		dto.Name = scdp.GetText(ctx, `.ProductReview-customerInfo .ProductReviewCustomerInfo-name`, chromedp.FromNode(reviewNode))
		dto.Badge = scdp.GetText(ctx, `.ProductReview-customerInfo .ProductReviewerComplianceBadge-tooltip`, chromedp.FromNode(reviewNode))
		helpful := scdp.GetText(ctx, `.ProductReview-helpfulButton`, chromedp.FromNode(reviewNode))
		dto.Helpful = getHelpful(helpful)

		r = append(r, dto)
	}

	return
}

func getHelpful(str string) int {
	helpfulStr := _helpfulRegex.FindString(str)
	if helpfulStr != "" {
		r, _ := strconv.Atoi(helpfulStr)
		return r
	}
	return 0
}

func loadReviews(ctx context.Context, from time.Time) bool {
	reviewNodes := scdp.GetNodes(ctx, "article.ProductReview")

	count := len(reviewNodes)

	if count > 0 {
		lastReviewNode := reviewNodes[count-1]
		dateText := scdp.GetText(ctx, `.ProductReview-reviewDetails p[data-hb-id="pl-text"]`, chromedp.FromNode(lastReviewNode))
		date, _ := time.Parse("01/02/2006", dateText)
		if date.After(from) {
			return loadMore(ctx)
		}
	}

	return false
}

func loadMore(ctx context.Context) bool {
	loadButtons := scdp.GetNodes(ctx, "div.ProductReviewList-links button")
	if len(loadButtons) > 0 {
		// 有加载按钮
		buttonText := scdp.GetText(ctx, ".Button-content", chromedp.FromNode(loadButtons[0]))
		if strings.Contains(buttonText, "More") { // 是加载更多按钮
			scdp.Click(ctx, `div.ProductReviewList-links button`)
			return true
		} else {
			return false
		}
	} else {
		// 没有加载按钮
		return false
	}
}
