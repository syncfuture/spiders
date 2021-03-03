package wayfair

import (
	"context"
	"testing"
	"time"

	"github.com/chromedp/cdproto/cdp"
	"github.com/chromedp/cdproto/fetch"
	"github.com/chromedp/chromedp"
	"github.com/olivere/elastic/v7"
	"github.com/stretchr/testify/assert"
	"github.com/syncfuture/go/sconfig"
	log "github.com/syncfuture/go/slog"
	"github.com/syncfuture/go/u"
	"github.com/syncfuture/scraper/scdp"
	"github.com/syncfuture/scraper/store/webshare"
	"github.com/syncfuture/spiders/wayfair/dal/es"
	"github.com/syncfuture/spiders/wayfair/model"
)

func TestGetReviews(t *testing.T) {
	cp := sconfig.NewJsonConfigProvider()
	log.Init(cp)

	store := webshare.NewDefaultWebShareProxyStore()
	t.Log(store)

	itemDAL, err := es.NewESItemDAL(
		elastic.SetURL("http://sa:Famous901@localhost:9200"),
		elastic.SetSniff(false),
	)
	if u.LogError(err) {
		return
	}

	reviewDAL, err := es.NewESReviewDAL(
		elastic.SetURL("http://sa:Famous901@localhost:9200"),
		elastic.SetSniff(false),
	)
	if u.LogError(err) {
		return
	}

	query := &model.ItemQuery{}
	rs, err := itemDAL.GetAllItems(query)
	if u.LogError(err) {
		return
	}

	assert.NotEmpty(t, rs)

	client := scdp.NewChromeDPClient(cp, store)
	scraper := NewReviewsScraper(client)
	from := time.Now().Add(time.Hour * 24 * -30)

	// scraper.FetchReviews(&model.ItemDTO{SKU: "fssx4862"}, from) // fssx4862 截至没有评论3/3/2021

	for _, item := range rs.Items {
		reviews, err := scraper.FetchReviews(item, from)
		if !u.LogError(err) {
			item.Status = 1
			err = reviewDAL.SaveReviews(reviews...)
			u.LogError(err)
		} else {
			item.Status = -1
		}
		itemDAL.SaveItems(item)
	}
}

func Test123(t *testing.T) {
	cp := sconfig.NewJsonConfigProvider()
	log.Init(cp)

	err := run(aaa)
	u.LogError(err)
}

func aaa(ctx context.Context) error {
	var txt string
	err := chromedp.Run(ctx,
		fetch.Enable().WithHandleAuthRequests(true),
		chromedp.Navigate("http://ipv4.webshare.io/"),
		chromedp.Text("html", &txt, chromedp.ByQuery),
	)
	return err
}

func run(action func(ctx context.Context) error) error {
	store := webshare.NewDefaultWebShareProxyStore()
	proxy := store.Lease()
	defer store.Return(proxy)

	timeoutCtx, cancel1 := context.WithTimeout(context.Background(), time.Second*60)
	defer cancel1()

	execCtx, cancel2 := chromedp.NewExecAllocator(
		timeoutCtx,
		chromedp.NoFirstRun,
		chromedp.NoDefaultBrowserCheck,
		// chromedp.Headless,
		chromedp.Flag("disable-background-networking", true),
		chromedp.Flag("enable-features", "NetworkService,NetworkServiceInProcess"),
		chromedp.Flag("disable-background-timer-throttling", true),
		chromedp.Flag("disable-backgrounding-occluded-windows", true),
		chromedp.Flag("disable-breakpad", true),
		chromedp.Flag("disable-client-side-phishing-detection", true),
		chromedp.Flag("disable-default-apps", true),
		chromedp.Flag("disable-dev-shm-usage", true),
		chromedp.Flag("disable-extensions", true),
		chromedp.Flag("disable-features", "site-per-process,TranslateUI,BlinkGenPropertyTrees"),
		chromedp.Flag("disable-hang-monitor", true),
		chromedp.Flag("disable-ipc-flooding-protection", true),
		chromedp.Flag("disable-popup-blocking", true),
		chromedp.Flag("disable-prompt-on-repost", true),
		chromedp.Flag("disable-renderer-backgrounding", true),
		chromedp.Flag("disable-sync", true),
		chromedp.Flag("force-color-profile", "srgb"),
		chromedp.Flag("metrics-recording-only", true),
		chromedp.Flag("safebrowsing-disable-auto-update", true),
		chromedp.Flag("enable-automation", true),
		chromedp.Flag("password-store", "basic"),
		chromedp.Flag("use-mock-keychain", true),

		chromedp.Flag("incognito", true),
		chromedp.ProxyServer(proxy.Host),
		chromedp.Flag("blink-settings", "imagesEnabled=false"),
	)
	defer cancel2()

	mainCtx, cancel3 := chromedp.NewContext(execCtx)
	defer cancel3()

	chromedp.ListenTarget(mainCtx, func(ev interface{}) {
		go func() {
			switch ev := ev.(type) {
			case *fetch.EventAuthRequired:
				c := chromedp.FromContext(mainCtx)
				execCtx := cdp.WithExecutor(mainCtx, c.Target)

				resp := &fetch.AuthChallengeResponse{
					Response: fetch.AuthChallengeResponseResponseProvideCredentials,
					Username: proxy.Username,
					Password: proxy.Password,
				}

				err := fetch.ContinueWithAuth(ev.RequestID, resp).Do(execCtx)
				u.LogError(err)

			case *fetch.EventRequestPaused:
				c := chromedp.FromContext(mainCtx)
				execCtx := cdp.WithExecutor(mainCtx, c.Target)
				err := fetch.ContinueRequest(ev.RequestID).Do(execCtx)
				u.LogError(err)
			}
		}()
	})

	return action(mainCtx)
}
