package body

import (
	"golang.org/x/text/cases"
	"golang.org/x/text/language"

	"github.com/charmbracelet/huh"
	"github.com/nantli/goodcommit/pkg/commit"
	"github.com/nantli/goodcommit/pkg/module"
)

const MODULE_NAME = "body"

type Body struct {
	config module.Config
}

func (b *Body) LoadConfig() error {
	// No configuration to load for this module.
	return nil
}

func (b *Body) NewField(commit *commit.Config) (huh.Field, error) {
	return huh.NewText().
		Title("📖・Write the Commit Body").
		Description("Provide a more detailed description of the changes (ctrl+j creates a new line).").
		Value(&commit.Body).
		Editor("vim"), nil
}

func (b *Body) PostProcess(commit *commit.Config) error {
	// Capitalize first letter of body
	caser := cases.Title(language.English)
	commit.Body = caser.String(commit.Body[:1]) + commit.Body[1:]
	// Add a period at the end of the body if it doesn't have one
	if commit.Body[len(commit.Body)-1] != '.' {
		commit.Body += "."
	}
	return nil
}

func (b *Body) GetName() string {
	return MODULE_NAME
}

func (b *Body) GetConfig() module.Config {
	return b.config
}

func (b *Body) SetConfig(config module.Config) {
	b.config = config
}

func (b *Body) InitCommitInfo(commit *commit.Config) error {
	// No initialization needed for this module.
	return nil
}

func (b *Body) IsActive() bool {
	return b.config.Active
}

func New() module.Module {
	return &Body{config: module.Config{Name: MODULE_NAME}}
}
