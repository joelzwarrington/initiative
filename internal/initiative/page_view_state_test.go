package initiative

import (
	"strings"
	"testing"

	"github.com/joelzwarrington/initiative/dnd"
	"github.com/termkit/skeleton"
)

// Test that View() method correctly handles all states and always includes help
func TestPage_ViewAlwaysIncludesHelp(t *testing.T) {
	tests := []struct {
		name     string
		pageType string
		setup    func() (interface{ View() string }, []string)
	}{
		{
			name:     "encounter page empty state",
			pageType: "encounter",
			setup: func() (interface{ View() string }, []string) {
				s := skeleton.NewSkeleton()
				page := newEncounterPage(s, nil, nil)
				page.width = 80
				page.height = 24
				return page, []string{"?", "help"}
			},
		},
		{
			name:     "encounter page form state",
			pageType: "encounter",
			setup: func() (interface{ View() string }, []string) {
				s := skeleton.NewSkeleton()
				page := newEncounterPage(s, nil, nil)
				page.width = 80
				page.height = 24
				page.encounterForm = newEncounterForm(nil, nil, 80, 24)
				return page, []string{"esc", "cancel"}
			},
		},
		{
			name:     "character page list state",
			pageType: "character",
			setup: func() (interface{ View() string }, []string) {
				s := skeleton.NewSkeleton()
				characters := map[string]*dnd.Character{
					"test": dnd.NewCharacter("Test"),
				}
				page := newCharacterPage(s, characters)
				page.width = 80
				page.height = 24
				return page, []string{"n", "new character"}
			},
		},
		{
			name:     "character page form state",
			pageType: "character",
			setup: func() (interface{ View() string }, []string) {
				s := skeleton.NewSkeleton()
				page := newCharacterPage(s, nil)
				page.width = 80
				page.height = 24
				page.characterForm = newCharacterForm("", nil, 80, 24)
				return page, []string{"esc", "cancel"}
			},
		},
		{
			name:     "sources page list state",
			pageType: "sources",
			setup: func() (interface{ View() string }, []string) {
				s := skeleton.NewSkeleton()
				sources := map[string]*dnd.Source{
					"test": {Meta: dnd.SourceMeta{Key: "test", Name: "Test Source", Version: "1.0"}},
				}
				page := newSourcesPage(s, sources)
				page.width = 80
				page.height = 24
				return page, []string{"↑", "↓", "enter"}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			viewer, expectedHelp := tt.setup()

			output := viewer.View()

			if output == "" {
				t.Errorf("%s should return non-empty view", tt.pageType)
			}

			// Verify help component is present
			hasHelp := false
			for _, helpPattern := range expectedHelp {
				if strings.Contains(output, helpPattern) {
					hasHelp = true
					break
				}
			}

			if !hasHelp {
				t.Errorf("%s view should contain help component with one of: %v\nFull output:\n%s",
					tt.pageType, expectedHelp, output)
			}
		})
	}
}

func TestEncounterPage_ViewStateTransitions(t *testing.T) {
	s := skeleton.NewSkeleton()
	page := newEncounterPage(s, nil, nil)
	page.width = 80
	page.height = 24

	// Test empty state
	output1 := page.View()
	if !strings.Contains(output1, "No encounter") {
		t.Error("Expected empty state to show 'No encounter' message")
	}

	// Transition to form state
	page.encounterForm = newEncounterForm(nil, nil, 80, 24)
	output2 := page.View()
	if strings.Contains(output2, "No encounter") {
		t.Error("Form state should not show 'No encounter' message")
	}

	// Both outputs should have help (different help)
	if !strings.Contains(output1, "n") || !strings.Contains(output1, "new") {
		t.Error("Empty state should show 'new' key in help")
	}
	if !strings.Contains(output2, "esc") {
		t.Error("Form state should show 'esc' key in help")
	}
}

func TestCharacterPage_ViewStateTransitions(t *testing.T) {
	s := skeleton.NewSkeleton()
	characters := map[string]*dnd.Character{
		"test": dnd.NewCharacter("Test Character"),
	}
	page := newCharacterPage(s, characters)
	page.width = 80
	page.height = 24

	// Test list state
	output1 := page.View()
	if !strings.Contains(output1, "Test Character") {
		t.Error("List state should show character name")
	}

	// Transition to detail state
	page.currentCharacter = "test"
	output2 := page.View()
	if !strings.Contains(output2, "Viewing character: Test Character") {
		t.Error("Detail state should show 'Viewing character' message")
	}

	// Transition to form state
	page.currentCharacter = ""
	page.characterForm = newCharacterForm("", nil, 80, 24)
	output3 := page.View()
	if strings.Contains(output3, "Test Character") {
		t.Error("Form state should not show character list")
	}

	// All states should have help
	helpStates := []string{output1, output2, output3}
	for i, output := range helpStates {
		hasHelp := strings.Contains(output, "↑") || strings.Contains(output, "esc") ||
			strings.Contains(output, "n") || strings.Contains(output, "?")
		if !hasHelp {
			t.Errorf("State %d should have help component", i+1)
		}
	}
}
