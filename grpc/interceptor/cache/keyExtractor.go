package cache

import (
	"fmt"
	"regexp"
)

const keyPattern = `^/(([^.]+)\.([^/]+)/([^:]+)):`

var keyRegex *regexp.Regexp

type KeyParts struct {
	PackageName string
	ServiceName string
	MethodName  string
	MethodLongName  string
}

type KeyExtractor interface {
	Extract(key interface{}) (*KeyParts, error)
}

type DefaultKeyExtractorImpl struct{}

func (d *DefaultKeyExtractorImpl) Extract(key interface{}) (*KeyParts, error) {
	if keyAsString, ok := keyIsString(key); !ok {
		return nil, fmt.Errorf("expected key type is []byte but got %T", key)
	} else if parts := keyRegex.FindStringSubmatch(keyAsString); len(parts) != 5 {
		return nil, fmt.Errorf("key expected format is '%s' but got '%s'", keyPattern, keyAsString)
	} else {
		return &KeyParts{
			MethodLongName:parts[1],
			PackageName: parts[2],
			ServiceName: parts[3],
			MethodName:  parts[4],
		}, nil
	}
}

func keyIsString(key interface{}) (string, bool) {
	if bytes, ok := key.([]byte); ok {
		return string(bytes), true
	}

	return "", false
}

func NewDefaultKeyExtractorImpl() *DefaultKeyExtractorImpl {
	return &DefaultKeyExtractorImpl{}
}

func init() {
	keyRegex = regexp.MustCompile(keyPattern)
}
