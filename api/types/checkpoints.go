package types

// CheckpointCreateOptions holds parameters to create a checkpoint from a container
type CheckpointCreateOptions struct {
	CheckpointID  string
	CheckpointDir string
	Exit          bool
}

// CheckpointListOptions holds parameters to list checkpoints for a container
type CheckpointListOptions struct {
	CheckpointDir string
}

// CheckpointDeleteOptions holds parameters to delete a checkpoint from a container
type CheckpointDeleteOptions struct {
	CheckpointID  string
	CheckpointDir string
}

// Checkpoint represents the details of a checkpoint
type Checkpoint struct {
	Name string // Name is the name of the checkpoint
}
