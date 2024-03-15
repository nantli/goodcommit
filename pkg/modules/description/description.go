package description

import (
	"strings"

	"github.com/charmbracelet/huh"
	"github.com/nantli/goodcommit/pkg/module"
)

const MODULE_NAME = "description"

type Description struct {
	config module.Config
}

func (d *Description) LoadConfig() error {
	// No configuration to load for this module.
	return nil
}

func (d *Description) NewField(commit *module.CommitInfo) (huh.Field, error) {
	return huh.NewInput().
		Title("🐵・Write the Commit Description").
		Description("Briefly describe the changes in this commit (max 50 chars).").
		CharLimit(50).
		Value(&commit.Description), nil
}

func (d *Description) PostProcess(commit *module.CommitInfo) error {
	// lowercase first letter of description
	commit.Description = strings.ToLower(commit.Description[:1]) + commit.Description[1:]
	return nil
}

func (d *Description) GetConfig() module.Config {
	return d.config
}

func (d *Description) SetConfig(config module.Config) {
	d.config = config
}

func (d *Description) Debug() error {
	// Optionally implement debugging information.
	return nil
}

func (d *Description) GetName() string {
	return MODULE_NAME
}

func New() module.Module {
	return &Description{config: module.Config{Name: MODULE_NAME}}
}
