package runner

import (
	"fmt"
	"testing"

	"github.com/docker/docker/pkg/reexec"
	"github.com/go-check/check"

	// Trigger inclusion of test suites
	"github.com/docker/docker/integration-cli"
)

// Test starts the integration test suites
func Test(t *testing.T) {
	reexec.Init() // This is required for external graphdriver tests

	if !integration.IsLocalDaemon() {
		fmt.Println("INFO: Testing against a remote daemon")
	} else {
		fmt.Println("INFO: Testing against a local daemon")
	}

	integration.AddSuites()
	check.TestingT(t)
}
