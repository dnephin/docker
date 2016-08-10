package fixtures

import (
	"testing"

	"github.com/docker/docker/pkg/testutil/assert"
)

type MockFixture struct {
	calledCleanup bool
}

func (m *MockFixture) Cleanup() error {
	m.calledCleanup = true
	return nil
}

func mockInit(t TestingT) (Fixture, error) {
	return &MockFixture{}, nil
}

func TestRegisterFixture(t *testing.T) {
	fixture := Register(t, mockInit, Suite)

	// Register with expanded scope should return the same object
	assert.Equal(t, fixture, Register(t, mockInit, Global))

	name := "github.com/docker/docker/pkg/testutil/fixtures.mockInit"
	_, exists := reg.active[Suite][name]
	assert.Equal(t, exists, true)
}

func TestCleanupScope(t *testing.T) {
	Register(t, mockInit, Suite)

	Cleanup(t, Suite)
	assert.Equal(t, len(reg.active[Suite]), 0)
}
