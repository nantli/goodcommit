package main

import (
	"errors"
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/charmbracelet/huh"
	"github.com/charmbracelet/lipgloss"
)

type CommitInfo struct {
	Type         string
	Scope        string
	Description  string
	Body         string
	Footer       string
	Breaking     bool
	CoAuthoredBy string
}

func main() {
	var commit CommitInfo

	commitTypes := []string{"feat", "fix", "chore"}
	commitScopes := []string{"core", "ui", "docs"}
	contributors := []string{"Diego", "Arturo", "Marcos"}

	// Should we run in accessible mode?
	accessible, _ := strconv.ParseBool(os.Getenv("ACCESSIBLE"))

	form := huh.NewForm(
		huh.NewGroup(huh.NewNote().
			Title("GOODCOMMIT").
			Description("Helping you craft the perfect commit message.")),

		// Commit type and scope group.
		huh.NewGroup(
			// Select commit type.
			huh.NewSelect[string]().
				Options(huh.NewOptions(commitTypes...)...).
				Title("Commit type").
				Description("Select the type of change that you're committing.").
				Value(&commit.Type),

			// Select commit scope.
			huh.NewSelect[string]().
				Options(huh.NewOptions(commitScopes...)...).
				Title("Commit scope").
				Description("Select the scope of the change.").
				Value(&commit.Scope),
		),

		// Commit description and body group.
		huh.NewGroup(
			// Input commit description.
			huh.NewInput().
				Value(&commit.Description).
				Title("Short description").
				Placeholder("Add a short description of the change.").
				Validate(func(s string) error {
					if len(s) == 0 {
						return errors.New("description is required")
					}
					return nil
				}),

			// Input commit body.
			huh.NewText().
				Value(&commit.Body).
				Title("Commit body").
				Description("Provide a more detailed description of the change.").
				Lines(5),
		),

		// Commit Footer group.
		huh.NewGroup(
			// Input commit footer.
			huh.NewText().
				Value(&commit.Footer).
				Title("Commit footer").
				Description("Add any references to issues or note breaking changes.").
				Lines(3),

			// Confirm if the change is breaking.
			huh.NewConfirm().
				Title("Is this a BREAKING CHANGE?").
				Value(&commit.Breaking).
				Affirmative("Yes").
				Negative("No"),

			// Select co-author.
			huh.NewSelect[string]().
				Options(huh.NewOptions(contributors...)...).
				Title("Co-authored by").
				Description("Select a contributor if this commit was co-authored.").
				Value(&commit.CoAuthoredBy),
		),
	).WithAccessible(accessible)

	err := form.Run()

	if err != nil {
		fmt.Println("Uh oh:", err)
		os.Exit(1)
	}

	// Print commit summary.
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
