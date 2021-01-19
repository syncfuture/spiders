package model

import (
	"net/http"
	"net/url"
)

func (x *Proxy) ToProxyURL() func(*http.Request) (*url.URL, error) {
	if x.Username != "" && x.Password != "" {
		return http.ProxyURL(&url.URL{
			Scheme: x.Scheme,
			Host:   x.Host,
			User:   url.UserPassword(x.Username, x.Password),
		})
	}

	return http.ProxyURL(&url.URL{
		Scheme: x.Scheme,
		Host:   x.Host,
	})
}
