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

func (g *GoodCommit) Execute(accessible bool) (string, error) {
	if err := g.commiter.RunForm(accessible); err != nil {
		return "", err
	}
	if err := g.commiter.RunPostProcessing(); err != nil {
		return "", err
	}
	g.commiter.PreviewCommit()
	return g.commiter.RenderMessage(), nil
}
