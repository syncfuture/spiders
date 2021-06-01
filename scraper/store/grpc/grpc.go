package grpc

import (
	"context"

	"github.com/SyncSoftInc/proxy/protoc/proxy"
	log "github.com/syncfuture/go/slog"
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
	mr, err := x.client.Return(context.Background(), &proxy.ReturnCommand{
		Proxy: p,
	})
	u.LogError(err)
	if mr.MsgCode != "" {
		log.Error(mr.MsgCode)
	}
}
