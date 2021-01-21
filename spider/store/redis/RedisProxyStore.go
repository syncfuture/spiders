package redis

import (
	"encoding/json"
	"errors"

	"github.com/go-redis/redis/v7"
	"github.com/syncfuture/go/sid"
	"github.com/syncfuture/go/srand"
	"github.com/syncfuture/go/sredis"
	"github.com/syncfuture/spiders/spider/model"
	"github.com/syncfuture/spiders/spider/store"
)

type RedisProxyStore struct {
	key         string
	client      redis.UniversalClient
	idGenerator sid.IIDGenerator
}

func NewRedisProxyStore(key string, config *sredis.RedisConfig) store.IProxyStore {
	return &RedisProxyStore{
		key:         key,
		client:      sredis.NewClient(config),
		idGenerator: sid.NewSonyflakeIDGenerator(),
	}
}

func (x *RedisProxyStore) SaveProxy(proxy *model.Proxy) error {
	if proxy.ID == "" {
		proxy.ID = x.idGenerator.GenerateString()
	}
	data, err := json.Marshal(proxy)
	if err != nil {
		return err
	}
	err = x.client.HSet(x.key, proxy.ID, string(data)).Err()
	return err
}

func (x *RedisProxyStore) GetProxy(id string) (r *model.Proxy, err error) {
	j, err := x.client.HGet(x.key, id).Bytes()
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(j, &r)
	if err != nil {
		return nil, err
	}

	return
}

func (x *RedisProxyStore) GetRandomProxy() (r *model.Proxy, err error) {
	proxies, err := x.GetProxies()
	if err != nil {
		return nil, err
	}

	if len(proxies) == 0 {
		return nil, errors.New("no available proxies")
	}

	return proxies[srand.IntRange(0, len(proxies)-1)], err
}

func (x *RedisProxyStore) GetProxies() ([]*model.Proxy, error) {
	m, err := x.client.HGetAll(x.key).Result()
	if err != nil {
		return nil, err
	}

	r := make([]*model.Proxy, 0, len(m))
	for _, v := range m {
		var p *model.Proxy
		err = json.Unmarshal([]byte(v), &p)
		if err != nil {
			return nil, err
		}
		r = append(r, p)
	}
	return r, err
}
