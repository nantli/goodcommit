package scopes

import (
	"encoding/json"
	"fmt"
	"os"
	"slices"
	"strings"

	"github.com/charmbracelet/huh"
	"github.com/nantli/goodcommit/pkg/module"
)

const MODULE_NAME = "scopes"

type commitScope struct {
	Name        string   `json:"name"`
	Emoji       string   `json:"emoji"`
	Description string   `json:"description"`
	Types       []string `json:"types"`
}

type Scopes struct {
	config module.Config
	Items  []commitScope `json:"scopes"`
}

func (s *Scopes) Load() error {

	if s.config.Path == "" {
		return nil
	}

	raw, err := os.ReadFile(s.config.Path)
	if err != nil {
		fmt.Println("Error occurred while reading config:", err)
		os.Exit(1)
	}
	err = json.Unmarshal(raw, &s)
	if err != nil {
		fmt.Println("Error occurred while parsing config:", err)
		os.Exit(1)
	}

	return nil
}

func (s *Scopes) NewField(commit *module.CommitInfo) (huh.Field, error) {

	var typeOptions []huh.Option[string]
	for _, i := range s.Items {
		if slices.Contains(i.Types, commit.Type) {
			typeOptions = append(typeOptions, huh.NewOption(i.Emoji+"\t- "+i.Name, i.Name))
		}
	}

	return huh.NewSelect[string]().
		Options(typeOptions...).
		Title("Commit scope").
		Description("Select the scope for these changes.").
		Value(&commit.Scope), nil
}

func (s *Scopes) PostProcess(commit *module.CommitInfo) error {
	if commit.Scope == "" {
		return fmt.Errorf("commit scope is required")
	}
	commit.Scope = strings.ToLower(commit.Scope)
	return nil
}

func (s *Scopes) GetConfig() module.Config {
	return s.config
}

func New(config module.Config) (module.Module, error) {
	return &Scopes{config, []commitScope{}}, nil
}
