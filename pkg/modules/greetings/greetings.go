package greetings

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"

	"github.com/charmbracelet/huh"
	"github.com/nantli/goodcommit/pkg/module"
)

const MODULE_NAME = "greetings"

type Greetings struct {
	config module.Config
}

func (g *Greetings) LoadConfig() error {
	// Load any necessary configuration. For this module, there might not be any config to load.
	return nil
}

func (g *Greetings) NewField(commit *module.CommitInfo) (huh.Field, error) {
	stagedFiles, err := g.getStagedFiles()
	if err != nil {
		return nil, fmt.Errorf("error getting staged files: %w", err)
	}

	if len(stagedFiles) == 0 {
		return nil, fmt.Errorf("no staged files found")
	}

	// Display staged files as a simple text field in the form
	stagedFilesText := fmt.Sprintf("Staged Files:\n%s", stagedFiles)
	return huh.NewConfirm().Title("🐒・Do you want to commit these files?").Description(stagedFilesText).Validate(
		func(b bool) error {
			if b {
				return nil
			}
			os.Exit(1)
			return nil
		},
	), nil
}

func (g *Greetings) PostProcess(commit *module.CommitInfo) error {
	// This module does not modify the commit info, so no post-processing is needed.
	return nil
}

func (g *Greetings) GetConfig() module.Config {
	return g.config
}

func (g *Greetings) SetConfig(config module.Config) {
	g.config = config
}

func (g *Greetings) Debug() error {
	// Optionally implement debugging information
	return nil
}

func (g *Greetings) GetName() string {
	return MODULE_NAME
}

func (g *Greetings) getStagedFiles() (string, error) {
	cmd := exec.Command("git", "diff", "--cached", "--name-only")
	var out bytes.Buffer
	cmd.Stdout = &out
	err := cmd.Run()
	if err != nil {
		return "", err
	}
	return out.String(), nil
}

func New() module.Module {
	return &Greetings{config: module.Config{Name: MODULE_NAME}}
}
