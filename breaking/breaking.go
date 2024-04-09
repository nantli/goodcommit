// Package breaking provides a github.com/nantli/goodcommit module for indicating breaking changes.
// It prompts the user to indicate whether the commit introduces a breaking change.
package breaking

import (
	"github.com/charmbracelet/huh"
	gc "github.com/nantli/goodcommit"
)

// MODULE_NAME is the name of the module and should be used as the name of the module in the config.json file.
const MODULE_NAME = "breaking"

type breaking struct {
	config gc.ModuleConfig
}

func (b *breaking) LoadConfig() error {
	return nil
}

// NewField returns a new huh.Confirm field for indicating breaking changes.
// Only appears if the commit type is "feat" or "fix".
func (b *breaking) NewField(commit *gc.Commit) (huh.Field, error) {

	if commit.Type != "feat" && commit.Type != "fix" {
		return nil, nil
	}

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
	// No initialization of the commit is done by the breaking module.
	return nil
}

func (b *breaking) IsActive() bool {
	return b.config.Active
}

// New returns a new instance of the breaking module.
// The breaking module is a github.com/nantli/goodcommit module that is used to indicate breaking changes introduced by a commit.
func New() gc.Module {
	return &breaking{config: gc.ModuleConfig{Name: MODULE_NAME}}
}
