package types

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/charmbracelet/huh"
	"github.com/nantli/goodcommit/pkg/module"
)

const MODULE_NAME = "types"

type Types struct {
	config module.Config
	Items  []module.Item `json:"types"`
}

func (s *Types) LoadConfig() error {

	if s.config.Path == "" {
		return nil
	}

	raw, err := os.ReadFile(s.config.Path)
	if err != nil {
		fmt.Println("Error occurred while reading types config:", err)
		os.Exit(1)
	}
	err = json.Unmarshal(raw, &s)
	if err != nil {
		fmt.Println("Error occurred while parsing types config:", err)
		os.Exit(1)
	}

	return nil
}

func (s *Types) NewField(commit *module.CommitInfo) (huh.Field, error) {

	var typeOptions []huh.Option[string]
	for _, i := range s.Items {
		typeOptions = append(typeOptions, huh.NewOption(i.Emoji+" "+i.Name+" - "+i.Title, i.Id))
	}
	return huh.NewSelect[string]().
		Options(typeOptions...).
		Title("🪰・Select a Commit Type").
		Description("Folowing the Conventional Commits specification.\n").
		Value(&commit.Type), nil
}

func (s *Types) PostProcess(commit *module.CommitInfo) error {
	if commit.Type == "" && s.IsActive() {
		return fmt.Errorf("commit type is required")
	}
	return nil
}

func (s *Types) GetConfig() module.Config {
	return s.config
}

func (s *Types) SetConfig(config module.Config) {
	s.config = config
}

func (s *Types) Debug() error {
	// print configuration and items in a human readable format
	fmt.Println(s.config)
	fmt.Println(s.Items)

	return nil
}

func (s *Types) GetName() string {
	return MODULE_NAME
}

func (s *Types) InitCommitInfo(commit *module.CommitInfo) error {
	return nil
}

func (s *Types) IsActive() bool {
	return s.config.Active
}

func New() module.Module {
	return &Types{config: module.Config{Name: MODULE_NAME}, Items: []module.Item{}}
}
