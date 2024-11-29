package interceptor

import (
	"encoding/json"
)

type KeyDerivator interface {
	Derive(method string, req interface{}) (interface{}, error)
}

type DefaultKeyDerivatorImpl struct{}

func (d *DefaultKeyDerivatorImpl) Derive(method string, req interface{}) (interface{}, error) {
	if raw, err := json.Marshal(req); err != nil {
		return nil, err
	} else {
		return append([]byte(method+":"), raw...), nil
	}
}

func NewDefaultKeyDerivatorImpl() KeyDerivator {
	return &DefaultKeyDerivatorImpl{}
}
