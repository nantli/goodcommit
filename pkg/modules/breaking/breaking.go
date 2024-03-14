package breaking

import (
	"github.com/charmbracelet/huh"
	"github.com/nantli/goodcommit/pkg/module"
)

const MODULE_NAME = "breaking"

type Breaking struct {
	config module.Config
}

func (s *Breaking) LoadConfig() error {
	return nil
}

func (s *Breaking) NewField(commit *module.CommitInfo) (huh.Field, error) {

	return huh.NewConfirm().
		Title("🙊・Does this commit introduce a Breaking Change?").
		Affirmative("Yes 🚨").
		Negative("No 🏖️").
		Value(&commit.Breaking), nil
}

func (s *Breaking) PostProcess(commit *module.CommitInfo) error {
	return nil
}

func (s *Breaking) GetConfig() module.Config {
	return s.config
}

func (s *Breaking) SetConfig(config module.Config) {
	s.config = config
}

func (s *Breaking) Debug() error {

	return nil
}

func (s *Breaking) GetName() string {
	return MODULE_NAME
}

func New() module.Module {
	return &Breaking{config: module.Config{Name: MODULE_NAME}}
}
