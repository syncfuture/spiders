package webshare

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/syncfuture/go/sid"
	log "github.com/syncfuture/go/slog"
	"github.com/syncfuture/go/srand"
	"github.com/syncfuture/spiders/spider/model"
	"github.com/syncfuture/spiders/spider/store"
)

const (
	_apiURL = "https://proxy.webshare.io/api/proxy/list/?page=1"
)

var (
	_idGenerator sid.IIDGenerator = sid.NewSonyflakeIDGenerator()
)

type (
	WebShareProxyStore struct {
		key     string
		proxies []*model.Proxy
	}

	APIResultDTO struct {
		Count    int                 `json:"count,omitempty"`
		Next     string              `json:"next,omitempty"`
		Previous string              `json:"previous,omitempty"`
		Results  []*WebShareProxyDTO `json:"results,omitempty"`
	}

	WebShareProxyDTO struct {
		Username              string     `json:"username,omitempty"`
		Password              string     `json:"password,omitempty"`
		ProxyAddress          string     `json:"proxy_address,omitempty"`
		CountryCode           string     `json:"country_code,omitempty"`
		CountryCodeConfidence float32    `json:"country_code_confidence,omitempty"`
		Valid                 bool       `json:"valid,omitempty"`
		LastVerification      *time.Time `json:"last_verification,omitempty"`
		Ports                 PortsDTO   `json:"ports,omitempty"`
	}

	PortsDTO struct {
		Http   interface{} `json:"http,omitempty"`
		Socks5 interface{} `json:"socks5,omitempty"`
	}
)

func NewWebShareProxyStore(key string) store.IProxyStore {
	return &WebShareProxyStore{
		key:     key,
		proxies: make([]*model.Proxy, 0, 20),
	}
}

func (x *WebShareProxyStore) SaveProxy(proxy *model.Proxy) error {
	for i, p := range x.proxies {
		if p.ID == proxy.ID {
			if proxy.Blocked {
				log.Debugf("proxy [%s] is been blocked, [%d] proxies left.", proxy.Host, len(x.proxies))
				x.proxies = append(x.proxies[:i], x.proxies[i+1:]...)
			} else {
				x.proxies[i] = proxy
			}
			return nil
		}
	}

	return nil
}

func (x *WebShareProxyStore) GetProxy(id string) (*model.Proxy, error) {
	for i, proxy := range x.proxies {
		if proxy.ID == proxy.ID {
			return x.proxies[i], nil
		}
	}
	return nil, errors.New("not found")
}

func (x *WebShareProxyStore) GetRandomProxy() (r *model.Proxy, err error) {
	proxies, err := x.GetProxies()
	if err != nil {
		return nil, err
	}

	if len(proxies) == 0 {
		return nil, errors.New("no available proxies")
	}

	randomIndex := srand.IntRange(0, len(proxies)-1)
	return proxies[randomIndex], err
}

func (x *WebShareProxyStore) GetProxies() ([]*model.Proxy, error) {
	var err error
	if len(x.proxies) == 0 {
		x.proxies, err = x.getProxiesFromAPI()
	}
	return x.proxies, err
}

func (x *WebShareProxyStore) getProxiesFromAPI() ([]*model.Proxy, error) {
	log.Debug("fetching proxies from webshare api...")
	msg, _ := http.NewRequest("GET", _apiURL, nil)
	msg.Header.Set("Authorization", "Token "+x.key)
	resp, err := http.DefaultClient.Do(msg)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var rs *APIResultDTO
	err = json.Unmarshal(data, &rs)
	if err != nil {
		return nil, err
	}

	if len(rs.Results) == 0 {
		log.Warn(string(data))
		return nil, err
	}
	log.Debugf("[%d] proxies fetched", len(rs.Results))

	r := make([]*model.Proxy, 0, len(rs.Results))
	for _, dto := range rs.Results {
		r = append(r, dto.ToProxy())
	}

	return r, nil
}
func (x *WebShareProxyStore) Clear() error {
	x.proxies = nil
	log.Debug("proxies cleared")

	return nil
}

func (x *WebShareProxyDTO) ToProxy() (r *model.Proxy) {
	r = new(model.Proxy)

	r.ID = x.ProxyAddress
	r.Host = fmt.Sprintf("%s:%v", x.ProxyAddress, x.Ports.Http)
	r.Scheme = "http"
	r.Username = x.Username
	r.Password = x.Password

	return
}
