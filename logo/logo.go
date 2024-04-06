package logo

import (
	"github.com/charmbracelet/huh"
	gc "github.com/nantli/goodcommit"
)

const MODULE_NAME = "logo"

type logo struct {
	config gc.ModuleConfig
}

func (l *logo) LoadConfig() error {
	return nil
}

func (l *logo) NewField(commit *gc.Commit) (huh.Field, error) {
	asciiArt := `                          
		 ____ ____ ____ ____ ____ ____      
		||N |||a |||n |||t |||l |||i ||     
		||__|||__|||__|||__|||__|||__||     
		|/__\|/__\|/__\|/__\|/__\|/__\|     
	┌─────────────────────────────────────┐ 
	│  You're gonna like this commit...   │ 
	└─────────────────────────────────────┘ 
	`
	return huh.NewNote().Title(asciiArt), nil
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
	return &logo{config: gc.ModuleConfig{Name: MODULE_NAME}}
}
