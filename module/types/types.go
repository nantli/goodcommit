package types

import (
	"github.com/charmbracelet/huh"
	"github.com/nantli/goodcommit/module"
)

const MODULE_NAME = "types"

type Types struct {
	config module.Config
}

func (s *Types) NewField(commit *module.CommitInfo) (huh.Field, error) {
	return huh.NewNote().Title(MODULE_NAME).Next(true), nil
}

func New(config module.Config) (module.Module, error) {
	return &Types{config}, nil
}
