// Package body provides a github.com/nantli/goodcommit module for writing the commit body.
// It presents the user with a free-form text box in which they can write
// a detailed description of the changes made in the commit.
package body

import (
	"golang.org/x/text/cases"
	"golang.org/x/text/language"

	"github.com/charmbracelet/huh"
	gc "github.com/nantli/goodcommit"
)

// MODULE_NAME is the name of the module and should be used as the name of the module in the config.json file.
const MODULE_NAME = "body"

type body struct {
	config gc.ModuleConfig
}

func (b *body) LoadConfig() error {
	// No configuration to load for this module.
	return nil
}

// NewField returns a huh.Text field that will be used to write the commit body.
func (b *body) NewField(commit *gc.Commit) (huh.Field, error) {
	return huh.NewText().
		Title("📖・Write the Commit Body").
		Description("Provide a more detailed description of the changes (ctrl+j creates a new line).").
		Value(&commit.Body).
		Editor("vim"), nil
}

func (b *body) PostProcess(commit *gc.Commit) error {
	if commit.Body == "" {
		return nil
	}

	// Capitalize first letter of body
	caser := cases.Title(language.English)
	commit.Body = caser.String(commit.Body[:1]) + commit.Body[1:]

	// Add a period at the end of the body if it doesn't have one
	if commit.Body[len(commit.Body)-1] != '.' {
		commit.Body += "."
	}

	return nil
}

func (b *body) Name() string {
	return MODULE_NAME
}

func (b *body) Config() gc.ModuleConfig {
	return b.config
}

func (b *body) SetConfig(config gc.ModuleConfig) {
	b.config = config
}

func (b *body) InitCommitInfo(commit *gc.Commit) error {
	// No initialization of the commit is done by the body module.
	return nil
}

func (b *body) IsActive() bool {
	return b.config.Active
}

// New returns a new instance of the body module.
// The body module is a github.com/nantli/goodcommit module that is used to write the commit body.
func New() gc.Module {
	return &body{config: gc.ModuleConfig{Name: MODULE_NAME}}
}
