package cache

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDefaultKeyExtractorImpl_Extract(t *testing.T) {
	e := NewDefaultKeyExtractorImpl()
	expectedPackage := "catalog"
	expectedService := "Catalog"
	expectedMethod := "GetCustomVideosByProgramId"
	key := []byte(fmt.Sprintf(`/%s.%s/%s:{"some json":1}`, expectedPackage, expectedService, expectedMethod))

	parts, err := e.Extract(key)

	assert.NoError(t, err)
	assert.EqualValues(t, expectedPackage, parts.PackageName)
	assert.EqualValues(t, expectedService, parts.ServiceName)
	assert.EqualValues(t, expectedMethod, parts.MethodName)
}
