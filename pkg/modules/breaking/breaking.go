package breaking

import (
	"github.com/charmbracelet/huh"
	"github.com/nantli/goodcommit/pkg/module"
)

const MODULE_NAME = "breaking"

type Breaking struct {
	config module.Config
}

func (s *Breaking) Load() error {
	return nil
}

func (s *Breaking) NewField(commit *module.CommitInfo) (huh.Field, error) {

	return huh.NewConfirm().
		Title("Are you sure?").
		Affirmative("Yes!").
		Negative("No.").
		Value(&commit.Breaking), nil
}

func (s *Breaking) PostProcess(commit *module.CommitInfo) error {
	return nil
}

func (s *Breaking) GetConfig() module.Config {
	return s.config
}

func New(config module.Config) (module.Module, error) {
	return &Breaking{config}, nil
}
