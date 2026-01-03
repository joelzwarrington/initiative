package initiative

import (
	"strings"
	"testing"

	"github.com/joelzwarrington/initiative/dnd"
	"github.com/termkit/skeleton"
)

func TestCharacterPage_HelpIntegration(t *testing.T) {
	tests := []struct {
		name     string
		setup    func(*characterPage)
		contains []string
	}{
		{
			name: "character form shows form help keys",
			setup: func(p *characterPage) {
				p.characterForm = newCharacterForm("", nil, 80, 24)
			},
			contains: []string{"esc", "cancel"},
		},
		{
			name: "character list shows list help keys",
			setup: func(p *characterPage) {
				p.characters = map[string]*dnd.Character{
					"test": dnd.NewCharacter("Test"),
				}
			},
			contains: []string{"n", "new character", "↑", "↓"},
		},
		{
			name: "empty character list shows help",
			setup: func(p *characterPage) {
				p.characters = map[string]*dnd.Character{}
			},
			contains: []string{"n", "new character"},
		},
		{
			name: "character detail shows back key",
			setup: func(p *characterPage) {
				p.characters = map[string]*dnd.Character{
					"test": dnd.NewCharacter("Test"),
				}
				p.currentCharacter = "test"
			},
			contains: []string{"esc", "back"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := skeleton.NewSkeleton()
			page := newCharacterPage(s, nil)
			page.width = 80
			page.height = 24

			tt.setup(page)

			output := page.View()

			for _, pattern := range tt.contains {
				if !strings.Contains(output, pattern) {
					t.Errorf("Expected help to contain %q, but it didn't.\nFull output:\n%s", pattern, output)
				}
			}
		})
	}
}

func TestCharacterPage_ViewRendersHelpInAllStates(t *testing.T) {
	s := skeleton.NewSkeleton()
	characters := map[string]*dnd.Character{
		"test": dnd.NewCharacter("Test Character"),
	}
	page := newCharacterPage(s, characters)
	page.width = 80
	page.height = 24

	states := []struct {
		name  string
		setup func()
	}{
		{
			name: "viewing list",
			setup: func() {
				page.currentCharacter = ""
				page.characterForm = nil
			},
		},
		{
			name: "viewing character detail",
			setup: func() {
				page.currentCharacter = "test"
				page.characterForm = nil
			},
		},
		{
			name: "editing character",
			setup: func() {
				page.currentCharacter = ""
				page.characterForm = newCharacterForm("test", characters["test"], 80, 24)
			},
		},
		{
			name: "adding new character",
			setup: func() {
				page.currentCharacter = ""
				page.characterForm = newCharacterForm("", nil, 80, 24)
			},
		},
	}

	for _, state := range states {
		t.Run(state.name, func(t *testing.T) {
			state.setup()

			output := page.View()

			// Every state should contain help component (at minimum the ? key)
			if !strings.Contains(output, "?") && !strings.Contains(output, "help") &&
				!strings.Contains(output, "enter") && !strings.Contains(output, "tab") &&
				!strings.Contains(output, "esc") {
				t.Errorf("Expected help component in %s state, but no help keys found.\nFull output:\n%s", state.name, output)
			}
		})
	}
}
