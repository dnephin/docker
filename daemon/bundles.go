package daemon

import (
	"fmt"
	"path"
	"sort"
	"time"

	"github.com/docker/distribution/digest"
	"github.com/docker/docker/image/bundle"
	"github.com/docker/docker/reference"
	"github.com/docker/engine-api/types"
	"github.com/docker/engine-api/types/filters"
)

// CreateBundle(src, repository, tag string, inConfig io.ReadCloser, outStream io.Writer) error
// PullBundle(ctx context.Context, bundle, tag string, metaHeaders map[string][]string, authConfig *types.AuthConfig, outStream io.Writer) error
// PushBundle(ctx context.Context, bundle, tag string, metaHeaders map[string][]string, authConfig *types.AuthConfig, outStream io.Writer) error
// BundleDelete(bundleRef string, force, prune bool) ([]types.BundleDelete, error)
// Bundles(filterArgs string, filter string, all bool) ([]*types.Bundle, error)
// LookupBundle(name string) (*types.BundleInspect, error)
// TagBundle(bundleName, repository, tag string) error

var acceptedBundleFilterTags = map[string]bool{
	"label":  true,
	"before": true,
	"since":  true,
}

// bundleByCreated is a temporary type used to sort a list of bundles by creation
// time.
type bundleByCreated []*types.Bundle

func (r bundleByCreated) Len() int           { return len(r) }
func (r bundleByCreated) Swap(i, j int)      { r[i], r[j] = r[j], r[i] }
func (r bundleByCreated) Less(i, j int) bool { return r[i].Created < r[j].Created }

// Bundles returns a filtered list of bundles. filterArgs is a JSON-encoded set
// of filter arguments which will be interpreted by api/types/filters.
// filter is a shell glob string applied to repository names. The argument
// named all controls whether all bundles in the graph are filtered, or just
// the heads.
func (daemon *Daemon) Bundles(filterArgs, filter string) ([]*types.Bundle, error) {
	bundleFilters, err := filters.FromParam(filterArgs)
	if err != nil {
		return nil, err
	}
	if err := bundleFilters.Validate(acceptedBundleFilterTags); err != nil {
		return nil, err
	}

	bundles := make([]*types.Bundle, 0)

	var beforeFilter, sinceFilter *bundle.Bundle
	err = bundleFilters.WalkValues("before", func(value string) error {
		beforeFilter, err = daemon.GetBundle(value)
		return err
	})
	if err != nil {
		return nil, err
	}

	err = bundleFilters.WalkValues("since", func(value string) error {
		sinceFilter, err = daemon.GetBundle(value)
		return err
	})
	if err != nil {
		return nil, err
	}

	var filterTagged bool
	if filter != "" {
		filterRef, err := reference.ParseNamed(filter)
		if err == nil { // parse error means wildcard repo
			if _, ok := filterRef.(reference.NamedTagged); ok {
				filterTagged = true
			}
		}
	}

	for id, b := range daemon.bundleStore.Map() {
		if beforeFilter != nil {
			if b.Created.Equal(beforeFilter.Created) || b.Created.After(beforeFilter.Created) {
				continue
			}
		}

		if sinceFilter != nil {
			if b.Created.Equal(sinceFilter.Created) || b.Created.Before(sinceFilter.Created) {
				continue
			}
		}

		if bundleFilters.Include("label") {
			if !bundleFilters.MatchKVList("label", b.Labels) {
				continue
			}
		}

		newBundle := newAPIBundle(b)

		for _, ref := range daemon.bundleReferenceStore.References(digest.Digest(id)) {
			if filter != "" { // filter by tag/repo name
				if filterTagged { // filter by tag, require full ref match
					if ref.String() != filter {
						continue
					}
				} else if matched, err := path.Match(filter, ref.Name()); !matched || err != nil { // name only match, FIXME: docs say exact
					continue
				}
			}
			if _, ok := ref.(reference.Canonical); ok {
				newBundle.RepoDigests = append(newBundle.RepoDigests, ref.String())
			}
			if _, ok := ref.(reference.NamedTagged); ok {
				newBundle.RepoTags = append(newBundle.RepoTags, ref.String())
			}
		}

		bundles = append(bundles, newBundle)
	}

	sort.Sort(sort.Reverse(bundleByCreated(bundles)))

	return bundles, nil
}

func newAPIBundle(bundle *bundle.Bundle) *types.Bundle {
	newBundle := new(types.Bundle)
	newBundle.ID = bundle.ID().String()
	newBundle.Created = bundle.Created.Unix()
	newBundle.Labels = bundle.Labels
	return newBundle
}

// GetBundleID returns an bundle ID corresponding to the bundle referred to by
// refOrID.
func (daemon *Daemon) GetBundleID(refOrID string) (bundle.ID, error) {
	id, ref, err := reference.ParseIDOrReference(refOrID)
	if err != nil {
		return "", err
	}
	if id != "" {
		if _, err := daemon.bundleStore.Get(bundle.ID(id)); err != nil {
			return "", ErrRefDoesNotExist{refOrID}
		}
		return bundle.ID(id), nil
	}

	if id, err := daemon.bundleReferenceStore.Get(ref); err == nil {
		return bundle.ID(id), nil
	}
	if tagged, ok := ref.(reference.NamedTagged); ok {
		if id, err := daemon.bundleStore.Search(tagged.Tag()); err == nil {
			for _, namedRef := range daemon.bundleReferenceStore.References(digest.Digest(id)) {
				if namedRef.Name() == ref.Name() {
					return id, nil
				}
			}
		}
	}

	// Search based on ID
	if id, err := daemon.bundleStore.Search(refOrID); err == nil {
		return id, nil
	}

	return "", ErrRefDoesNotExist{refOrID}
}

// GetBundle returns an bundle corresponding to the bundle referred to by refOrID.
func (daemon *Daemon) GetBundle(refOrID string) (*bundle.Bundle, error) {
	imgID, err := daemon.GetBundleID(refOrID)
	if err != nil {
		return nil, err
	}
	return daemon.bundleStore.Get(imgID)
}

// LookupBundle looks up an Bundle by name and returns it as an BundleInspect
// structure.
func (daemon *Daemon) LookupBundle(name string) (*types.BundleInspect, error) {
	bundle, err := daemon.GetBundle(name)
	if err != nil {
		return nil, fmt.Errorf("no such bundle: %s", name)
	}

	// todo(tonistiigi): separate to func
	refs := daemon.bundleReferenceStore.References(digest.Digest(bundle.ID()))
	repoTags := []string{}
	repoDigests := []string{}
	for _, ref := range refs {
		switch ref.(type) {
		case reference.NamedTagged:
			repoTags = append(repoTags, ref.String())
		case reference.Canonical:
			repoDigests = append(repoDigests, ref.String())
		}
	}

	bundleInspect := &types.BundleInspect{
		ID:            bundle.ID().String(),
		RepoTags:      repoTags,
		RepoDigests:   repoDigests,
		Created:       bundle.Created.Format(time.RFC3339Nano),
		DockerVersion: bundle.DockerVersion,
	}

	for _, s := range bundle.Services {
		img, err := daemon.LookupImage(string(s.Image))
		if err != nil {
			return nil, err
		}
		sInspect := &types.BundleService{
			Name:  s.Name,
			Image: img,
		}
		bundleInspect.Services = append(bundleInspect.Services, sInspect)
	}

	return bundleInspect, nil
}
