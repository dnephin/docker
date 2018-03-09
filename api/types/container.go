package types

import "github.com/docker/docker/api/types/filters"

// ContainerListOptions holds parameters to list containers with.
type ContainerListOptions struct {
	Quiet   bool
	Size    bool
	All     bool
	Latest  bool
	Since   string
	Before  string
	Limit   int
	Filters filters.Args
}
