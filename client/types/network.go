package types // import "github.com/docker/docker/client/types"

// NetworkInspectOptions holds parameters to inspect network
type NetworkInspectOptions struct {
	Scope   string
	Verbose bool
}
