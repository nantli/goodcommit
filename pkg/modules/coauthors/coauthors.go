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

func (c *CoAuthors) NewField(commit *module.CommitInfo) (huh.Field, error) {
	var coAuthorOptions []huh.Option[string]
	for _, item := range c.Items {
		coAuthorOptions = append(coAuthorOptions, huh.NewOption(item.Name, item.Id))
	}

	return huh.NewMultiSelect[string]().
		Title("👥・Select Co-Authors").
		Description("Choose co-authors for this commit.").
		Options(coAuthorOptions...).
		Value(&commit.CoAuthoredBy), nil
}

func (c *CoAuthors) PostProcess(commit *module.CommitInfo) error {
	// Additional processing after form submission if needed
	return nil
}

func (c *CoAuthors) GetConfig() module.Config {
	return c.config
}

func (c *CoAuthors) SetConfig(config module.Config) {
	c.config = config
}

func (c *CoAuthors) Debug() error {
	// Optionally implement debugging information
	return nil
}

func (c *CoAuthors) GetName() string {
	return MODULE_NAME
}

func New() module.Module {
	return &CoAuthors{config: module.Config{Name: MODULE_NAME}}
}
