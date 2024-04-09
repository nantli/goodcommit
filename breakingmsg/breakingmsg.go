// Package breakingmsg provides a github.com/nantli/goodcommit module for writing the breaking message.
// It presents the user with a free-form text box in which they can write
// a detailed description of the breaking changes introduced by the commit.
// It only appears if the commit is previously marked as containing a breaking change (by the breaking module for example).
package breakingmsg

import (
	"fmt"

	"github.com/charmbracelet/huh"
	gc "github.com/nantli/goodcommit"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

// MODULE_NAME is the name of the module and should be used as the name of the module in the config.json file.
const MODULE_NAME = "breakingmsg"

type breakingMsg struct {
	config gc.ModuleConfig
}

func (bm *breakingMsg) LoadConfig() error {
	// No configuration needed for this module.
	return nil
}

// NewField returns a huh.Text field for the breaking message.
// It only appears if the commit is marked as breaking.
func (bm *breakingMsg) NewField(commit *gc.Commit) (huh.Field, error) {
	// Only show this field if the commit is marked as breaking and not a chore
	if commit.Breaking {
		return huh.NewText().
			Title("💥・Breaking Changes Details").
			Description("Provide detailed information about the breaking changes.\n").
			Value(commit.Extras["breakingmsg"]).
			Editor("vim"), nil
	}
	return nil, nil
}

func (bm *breakingMsg) PostProcess(commit *gc.Commit) error {
	if commit.Extras["breakingmsg"] == nil || *commit.Extras["breakingmsg"] == "" {
		return nil
	}
	// Capitalize first letter of breaking message
	caser := cases.Title(language.English)
	*commit.Extras["breakingmsg"] = caser.String((*commit.Extras["breakingmsg"])[:1]) + (*commit.Extras["breakingmsg"])[1:]
	// Add a period at the end of the breaking message if it doesn't have one
	if (*commit.Extras["breakingmsg"])[len(*commit.Extras["breakingmsg"])-1] != '.' {
		*commit.Extras["breakingmsg"] += "."
	}
	// At the end of the body, add a new line and the breaking message
	commit.Body = fmt.Sprintf("%s\n\nBREAKING CHANGE: %s", commit.Body, *commit.Extras["breakingmsg"])
	return nil
}

func (bm *breakingMsg) Config() gc.ModuleConfig {
	return bm.config
}

func (bm *breakingMsg) SetConfig(config gc.ModuleConfig) {
	bm.config = config
}

func (bm *breakingMsg) Name() string {
	return MODULE_NAME
}

func (bm *breakingMsg) InitCommitInfo(commit *gc.Commit) error {
	placeholder := ""
	commit.Extras["breakingmsg"] = &placeholder
	return nil
}

func (bm *breakingMsg) IsActive() bool {
	return bm.config.Active
}

// New returns a new instance of the breakingmsg module.
// The breakingmsg module is a github.com/nantli/goodcommit module that is used to write the breaking changes message.
func New() gc.Module {
	return &breakingMsg{config: gc.ModuleConfig{Name: MODULE_NAME}}
}
