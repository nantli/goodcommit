// Package types provides a github.com/nantli/goodcommit module that can be used to select the type of the commit.
// It presents the user with a selection of types and allows them to select one.
// The selected type is then added to the commit title.
package types

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"github.com/charmbracelet/huh"
	gc "github.com/nantli/goodcommit"
)

// item is the structure for each entry in the types configuration file.
type item struct {
	Id          string `json:"id"`
	Name        string `json:"name"`
	Title       string `json:"title"`
	Description string `json:"description"`
	Emoji       string `json:"emoji"`
}

// MODULE_NAME is the name of the module and should be used as the name of the module in the config.json file.
const MODULE_NAME = "types"

type types struct {
	config gc.ModuleConfig
	Items  []item `json:"types"`
}

// LoadConfig loads the types configuration file.
// Example config file:
//
//	{
//		"types": [
//			{
//				"id": "feat",
//				"name": "Feature",
//				"title": "A new feature",
//				"description": "A new feature",
//				"emoji": "✨"
//			}
//		]
//	}
func (t *types) LoadConfig() error {

	if t.config.Path == "" {
		return nil
	}

	raw, err := os.ReadFile(t.config.Path)
	if err != nil {
		fmt.Println("Error occurred while reading types config:", err)
		os.Exit(1)
	}
	err = json.Unmarshal(raw, &t)
	if err != nil {
		fmt.Println("Error occurred while parsing types config:", err)
		os.Exit(1)
	}

	return nil
}

// NewField returns a huh.Select field that allows the user to select the type of the commit.
func (t *types) NewField(commit *gc.Commit) (huh.Field, error) {

	var typeOptions []huh.Option[string]
	for _, i := range t.Items {
		typeOptions = append(typeOptions, huh.NewOption(i.Emoji+" "+i.Name+" - "+i.Title, i.Id))
	}
	return huh.NewSelect[string]().
		Options(typeOptions...).
		Title("🪰・Select a Commit Type").
		Description("Folowing the Conventional Commits specification.\n").
		Value(&commit.Type), nil
}

func (t *types) PostProcess(commit *gc.Commit) error {
	if commit.Type == "" && t.IsActive() {
		return fmt.Errorf("commit type is required")
	}
	commit.Type = strings.ToLower(commit.Type)
	return nil
}

func (t *types) Config() gc.ModuleConfig {
	return t.config
}

func (t *types) SetConfig(config gc.ModuleConfig) {
	t.config = config
}

func (t *types) Name() string {
	return MODULE_NAME
}

func (t *types) InitCommitInfo(commit *gc.Commit) error {
	return nil
}

func (t *types) IsActive() bool {
	return t.config.Active
}

// New returns a new instance of the types module.
// The types module is a github.com/nantli/goodcommit module that can be used to select the type of the commit.
// The selected type is then added to the commit title.
func New() gc.Module {
	return &types{config: gc.ModuleConfig{Name: MODULE_NAME}, Items: []item{}}
}
