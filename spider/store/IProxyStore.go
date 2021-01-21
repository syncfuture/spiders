package store

import "github.com/syncfuture/spiders/spider/model"

type IProxyStore interface {
	SaveProxy(proxy *model.Proxy) error
	GetProxy(id string) (*model.Proxy, error)
	GetRandomProxy() (*model.Proxy, error)
	GetProxies() ([]*model.Proxy, error)
}
