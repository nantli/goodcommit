// Package coauthors provides a github.com/nantli/goodcommit module for selecting co-authors.
// It presents the user with a multi-select field for selecting co-authors from a predefined list.
// The selected co-authors are then added to the commit body.
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

// item is the structure for each entry in the co-authors configuration file.
type item struct {
	Id    string `json:"id"`
	Name  string `json:"name"`
	Emoji string `json:"emoji"`
}

// MODULE_NAME is the name of the module and should be used as the name of the module in the config.json file.
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

// LoadConfig loads the co-authors configuration file.
// Example config file:
//
//	{
//	    "coauthors": [
//	        {
//	            "id": "nantli@nantli.dev",
//	            "name": "Nantli",
//	            "emoji": "🤓"
//	        }
//	    ]
//	}
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

// NewField returns a huh.MultiSelect field with options for each co-author.
// The commit author is excluded from the list of co-authors.
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
	// Build the co-authors string
	coAuthors := commit.CoAuthoredBy
	for i, coAuthor := range coAuthors {
		coAuthors[i] = c.item(coAuthor).Name + " <" + c.item(coAuthor).Id + ">"
	}
	commit.CoAuthoredBy = coAuthors

	// Sign the commit body with the author and co-authors Emojis
	emojis := []string{}
	for _, coAuthor := range coAuthors {
		emojis = append(emojis, c.item(coAuthor).Emoji)
	}

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
	// No initialization of the commit is done by this module
	return nil
}

func (c *coAuthors) IsActive() bool {
	return c.config.Active
}

// New returns a new instance of the co-authors module.
// The coauthors module is a github.com/nantli/goodcommit module that allows the user to select co-authors for the commit.
func New() gc.Module {
	return &coAuthors{config: gc.ModuleConfig{Name: MODULE_NAME}}
}
