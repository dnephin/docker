// +build experimental

package bundle

import (
	"fmt"

	"github.com/docker/docker/api/client"
	"github.com/docker/docker/cli"
	"github.com/spf13/cobra"
)

// NewBundleCommand returns a cobra command for `bundle` subcommands
func NewBundleCommand(dockerCli *client.DockerCli) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "bundle",
		Short: "Manage Docker bundles",
		Args:  cli.NoArgs,
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Fprintf(dockerCli.Err(), "\n"+cmd.UsageString())
		},
	}
	cmd.AddCommand(
		newListCommand(dockerCli),
		newRemoveCommand(dockerCli),
	)
	return cmd
}
