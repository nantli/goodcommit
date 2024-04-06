package goodcommit

import "github.com/charmbracelet/huh"

type Commit struct {
	Type         string
	Scope        string
	Description  string
	Body         string
	Footer       string
	Breaking     bool
	CoAuthoredBy []string
	Extras       map[string]*string
}

type ModuleConfig struct {
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

type Module interface {
	LoadConfig() error
	NewField(commit *Commit) (huh.Field, error)
	PostProcess(commit *Commit) error
	Config() ModuleConfig
	Name() string
	SetConfig(config ModuleConfig)
	InitCommitInfo(commit *Commit) error
	IsActive() bool
}
