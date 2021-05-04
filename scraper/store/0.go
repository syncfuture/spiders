package store

import "github.com/SyncSoftInc/proxy/protoc/proxy"

type IProxyStore interface {
	Rent() *proxy.Proxy
	Return(*proxy.Proxy)
}
