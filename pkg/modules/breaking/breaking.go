// Package breaking provides a module for goodcommit that prompts the user to indicate
// whether the commit introduces a breaking change.
package breaking

import (
	"github.com/charmbracelet/huh"
	"github.com/nantli/goodcommit/pkg/commit"
	"github.com/nantli/goodcommit/pkg/module"
)

const MODULE_NAME = "breaking"

type Breaking struct {
	config module.Config
}

func (b *Breaking) LoadConfig() error {
	return nil
}

// NewField returns a new Confirm field for the user to indicate whether the commit introduces a breaking change.
func (b *Breaking) NewField(commit *commit.Config) (huh.Field, error) {

	return huh.NewConfirm().
		Title("☎️・Does this commit introduce a Breaking Change?").
		Affirmative("Yes 🚨").
		Negative("No 🏖️").
		Value(&commit.Breaking), nil
}

// PostProcess adds ! symbol to commit type if the commit is a breaking change
func (b *Breaking) PostProcess(commit *commit.Config) error {

	return nil
}

func (b *Breaking) GetConfig() module.Config {
	return b.config
}

func (b *Breaking) SetConfig(config module.Config) {
	b.config = config
}

func (b *Breaking) GetName() string {
	return MODULE_NAME
}

func (b *Breaking) InitCommitInfo(commit *commit.Config) error {
	// No initialization needed for this module.
	return nil
}

func (b *Breaking) IsActive() bool {
	return b.config.Active
}

func New() module.Module {
	return &Breaking{config: module.Config{Name: MODULE_NAME}}
}
