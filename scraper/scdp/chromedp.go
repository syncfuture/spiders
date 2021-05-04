package scdp

import (
	"context"
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
	"net/url"
	"runtime"
	"time"

	"github.com/chromedp/cdproto/cdp"
	"github.com/chromedp/cdproto/fetch"
	"github.com/chromedp/chromedp"

	task "github.com/syncfuture/go/stask"
	"github.com/syncfuture/spiders/scraper/store"

	"os/exec"

	config "github.com/syncfuture/go/sconfig"
	log "github.com/syncfuture/go/slog"
	"github.com/syncfuture/go/u"
)

type ChromeDPClient struct {
	ConfigProvider       config.IConfigProvider
	ChromeCMD            string
	ExcelFile            string
	WebSocketDebuggerURL string
	WaitChromeDelay      int
	BatchSize            int
	BatchInterval        int
	Timeout              time.Duration
	SliceScheduler       task.ISliceScheduler
	ProxyStore           store.IProxyStore
}

func NewChromeDPClient(cp config.IConfigProvider, proxyStore store.IProxyStore, batchEvents ...func(int, int)) (r *ChromeDPClient) {
	r = new(ChromeDPClient)
	r.ConfigProvider = cp
	if runtime.GOOS == "windows" {
		r.ChromeCMD = r.ConfigProvider.GetString("WinChromeCMD")
	} else if runtime.GOOS == "darwin" {
		r.ChromeCMD = r.ConfigProvider.GetString("MacChromeCMD")
	}
	r.ProxyStore = proxyStore
	r.WaitChromeDelay = r.ConfigProvider.GetIntDefault("WaitChromeDelay", 2000)
	r.BatchSize = r.ConfigProvider.GetIntDefault("BatchSize", 4)
	r.BatchInterval = r.ConfigProvider.GetIntDefault("BatchInterval", 500)
	timeout := r.ConfigProvider.GetIntDefault("Timeout", 10000)
	r.Timeout = time.Duration(timeout) * time.Millisecond
	r.SliceScheduler = task.NewBatchScheduler(r.BatchSize, r.BatchInterval, batchEvents...)
	return r
}

func (x *ChromeDPClient) startHead() {
	portOpen := isPortOpen(9222)

	if !portOpen {
		log.Info("Starting chrome...")

		cmd := exec.Command(x.ChromeCMD, "--remote-debugging-port=9222")
		err := cmd.Start()
		time.Sleep(time.Duration(x.WaitChromeDelay) * time.Millisecond)
		if u.LogError(err) {
			return
		}
	}

	resp, err := http.Get("http://localhost:9222/json/version")
	if u.LogError(err) {
		return
	}
	defer resp.Body.Close()
	configJson, err := ioutil.ReadAll(resp.Body)
	if u.LogError(err) {
		return
	}

	config := make(map[string]string)
	json.Unmarshal(configJson, &config)

	x.WebSocketDebuggerURL = config["webSocketDebuggerUrl"]
	if x.WebSocketDebuggerURL == "" {
		log.Fatal("get webSocketDebuggerUrl failed")
	}
	log.Debug("Connecting to ", x.WebSocketDebuggerURL)
}

func (x *ChromeDPClient) FetchWithHead(action func(ctx context.Context) error) error {
	x.startHead()

	ctx := context.Background()
	timeoutCtx, cancel1 := context.WithTimeout(ctx, x.Timeout)
	defer cancel1()

	allocCtx, cancel2 := chromedp.NewRemoteAllocator(timeoutCtx, x.WebSocketDebuggerURL)
	defer cancel2()

	taskCtx, cancel3 := chromedp.NewContext(allocCtx)
	defer cancel3()

	return action(taskCtx)
}

func (x *ChromeDPClient) FetchWithProxy(proxyURL *url.URL, action func(ctx context.Context) error) error {
	if x.ProxyStore == nil {
		return errors.New("ProxyStore is nil")
	}

	ctx := context.Background()
	timeoutCtx, cancel1 := context.WithTimeout(ctx, x.Timeout)
	defer cancel1()

	// options := append(
	// 	chromedp.DefaultExecAllocatorOptions[:],
	// 	chromedp.ProxyServer(proxy),
	// 	chromedp.Flag("blink-settings", "imagesEnabled=false"),
	// )

	// allocCtx, cancel2 := chromedp.NewExecAllocator(timeoutCtx, options...)
	allocCtx, cancel2 := chromedp.NewExecAllocator(
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
		chromedp.ProxyServer(proxyURL.Host),
		chromedp.Flag("blink-settings", "imagesEnabled=false"),
	)
	defer cancel2()

	mainCtx, cancel3 := chromedp.NewContext(allocCtx)
	defer cancel3()

	username := proxyURL.User.Username()
	password, _ := proxyURL.User.Password()

	if username != "" && password != "" {
		chromedp.ListenTarget(mainCtx, func(ev interface{}) {
			go func() {
				switch ev := ev.(type) {
				case *fetch.EventAuthRequired:
					c := chromedp.FromContext(mainCtx)
					execCtx := cdp.WithExecutor(mainCtx, c.Target)

					resp := &fetch.AuthChallengeResponse{
						Response: fetch.AuthChallengeResponseResponseProvideCredentials,
						Username: username,
						Password: password,
					}

					err := fetch.ContinueWithAuth(ev.RequestID, resp).Do(execCtx)
					u.LogError(err)

				case *fetch.EventRequestPaused:
					c := chromedp.FromContext(mainCtx)
					execCtx := cdp.WithExecutor(mainCtx, c.Target)
					err := fetch.ContinueRequest(ev.RequestID).Do(execCtx)
					if err != nil {
						log.Debug(err)
					}
				}
			}()
		})
	}

	return action(mainCtx)
}
