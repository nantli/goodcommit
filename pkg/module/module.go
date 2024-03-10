package module

import "github.com/charmbracelet/huh"

type Config struct {
	Page     int    `json:"page"`
	Position int    `json:"position"`
	Name     string `json:"name"`
	Active   bool   `json:"active"`
	Path     string `json:"path,omitempty"`
	Priority int    `json:"priority"`
}

type Module interface {
	Load() error
	NewField(commit *CommitInfo) (huh.Field, error)
	PostProcess(commit *CommitInfo) error
	GetConfig() Config
}

type CommitInfo struct {
	Type         string
	Scope        string
	Description  string
	Body         string
	Footer       string
	Breaking     bool
	CoAuthoredBy string
}
