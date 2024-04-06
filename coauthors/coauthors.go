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

const MODULE_NAME = "coauthors"

type coAuthors struct {
	config gc.ModuleConfig
	Items  []item `json:"coauthors"`
}

func (c *coAuthors) item(id string) item {
	for _, i := range c.Items {
		if i.Id == id {
			return i
		}
	}
	return item{}
}

func (c *coAuthors) LoadConfig() error {
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
func (c *coAuthors) NewField(commit *gc.Commit) (huh.Field, error) {

	// Get the user's email
	emailCmd := exec.Command("git", "config", "--get", "user.email")
	email, err := emailCmd.Output()
	if err != nil {
		return nil, fmt.Errorf("error getting user email: %w", err)
	}
	userEmail := strings.TrimSpace(string(email))

	// Filter out the user's email from the co-authors
	coAuthors := []item{}
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

func (c *coAuthors) PostProcess(commit *gc.Commit) error {
	coAuthors := commit.CoAuthoredBy
	for i, coAuthor := range coAuthors {
		coAuthors[i] = c.item(coAuthor).Name + " <" + c.item(coAuthor).Id + ">"
	}
	commit.CoAuthoredBy = coAuthors
	// Sign the commit body with the co-authors Emojis
	emojis := []string{}
	for _, coAuthor := range coAuthors {
		emojis = append(emojis, c.item(coAuthor).Emoji)
	}

	// add emoji from author using emailCmd := exec.Command("git", "config", "--get", "user.email") to get mail and the serching with id
	emailCmd := exec.Command("git", "config", "--get", "user.email")
	email, err := emailCmd.Output()
	if err != nil {
		return fmt.Errorf("error getting user email: %w", err)
	}
	authorId := strings.TrimSpace(string(email))

	commit.Body = commit.Body + "\n\n" + c.item(authorId).Emoji + " " + strings.Join(emojis, " ")
	return nil
}

func (c *coAuthors) Config() gc.ModuleConfig {
	return c.config
}

func (c *coAuthors) SetConfig(config gc.ModuleConfig) {
	c.config = config
}

func (c *coAuthors) Name() string {
	return MODULE_NAME
}

func (c *coAuthors) InitCommitInfo(commit *gc.Commit) error {
	// No initialization needed for this module.
	return nil
}

func (c *coAuthors) IsActive() bool {
	return c.config.Active
}

func New() gc.Module {
	return &coAuthors{config: gc.ModuleConfig{Name: MODULE_NAME}}
}
