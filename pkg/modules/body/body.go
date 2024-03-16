package body

import (
	"github.com/charmbracelet/huh"
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

func (b *Body) NewField(commit *module.CommitInfo) (huh.Field, error) {
	return huh.NewText().
		Title("📖・Write the Commit Body").
		Description("Provide a more detailed description of the changes (ctrl+j creates a new line).").
		Value(&commit.Body).
		Editor("vim"), nil
}

func (b *Body) PostProcess(commit *module.CommitInfo) error {
	// No post-processing needed for this module.
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

func New() module.Module {
	return &Body{config: module.Config{Name: MODULE_NAME}}
}
