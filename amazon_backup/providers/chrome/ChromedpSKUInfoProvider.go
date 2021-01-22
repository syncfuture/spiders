package chrome

import (
	"context"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
	"github.com/chromedp/cdproto/cdp"
	"github.com/chromedp/chromedp"
	log "github.com/syncfuture/go/slog"
	"github.com/syncfuture/go/u"
	"github.com/syncfuture/spiders/amazon/core"
	"github.com/syncfuture/spiders/amazon/protoc/product"
)

const (
	_baseURL = "https://www.amazon.com/product/dp/"
)

var (
	cleanerRegex = regexp.MustCompile(`\([^\)]+\)|[\r\n]`)
	//增加选项，允许chrome窗口显示出来
	options = []chromedp.ExecAllocatorOption{
		// chromedp.Flag("headless", false),
		// chromedp.Flag("hide-scrollbars", false),
		// chromedp.Flag("mute-audio", false),
		chromedp.UserAgent(`Mozilla/5.0 (Windows NT 6.3; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/73.0.3683.103 Safari/537.36`),
	}
)

func init() {
	options = append(chromedp.DefaultExecAllocatorOptions[:], options...)
}

type ChromedpSKUInfoProvider struct{}

func (x *ChromedpSKUInfoProvider) Get(asin string) (r *product.ProductSKU) {
	asin = strings.TrimSpace(strings.ToUpper(asin))
	r = &product.ProductSKU{
		ID:           core.GenerateID(),
		ASIN:         asin,
		Price:        -1,
		CreatedOnUTC: time.Now().UTC().Format(time.RFC3339),
	}

	//创建chrome窗口
	baseCtx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()
	allocCtx, cancel := chromedp.NewExecAllocator(baseCtx, options...)
	defer cancel()
	ctx, cancel := chromedp.NewContext(allocCtx, chromedp.WithLogf(log.Infof))
	defer cancel()

	var priceStr, sellerStr, stockStr, rankStr string
	var bulletStr string
	var outOfStockBoxExists, priceBoxExists, sellerBoxExists, stockBoxExists, bulletBoxExists bool

	// 开始抓取并判断元素是否存在
	if err := chromedp.Run(ctx,
		// 跳转去目标页面
		chromedp.Navigate(_baseURL+asin),
		// 等待buybox可见
		chromedp.WaitVisible("#buybox", chromedp.ByID),
		// 查询缺货节点是否存在
		chromedp.Query("#outOfStockBuyBox_feature_div", chromedp.AtLeast(0), chromedp.After(func(i context.Context, n ...*cdp.Node) error {
			outOfStockBoxExists = len(n) > 0
			return nil
		})),
		// 查询价格节点是否存在
		chromedp.Query("#priceblock_ourprice", chromedp.AtLeast(0), chromedp.After(func(i context.Context, n ...*cdp.Node) error {
			priceBoxExists = len(n) > 0
			return nil
		})),
		// 查询卖家节点是否存在
		chromedp.Query("#shipsFromSoldByInsideBuyBox_feature_div", chromedp.AtLeast(0), chromedp.After(func(i context.Context, n ...*cdp.Node) error {
			sellerBoxExists = len(n) > 0
			return nil
		})),
		// 查询库存节点是否存在
		chromedp.Query("#availabilityInsideBuyBox_feature_div", chromedp.AtLeast(0), chromedp.After(func(i context.Context, n ...*cdp.Node) error {
			stockBoxExists = len(n) > 0
			return nil
		})),
		// 查询信息表是否存在
		chromedp.Query("#productDetails_detailBullets_sections1", chromedp.AtLeast(0), chromedp.After(func(i context.Context, n ...*cdp.Node) error {
			bulletBoxExists = len(n) > 0
			return nil
		})),
	); u.LogError(err) {
		// buybox未找到，返回Active为false
		return r
	}

	if outOfStockBoxExists {
		// 没库存
		r.Active = true
		return r
	}

	tasks := make(chromedp.Tasks, 0, 4)
	// 获取价格任务
	if priceBoxExists {
		tasks = append(tasks, chromedp.Text("#priceblock_ourprice", &priceStr, chromedp.ByID))
	}
	// 获取卖家任务
	if sellerBoxExists {
		tasks = append(tasks, chromedp.Text("#shipsFromSoldByInsideBuyBox_feature_div", &sellerStr, chromedp.ByID))
	}
	// 获取库存状态任务
	if stockBoxExists {
		tasks = append(tasks, chromedp.Text("#availabilityInsideBuyBox_feature_div", &stockStr, chromedp.ByID))
	}
	// 获取评级任务
	if bulletBoxExists {
		tasks = append(tasks, chromedp.OuterHTML("#productDetails_detailBullets_sections1", &bulletStr, chromedp.ByID))
	}
	// 开始查询页面
	if err := chromedp.Run(ctx, tasks); err != nil {
		if u.LogError(err) {
			return r
		}
	}

	if bulletBoxExists {
		bulletTab, err := goquery.NewDocumentFromReader(strings.NewReader(bulletStr))
		if err == nil {
			bulletTab.Find("tr").Each(func(i int, tr *goquery.Selection) {
				th := strings.TrimSpace(tr.Find("th").Text())
				if th == "Best Sellers Rank" {
					rankStr = tr.Find("td").Text()
					// rankStr = newlineRegex.ReplaceAllString(rankStr, "")
					rankStr = format(&rankStr)
				}
			})
		} else {
			u.LogError(err)
		}
	}

	priceStr = strings.TrimLeft(format(&priceStr), "$")
	sellerStr = format(&sellerStr)
	stockStr = format(&stockStr)

	instock := !strings.Contains(stockStr, "Out") && !strings.Contains(stockStr, "In stock on")
	price, err := strconv.ParseFloat(priceStr, 32)
	u.LogError(err)

	r.Active = true
	r.Seller = sellerStr
	r.Price = float32(price)
	r.Instock = instock
	return r
}

func format(str *string) string {
	if str == nil {
		return ""
	}
	return strings.TrimSpace(strings.Trim(cleanerRegex.ReplaceAllString(*str, ""), "."))
}
