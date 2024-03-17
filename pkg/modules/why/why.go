// Package why provides a module for goodcommit that prompts the user
// to explain why the change was needed.
package why

import (
	"fmt"

	"github.com/charmbracelet/huh"
	"github.com/nantli/goodcommit/pkg/module"
)

const MODULE_NAME = "why"

type Why struct {
	config module.Config
}

func (w *Why) LoadConfig() error {
	// Load any necessary configuration for the Why module.
	return nil
}

// NewField returns a new Input field for the user to explain why the change was needed.
func (w *Why) NewField(commit *module.CommitInfo) (huh.Field, error) {
	return huh.NewInput().
		Title("❔・Why was this change needed?").
		Description("Explain the reason for this change (max 100 chars).").
		CharLimit(100).
		Value(commit.Extras["why"]), nil
}

// PostProcess prepends the value of the Why field to the commit body
func (w *Why) PostProcess(commit *module.CommitInfo) error {
	if commit.Extras["why"] == nil || *commit.Extras["why"] == "" {
		return nil
	}

	commit.Body = fmt.Sprintf("WHY: %s\n\n%s", *commit.Extras["why"], commit.Body)
	return nil
}

func (w *Why) GetConfig() module.Config {
	return w.config
}

func (w *Why) SetConfig(config module.Config) {
	w.config = config
}

func (w *Why) GetName() string {
	return MODULE_NAME
}

func (w *Why) InitCommitInfo(commit *module.CommitInfo) error {
	placeholder := ""
	commit.Extras["why"] = &placeholder
	return nil
}

func (w *Why) IsActive() bool {
	return w.config.Active
}

func New() module.Module {
	return &Why{config: module.Config{Name: MODULE_NAME}}
}
