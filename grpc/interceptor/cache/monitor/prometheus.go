package monitor

import (
	"time"

	"github.com/nsitbon/grpc-cache-interceptor/grpc/interceptor/cache"
)

type Option func(*PrometheusMonitor)

type Recorder interface {
	IncCacheKeyMetric(service string, key string, metric string)
}

type PrometheusMonitor struct {
	recorder     Recorder
	cache        cache.Cache
	keyExtractor cache.KeyExtractor
	onError      func(err error)
}

func (d *PrometheusMonitor) Get(key interface{}) (interface{}, bool) {
	ret, found := d.cache.Get(key)

	if parts, err := d.keyExtractor.Extract(key); err == nil {
		d.recorder.IncCacheKeyMetric(parts.PackageName, parts.MethodName+"Cache", getMetric(found))
	} else if d.onError != nil {
		d.onError(err)
	}

	return ret, found
}

func getMetric(found bool) string {
	if found {
		return "hit_count"
	} else {
		return "miss_count"
	}
}

func (d *PrometheusMonitor) Set(key interface{}, value interface{}) {
	d.cache.Set(key, value)
}

func (d *PrometheusMonitor) SetWithTTL(key interface{}, value interface{}, ttl time.Duration) {
	d.cache.SetWithTTL(key, value, ttl)
}

func NewPrometheusMonitor(recorder Recorder, innerCache cache.Cache, keyExtractor cache.KeyExtractor, opts ...Option) cache.Cache {
	m := &PrometheusMonitor{recorder: recorder, cache: innerCache, keyExtractor: keyExtractor}

	for _, opt := range opts {
		opt(m)
	}

	return m
}

func WithErrorFunc(onError func(error)) Option {
	return func(that *PrometheusMonitor) {
		that.onError = onError
	}
}
