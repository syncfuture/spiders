package wayfair

import (
	"context"
	"encoding/json"
	"fmt"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/chromedp/cdproto/cdp"
	"github.com/chromedp/cdproto/network"
	"github.com/chromedp/chromedp"

	log "github.com/syncfuture/go/slog"
	"github.com/syncfuture/go/u"
	"github.com/syncfuture/scraper/scdp"
	"github.com/syncfuture/spiders/wayfair/model"
)

const (
	// _expError   = "ERR_INVALID_AUTH_CREDENTIALS"
	// _retryError = "ERR_TOO_MANY_RETRIES"
	_urlFormat = "https://www.wayfair.com/category/pdp/product-%s.html"
	_ajaxURL   = "https://www.wayfair.com/graphql?hash="
)

var (
	_reviewCountRegex = regexp.MustCompile(`(\d+) Reviews`)
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
	proxy := x.CDPClient.ProxyStore.Lease()

	err = x.CDPClient.FetchWithProxy(proxy.ToURL(true), func(mainCtx context.Context) error {
		// 监听ajax事件，获取评论
		chromedp.ListenTarget(mainCtx, func(ev interface{}) {
			go captureReviewsFromAjax(ev, mainCtx, item, from, &dic)
		})
		// chromedp.ListenTarget(mainCtx, func(ev interface{}) {
		// 	go func() {
		// 		switch ev := ev.(type) {

		// 		case *network.EventResponseReceived:
		// 			if ev.Response.Status != 200 || !strings.Contains(ev.Response.URL, _ajaxURL) {
		// 				return
		// 			}
		// 			c := chromedp.FromContext(mainCtx)
		// 			data, err := network.GetResponseBody(ev.RequestID).Do(cdp.WithExecutor(mainCtx, c.Target))
		// 			if u.LogError(err) {
		// 				return
		// 			}
		// 			var resp *model.ReviewResp
		// 			err = json.Unmarshal(data, &resp)
		// 			if u.LogError(err) {
		// 				return
		// 			}

		// 			if resp != nil && resp.Data != nil && resp.Data.Product != nil && resp.Data.Product.CustomerReviews != nil && len(resp.Data.Product.CustomerReviews.Reviews) > 0 {
		// 				// reviews[]
		// 				for _, review := range resp.Data.Product.CustomerReviews.Reviews {
		// 					date, err := time.Parse("01/02/2006", review.Date)
		// 					if err != nil {
		// 						log.Errorf("%d parse date failed: %s", review.ReviewID, err.Error())
		// 						continue
		// 					}
		// 					if date.After(from) {
		// 						review.SKU = item.SKU
		// 						review.Items = item.Items
		// 						dic[review.ReviewID] = review
		// 					}
		// 				}
		// 			}
		// 		}
		// 	}()
		// 	// other needed network Event
		// })

		// client.FetchWithHead(func(ctx context.Context) {
		// 跳转去产品页
		var url string
		if item.URL == "" {
			url = fmt.Sprintf(_urlFormat, item.SKU)
		} else {
			url = item.URL
		}
		err := scdp.NavigateWithAuth(mainCtx, url, 1920, 936)
		if err != nil {
			// if strings.Contains(err.Error(), _expError) || strings.Contains(err.Error(), _retryError) {
			// 	proxy.Expired = true
			// }
			proxy.Expired = true // 导航失败，直接视作代理池过期
			return err
		}

		// 检查是否被人机验证阻止
		captchaNodes := scdp.GetNodesWithAuth(mainCtx, "h1.Captcha-title")
		if len(captchaNodes) > 0 {
			log.Warnf("%s blocked by captcha", proxy.Host)
			proxy.Blocked = true
			return err
		}

		var reviewStats string
		err = chromedp.Run(mainCtx,
			chromedp.Location(&item.URL), // 更新item的URL
			chromedp.WaitVisible(".ReviewStats", chromedp.ByQuery),
			chromedp.Text(".ReviewStats", &reviewStats, chromedp.ByQuery), // 获取评论数
		)
		if u.LogError(err) {
			return err
		}
		reviewStatsMatches := _reviewCountRegex.FindStringSubmatch(reviewStats)
		if len(reviewStatsMatches) < 2 {
			// 没有评论，不做后续处理
			return nil
		}
		count, err := strconv.Atoi(reviewStatsMatches[1])
		if u.LogError(err) {
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

		log.Debugf("%s got %d reviews", url, len(dic))

		return err
	})

	x.CDPClient.ProxyStore.Return(proxy)

	if proxy.Blocked || proxy.Expired {
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
		if u.LogError(err) {
			return
		}
		var resp *model.ReviewResp
		err = json.Unmarshal(data, &resp) // 反序列化
		if u.LogError(err) {
			return
		}

		if resp != nil && resp.Data != nil && resp.Data.Product != nil && resp.Data.Product.CustomerReviews != nil && len(resp.Data.Product.CustomerReviews.Reviews) > 0 { // 再次排除非Review Ajax请求
			for _, review := range resp.Data.Product.CustomerReviews.Reviews {
				date, err := time.Parse("01/02/2006", review.Date)
				if err != nil {
					log.Errorf("%d parse date failed: %s", review.ReviewID, err.Error())
					continue
				}
				review.CreatedOnUTC = date.UTC()
				if review.CreatedOnUTC.After(from.UTC()) {
					review.SKU = item.SKU
					review.Items = item.Items
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
		date, _ := time.Parse("01/02/2006", dateText)
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
