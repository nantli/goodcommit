// Package scopes provides a github.com/nantli/goodcommit module that allows the user to select scopes for the commit.
// It presents a multi-select menu with the available scopes.
// The selected scopes are then added to the commit title and body.
package scopes

import (
	"encoding/json"
	"fmt"
	"os"
	"slices"

	"github.com/charmbracelet/huh"
	gc "github.com/nantli/goodcommit"
)

// item is the structure for each entry in the scopes configuration file.
type item struct {
	Id          string   `json:"id"`
	Name        string   `json:"name"`
	Title       string   `json:"title"`
	Description string   `json:"description"`
	Emoji       string   `json:"emoji"`
	Conditional []string `json:"conditional"` // The types of commits that this scope is valid for.
}

// MODULE_NAME is the name of the module and should be used as the name of the module in the config.json file.
const MODULE_NAME = "scopes"

type scopes struct {
	config gc.ModuleConfig
	Items  []item `json:"scopes"`
}

func (s *scopes) item(id string) item {
	for _, i := range s.Items {
		if i.Id == id {
			return i
		}
	}
	return item{}
}

// LoadConfig loads the scopes configuration file.
// Example:
//
//	{
//		"scopes": [
//			{
//				"id": "modules",
//				"emoji": "📦",
//				"name": "Modules",
//				"description": "Use this when changes are made to modules",
//				"conditional": ["feat", "fix", "chore"]
//			}
//		]
//	}
func (s *scopes) LoadConfig() error {

	if s.config.Path == "" {
		return nil
	}

	raw, err := os.ReadFile(s.config.Path)
	if err != nil {
		fmt.Println("Error occurred while reading scopes config:", err)
		os.Exit(1)
	}
	err = json.Unmarshal(raw, &s)
	if err != nil {
		fmt.Println("Error occurred while parsing scopes config:", err)
		os.Exit(1)
	}

	return nil
}

// NewField returns a huh.MultiSelect field that allows the user to select the scopes for the commit.
// The options are built based on the selected commit type.
func (s *scopes) NewField(commit *gc.Commit) (huh.Field, error) {

	var typeOptions []huh.Option[string]
	for _, i := range s.Items {
		if slices.Contains(i.Conditional, commit.Type) {
			typeOptions = append(typeOptions, huh.NewOption(commit.Type+"("+i.Emoji+"): "+i.Name+" - "+i.Description, i.Id))
		}
	}

	if len(typeOptions) == 0 {
		return nil, fmt.Errorf("no valid scope options found for commit type: %s", commit.Type)
	}

	return huh.NewMultiSelect[string]().
		Options(typeOptions...).
		Title("🪱・Select Commit Scopes").
		Description("Additional contextual information about the changes. Multiple selections allowed.\n").
		Value(&commit.Scopes), nil // commit.Scopes should be a slice of strings
}

func (s *scopes) PostProcess(commit *gc.Commit) error {
	scopeHeader := "SCOPE: "
	scopeEmojis := ""
	if len(commit.Scopes) == 0 && s.IsActive() {
		commit.Scope = ""
		return nil
	}
	if len(commit.Scopes) > 1 {
		scopeHeader = "SCOPES: "
	}
	for _, scopeId := range commit.Scopes {
		if scopeId != "empty" {
			scopeHeader += s.item(scopeId).Name + " "
			scopeEmojis += s.item(scopeId).Emoji
		}
	}
	commit.Scope = scopeEmojis
	commit.Body = scopeHeader + "\n" + commit.Body

	return nil
}

func (s *scopes) Config() gc.ModuleConfig {
	return s.config
}

func (s *scopes) SetConfig(config gc.ModuleConfig) {
	s.config = config
}

func (s *scopes) Name() string {
	return MODULE_NAME
}

func (s *scopes) InitCommitInfo(commit *gc.Commit) error {
	return nil
}

func (s *scopes) IsActive() bool {
	return s.config.Active
}

// New returns a new instance of the scopes module.
// The scopes module is a github.com/nantli/goodcommit module that allows the user to select scopes for the commit.
// The selected scopes are then added to the commit title and body.
func New() gc.Module {
	return &scopes{config: gc.ModuleConfig{Name: MODULE_NAME}, Items: []item{}}
}
