// +build !experimental

package bundle

import (
	"github.com/docker/docker/api/client"
	"github.com/spf13/cobra"
)

// NewBundleCommand returns no command
func NewBundleCommand(dockerCli *client.DockerCli) *cobra.Command {
	return &cobra.Command{}
}
