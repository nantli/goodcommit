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
	gc "github.com/nantli/goodcommit"
)

type goodCommiter struct {
	modules []gc.Module
	commit  gc.Commit
}

func (c *goodCommiter) RunForm(accessible bool) error {
	modulesByPage := make(map[int][]*gc.Module)
	pinnedModules := make(map[*gc.Module][]int)
	var maxPage int // Track the maximum page number

	// First pass: Iterate over the modules to populate modulesByPage and track pinned modules
	for _, m := range c.modules {
		if m.IsActive() {
			page := m.Config().Page
			if page > maxPage {
				maxPage = page // Update maxPage if the current page is higher
			}
			modulesByPage[page] = append(modulesByPage[page], &m)
			if m.Config().Pinned {
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
			if mi.Config().Page == mj.Config().Page {
				if mi.Config().Pinned == mj.Config().Pinned {
					return mi.Config().Position < mj.Config().Position
				}
				return mi.Config().Pinned && !mj.Config().Pinned
			}
			return mi.Config().Page < mj.Config().Page
		})

		var fields []huh.Field
		// Add the fields from the modules to the group
		for _, m := range modulesByPage[page] {

			field, err := (*m).NewField(&c.commit)
			if err != nil {
				return err
			}

			if !(*m).Config().Pinned && field != nil && (*m).Config().Active {
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
			if (*m).Config().Checkpoint {
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

func (c *goodCommiter) RunPostProcessing() error {
	for i := 0; i < 100; i++ {
		for _, m := range c.modules {
			if m.Config().Priority != i {
				continue
			}
			if err := m.PostProcess(&c.commit); err != nil {
				return err
			}
		}
	}
	return nil
}

func (c *goodCommiter) PreviewCommit() {
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

	breakingChangeIndicator := ""
	if c.commit.Breaking {
		breakingChangeIndicator = "!"
	}

	// Use the determined style for the commit type
	fmt.Fprintf(&sb,
		"%s\n\nType: %s%s\nScope: %s\nDescription: %s\nBody:\n\n%s\n",
		lipgloss.NewStyle().Bold(true).Render("COMMIT SUMMARY 💎"),
		typeStyle.Render(c.commit.Type), // Apply the conditional styling here
		breakingChangeIndicator,
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

func (c *goodCommiter) RenderMessage() string {
	var scopeStr string
	if c.commit.Scope != "" {
		scopeStr = fmt.Sprintf("(%s)", c.commit.Scope)
	}

	breakingChangeIndicator := ""
	if c.commit.Breaking {
		breakingChangeIndicator = "!"
	}

	commitMsg := fmt.Sprintf("%s%s%s: %s\n\n%s\n", c.commit.Type, scopeStr, breakingChangeIndicator, c.commit.Description, c.commit.Body)

	if len(c.commit.CoAuthoredBy) > 0 {
		var coauthors string
		// build coauthors to gather all entries in CoAuthoredBy
		for _, coauthor := range c.commit.CoAuthoredBy {
			coauthors += fmt.Sprintf("\nCo-authored-by: %s", coauthor)
		}
		commitMsg += coauthors
	}

	if c.commit.Footer != "" {
		commitMsg += c.commit.Footer
	}

	return commitMsg
}

func (c *goodCommiter) LoadModules(modules []gc.Module) error {
	// run InitCommitInfo from all modules in priority order
	for i := 0; i < 100; i++ {
		for _, m := range modules {
			if m.Config().Priority > i {
				continue
			}
			if err := m.InitCommitInfo(&c.commit); err != nil {
				return err
			}
		}
	}
	c.modules = modules
	return nil
}

func New() (*goodCommiter, error) {
	commit := gc.Commit{Extras: make(map[string]*string)}

	return &goodCommiter{modules: []gc.Module{}, commit: commit}, nil
}
