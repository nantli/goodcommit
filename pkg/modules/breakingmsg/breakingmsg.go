package breakingmsg

import (
	"fmt"

	"github.com/charmbracelet/huh"
	"github.com/nantli/goodcommit/pkg/commit"
	"github.com/nantli/goodcommit/pkg/module"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

const MODULE_NAME = "breakingmsg"

type BreakingMsg struct {
	config module.Config
}

func (bm *BreakingMsg) LoadConfig() error {
	// Load configuration if necessary
	return nil
}

func (bm *BreakingMsg) NewField(commit *commit.Config) (huh.Field, error) {
	// Only show this field if the commit is marked as breaking and not a chore
	if commit.Breaking && commit.Type != "chore" {
		return huh.NewText().
			Title("💥・Breaking Changes Details").
			Description("Provide detailed information about the breaking changes.\n").
			Value(commit.Extras["breakingmsg"]).
			Editor("vim"), nil
	}
	return nil, nil
}

func (bm *BreakingMsg) PostProcess(commit *commit.Config) error {
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
	// At at the end of the Body, add a new line and the breaking message
	commit.Body = fmt.Sprintf("%s\n\nBREAKING CHANGE: %s", commit.Body, *commit.Extras["breakingmsg"])
	return nil
}

func (bm *BreakingMsg) GetConfig() module.Config {
	return bm.config
}

func (bm *BreakingMsg) SetConfig(config module.Config) {
	bm.config = config
}

func (bm *BreakingMsg) GetName() string {
	return MODULE_NAME
}

func (bm *BreakingMsg) InitCommitInfo(commit *commit.Config) error {
	placeholder := ""
	commit.Extras["breakingmsg"] = &placeholder
	return nil
}

func (bm *BreakingMsg) IsActive() bool {
	return bm.config.Active
}

func New() module.Module {
	return &BreakingMsg{config: module.Config{Name: MODULE_NAME}}
}
