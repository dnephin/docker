package stack

// Backend is all the methods that need to be implemented
// to provide stack specific functionality.
type Backend interface {
	CreateStack(name, bundle string) (string, error) // TODO: add config
}
