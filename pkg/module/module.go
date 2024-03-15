package module

import "github.com/charmbracelet/huh"

type Config struct {
	Page       int    `json:"page"`
	Position   int    `json:"position"`
	Name       string `json:"name"`
	Active     bool   `json:"active"`
	Path       string `json:"path,omitempty"`
	Priority   int    `json:"priority"`
	Checkpoint bool   `json:"checkpoint"`
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
	NewField(commit *CommitInfo) (huh.Field, error)
	PostProcess(commit *CommitInfo) error
	GetConfig() Config
	GetName() string
	SetConfig(config Config)
	Debug() error
}

type CommitInfo struct {
	Type         string
	Scope        string
	Description  string
	Body         string
	Footer       string
	Breaking     bool
	CoAuthoredBy []string
}
