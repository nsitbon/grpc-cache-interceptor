package cache

import (
	"time"
)

var NullCache Cache = &nullCache{}

type nullCache struct {}

func (n nullCache) Get(_ interface{}) (interface{}, bool) {
	return nil, false
}

func (n nullCache) Set(_ interface{}, _ interface{}) {}

func (n nullCache) SetWithTTL(_ interface{}, _ interface{}, _ time.Duration) {}

type Cache interface {
	Get(key interface{}) (interface{}, bool)
	Set(key interface{}, value interface{})
	SetWithTTL(key interface{}, value interface{}, ttl time.Duration)
}
