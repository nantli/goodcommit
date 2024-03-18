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
	pinnedModules := make(map[*module.Module][]int)
	var maxPage int // Track the maximum page number

	// First pass: Iterate over the modules to populate modulesByPage and track pinned modules
	for _, m := range c.modules {
		if m.IsActive() {
			page := m.GetConfig().Page
			if page > maxPage {
				maxPage = page // Update maxPage if the current page is higher
			}
			modulesByPage[page] = append(modulesByPage[page], &m)
			if m.GetConfig().Pinned {
				// If the module is pinned, track it for addition to subsequent pages
				for p := page + 1; p <= 30; p++ { // Assuming a max of 20 pages for simplicity
					pinnedModules[&m] = append(pinnedModules[&m], p)
				}
			}
		}
	}

	// Second pass: Add pinned modules to their respective pages, up to maxPage
	for m, pages := range pinnedModules {
		for _, page := range pages {
			if page <= maxPage { // Only add to pages within the maxPage limit
				modulesByPage[page] = append(modulesByPage[page], m)
			}
		}
	}

	var pages []int
	for page := range modulesByPage {
		pages = append(pages, page)
	}
	sort.Ints(pages) // Sort the pages

	var groups []*huh.Group
	pageHasNonPinned := make(map[int]bool) // Track if a page has non-pinned modules

	for _, page := range pages { // Iterate over sorted pages
		// Sort the modules by position
		sort.Slice(modulesByPage[page], func(i, j int) bool {
			mi, mj := *modulesByPage[page][i], *modulesByPage[page][j]
			if mi.GetConfig().Page == mj.GetConfig().Page {
				if mi.GetConfig().Pinned == mj.GetConfig().Pinned {
					return mi.GetConfig().Position < mj.GetConfig().Position
				}
				return mi.GetConfig().Pinned && !mj.GetConfig().Pinned
			}
			return mi.GetConfig().Page < mj.GetConfig().Page
		})

		var fields []huh.Field
		// Add the fields from the modules to the group
		for _, m := range modulesByPage[page] {

			field, err := (*m).NewField(&c.commit)
			if err != nil {
				return err
			}

			if !(*m).GetConfig().Pinned && field != nil && (*m).GetConfig().Active {
				pageHasNonPinned[page] = true // Mark page as having non-pinned modules
			}

			if field != nil {
				fields = append(fields, field)
			}
		}

		// Only create a group if the page has non-pinned modules
		if pageHasNonPinned[page] {
			group := huh.NewGroup(fields...)
			groups = append(groups, group)
		}

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
			if m.GetConfig().Priority != i {
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
	alertStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("#FF0000"))
	footerStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("#00D4F4"))

	// Determine the style to use based on whether the commit type includes an exclamation mark
	var typeStyle lipgloss.Style
	if strings.Contains(c.commit.Type, "!") {
		typeStyle = alertStyle
	} else {
		typeStyle = keywordStyle
	}

	// Use the determined style for the commit type
	fmt.Fprintf(&sb,
		"%s\n\nType: %s\nScope: %s\nDescription: %s\nBody:\n\n%s\n",
		lipgloss.NewStyle().Bold(true).Render("COMMIT SUMMARY 💎"),
		typeStyle.Render(c.commit.Type), // Apply the conditional styling here
		keywordStyle.Render(c.commit.Scope),
		keywordStyle.Render(c.commit.Description),
		lipgloss.NewStyle().Italic(true).Render(c.commit.Body),
	)

	if len(c.commit.CoAuthoredBy) > 0 {
		var coauthors string
		// build coauthors to gather all entries in CoAuthoredBy
		for _, coauthor := range c.commit.CoAuthoredBy {
			coauthors += fmt.Sprintf("\nCo-authored-by: %s", coauthor)
		}
		fmt.Fprintf(&sb, "%s", footerStyle.Render(coauthors))
	}

	if c.commit.Footer != "" {
		fmt.Fprintf(&sb, "%s", footerStyle.Render(c.commit.Footer))
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

func New(modules []module.Module) (*GoodCommiter, error) {
	commit := module.CommitInfo{Extras: make(map[string]*string)}
	// run InitCommitInfo from all modules in priority order
	for i := 0; i < 100; i++ {
		for _, m := range modules {
			if m.GetConfig().Priority > i {
				continue
			}
			if err := m.InitCommitInfo(&commit); err != nil {
				return nil, err
			}
		}
	}

	return &GoodCommiter{modules, commit}, nil
}
