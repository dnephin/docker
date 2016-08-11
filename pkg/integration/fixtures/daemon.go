package fixtures

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/docker/docker/pkg/testutil/fixtures"
)

// RunDaemon initializes the DaemonFixture by starting a docker daemon for tests
// TODO: use an env variable to optionally skip this fixture
func RunDaemon(t fixtures.TestingT) (fixtures.Fixture, error) {
	// TODO: use icmd.RunCommand (only print full stdout/stderr on failure)
	cmd := exec.Command("bash", "-exc", fmt.Sprintf("%s; %s",
		bundle(".integration-daemon-start"),
		bundle(".integration-daemon-setup")))
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		return nil, err
	}

	// TODO: remove dependency on hack/make.sh
	path, err := filepath.Abs(os.Getenv("DEST"))
	if err != nil {
		return nil, err
	}

	host := fmt.Sprintf("unix://%s/docker.sock", path)
	if err := os.Setenv("DOCKER_HOST", host); err != nil {
		return nil, err
	}
	return &DaemonFixture{Host: host}, nil
}

// DaemonFixture holds the state for a running Docker Daemon
type DaemonFixture struct {
	Host string
}

// Cleanup stops the Daemon and performs some cleanup
func (f *DaemonFixture) Cleanup() error {
	// TODO: use icmd.RunCommand (only print full stdout/stderr on failure)
	cmd := exec.Command("bash", "-exc", bundle(".integration-daemon-stop"))
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

func bundle(name string) string {
	// TODO: 'go_test_dir' changes the directory to integration-cli. Change this
	// path once 'go_test_dir' is removed
	return fmt.Sprintf("source ../hack/make/%s", name)
}
