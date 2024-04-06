package signedoffby

import (
	"bytes"
	"fmt"
	"os/exec"

	"github.com/charmbracelet/huh"
	gc "github.com/nantli/goodcommit"
)

const MODULE_NAME = "signedoffby"

type signedOffBy struct {
	config gc.ModuleConfig
}

func (s *signedOffBy) LoadConfig() error {
	// Load any necessary configuration. For this module, it might be inactive or active.
	return nil
}

func (s *signedOffBy) NewField(commit *gc.Commit) (huh.Field, error) {
	// This module does not require input from the user.
	return nil, nil
}

func (s *signedOffBy) PostProcess(commit *gc.Commit) error {
	// Execute the command to get the user's name from Git config
	nameCmd := exec.Command("git", "config", "--get", "user.name")
	var nameOut bytes.Buffer
	nameCmd.Stdout = &nameOut
	if err := nameCmd.Run(); err != nil {
		return fmt.Errorf("failed to get git user name: %w", err)
	}
	authorName := nameOut.String()
	authorName = authorName[:len(authorName)-1] // Remove the newline at the end

	// Execute the command to get the user's email from Git config
	emailCmd := exec.Command("git", "config", "--get", "user.email")
	var emailOut bytes.Buffer
	emailCmd.Stdout = &emailOut
	if err := emailCmd.Run(); err != nil {
		return fmt.Errorf("failed to get git user email: %w", err)
	}
	authorEmail := emailOut.String()
	authorEmail = authorEmail[:len(authorEmail)-1] // Remove the newline at the end

	// Append "Signed-off-by" to the commit footer with the gathered info
	commit.Footer += fmt.Sprintf("\nSigned-off-by: %s <%s>", authorName, authorEmail)
	return nil
}

func (s *signedOffBy) Config() gc.ModuleConfig {
	return s.config
}

func (s *signedOffBy) SetConfig(config gc.ModuleConfig) {
	s.config = config
}

func (s *signedOffBy) Name() string {
	return MODULE_NAME
}

func (s *signedOffBy) IsActive() bool {
	return s.config.Active
}

func (s *signedOffBy) InitCommitInfo(commit *gc.Commit) error {
	// Initialize any necessary fields in CommitInfo.
	return nil
}

func New() gc.Module {
	return &signedOffBy{}
}
