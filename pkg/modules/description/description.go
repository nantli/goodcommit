// Package description provides a module for goodcommit that prompts the user
// for a brief description of the commit.
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
	return nil
}

// NewField returns a new Input field for the user to write a brief description of the commit (max 50 chars).
func (d *Description) NewField(commit *module.CommitInfo) (huh.Field, error) {
	return huh.NewInput().
		Title("✏️・Write the Commit Description").
		Description("Briefly describe the changes in this commit (max 50 chars).").
		CharLimit(50).
		Value(&commit.Description), nil
}

// PostProcess lowercases the first letter of the commit description.
func (d *Description) PostProcess(commit *module.CommitInfo) error {
	if commit.Description == "" {
		return nil
	}
	commit.Description = strings.ToLower(commit.Description[:1]) + commit.Description[1:]
	return nil
}

func (d *Description) GetConfig() module.Config {
	return d.config
}

func (d *Description) SetConfig(config module.Config) {
	d.config = config
}

func (d *Description) GetName() string {
	return MODULE_NAME
}

func (d *Description) InitCommitInfo(commit *module.CommitInfo) error {
	return nil
}

func (d *Description) IsActive() bool {
	return d.config.Active
}

func New() module.Module {
	return &Description{config: module.Config{Name: MODULE_NAME}}
}
