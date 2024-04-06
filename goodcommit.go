package goodcommit

type Commiter interface {
	LoadModules(modules []Module) error
	RunForm(accessible bool) error
	RunPostProcessing() error
	PreviewCommit()
	RenderMessage() string
}

type goodCommit struct {
	commiter Commiter
}

func (g *goodCommit) Execute(accessible bool) (string, error) {
	if err := g.commiter.RunForm(accessible); err != nil {
		return "", err
	}
	if err := g.commiter.RunPostProcessing(); err != nil {
		return "", err
	}
	g.commiter.PreviewCommit()
	return g.commiter.RenderMessage(), nil
}

func New(c Commiter) *goodCommit {
	return &goodCommit{commiter: c}
}
