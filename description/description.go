// Package description provides a module for goodcommit that prompts the user
// for a brief description of the commit.
package description

import (
	"strings"

	"github.com/charmbracelet/huh"
	gc "github.com/nantli/goodcommit"
)

const MODULE_NAME = "description"

type description struct {
	config gc.ModuleConfig
}

func (d *description) LoadConfig() error {
	return nil
}

// NewField returns a new Input field for the user to write a brief description of the commit (max 50 chars).
func (d *description) NewField(commit *gc.Commit) (huh.Field, error) {
	return huh.NewInput().
		Title("✏️・Write the Commit Description").
		Description("Briefly describe the changes in this commit (max 50 chars).").
		CharLimit(50).
		Value(&commit.Description), nil
}

// PostProcess lowercases the first letter of the commit description.
func (d *description) PostProcess(commit *gc.Commit) error {
	if commit.Description == "" {
		return nil
	}
	commit.Description = strings.TrimSuffix(commit.Description, ".")
	commit.Description = strings.ToLower(commit.Description[:1]) + commit.Description[1:]
	return nil
}

func (d *description) Config() gc.ModuleConfig {
	return d.config
}

func (d *description) SetConfig(config gc.ModuleConfig) {
	d.config = config
}

func (d *description) Name() string {
	return MODULE_NAME
}

func (d *description) InitCommitInfo(commit *gc.Commit) error {
	return nil
}

func (d *description) IsActive() bool {
	return d.config.Active
}

func New() gc.Module {
	return &description{config: gc.ModuleConfig{Name: MODULE_NAME}}
}
