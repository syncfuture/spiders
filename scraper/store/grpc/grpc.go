package grpc

import (
	"context"

	"github.com/SyncSoftInc/proxy/protoc/proxy"
	"github.com/syncfuture/go/u"
	"github.com/syncfuture/spiders/scraper/store"
	"google.golang.org/grpc"
)

func NewGRPCProxyStore(addr, provider string) store.IProxyStore {
	conn, err := grpc.Dial(addr, grpc.WithInsecure())
	u.LogFaltal(err)
	return &GRPCProxyStore{
		provider: provider,
		client:   proxy.NewProxyServiceClient(conn),
	}
}

type GRPCProxyStore struct {
	provider string
	client   proxy.ProxyServiceClient
}

func (x *GRPCProxyStore) Rent() *proxy.Proxy {
	p, err := x.client.Rent(context.Background(), &proxy.RentCommand{
		Provider: x.provider,
	})

	if u.LogError(err) {
		return nil
	}

	return p
}
func (x *GRPCProxyStore) Return(p *proxy.Proxy) {
	x.client.Return(context.Background(), &proxy.ReturnCommand{
		Provider: x.provider,
		Proxy:    p,
	})
}
