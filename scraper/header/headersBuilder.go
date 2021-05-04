package header

import "github.com/syncfuture/go/srand"

type headersBuilder struct {
	Headers map[string]string
}

func NewHeadersBuilder() *headersBuilder {
	r := &headersBuilder{
		Headers: make(map[string]string, 10),
	}
	r.build()
	return r
}

var _accepts = []string{
	"text/html,application/xhtml+xml,application/xml;q=0.9,image/webp,image/apng,*/*;q=0.8,application/signed-exchange;v=b3;q=0.9",
	// "text/html,application/xhtml+xml,application/xml;q=0.9,image/webp,image/apng,*/*;q=0.8,application/signed-exchange;v=b3;q=0.9",
	// "text/html,application/xhtml+xml,application/xml;q=0.9,image/webp,image/apng,*/*;q=0.8,application/signed-exchange;v=b3;q=0.9",
	// "text/html,application/xhtml+xml,application/xml;q=0.9,image/webp,image/apng,*/*;q=0.8,application/signed-exchange;v=b3;q=0.9",
	// "text/html,application/xhtml+xml,application/xml;q=0.9,image/webp,image/apng,*/*;q=0.8,application/signed-exchange;v=b3;q=0.9",
	// "text/html,application/xhtml+xml,application/xml;q=0.9,image/webp,image/apng,*/*;q=0.8,application/signed-exchange;v=b3;q=0.9",
	// "text/html,application/xhtml+xml,application/xml;q=0.9,image/webp,image/apng,*/*;q=0.8,application/signed-exchange;v=b3;q=0.9",
	// "text/html,application/xml",
	// "text/html",
	// "text/html,*/*",
	// "text/html,application/xhtml+xml,application/xml",
	// "text/html,application/xhtml+xml, */*",
}

func (x *headersBuilder) buildAccept() {
	x.Headers["Accept"] = _accepts[srand.IntRange(0, len(_accepts)-1)]
}

var _acceptEncodings = []string{
	"gzip, deflate, br",
}

func (x *headersBuilder) buildAcceptEncoding() {
	a := srand.IntRange(0, 1)
	if a == 1 {
		x.Headers["Accept-Encoding"] = _acceptEncodings[srand.IntRange(0, len(_acceptEncodings)-1)]
	}
}

var _acceptLanguages = []string{
	// "en-US,en;q=0.9",
	// "en-US,en;q=0.8",
	// "en-US,en",
	// "en-UK,en;",
	// "zh-CN,fr-FR;q=0.5",
	"en-US,en;q=0.8,zh-Hans-CN;q=0.5,zh-Hans;q=0.3",
	"en-US,en;q=0.8,zh-Hans-CN;q=0.5",
	"en-US,en;q=0.8,zh-Hans-CN;q=0.5",
	"en-US,en;q=0.9,zh-Hans;q=0.4",
	"en-US,en;q=0.9,zh-CN;q=0.8,zh;q=0.7,zh-TW;q=0.6,ja;q=0.5",
}

func (x *headersBuilder) buildAcceptLanguage() {
	x.Headers["Accept-Language"] = _acceptLanguages[srand.IntRange(0, len(_acceptLanguages)-1)]
}

var _userAgents = []string{
	"Mozilla/5.0 (Windows NT 6.3; WOW64; Trident/7.0; rv:11.0) like Gecko",
	"Mozilla/5.0 (Windows NT 5.1) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/28.0.1500.95 Safari/537.36",
	"Mozilla/5.0 (Windows NT 6.1; WOW64; Trident/7.0; SLCC2; .NET CLR 2.0.50727; .NET CLR 3.5.30729; .NET CLR 3.0.30729; Media Center PC 6.0; .NET4.0C; rv:11.0) like Gecko)",
	"Mozilla/5.0 (Windows; U; Windows NT 5.2) Gecko/2008070208 Firefox/3.0.1",
	"Mozilla/5.0 (Windows; U; Windows NT 5.1) Gecko/20070309 Firefox/2.0.0.3",
	"Mozilla/5.0 (Windows; U; Windows NT 5.1) Gecko/20070803 Firefox/1.5.0.12",
	"Opera/9.27 (Windows NT 5.2; U; zh-cn)",
	"Mozilla/5.0 (Macintosh; PPC Mac OS X; U; en) Opera 8.0",
	"Opera/8.0 (Macintosh; PPC Mac OS X; U; en)",
	"Mozilla/5.0 (Windows; U; Windows NT 5.1; en-US; rv:1.8.1.12) Gecko/20080219 Firefox/2.0.0.12 Navigator/9.0.0.6",
	"Mozilla/4.0 (compatible; MSIE 8.0; Windows NT 6.1; Win64; x64; Trident/4.0)",
	"Mozilla/4.0 (compatible; MSIE 8.0; Windows NT 6.1; Trident/4.0)",
	"Mozilla/5.0 (compatible; MSIE 10.0; Windows NT 6.1; WOW64; Trident/6.0; SLCC2; .NET CLR 2.0.50727; .NET CLR 3.5.30729; .NET CLR 3.0.30729; Media Center PC 6.0; InfoPath.2; .NET4.0C; .NET4.0E)",
	"Mozilla/5.0 (Windows NT 6.1; WOW64) AppleWebKit/537.1 (KHTML, like Gecko) Maxthon/4.0.6.2000 Chrome/26.0.1410.43 Safari/537.1 ",
	"Mozilla/5.0 (compatible; MSIE 10.0; Windows NT 6.1; WOW64; Trident/6.0; SLCC2; .NET CLR 2.0.50727; .NET CLR 3.5.30729; .NET CLR 3.0.30729; Media Center PC 6.0; InfoPath.2; .NET4.0C; .NET4.0E; QQBrowser/7.3.9825.400)",
	"Mozilla/5.0 (Windows NT 6.1; WOW64; rv:21.0) Gecko/20100101 Firefox/21.0 ",
	"Mozilla/5.0 (Windows NT 6.1; WOW64) AppleWebKit/537.1 (KHTML, like Gecko) Chrome/21.0.1180.92 Safari/537.1 LBBROWSER",
	"Mozilla/5.0 (compatible; MSIE 10.0; Windows NT 6.1; WOW64; Trident/6.0; BIDUBrowser 2.x)",
	"Mozilla/5.0 (Windows NT 6.1; WOW64) AppleWebKit/536.11 (KHTML, like Gecko) Chrome/20.0.1132.11 TaoBrowser/3.0 Safari/536.11",
	"Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/87.0.4280.141 Safari/537.36",
	"Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/87.0.4280.141 Safari/537.36 Edg/87.0.664.75",
	"Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/87.0.4280.67 Safari/537.36 Edg/87.0.664.55",
	"Mozilla/5.0 (Windows NT 10.0; WOW64; Trident/7.0; rv:11.0) like Gecko",
}

func (x *headersBuilder) buildUserAgent() {
	x.Headers["User-Agent"] = _userAgents[srand.IntRange(0, len(_userAgents)-1)]
}

func (x *headersBuilder) buildDNT() {
	x.Headers["DNT"] = "1"
	// a := srand.IntRange(0, 1)
	// if a == 1 {
	// 	x.Headers["dnt"] = "1"
	// }
}

func (x *headersBuilder) buildUpgradeInsecureRequests() {
	x.Headers["Upgrade-Insecure-Requests"] = "1"
	// 	a := srand.IntRange(0, 1)
	// if a == 1 {
	// 	x.Headers["Upgrade-Insecure-Requests"] = "1"
	// }
}

var _cacheControls = []string{
	"no-cache",
	"max-age=0",
}

func (x *headersBuilder) buildCacheControl() {
	x.Headers["Cache-Control"] = _cacheControls[srand.IntRange(0, len(_cacheControls)-1)]
	// a := srand.IntRange(0, 1)
	// if a == 1 {
	// 	x.Headers["Cache-Control"] = _cacheControls[srand.IntRange(0, len(_cacheControls)-1)]
	// }
}

var _connections = []string{
	"keep-alive",
	"close",
}

func (x *headersBuilder) buildConnection() {
	x.Headers["Connection"] = _connections[srand.IntRange(0, len(_connections)-1)]
	// 	a := srand.IntRange(0, 1)
	// if a == 1 {
	// 	x.Headers["Connection"] = _connections[srand.IntRange(0, len(_connections)-1)]
	// }
}

func (x *headersBuilder) buildSec() {
	a := srand.IntRange(0, 1)
	if a == 1 {
		x.Headers["Sec-Fetch-Site"] = "none"
		x.Headers["Sec-Fetch-Mode"] = "navigate"
		x.Headers["Sec-Fetch-Dest"] = "document"
	}
}

func (x *headersBuilder) build() {
	x.buildAccept()
	x.buildAcceptEncoding()
	x.buildAcceptLanguage()
	x.buildUserAgent()
	x.buildDNT()
	x.buildUpgradeInsecureRequests()
	x.buildCacheControl()
	x.buildConnection()
	x.buildSec()
}
