package interceptor

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

type MockGrpcRequest struct {
	AString   string
	AnInteger int
}

func TestItConcatTheMethodAColonAndTheMArshalledRequest(t *testing.T) {
	kd := NewDefaultKeyDerivatorImpl()
	method := "/catalog/MyMethod"
	req := MockGrpcRequest{"a string", 1}

	key, err := kd.Derive(method, &req)

	assert.NoError(t, err)
	assert.Equal(t, `/catalog/MyMethod:{"AString":"a string","AnInteger":1}`, string(key.([]byte)))
}
