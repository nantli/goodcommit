// Package greetings provides a module for goodcommit that handles initial greetings.
//
// The greetings module is intended to be the first module in goodcommit. It displays a greeting
// message and displays staged files to the user, asking for confirmation to proceed with the commit.
package greetings

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"

	"github.com/charmbracelet/huh"
	gc "github.com/nantli/goodcommit"
)

const MODULE_NAME = "greetings"

type greetings struct {
	config gc.ModuleConfig
}

func (g *greetings) LoadConfig() error {
	return nil
}

func (g *greetings) NewField(commit *gc.Commit) (huh.Field, error) {
	stagedFiles, err := g.getStagedFiles()
	if err != nil {
		return nil, fmt.Errorf("error getting staged files: %w", err)
	}

	if len(stagedFiles) == 0 {
		return nil, fmt.Errorf("no staged files found")
	}

	// Display staged files as a simple text field in the form
	stagedFilesText := fmt.Sprintf("\nStaged Files:\n%s", stagedFiles)
	return huh.NewConfirm().Title("🐝・Do you want to commit these files?").Description(stagedFilesText).Validate(
		func(b bool) error {
			if b {
				return nil
			}
			os.Exit(1)
			return nil
		},
	), nil
}

func (g *greetings) PostProcess(commit *gc.Commit) error {
	// This module does not modify the commit info, so no post-processing is needed.
	return nil
}

func (g *greetings) Config() gc.ModuleConfig {
	return g.config
}

func (g *greetings) SetConfig(config gc.ModuleConfig) {
	g.config = config
}

func (g *greetings) Name() string {
	return MODULE_NAME
}

func (g *greetings) getStagedFiles() (string, error) {
	cmd := exec.Command("git", "diff", "--cached", "--name-only")
	var out bytes.Buffer
	cmd.Stdout = &out
	err := cmd.Run()
	if err != nil {
		return "", err
	}
	return out.String(), nil
}

func (g *greetings) InitCommitInfo(commit *gc.Commit) error {
	return nil
}

func (g *greetings) IsActive() bool {
	return g.config.Active
}

func New() gc.Module {
	return &greetings{config: gc.ModuleConfig{Name: MODULE_NAME}}
}
