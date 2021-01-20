package cdp

import (
	"context"
	"testing"

	"github.com/chromedp/cdproto/dom"
	"github.com/chromedp/chromedp"
	"github.com/syncfuture/go/sconfig"
	log "github.com/syncfuture/go/slog"
)

func TestAAAA(t *testing.T) {
	cp := sconfig.NewJsonConfigProvider()

	cdp := NewChromeDPWithHead(cp)

	cdp.Init()
	cdp.Fetch(func(c context.Context) {
		// ctx := context.Background()
		// timeoutCtx, cancel := context.WithTimeout(ctx, cdp.Timeout)
		// defer cancel()

		// allocCtx, cancel := chromedp.NewRemoteAllocator(timeoutCtx, cdp.WebSocketDebuggerURL)
		// defer cancel()

		// taskCtx, cancel := chromedp.NewContext(allocCtx)
		// defer cancel()

		chromedp.Run(c,
			chromedp.Navigate("https://www.amazon.com/gp/offer-listing/B08164VTWH//ref=olp_f_new?f_primeEligible=true&f_new=true"),
			chromedp.WaitVisible("#olpOfferListColumn", chromedp.ByID),
			chromedp.ActionFunc(func(ctx context.Context) error {
				node, err := dom.GetDocument().Do(ctx)
				if err != nil {
					return err
				}
				str, err := dom.GetOuterHTML().WithNodeID(node.NodeID).Do(ctx)
				log.Info(str)
				return err
			}),
		)

		t.Log(c)
	})
}
