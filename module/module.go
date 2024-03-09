package module

import "github.com/charmbracelet/huh"

type Config struct {
	Page     int    `json:"page"`
	Position int    `json:"position"`
	Name     string `json:"name"`
	Active   bool   `json:"active"`
}

type Module interface {
	NewField(commit *CommitInfo) (huh.Field, error)
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
