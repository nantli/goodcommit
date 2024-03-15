package goodcommit

import "github.com/nantli/goodcommit/pkg/commiter"

// GoodCommit handles the execution of a commiter.
type GoodCommit struct {
	commiter commiter.Commiter
}

// New creates a new GoodCommit instance with the provided commiter.
func New(c commiter.Commiter) *GoodCommit {
	return &GoodCommit{commiter: c}
}

// Execute runs the commiter's methods in the desired order.
func (g *GoodCommit) Execute(accessible bool) error {
	return g.commiter.Execute(accessible)
}
