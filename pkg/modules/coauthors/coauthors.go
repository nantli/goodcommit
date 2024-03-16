// Package coauthors provides a module for goodcommit that allows the user to select
// co-authors for the commit from a predefined list.
package coauthors

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/charmbracelet/huh"
	"github.com/nantli/goodcommit/pkg/module"
)

const MODULE_NAME = "coauthors"

type CoAuthors struct {
	config module.Config
	Items  []module.Item `json:"coauthors"`
}

func (c *CoAuthors) LoadConfig() error {
	if c.config.Path == "" {
		return nil
	}

	raw, err := os.ReadFile(c.config.Path)
	if err != nil {
		return fmt.Errorf("error reading config: %w", err)
	}
	err = json.Unmarshal(raw, c)
	if err != nil {
		return fmt.Errorf("error parsing config: %w", err)
	}

	return nil
}

// NewField returns a MultiSelect field with options for each co-author.
func (c *CoAuthors) NewField(commit *module.CommitInfo) (huh.Field, error) {
	var coAuthorOptions []huh.Option[string]
	for _, item := range c.Items {
		coAuthorOptions = append(coAuthorOptions, huh.NewOption(item.Name+" - "+item.Id, item.Name+" <"+item.Id+">"))
	}

	return huh.NewMultiSelect[string]().
		Title("👥・Select Co-Authors").
		Description("Choose co-authors for this commit.").
		Options(coAuthorOptions...).
		Value(&commit.CoAuthoredBy), nil
}

func (c *CoAuthors) PostProcess(commit *module.CommitInfo) error {
	return nil
}

func (c *CoAuthors) GetConfig() module.Config {
	return c.config
}

func (c *CoAuthors) SetConfig(config module.Config) {
	c.config = config
}

func (c *CoAuthors) GetName() string {
	return MODULE_NAME
}

func (c *CoAuthors) InitCommitInfo(commit *module.CommitInfo) error {
	// No initialization needed for this module.
	return nil
}

func (c *CoAuthors) IsActive() bool {
	return c.config.Active
}

func New() module.Module {
	return &CoAuthors{config: module.Config{Name: MODULE_NAME}}
}
