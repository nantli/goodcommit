package scopes

import (
	"encoding/json"
	"fmt"
	"os"
	"slices"

	"github.com/charmbracelet/huh"
	gc "github.com/nantli/goodcommit"
)

type item struct {
	Id          string   `json:"id"`
	Name        string   `json:"name"`
	Title       string   `json:"title"`
	Description string   `json:"description"`
	Emoji       string   `json:"emoji"`
	Conditional []string `json:"conditional"`
}

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

	// Adjusted to use MultiSelect and to work with a slice of strings for multiple selections
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

func (s *scopes) Debug() error {
	// print configuration and items in a human readable format
	fmt.Println(s.config)
	fmt.Println(s.Items)

	return nil
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

func New() gc.Module {
	return &scopes{config: gc.ModuleConfig{Name: MODULE_NAME}, Items: []item{}}
}
