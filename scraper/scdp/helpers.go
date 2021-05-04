package scdp

import (
	"context"
	"net"
	"strconv"
	"time"

	"github.com/chromedp/cdproto/cdp"
	"github.com/chromedp/cdproto/emulation"
	"github.com/chromedp/cdproto/fetch"
	"github.com/chromedp/chromedp"
	log "github.com/syncfuture/go/slog"
)

func isPortOpen(port int) bool {
	_, err := net.DialTimeout("tcp", "localhost:"+strconv.Itoa(port), time.Millisecond*500)
	if err != nil {
		return false
	}
	return true
}

func GetText(ctx context.Context, selector string, opts ...chromedp.QueryOption) (r string) {
	opts = append(opts, chromedp.ByQuery, chromedp.AtLeast(0))
	err := chromedp.Run(ctx,
		chromedp.Text(selector, &r, opts...),
	)
	if err != nil {
		log.Debug(err)
	}
	return
}

func GetTextWithAuth(ctx context.Context, selector string, opts ...chromedp.QueryOption) (r string) {
	opts = append(opts, chromedp.ByQuery, chromedp.AtLeast(0))
	err := chromedp.Run(ctx,
		fetch.Enable().WithHandleAuthRequests(true),
		chromedp.Text(selector, &r, opts...),
	)
	if err != nil {
		log.Debug(err)
	}
	return
}

func GetNodes(ctx context.Context, selector string, opts ...chromedp.QueryOption) (r []*cdp.Node) {
	opts = append(opts, chromedp.ByQueryAll, chromedp.AtLeast(0))
	err := chromedp.Run(ctx,
		chromedp.Nodes(selector, &r, opts...),
	)
	if err != nil {
		log.Debug(err)
	}
	if r == nil {
		r = make([]*cdp.Node, 0)
	}
	return
}

func GetNodesWithAuth(ctx context.Context, selector string, opts ...chromedp.QueryOption) (r []*cdp.Node) {
	opts = append(opts, chromedp.ByQueryAll, chromedp.AtLeast(0))
	err := chromedp.Run(ctx,
		fetch.Enable().WithHandleAuthRequests(true),
		chromedp.Nodes(selector, &r, opts...),
	)
	if err != nil {
		log.Debug(err)
	}
	if r == nil {
		r = make([]*cdp.Node, 0)
	}
	return
}

func GetAttributes(ctx context.Context, selector, attrName string, opts ...chromedp.QueryOption) (r []string) {
	nodes := GetNodes(ctx, selector, opts...)
	r = make([]string, 0, len(nodes))
	for _, node := range nodes {
		r = append(r, node.AttributeValue(attrName))
	}
	return
}

func GetAttributesWithAuth(ctx context.Context, selector, attrName string, opts ...chromedp.QueryOption) (r []string) {
	nodes := GetNodesWithAuth(ctx, selector, opts...)
	r = make([]string, 0, len(nodes))
	for _, node := range nodes {
		r = append(r, node.AttributeValue(attrName))
	}
	return
}

func Click(ctx context.Context, selector string, opts ...chromedp.QueryOption) bool {
	opts = append(opts, chromedp.ByQuery, chromedp.AtLeast(0))
	err := chromedp.Run(ctx,
		chromedp.Click(selector, opts...),
	)
	if err != nil {
		log.Debug(err)
		return false
	}

	return true
}

func ClickWithAuth(ctx context.Context, selector string, opts ...chromedp.QueryOption) bool {
	opts = append(opts, chromedp.ByQuery, chromedp.AtLeast(0))
	err := chromedp.Run(ctx,
		fetch.Enable().WithHandleAuthRequests(true),
		chromedp.Click(selector, opts...),
	)
	if err != nil {
		log.Debug(err)
		return false
	}

	return true
}

func Navigate(ctx context.Context, url string, width, height int64) error {
	err := chromedp.Run(ctx,
		emulation.SetDeviceMetricsOverride(width, height, 1.0, false), // 设置屏幕尺寸，防止自适应
		chromedp.Navigate(url),
	)

	return err
}

func NavigateWithAuth(ctx context.Context, url string, width, height int64) error {
	err := chromedp.Run(ctx,
		fetch.Enable().WithHandleAuthRequests(true),
		emulation.SetDeviceMetricsOverride(width, height, 1.0, false), // 设置屏幕尺寸，防止自适应
		chromedp.Navigate(url),
	)

	return err
}
