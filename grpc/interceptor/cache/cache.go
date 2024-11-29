package cache

import (
	"context"

	"github.com/nsitbon/grpc-cache-interceptor/grpc/interceptor"
	"google.golang.org/grpc"
)

type CachingInterceptorKeyDerivator func(string, interface{}) (interface{}, error)

type CachingInterceptorOption func(interceptor *CachingInterceptor)

func WithMemMove(memMove MemMoveFn) CachingInterceptorOption {
	return func(i *CachingInterceptor) {
		i.memMove = memMove
	}
}

type MemMoveFn func(dest, src interface{}) error

type CachingInterceptor struct {
	memMove      MemMoveFn
	keyDerivator interceptor.KeyDerivator
	cache        Cache
}

func NewCachingInterceptor(cache Cache, options ...CachingInterceptorOption) *CachingInterceptor {
	i := &CachingInterceptor{cache: cache, keyDerivator: interceptor.NewDefaultKeyDerivatorImpl()}

	for _, opt := range options {
		opt(i)
	}

	if i.memMove == nil {
		i.memMove = interceptor.MemMove
	}

	return i
}

func (i *CachingInterceptor) Intercept(ctx context.Context, method string, req, reply interface{}, cc *grpc.ClientConn, invoker grpc.UnaryInvoker, opts ...grpc.CallOption) error {
	if key, err := i.keyDerivator.Derive(method, req); err != nil {
		return err
	} else if i.getFromCacheAndSetReply(method, key, reply) {
		return nil
	} else if err = invoker(ctx, method, req, reply, cc, opts...); err == nil {
		i.cache.Set(key, reply)
		return nil
	} else {
		return err
	}
}

/*
v and reply are both interface whose dynamic types are equal and of kind pointer to T.
Golang doesn't support generics so we implement genericity by simply copying the content of v (of type T)
to the content of reply (also of type T). To do that we "treat" the content as raw bytes and create a byte slice
backed by the content (of type T).
*/
func (i *CachingInterceptor) getFromCacheAndSetReply(method string, key interface{}, reply interface{}) bool {
	if v, ok := i.cache.Get(key); ok {
		return i.memMove(reply, v) == nil
	}

	return false
}
