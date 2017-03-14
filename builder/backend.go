package builder

import (
	"context"
	"io"

	"github.com/docker/docker/api/types"
)

type backend struct {
}

type GetImageOptions struct {
	PullParent  bool
	AuthConfigs map[string]types.AuthConfig
	Output      io.Writer
}

// GetImage returns the image referenced by name
func (b *backend) GetImage(ctx context.Context, name string, opts GetImageOptions) (Image, error) {
	var image Image
	var err error

	if !opts.PullParent {
		// TODO: don't use `name`, instead resolve it to a digest
		image, _ = b.docker.GetImageOnBuild(name)
		// TODO: shouldn't we error out if error is different from "not found" ?
	}
	if image == nil {
		image, err = b.docker.PullOnBuild(ctx, name, opts.AuthConfigs, opts.Output)
	}
	return image, err
}
