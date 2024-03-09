package main

import (
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/charmbracelet/huh"
	"github.com/charmbracelet/lipgloss"
	"github.com/nantli/goodcommit/config"
	"github.com/nantli/goodcommit/module"
)

func main() {
	var commit module.CommitInfo

	// Load modules from configuration
	modules := config.LoadConfig()
	var groups []*huh.Group

	for _, m := range modules {
		field, err := m.NewField(&commit)
		if err != nil {
			fmt.Println("Uh oh:", err)
			os.Exit(1)
		}

		groups = append(groups, huh.NewGroup(field))
	}

	// Should we run in accessible mode?
	accessible, _ := strconv.ParseBool(os.Getenv("ACCESSIBLE"))

	form := huh.NewForm(
		groups...,
	).WithAccessible(accessible)

	err := form.Run()

	if err != nil {
		fmt.Println("Uh oh:", err)
		os.Exit(1)
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
