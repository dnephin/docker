package types // import "github.com/docker/docker/client/types"

import "github.com/docker/docker/api/types/filters"

// SecretListOptions holds parameters to list secrets
type SecretListOptions struct {
	Filters filters.Args
}

// ConfigListOptions holds parameters to list configs
type ConfigListOptions struct {
	Filters filters.Args
}
