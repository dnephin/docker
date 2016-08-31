// +build experimental

package bundle

import (
	"github.com/docker/docker/api/client"
	"github.com/docker/docker/cli"
	"github.com/docker/docker/opts"
	"github.com/spf13/cobra"
)

type listOptions struct {
	quiet  bool
	filter opts.FilterOpt
}

func newListCommand(dockerCli *client.DockerCli) *cobra.Command {
	opts := listOptions{filter: opts.NewFilterOpt()}

	cmd := &cobra.Command{
		Use:     "ls [OPTIONS]",
		Aliases: []string{"list"},
		Short:   "List bundles",
		Args:    cli.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			return runList(dockerCli, opts)
		},
	}

	flags := cmd.Flags()
	flags.BoolVarP(&opts.quiet, "quiet", "q", false, "Only display IDs")
	flags.VarP(&opts.filter, "filter", "f", "Filter output based on conditions provided")

	return cmd
}

func runList(dockerCli *client.DockerCli, opts listOptions) error {
	//	ctx := context.Background()
	//	client := dockerCli.Client()
	//
	//	bundles, err := client.BundleList(ctx, ...)
	//	if err != nil {
	//		return err
	//	}
	//
	//	out := dockerCli.Out()
	//	if opts.quiet {
	//		PrintQuiet(out, bundles)
	//	} else {
	//		taskFilter := filters.NewArgs()
	//		for _, service := range services {
	//			taskFilter.Add("service", service.ID)
	//		}
	//		PrintNotQuiet(out, services, nodes, tasks)
	//	}
	return nil
}

//func printQuiet(out io.Writer, bundles []types.Bundles) {
//	for _, service := range services {
//		fmt.Fprintln(out, service.ID)
//	}
//}
