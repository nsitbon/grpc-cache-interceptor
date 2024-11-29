package cache

import (
	"encoding/json"
	"testing"
	"time"

	grpc_testing "github.com/grpc-ecosystem/go-grpc-middleware/testing"
	pb_testproto "github.com/grpc-ecosystem/go-grpc-middleware/testing/testproto"
	"google.golang.org/grpc"
)

type GetCall struct {
	Key   interface{}
	Value interface{}
	Found bool
}

type SetCall struct {
	Key   interface{}
	Value interface{}
}

type MockingCache struct {
	values   map[interface{}]interface{}
	GetCalls []GetCall
	SetCalls []SetCall
}

func (m *MockingCache) SetWithTTL(key interface{}, value interface{}, ttl time.Duration) {}

func NewMockingCache() *MockingCache {
	return &MockingCache{values: make(map[interface{}]interface{})}
}

func (m *MockingCache) Get(key interface{}) (interface{}, bool) {
	v, ok := m.values[string(key.([]byte))]
	m.GetCalls = append(m.GetCalls, GetCall{Key: key, Value: v, Found: ok})
	return v, ok
}

func (m *MockingCache) Set(key interface{}, value interface{}) {
	m.values[string(key.([]byte))] = value
	m.SetCalls = append(m.SetCalls, SetCall{Key: key, Value: value})
}

func deriveKey(key interface{}) string {
	if v, err := json.Marshal(key); err != nil {
		panic(err)
	} else {
		return "/mwitkow.testproto.TestService/Ping:" + string(v)
	}
}

func TestItCallsTheInvokerOnCacheMiss(t *testing.T) {
	cache := NullCache
	its := getSuite(cache, t)
	defer its.TearDownSuite()

	request := &pb_testproto.PingRequest{Value: "my_fake_ping_payload"}
	resp, err := its.Client.Ping(its.SimpleCtx(), request)

	its.NoError(err)
	its.Equal("my_fake_ping_payload", resp.Value)
	its.Equal(int32(42), resp.Counter)
}

func TestItFirstTryToGetValueFromCache(t *testing.T) {
	cache := NewMockingCache()
	its := getSuite(cache, t)
	defer its.TearDownSuite()

	request := &pb_testproto.PingRequest{Value: "my_fake_ping_payload"}
	_, _ = its.Client.Ping(its.SimpleCtx(), request)

	its.Len(cache.GetCalls, 1)
	its.EqualValues(string(cache.GetCalls[0].Key.([]byte)), deriveKey(request))
	its.False(cache.GetCalls[0].Found)
}

func TestItSetsTheCacheAfterCallingTheInvoker(t *testing.T) {
	cache := NewMockingCache()
	its := getSuite(cache, t)
	defer its.TearDownSuite()

	request := &pb_testproto.PingRequest{Value: "my_fake_ping_payload"}
	resp, _ := its.Client.Ping(its.SimpleCtx(), request)

	its.Len(cache.SetCalls, 1)
	its.EqualValues(cache.SetCalls[0].Key, deriveKey(request))
	its.EqualValues(cache.SetCalls[0].Value, resp)
}

func TestItDoesntCallTheInvokerOnCacheHit(t *testing.T) {
	cache := NewMockingCache()
	its := getSuite(cache, t)
	defer its.TearDownSuite()

	request := &pb_testproto.PingRequest{Value: "my_fake_ping_payload"}

	expectedResponse := &pb_testproto.PingResponse{Value: "expected value", Counter: 1}
	cache.values[deriveKey(request)] = expectedResponse
	response, err := its.Client.Ping(its.SimpleCtx(), request)

	its.NoError(err)
	its.Equal(expectedResponse.Value, response.Value)
	its.Equal(expectedResponse.Counter, response.Counter)

	its.Len(cache.GetCalls, 1)
	its.Len(cache.SetCalls, 0)
}

func getSuite(cache Cache, t *testing.T) *grpc_testing.InterceptorTestSuite {
	its := &grpc_testing.InterceptorTestSuite{ClientOpts: []grpc.DialOption{
		grpc.WithUnaryInterceptor(NewCachingInterceptor(cache).Intercept)},
	}
	its.Suite.SetT(t)
	its.SetupSuite()
	return its
}
