/*
fixtures is a framework for registering test fixtures during unit or integration
testing.

Fixtures can be created with a scope, and will be cleaned up when that scope
exits. If a second test requires the same fixture, and uses the same scope
a fixture is re-used.

This package was inspired by http://doc.pytest.org/en/latest/fixture.html
*/

package fixtures

import (
	"fmt"
	"reflect"
	"runtime"
)

// Scope is an enumeration of scopes for a fixture.
type Scope string

const (
	// Global scope is for the life of the process
	Global Scope = "global"
	// Suite scope is for duration of the test suite
	Suite Scope = "suite"
	// Test scope is for the duration of a single test case
	Test Scope = "test"
)

// registry stores active fixtures
type registry struct {
	active map[Scope]map[string]Fixture
}

func (r *registry) add(name string, fixture Fixture, scope Scope) {
	r.active[scope][name] = fixture
}

func (r *registry) getActive(name string, scope Scope) Fixture {
	for _, scope := range expandScope(scope) {
		if fixture, ok := r.active[scope][name]; ok {
			return fixture
		}
	}
	return nil
}

func expandScope(scope Scope) []Scope {
	switch scope {
	case Test:
		return []Scope{Test}
	case Suite:
		return []Scope{Suite, Test}
	case Global:
		return []Scope{Global, Suite, Test}
	default:
		panic(fmt.Sprintf("Invalid Scope %q", scope))
	}
}

func newRegistry() *registry {
	active := make(map[Scope]map[string]Fixture)
	for _, scope := range []Scope{Global, Suite, Test} {
		active[scope] = make(map[string]Fixture)
	}
	return &registry{active: active}
}

// TestingT is the testing.T interface required by fixtures
type TestingT interface {
	Fatalf(string, ...interface{})
	Errorf(string, ...interface{})
}

// FixtureInit if a function which creates a test Fixture
type FixtureInit func(TestingT) (Fixture, error)

// Fixture is a test fixture which provides something required by a test case
type Fixture interface {
	Cleanup() error
}

// EmptyFixture is a struct that implements the Fixture interface. It can be
// used for fixtures that don't have any in-process state. An InitFixture
// function can return this struct to satisfy the interface.
type EmptyFixture struct{}

// Cleanup does nothing
func (f *EmptyFixture) Cleanup() error {
	return nil
}

var (
	reg *registry
)

func init() {
	reg = newRegistry()
}

func getFuncName(i FixtureInit) string {
	return runtime.FuncForPC(reflect.ValueOf(i).Pointer()).Name()
}

// Register a new fixture. If there is already an active fixture in the
// scope, the existing fixture will be returned.
func Register(t TestingT, init FixtureInit, scope Scope) Fixture {
	name := getFuncName(init)
	fixture := reg.getActive(name, scope)
	if fixture != nil {
		return fixture
	}

	fixture, err := init(t)
	if err != nil {
		t.Fatalf("Fixture %q failed: %s", name, err)
	}

	reg.add(name, fixture, scope)
	return fixture
}

// Cleanup runs cleanup on all fixtures in the scope and removes them from
// active fixtures
func Cleanup(t TestingT, scope Scope) {
	for name, fixture := range reg.active[scope] {
		if err := fixture.Cleanup(); err != nil {
			t.Errorf("Fixture %q failed cleanup: %s", name, err)
		}
		delete(reg.active[scope], name)
	}
}
