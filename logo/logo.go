package logo

import (
	"fmt"
	"os"

	"github.com/charmbracelet/huh"
	gc "github.com/nantli/goodcommit"
)

const MODULE_NAME = "logo"

type logo struct {
	config   gc.ModuleConfig
	asciiArt string // Add this line
}

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
	// No post-processing needed for the Logo module.
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

func New() gc.Module {
	return &logo{
		config:   gc.ModuleConfig{Name: MODULE_NAME},
		asciiArt: "", // Initialize with an empty string
	}
}
