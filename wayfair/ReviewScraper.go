package wayfair

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/url"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/chromedp/cdproto"
	"github.com/chromedp/cdproto/cdp"
	"github.com/chromedp/cdproto/network"
	"github.com/chromedp/chromedp"

	"github.com/syncfuture/go/serr"
	log "github.com/syncfuture/go/slog"
	"github.com/syncfuture/go/u"
	"github.com/syncfuture/scraper/scdp"
	"github.com/syncfuture/spiders/wayfair/model"
)

const (
	// _expError   = "ERR_INVALID_AUTH_CREDENTIALS"
	// _retryError = "ERR_TOO_MANY_RETRIES"
	_urlFormat  = "https://www.wayfair.com/category/pdp/product-%s.html"
	_ajaxURL    = "https://www.wayfair.com/graphql?hash="
	_dateFormat = "01/02/2006"
)

var (
	_reviewCountRegex = regexp.MustCompile(`(\d+) Reviews`)
	_errBlocked       = errors.New("BLOCKED")
	_errNotFound      = errors.New("NOT FOUND")
)

type skuBase struct {
	// Item      *model.ItemDTO
	CDPClient *scdp.ChromeDPClient
}

type ReviewsScraper struct {
	skuBase
	// PageInfo *PageInfo
}

func NewReviewsScraper(client *scdp.ChromeDPClient) (r *ReviewsScraper) {
	r = new(ReviewsScraper)
	// r.Item = item
	r.CDPClient = client
	return
}

func (x *ReviewsScraper) FetchReviews(item *model.ItemDTO, from time.Time) (r []*model.ReviewDTO, err error) {
	dic := make(map[int]*model.ReviewDTO)
	proxy := x.CDPClient.ProxyStore.Rent()
	proxyURL, err := url.Parse(proxy.URI)
	if err != nil {
		return nil, serr.WithStack(err)
	}

	err = x.CDPClient.FetchWithProxy(proxyURL, func(mainCtx context.Context) error {
		// 监听ajax事件，获取评论
		chromedp.ListenTarget(mainCtx, func(ev interface{}) {
			go captureReviewsFromAjax(ev, mainCtx, item, from, &dic)
		})

		// client.FetchWithHead(func(ctx context.Context) {
		// 跳转去产品页
		var url string
		if item.URL == "" {
			url = generateURL(item)
		} else {
			if notFound(item.URL) { // 跳转去了列表页，直接返回错误
				return _errNotFound
			}
			url = item.URL
		}
		err := scdp.NavigateWithAuth(mainCtx, url, 1920, 936)
		if err != nil {
			// if strings.Contains(err.Error(), _expError) || strings.Contains(err.Error(), _retryError) {
			// 	proxy.Expired = true
			// }
			// proxy.Expired = true // 导航失败，直接视作代理池过期
			proxy.Score = -1
			return err
		}

		// 检查是否被人机验证阻止
		// captchaNodes := scdp.GetNodesWithAuth(mainCtx, "h1.Captcha-title")
		// if len(captchaNodes) > 0 {
		// 	log.Warnf("%s blocked by captcha", proxy.Host)
		// 	proxy.Blocked = true
		// 	return err
		// }

		err = chromedp.Run(mainCtx,
			chromedp.Location(&item.URL), // 更新item的URL
		)
		if err != nil {
			return err
		}

		if blocked(item.URL) {
			item.URL = url // 恢复URL
			// log.Debugf("%s blocked by captcha", proxy.Host)
			// proxy.Blocked = true
			proxy.Score = 0
			return _errBlocked
		} else if notFound(item.URL) { // 跳转去了列表页，直接返回错误
			return _errNotFound
		}

		var reviewStats string
		err = chromedp.Run(mainCtx,
			chromedp.WaitVisible(".ReviewStats", chromedp.ByQuery),
			chromedp.Text(".ReviewStats", &reviewStats, chromedp.ByQuery), // 获取评论数
		)
		if err != nil {
			return err
		}
		reviewStatsMatches := _reviewCountRegex.FindStringSubmatch(reviewStats)
		if len(reviewStatsMatches) < 2 {
			// 没有评论，不做后续处理
			return nil
		}
		count, err := strconv.Atoi(reviewStatsMatches[1])
		if err != nil {
			return err
		} else if count <= 0 {
			// 没有评论，不做后续处理
			return nil
		}

		// 选择评论按时间排序，等待评论加载完毕
		err = chromedp.Run(mainCtx,
			chromedp.WaitVisible(`.ReviewsSearchSortFilter-dropdown`, chromedp.ByQuery),
			chromedp.Click(`.ReviewsSearchSortFilter-dropdown`, chromedp.ByQuery),
			chromedp.WaitVisible(`#productReviewsSubheaderDropdown-item-2`, chromedp.ByQuery),
			chromedp.Click(`#productReviewsSubheaderDropdown-item-2`, chromedp.ByQuery),
			chromedp.WaitVisible(`div.ProductReviewList-links`, chromedp.ByQuery),
		)
		if err != nil {
			log.Debug(err)
			return nil // 排序下拉列表相关错误不需要返回
		}

		for loadReviews(mainCtx, from) { // 循环点击加载更多按钮
		}

		log.Debugf("[%s] %s got %d reviews", proxy.URI, item.URL, len(dic))
		return nil
	})

	x.CDPClient.ProxyStore.Return(proxy)

	// if proxy.Blocked || proxy.Expired {
	if proxy.Score <= 0 {
		// 重试
		return x.FetchReviews(item, from)
	} else if err == nil {
		r = make([]*model.ReviewDTO, 0, len(dic))
		for _, v := range dic {
			r = append(r, v)
		}
	} else {
		r = make([]*model.ReviewDTO, 0)
	}

	return
}

func captureReviewsFromAjax(ev interface{}, mainCtx context.Context, item *model.ItemDTO, from time.Time, dic *map[int]*model.ReviewDTO) {
	switch ev := ev.(type) {

	case *network.EventResponseReceived:
		if ev.Response.Status != 200 || !strings.Contains(ev.Response.URL, _ajaxURL) { // 排除不满足条件的Response
			return
		}
		c := chromedp.FromContext(mainCtx)
		data, err := network.GetResponseBody(ev.RequestID).Do(cdp.WithExecutor(mainCtx, c.Target)) // 获取Response Body数据
		if err != nil {
			var e1 *cdproto.Error
			if errors.As(err, &e1) {
				if e1.Code == -32000 {
					return
				}
			}
			log.Errorf("%s: %s", item.SKU, err.Error())
			return
		}
		var resp *model.ReviewResp
		err = json.Unmarshal(data, &resp) // 反序列化
		if u.LogError(err) {
			return
		}

		if resp != nil && resp.Data != nil && resp.Data.Product != nil && resp.Data.Product.CustomerReviews != nil && len(resp.Data.Product.CustomerReviews.Reviews) > 0 { // 再次排除非Review Ajax请求
			for _, review := range resp.Data.Product.CustomerReviews.Reviews {
				var err error
				review.CreatedOn, err = time.Parse(_dateFormat, review.Date)
				if err != nil {
					log.Errorf("%d parse date failed: %s", review.ReviewID, err.Error())
					continue
				}
				if review.CreatedOn.After(from) {
					review.SKU = item.SKU
					review.ItemNOs = item.ItemNOs
					(*dic)[review.ReviewID] = review // 添加进结果
				}
			}
		}
	}
}

func loadReviews(ctx context.Context, from time.Time) bool {
	reviewNodes := scdp.GetNodesWithAuth(ctx, "article.ProductReview")

	count := len(reviewNodes)

	if count > 0 {
		lastReviewNode := reviewNodes[count-1]
		dateText := scdp.GetTextWithAuth(ctx, `.ProductReview-reviewDetails p[data-hb-id="pl-text"]`, chromedp.FromNode(lastReviewNode))
		date, _ := time.Parse(_dateFormat, dateText)
		if date.After(from) {
			return loadMore(ctx)
		}
	}

	return false
}

func loadMore(ctx context.Context) bool {
	loadButtons := scdp.GetNodesWithAuth(ctx, "div.ProductReviewList-links button")
	if len(loadButtons) > 0 {
		// 有加载按钮
		buttonText := scdp.GetTextWithAuth(ctx, ".Button-content", chromedp.FromNode(loadButtons[0]))
		if strings.Contains(buttonText, "More") { // 是加载更多按钮
			scdp.ClickWithAuth(ctx, `div.ProductReviewList-links button`)
			return true
		} else {
			return false
		}
	} else {
		// 没有加载按钮
		return false
	}
}

func blocked(url string) bool {
	return strings.Contains(url, "/captcha/")
}

func notFound(url string) bool {
	// return strings.Contains(url, "/sb0/") || strings.Contains(url, "/cat/") || url == "https://www.wayfair.com/" || url == "https://www.wayfair.com/?"
	return !strings.Contains(url, "/pdp/")
}

func generateURL(item *model.ItemDTO) string {
	return fmt.Sprintf(_urlFormat, item.SKU)
}
