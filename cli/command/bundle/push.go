// +build experimental

package bundle

import (
	"io"

	"golang.org/x/net/context"

	"github.com/docker/docker/api/client"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/cli"
	apiclient "github.com/docker/docker/client"
	"github.com/docker/docker/pkg/jsonmessage"
	"github.com/docker/docker/reference"
	"github.com/docker/docker/registry"
	"github.com/spf13/cobra"
)

// NewPushCommand creates a new `docker bundle push` command
func newPushCommand(dockerCli *client.DockerCli) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "push [OPTIONS] NAME[:TAG]",
		Short: "Push a bundle to a registry",
		Args:  cli.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return runPush(dockerCli, args[0])
		},
	}

	flags := cmd.Flags()
	client.AddTrustedFlags(flags, true)

	return cmd
}

func runPush(dockerCli *client.DockerCli, remote string) error {
	ref, err := reference.ParseNamed(remote)
	if err != nil {
		return err
	}

	// Resolve the Repository name from fqn to RepositoryInfo
	repoInfo, err := registry.ParseRepositoryInfo(ref)
	if err != nil {
		return err
	}

	ctx := context.Background()

	// Resolve the Auth config relevant for this server
	authConfig := dockerCli.ResolveAuthConfig(ctx, repoInfo.Index)
	requestPrivilege := dockerCli.RegistryAuthenticationPrivilegedFunc(repoInfo.Index, "push")

	if client.IsTrusted() {
		return dockerCli.TrustedPush(ctx, repoInfo, ref, authConfig, requestPrivilege)
	}

	responseBody, err := bundlePushPrivileged(ctx, dockerCli.Client(), authConfig, ref.String(), requestPrivilege)
	if err != nil {
		return err
	}
	defer responseBody.Close()

	return jsonmessage.DisplayJSONMessagesStream(responseBody, dockerCli.Out(), dockerCli.OutFd(), dockerCli.IsTerminalOut(), nil)
}

func bundlePushPrivileged(ctx context.Context, apiclient apiclient.APIClient, authConfig types.AuthConfig, ref string, requestPrivilege types.RequestPrivilegeFunc) (io.ReadCloser, error) {
	encodedAuth, err := client.EncodeAuthToBase64(authConfig)
	if err != nil {
		return nil, err
	}
	options := types.BundlePushOptions{
		RegistryAuth:  encodedAuth,
		PrivilegeFunc: requestPrivilege,
	}

	return apiclient.BundlePush(ctx, ref, options)
}
