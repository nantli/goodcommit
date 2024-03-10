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

func (s *Types) Load() error {

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

func (s *Types) NewField(commit *module.CommitInfo) (huh.Field, error) {

	var typeOptions []huh.Option[string]
	for _, i := range s.Items {
		typeOptions = append(typeOptions, huh.NewOption(i.Name+" - "+i.Title, i.Name))
	}
	return huh.NewSelect[string]().
		Options(typeOptions...).
		Title("Commit type").
		Description("Select the type of change that you're committing.").
		Value(&commit.Type), nil
}

func (s *Types) PostProcess(commit *module.CommitInfo) error {
	if commit.Type == "" {
		return fmt.Errorf("commit type is required")
	}
	return nil
}

func (s *Types) GetConfig() module.Config {
	return s.config
}

func New(config module.Config) (module.Module, error) {
	return &Types{config, []module.Item{}}, nil
}
