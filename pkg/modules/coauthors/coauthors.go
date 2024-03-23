// Package coauthors provides a module for goodcommit that allows the user to select
// co-authors for the commit from a predefined list.
package coauthors

import (
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/charmbracelet/huh"
	"github.com/nantli/goodcommit/pkg/commit"
	"github.com/nantli/goodcommit/pkg/module"
)

const MODULE_NAME = "coauthors"

type CoAuthors struct {
	config module.Config
	Items  []module.Item `json:"coauthors"`
}

func (c *CoAuthors) getItem(id string) module.Item {
	for _, i := range c.Items {
		if i.Id == id {
			return i
		}
	}
	return module.Item{}
}

func (c *CoAuthors) LoadConfig() error {
	if c.config.Path == "" {
		return nil
	}

	raw, err := os.ReadFile(c.config.Path)
	if err != nil {
		return fmt.Errorf("error reading config: %w", err)
	}
	err = json.Unmarshal(raw, c)
	if err != nil {
		return fmt.Errorf("error parsing config: %w", err)
	}

	return nil
}

// NewField returns a MultiSelect field with options for each co-author.
func (c *CoAuthors) NewField(commit *commit.Config) (huh.Field, error) {

	// Get the user's email
	emailCmd := exec.Command("git", "config", "--get", "user.email")
	email, err := emailCmd.Output()
	if err != nil {
		return nil, fmt.Errorf("error getting user email: %w", err)
	}
	userEmail := strings.TrimSpace(string(email))

	// Filter out the user's email from the co-authors
	coAuthors := []module.Item{}
	for _, item := range c.Items {
		if item.Id != userEmail {
			coAuthors = append(coAuthors, item)
		}
	}

	if len(coAuthors) == 0 {
		return nil, nil
	}

	var coAuthorOptions []huh.Option[string]
	for _, item := range coAuthors {
		coAuthorOptions = append(coAuthorOptions, huh.NewOption(item.Name+" - "+item.Id, item.Id))
	}

	return huh.NewMultiSelect[string]().
		Title("👥・Select Co-Authors").
		Description("Choose co-authors for this commit.").
		Options(coAuthorOptions...).
		Value(&commit.CoAuthoredBy), nil
}

func (c *CoAuthors) PostProcess(commit *commit.Config) error {
	coAuthors := commit.CoAuthoredBy
	for i, coAuthor := range coAuthors {
		coAuthors[i] = c.getItem(coAuthor).Name + " <" + c.getItem(coAuthor).Id + ">"
	}
	commit.CoAuthoredBy = coAuthors
	// Sign the commit body with the co-authors Emojis
	emojis := []string{}
	for _, coAuthor := range coAuthors {
		emojis = append(emojis, c.getItem(coAuthor).Emoji)
	}

	// add emoji from author using emailCmd := exec.Command("git", "config", "--get", "user.email") to get mail and the serching with id
	emailCmd := exec.Command("git", "config", "--get", "user.email")
	email, err := emailCmd.Output()
	if err != nil {
		return fmt.Errorf("error getting user email: %w", err)
	}
	authorId := strings.TrimSpace(string(email))

	commit.Body = commit.Body + "\n\n" + c.getItem(authorId).Emoji + " " + strings.Join(emojis, " ")
	return nil
}

func (c *CoAuthors) GetConfig() module.Config {
	return c.config
}

func (c *CoAuthors) SetConfig(config module.Config) {
	c.config = config
}

func (c *CoAuthors) GetName() string {
	return MODULE_NAME
}

func (c *CoAuthors) InitCommitInfo(commit *commit.Config) error {
	// No initialization needed for this module.
	return nil
}

func (c *CoAuthors) IsActive() bool {
	return c.config.Active
}

func New() module.Module {
	return &CoAuthors{config: module.Config{Name: MODULE_NAME}}
}
