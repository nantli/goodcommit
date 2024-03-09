package scope

import (
	"github.com/charmbracelet/huh"
	"github.com/nantli/goodcommit/module"
)

const MODULE_NAME = "scope"

type Scope struct {
	config module.Config
}

func (s *Scope) NewField(commit *module.CommitInfo) (huh.Field, error) {
	return huh.NewNote().Title(MODULE_NAME).Next(true), nil
}

func New(config module.Config) (module.Module, error) {
	return &Scope{config}, nil
}
