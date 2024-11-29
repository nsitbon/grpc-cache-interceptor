package adapter

import (
	"time"

	"github.com/dgraph-io/ristretto"
	"github.com/nsitbon/grpc-cache-interceptor/grpc/interceptor/cache"
)

type RistrettoCacheAdapter struct {
	cache *ristretto.Cache
}

func (r *RistrettoCacheAdapter) Get(key interface{}) (interface{}, bool) {
	return r.cache.Get(key)
}

func (r *RistrettoCacheAdapter) Set(key interface{}, value interface{}) {
	r.cache.Set(key, value, 0)
}

func (r *RistrettoCacheAdapter) SetWithTTL(key interface{}, value interface{}, ttl time.Duration) {
	r.cache.SetWithTTL(key, value, 0, ttl)
}

func NewRistrettoCacheAdapter(cache *ristretto.Cache) cache.Cache {
	return &RistrettoCacheAdapter{cache: cache}
}
