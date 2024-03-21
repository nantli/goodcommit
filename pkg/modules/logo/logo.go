package logo

import (
	"github.com/charmbracelet/huh"
	"github.com/nantli/goodcommit/pkg/commit"
	"github.com/nantli/goodcommit/pkg/module"
)

const MODULE_NAME = "logo"

type Logo struct {
	config module.Config
}

func (l *Logo) LoadConfig() error {
	return nil
}

func (l *Logo) NewField(commit *commit.Config) (huh.Field, error) {
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

func (l *Logo) PostProcess(commit *commit.Config) error {
	// No post-processing needed for the Logo module.
	return nil
}

func (l *Logo) GetConfig() module.Config {
	return l.config
}

func (l *Logo) SetConfig(config module.Config) {
	l.config = config
}

func (l *Logo) GetName() string {
	return MODULE_NAME
}

func (l *Logo) IsActive() bool {
	return l.config.Active
}

func (l *Logo) InitCommitInfo(commit *commit.Config) error {
	return nil
}

func New() module.Module {
	return &Logo{config: module.Config{Name: MODULE_NAME}}
}
