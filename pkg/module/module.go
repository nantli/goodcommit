package module

import (
	"github.com/charmbracelet/huh"
	"github.com/nantli/goodcommit/pkg/commit"
)

type Config struct {
	Page         int      `json:"page"`
	Position     int      `json:"position"`
	Name         string   `json:"name"`
	Active       bool     `json:"active"`
	Path         string   `json:"path,omitempty"`
	Priority     int      `json:"priority"`
	Checkpoint   bool     `json:"checkpoint"`
	Pinned       bool     `json:"pinned"`
	Dependencies []string `json:"dependencies"`
}

type Item struct {
	Id          string   `json:"id"`
	Name        string   `json:"name"`
	Title       string   `json:"title"`
	Description string   `json:"description"`
	Emoji       string   `json:"emoji"`
	Conditional []string `json:"conditional"`
}

type Module interface {
	LoadConfig() error
	NewField(commit *commit.Config) (huh.Field, error)
	PostProcess(commit *commit.Config) error
	GetConfig() Config
	GetName() string
	SetConfig(config Config)
	InitCommitInfo(commit *commit.Config) error
	IsActive() bool
}
