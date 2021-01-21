package webshare

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestWebShareProxyStore_getProxiesFromAPI(t *testing.T) {
	x := &WebShareProxyStore{
		key: "be09d781115fe3491743fa205ea786852513f474",
	}
	rs, err := x.getProxiesFromAPI()
	assert.NoError(t, err)
	assert.NotEmpty(t, rs)
}
