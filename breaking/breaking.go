// Package breaking provides a module for goodcommit that prompts the user to indicate
// whether the commit introduces a breaking change.
package breaking

import (
	"github.com/charmbracelet/huh"
	gc "github.com/nantli/goodcommit"
)

const MODULE_NAME = "breaking"

type breaking struct {
	config gc.ModuleConfig
}

func (b *breaking) LoadConfig() error {
	return nil
}

// NewField returns a new Confirm field for the user to indicate whether the commit introduces a breaking change.
func (b *breaking) NewField(commit *gc.Commit) (huh.Field, error) {

	return huh.NewConfirm().
		Title("☎️・Does this commit introduce a Breaking Change?").
		Affirmative("Yes 🚨").
		Negative("No 🏖️").
		Value(&commit.Breaking), nil
}

// PostProcess adds ! symbol to commit type if the commit is a breaking change
func (b *breaking) PostProcess(commit *gc.Commit) error {

	return nil
}

func (b *breaking) Config() gc.ModuleConfig {
	return b.config
}

func (b *breaking) SetConfig(config gc.ModuleConfig) {
	b.config = config
}

func (b *breaking) Name() string {
	return MODULE_NAME
}

func (b *breaking) InitCommitInfo(commit *gc.Commit) error {
	// No initialization needed for this module.
	return nil
}

func (b *breaking) IsActive() bool {
	return b.config.Active
}

func New() gc.Module {
	return &breaking{config: gc.ModuleConfig{Name: MODULE_NAME}}
}
