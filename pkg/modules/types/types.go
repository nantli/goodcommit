package types

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"github.com/charmbracelet/huh"
	"github.com/nantli/goodcommit/pkg/commit"
	"github.com/nantli/goodcommit/pkg/module"
)

const MODULE_NAME = "types"

type Types struct {
	config module.Config
	Items  []module.Item `json:"types"`
}

func (t *Types) LoadConfig() error {

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

func (t *Types) NewField(commit *commit.Config) (huh.Field, error) {

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

func (t *Types) PostProcess(commit *commit.Config) error {
	if commit.Type == "" && t.IsActive() {
		return fmt.Errorf("commit type is required")
	}
	commit.Type = strings.ToLower(commit.Type)
	return nil
}

func (t *Types) GetConfig() module.Config {
	return t.config
}

func (t *Types) SetConfig(config module.Config) {
	t.config = config
}

func (s *Types) Debug() error {
	// print configuration and items in a human readable format
	fmt.Println(s.config)
	fmt.Println(s.Items)

	return nil
}

func (t *Types) GetName() string {
	return MODULE_NAME
}

func (t *Types) InitCommitInfo(commit *commit.Config) error {
	return nil
}

func (t *Types) IsActive() bool {
	return t.config.Active
}

func New() module.Module {
	return &Types{config: module.Config{Name: MODULE_NAME}, Items: []module.Item{}}
}
