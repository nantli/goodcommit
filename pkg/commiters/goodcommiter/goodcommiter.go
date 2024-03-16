// Package goodcommiter provides the default implementation of the Commiter interface
// for the goodcommit application. New commiters can be created to handle different
// commit flows.
package goodcommiter

import (
	"fmt"
	"sort"
	"strings"

	"github.com/charmbracelet/huh"
	"github.com/charmbracelet/lipgloss"
	"github.com/nantli/goodcommit/pkg/module"
)

type GoodCommiter struct {
	modules []module.Module
	commit  module.CommitInfo
}

func (c *GoodCommiter) runForm(accessible bool) error {
	modulesByPage := make(map[int][]*module.Module)
	// Iterate over the modules and add them to the map
	for _, m := range c.modules {
		page := m.GetConfig().Page
		modulesByPage[page] = append(modulesByPage[page], &m)
	}

	var pages []int
	for page := range modulesByPage {
		pages = append(pages, page)
	}
	sort.Ints(pages) // Sort the pages

	var groups []*huh.Group
	for _, page := range pages { // Iterate over sorted pages
		// Sort the modules by position
		sort.Slice(modulesByPage[page], func(i, j int) bool {
			return (*modulesByPage[page][i]).GetConfig().Position < (*modulesByPage[page][j]).GetConfig().Position
		})

		var fields []huh.Field
		// Add the fields from the modules to the group
		for _, m := range modulesByPage[page] {
			field, err := (*m).NewField(&c.commit)
			if err != nil {
				return err
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
					return err
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

	err := form.Run()

	return err
}

func (c *GoodCommiter) runPostProcessing() error {
	for i := 0; i < 100; i++ {
		for _, m := range c.modules {
			if m.GetConfig().Priority > i {
				continue
			}
			if err := m.PostProcess(&c.commit); err != nil {
				return err
			}
		}
	}
	return nil
}

func (c *GoodCommiter) previewCommit() {
	var sb strings.Builder
	keywordStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("#FFD700"))
	fmt.Fprintf(&sb,
		"%s\n\nType: %s\nScope: %s\n\n%s\n\n%s\n\n%s",
		lipgloss.NewStyle().Bold(true).Render("COMMIT SUMMARY 💎"),
		keywordStyle.Render(c.commit.Type),
		keywordStyle.Render(c.commit.Scope),
		keywordStyle.Render(c.commit.Description),
		lipgloss.NewStyle().Italic(true).Render(c.commit.Body),
		keywordStyle.Render(c.commit.Footer),
	)

	if c.commit.Breaking {
		fmt.Fprintf(&sb, "\nBREAKING CHANGE!")
	}

	if len(c.commit.CoAuthoredBy) > 0 {
		var coauthors string
		// build coauthors to gather all entries in CoAuthoredBy
		for _, coauthor := range c.commit.CoAuthoredBy {
			coauthors += fmt.Sprintf("Co-authored-by: %s\n", coauthor)
		}
		fmt.Fprintf(&sb, "\n\n%s", keywordStyle.Render(coauthors))
	}

	fmt.Fprintf(&sb, "\n%s", lipgloss.NewStyle().Bold(true).Render("He's alright, he's a GOODCOMMIT!"))

	fmt.Println(
		lipgloss.NewStyle().
			Width(60).
			BorderStyle(lipgloss.RoundedBorder()).
			BorderForeground(lipgloss.Color("#FFD700")).
			Padding(1, 2).
			Render(sb.String()),
	)
}

func (c *GoodCommiter) stringifyCommit() string {
	return ""
}

func (c *GoodCommiter) commitChanges() error {
	message := c.stringifyCommit()
	print(message)
	return nil
}

func (c *GoodCommiter) Execute(accessible bool) error {
	if err := c.runForm(accessible); err != nil {
		return err
	}
	if err := c.runPostProcessing(); err != nil {
		return err
	}
	c.previewCommit()
	return c.commitChanges()
}

func New(modules []module.Module) *GoodCommiter {
	return &GoodCommiter{modules, module.CommitInfo{}}
}
