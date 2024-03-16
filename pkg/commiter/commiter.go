// Package commiter defines the Commiter interface which outlines the methods
// required for a commiter implementation in the goodcommit application.
package commiter

type Commiter interface {
	Execute(accessible bool) error
}
