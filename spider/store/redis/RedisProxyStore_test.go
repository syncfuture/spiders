package redis

import (
	"fmt"
	"sync"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/syncfuture/go/sredis"
	"github.com/syncfuture/spiders/spider/model"
	"github.com/syncfuture/spiders/spider/store"
)

var (
	_store store.IProxyStore
)

func init() {
	_store = NewRedisProxyStore("ams:Proxies", &sredis.RedisConfig{
		Addrs:    []string{"localhost:6379"},
		Password: "Famous901",
	})
}

func TestRedisProxyStore_SaveProxy(t *testing.T) {
	max := 1
	wg := new(sync.WaitGroup)
	wg.Add(max)

	for i := 0; i < max; i++ {
		go func(idx int) {
			defer wg.Done()
			err := _store.SaveProxy(&model.Proxy{
				Scheme:   fmt.Sprintf("URL %d", idx),
				Host:     fmt.Sprintf("Host %d:Port %d", idx, idx),
				Username: fmt.Sprintf("Username %d", idx),
				Password: fmt.Sprintf("Password %d", idx),
			})
			if err != nil {
				t.Log(err)
			}
		}(i)
	}

	wg.Wait()
}

func TestRedisProxyStore_GetProxy(t *testing.T) {
	a, err := _store.GetProxy("4b110fc0a00bca6")
	assert.NoError(t, err)

	assert.NotNil(t, a)
	assert.NotEmpty(t, a.Host)
}

func TestRedisProxyStore_GetProxies(t *testing.T) {
	a, err := _store.GetProxies()
	assert.NoError(t, err)

	assert.NotEmpty(t, a)
}
