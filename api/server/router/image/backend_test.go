package image

import (
	"io"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/backend"
	"github.com/docker/docker/api/types/filters"
	"github.com/docker/docker/api/types/image"
	"github.com/docker/docker/api/types/registry"
	"golang.org/x/net/context"
)

type fakeBackend struct {
}

var _ Backend = &fakeBackend{}

func (f *fakeBackend) Commit(name string, config *backend.ContainerCommitConfig) (imageID string, err error) {
	return "", nil
}

func (f *fakeBackend) ImageDelete(imageRef string, force, prune bool) ([]types.ImageDeleteResponseItem, error) {
	return []types.ImageDeleteResponseItem{}, nil
}

func (f *fakeBackend) ImageHistory(imageName string) ([]*image.HistoryResponseItem, error) {
	return []*image.HistoryResponseItem{}, nil
}

func (f *fakeBackend) Images(imageFilters filters.Args, all bool, withExtraAttrs bool) ([]*types.ImageSummary, error) {
	return []*types.ImageSummary{}, nil
}

func (f *fakeBackend) LookupImage(name string) (*types.ImageInspect, error) {
	return nil, nil
}

func (f *fakeBackend) TagImage(imageName, repository, tag string) error {
	return nil
}

func (f *fakeBackend) ImagesPrune(ctx context.Context, pruneFilters filters.Args) (*types.ImagesPruneReport, error) {
	return nil, nil
}

func (f *fakeBackend) LoadImage(inTar io.ReadCloser, outStream io.Writer, quiet bool) error {
	return nil
}

func (f *fakeBackend) ImportImage(src string, repository, platform string, tag string, msg string, inConfig io.ReadCloser, outStream io.Writer, changes []string) error {
	return nil
}

func (f *fakeBackend) ExportImage(names []string, outStream io.Writer) error {
	return nil
}

func (f *fakeBackend) PullImage(ctx context.Context, image, tag, platform string, metaHeaders map[string][]string, authConfig *types.AuthConfig, outStream io.Writer) error {
	return nil
}

func (f *fakeBackend) PushImage(ctx context.Context, image, tag string, metaHeaders map[string][]string, authConfig *types.AuthConfig, outStream io.Writer) error {
	return nil
}

func (f *fakeBackend) SearchRegistryForImages(ctx context.Context, filtersArgs string, term string, limit int, authConfig *types.AuthConfig, metaHeaders map[string][]string) (*registry.SearchResults, error) {
	return nil, nil
}
