package main

import (
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/charmbracelet/huh"
	"github.com/charmbracelet/lipgloss"
	"github.com/nantli/goodcommit/internal/config"
	"github.com/nantli/goodcommit/pkg/module"
)

func main() {
	var commit module.CommitInfo
	var typeGroup *huh.Group
	var groups []*huh.Group

	// Load modules from configuration
	modules := config.LoadConfig()

	// Build type huh group, to be the first form to run
	for _, m := range modules {
		if m.GetConfig().Name != "types" {
			continue
		}
		field, err := m.NewField(&commit)
		if err != nil {
			fmt.Println("Uh oh:", err)
			os.Exit(1)
		}

		typeGroup = huh.NewGroup(field)
	}

	// Should we run in accessible mode?
	accessible, _ := strconv.ParseBool(os.Getenv("ACCESSIBLE"))

	preTypeForm := huh.NewForm(
		typeGroup,
	).
		WithTheme(huh.ThemeCharm()).
		WithAccessible(accessible)

	err := preTypeForm.Run()

	if err != nil {
		fmt.Println("Uh oh:", err)
		os.Exit(1)
	}

	// Build post type huh groups
	for _, m := range modules {
		if m.GetConfig().Name == "types" {
			continue
		}
		field, err := m.NewField(&commit)
		if err != nil {
			fmt.Println("Uh oh:", err)
			os.Exit(1)
		}

		groups = append(groups, huh.NewGroup(field))
	}

	postTypeForm := huh.NewForm(
		groups...,
	).
		WithTheme(huh.ThemeCharm()).
		WithAccessible(accessible)

	err = postTypeForm.Run()

	if err != nil {
		fmt.Println("Uh oh:", err)
		os.Exit(1)
	}

	// Post process commit from lower to higher priority
	for i := 0; i < 100; i++ {
		for _, m := range modules {
			if m.GetConfig().Priority > i {
				continue
			}
			if err := m.PostProcess(&commit); err != nil {
				fmt.Println("Uh oh:", err)
				os.Exit(1)
			}
		}
	}

	// Print commit summary
	{
		var sb strings.Builder
		keywordStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("#FFD700"))
		fmt.Fprintf(&sb,
			"%s\n\nType: %s\nScope: %s\n\n%s\n\n%s\n\n%s",
			lipgloss.NewStyle().Bold(true).Render("COMMIT SUMMARY"),
			keywordStyle.Render(commit.Type),
			keywordStyle.Render(commit.Scope),
			keywordStyle.Render(commit.Description),
			keywordStyle.Render(commit.Body),
			keywordStyle.Render(commit.Footer),
		)

		if commit.Breaking {
			fmt.Fprintf(&sb, "\n\nBREAKING CHANGE!")
		}

		if commit.CoAuthoredBy != "" {
			fmt.Fprintf(&sb, "\n\nCo-authored-by: %s", keywordStyle.Render(commit.CoAuthoredBy))
		}

		fmt.Fprintf(&sb, "\n\n%s", lipgloss.NewStyle().Bold(true).Render("He's alright, he's a GOODCOMMIT!"))

		fmt.Println(
			lipgloss.NewStyle().
				Width(60).
				BorderStyle(lipgloss.RoundedBorder()).
				BorderForeground(lipgloss.Color("#FFD700")).
				Padding(1, 2).
				Render(sb.String()),
		)
	}
}
