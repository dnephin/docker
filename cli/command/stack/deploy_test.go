package stack

import (
	"github.com/docker/docker/pkg/testutil/assert"
	"testing"
)

func TestConvertDeployModeGlobal(t *testing.T) {
	mode, err := convertDeployMode("global", nil)
	assert.NilError(t, err)
	assert.NotNil(t, mode.Global)
}

func TestConvertDeployModeGlobalWithReplicas(t *testing.T) {
	replicas := uint64(4)
	mode, err := convertDeployMode("global", &replicas)
	assert.Error(t, err, "used with replicated mode")
}

func TestConvertDeployModeReplicated(t *testing.T) {
	mode, err := convertDeployMode("replicated", nil)
	assert.NilError(t, err)
	assert.NotNil(t, mode.Replicated)
}
