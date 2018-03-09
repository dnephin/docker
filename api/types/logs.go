package types // import "github.com/docker/docker/api/types"

// ContainerLogsOptions holds parameters to filter logs with.
type ContainerLogsOptions struct {
	ShowStdout bool
	ShowStderr bool
	Since      string
	Until      string
	Timestamps bool
	Follow     bool
	Tail       string
	Details    bool
}
