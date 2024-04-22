// Package logo provides a github.com/nantli/goodcommit module that shows a logo.
// It allows for extra personalization by pinning a logo to the top of every page of the goodcommit flow.
package logo

import (
	"fmt"
	"os"

	"github.com/charmbracelet/huh"
	gc "github.com/nantli/goodcommit"
)

// MODULE_NAME is the name of the module and should be used as the name of the module in the config.json file.
const MODULE_NAME = "logo"

type logo struct {
	config   gc.ModuleConfig
	asciiArt string // The ascii art to display in the commit message.
}

// LoadConfig loads the ascii art from the config file.
// the config file can be any text file, there're no specific requirements.
func (l *logo) LoadConfig() error {
	if l.config.Path != "" {
		raw, err := os.ReadFile(l.config.Path)
		if err != nil {
			return fmt.Errorf("failed to read logo file: %w", err)
		}
		l.asciiArt = string(raw)
	}
	return nil
}

// NewField returns a huh.Note field that displays the ascii art.
func (l *logo) NewField(commit *gc.Commit) (huh.Field, error) {
	if l.asciiArt == "" {
		l.asciiArt = `
    ┌─────────────────────────────────────┐
    │  You're gonna like this commit...   │
    └─────────────────────────────────────┘` // default ascii art
	}
	return huh.NewNote().Title(l.asciiArt), nil
}

func (l *logo) PostProcess(commit *gc.Commit) error {
	return nil
}

func (l *logo) Config() gc.ModuleConfig {
	return l.config
}

func (l *logo) SetConfig(config gc.ModuleConfig) {
	l.config = config
}

func (l *logo) Name() string {
	return MODULE_NAME
}

func (l *logo) IsActive() bool {
	return l.config.Active
}

func (l *logo) InitCommitInfo(commit *gc.Commit) error {
	return nil
}

// New returns a new instance of the logo module.
// The logo module is a github.com/nantli/goodcommit module that can be used to pin a logo
// to the top of every page of the goodcommit flow.
func New() gc.Module {
	return &logo{
		config:   gc.ModuleConfig{Name: MODULE_NAME},
		asciiArt: "", // Initialize with an empty string
	}
}
