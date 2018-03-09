package types // import "github.com/docker/docker/api/types"

import "github.com/docker/docker/api/types/filters"

// SecretListOptions holds parameters to list secrets
type SecretListOptions struct {
	Filters filters.Args
}

// ConfigListOptions holds parameters to list configs
type ConfigListOptions struct {
	Filters filters.Args
}

// NodeListOptions holds parameters to list nodes with.
type NodeListOptions struct {
	Filters filters.Args
}

// ServiceListOptions holds parameters to list services with.
type ServiceListOptions struct {
	Filters filters.Args
}

// TaskListOptions holds parameters to list tasks with.
type TaskListOptions struct {
	Filters filters.Args
}

// ServiceCreateResponse contains the information returned to a client
// on the creation of a new service.
type ServiceCreateResponse struct {
	// ID is the ID of the created service.
	ID string
	// Warnings is a set of non-fatal warning messages to pass on to the user.
	Warnings []string `json:",omitempty"`
}

// Values for RegistryAuthFrom in ServiceUpdateOptions
const (
	RegistryAuthFromSpec         = "spec"
	RegistryAuthFromPreviousSpec = "previous-spec"
)

// ServiceUpdateOptions contains the options to be used for updating services.
type ServiceUpdateOptions struct {
	// EncodedRegistryAuth is the encoded registry authorization credentials to
	// use when updating the service.
	//
	// This field follows the format of the X-Registry-Auth header.
	EncodedRegistryAuth string

	// TODO(stevvooe): Consider moving the version parameter of ServiceUpdate
	// into this field. While it does open API users up to racy writes, most
	// users may not need that level of consistency in practice.

	// RegistryAuthFrom specifies where to find the registry authorization
	// credentials if they are not given in EncodedRegistryAuth. Valid
	// values are "spec" and "previous-spec".
	RegistryAuthFrom string

	// Rollback indicates whether a server-side rollback should be
	// performed. When this is set, the provided spec will be ignored.
	// The valid values are "previous" and "none". An empty value is the
	// same as "none".
	Rollback string

	// QueryRegistry indicates whether the service update requires
	// contacting a registry. A registry may be contacted to retrieve
	// the image digest and manifest, which in turn can be used to update
	// platform or other information about the service.
	QueryRegistry bool
}
