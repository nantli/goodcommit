// Package why provides a module for goodcommit that prompts the user
// to explain why the change was needed.
package why

import (
	"fmt"

	"github.com/charmbracelet/huh"
	"github.com/nantli/goodcommit/pkg/commit"
	"github.com/nantli/goodcommit/pkg/module"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
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
func (w *Why) NewField(commit *commit.Config) (huh.Field, error) {
	return huh.NewInput().
		Title("❔・Why was this change needed?").
		Description("Explain the reason for this change (max 100 chars).").
		CharLimit(100).
		Value(commit.Extras["why"]), nil
}

// PostProcess prepends the value of the Why field to the commit body
func (w *Why) PostProcess(commit *commit.Config) error {
	if commit.Extras["why"] == nil || *commit.Extras["why"] == "" {
		return nil
	}
	// Capitalize first letter of why
	caser := cases.Title(language.English)
	*commit.Extras["why"] = caser.String((*commit.Extras["why"])[:1]) + (*commit.Extras["why"])[1:]
	// Add a period at the end of the why if it doesn't have one
	if (*commit.Extras["why"])[len(*commit.Extras["why"])-1] != '.' {
		*commit.Extras["why"] += "."
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

func (w *Why) InitCommitInfo(commit *commit.Config) error {
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
