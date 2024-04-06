package types

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"

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

const MODULE_NAME = "types"

type types struct {
	config gc.ModuleConfig
	Items  []item `json:"types"`
}

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

func (s *types) Debug() error {
	// print configuration and items in a human readable format
	fmt.Println(s.config)
	fmt.Println(s.Items)

	return nil
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

func New() gc.Module {
	return &types{config: gc.ModuleConfig{Name: MODULE_NAME}, Items: []item{}}
}
