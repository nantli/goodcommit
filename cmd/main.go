package main

import (
	"fmt"
	"os"
	"sort"
	"strconv"
	"strings"

	"github.com/charmbracelet/huh"
	"github.com/charmbracelet/lipgloss"
	"github.com/nantli/goodcommit/internal/config"
	"github.com/nantli/goodcommit/pkg/module"
)

func main() {
	// Should we run in accessible mode?
	accessible, _ := strconv.ParseBool(os.Getenv("ACCESSIBLE"))
	var commit module.CommitInfo
	var groups []*huh.Group

	// Load modules from configuration
	modules := config.LoadConfig()

	modulesByPage := make(map[int][]*module.Module)
	// Iterate over the modules and add them to the map
	for _, m := range modules {
		page := m.GetConfig().Page
		modulesByPage[page] = append(modulesByPage[page], &m)
	}

	for page := range modulesByPage {
		// Sort the modules by position
		sort.Slice(modulesByPage[page], func(i, j int) bool {
			return (*modulesByPage[page][i]).GetConfig().Position < (*modulesByPage[page][j]).GetConfig().Position
		})

		var fields []huh.Field
		// Add the fields from the modules to the group
		for _, m := range modulesByPage[page] {
			field, err := (*m).NewField(&commit)
			if err != nil {
				fmt.Println("Uh oh:", err)
				os.Exit(1)
			}
			fields = append(fields, field)
		}

		group := huh.NewGroup(fields...)
		groups = append(groups, group)

		// Check if any module in the page has the checkpoint set to true
		for _, m := range modulesByPage[page] {
			if (*m).GetConfig().Checkpoint {
				// Create the form with the current groups
				form := huh.NewForm(groups...).
					WithTheme(huh.ThemeCharm()).
					WithAccessible(accessible)

				// Run the form and check for errors
				if err := form.Run(); err != nil {
					fmt.Println("Uh oh:", err)
					os.Exit(1)
				}

				// Start a new set of groups for the next pages
				groups = []*huh.Group{}
				break
			}
		}
	}

	// Create and run the form with the remaining groups
	form := huh.NewForm(groups...).
		WithTheme(huh.ThemeCharm()).
		WithAccessible(accessible)

	if err := form.Run(); err != nil {
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
