package signedoffby

import (
	"bytes"
	"fmt"
	"os/exec"

	"github.com/charmbracelet/huh"
	"github.com/nantli/goodcommit/pkg/commit"
	"github.com/nantli/goodcommit/pkg/module"
)

const MODULE_NAME = "signedoffby"

type SignedOffBy struct {
	config module.Config
}

func (s *SignedOffBy) LoadConfig() error {
	// Load any necessary configuration. For this module, it might be inactive or active.
	return nil
}

func (s *SignedOffBy) NewField(commit *commit.Config) (huh.Field, error) {
	// This module does not require input from the user.
	return nil, nil
}

func (s *SignedOffBy) PostProcess(commit *commit.Config) error {
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

func (s *SignedOffBy) GetConfig() module.Config {
	return s.config
}

func (s *SignedOffBy) SetConfig(config module.Config) {
	s.config = config
}

func (s *SignedOffBy) GetName() string {
	return MODULE_NAME
}

func (s *SignedOffBy) IsActive() bool {
	return s.config.Active
}

func (s *SignedOffBy) InitCommitInfo(commit *commit.Config) error {
	// Initialize any necessary fields in CommitInfo.
	return nil
}

func New() module.Module {
	return &SignedOffBy{}
}
