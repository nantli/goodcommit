package commiter

// Commiter defines the interface for a goodcommit handler.
type Commiter interface {
	Execute(accessible bool) error
}
