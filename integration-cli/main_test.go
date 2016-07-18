package integration

import (
	"fmt"
	"testing"

	"github.com/docker/docker/pkg/reexec"
	"github.com/go-check/check"
)

// Test starts the integration test suites
func Test(t *testing.T) {
	reexec.Init() // This is required for external graphdriver tests

	if !IsLocalDaemon() {
		fmt.Println("INFO: Testing against a remote daemon")
	} else {
		fmt.Println("INFO: Testing against a local daemon")
	}

	AddSuites()
	check.TestingT(t)
}
