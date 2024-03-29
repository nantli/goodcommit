package scopes

import (
	"encoding/json"
	"fmt"
	"os"
	"slices"

	"github.com/charmbracelet/huh"
	"github.com/nantli/goodcommit/pkg/commit"
	"github.com/nantli/goodcommit/pkg/module"
)

const MODULE_NAME = "scopes"

type Scopes struct {
	config module.Config
	Items  []module.Item `json:"scopes"`
}

func (s *Scopes) getItem(id string) module.Item {
	for _, i := range s.Items {
		if i.Id == id {
			return i
		}
	}
	return module.Item{}
}

func (s *Scopes) LoadConfig() error {

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

func (s *Scopes) NewField(commit *commit.Config) (huh.Field, error) {

	var typeOptions []huh.Option[string]
	for _, i := range s.Items {
		if slices.Contains(i.Conditional, commit.Type) {
			typeOptions = append(typeOptions, huh.NewOption(commit.Type+"("+i.Emoji+"): "+i.Name+" - "+i.Description, i.Id))
		}
	}

	if len(typeOptions) == 0 {
		return nil, fmt.Errorf("no valid scope options found for commit type: %s", commit.Type)
	}

	return huh.NewSelect[string]().
		Options(typeOptions...).
		Title("🪱・Select a Commit Scope").
		Description("Additional contextual information about the changes.\n").
		Value(&commit.Scope), nil
}

func (s *Scopes) PostProcess(commit *commit.Config) error {
	if commit.Scope == "" && s.IsActive() {
		return fmt.Errorf("commit scope is required")
	}
	scopeId := commit.Scope
	commit.Scope = s.getItem(scopeId).Emoji
	if scopeId != "empty" {
		commit.Body = fmt.Sprintf("SCOPE: %s\n%s", s.getItem(scopeId).Name, commit.Body)
	}

	return nil
}

func (s *Scopes) GetConfig() module.Config {
	return s.config
}

func (s *Scopes) SetConfig(config module.Config) {
	s.config = config
}

func (s *Scopes) Debug() error {
	// print configuration and items in a human readable format
	fmt.Println(s.config)
	fmt.Println(s.Items)

	return nil
}

func (s *Scopes) GetName() string {
	return MODULE_NAME
}

func (s *Scopes) InitCommitInfo(commit *commit.Config) error {
	return nil
}

func (s *Scopes) IsActive() bool {
	return s.config.Active
}

func New() module.Module {
	return &Scopes{config: module.Config{Name: MODULE_NAME}, Items: []module.Item{}}
}
