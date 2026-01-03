package initiative

import (
	"strings"
	"testing"
	"time"

	"github.com/joelzwarrington/initiative/dnd"
	"github.com/termkit/skeleton"
)

func TestEncounterPage_HelpIntegration(t *testing.T) {
	tests := []struct {
		name     string
		setup    func(*encounterPage)
		contains []string
	}{
		{
			name: "empty state shows help with new encounter key",
			setup: func(p *encounterPage) {
				// Default state is empty
			},
			contains: []string{"n", "new", "?", "help"},
		},
		{
			name: "encounter form shows form help keys",
			setup: func(p *encounterPage) {
				p.encounterForm = newEncounterForm(nil, nil, 80, 24)
			},
			contains: []string{"esc", "cancel"},
		},
		{
			name: "hit point form shows form help keys",
			setup: func(p *encounterPage) {
				p.hitPointForm = newHitPointForm(0, 0, "Test Creature", 80, 24, true)
			},
			contains: []string{"esc", "cancel"},
		},
		{
			name: "cancellation form shows form help keys",
			setup: func(p *encounterPage) {
				p.cancellationForm = newCancellationForm(80, 24)
			},
			contains: []string{"enter", "y", "n"},
		},
		{
			name: "encounter view shows encounter help keys",
			setup: func(p *encounterPage) {
				// Create encounter with a creature so we get full help keys
				char := dnd.NewCharacter("Test Hero")
				creature := dnd.Creature(char)
				groups := []*dnd.InitiativeGroup{
					dnd.NewInitiativeGroup(15, []*dnd.Creature{&creature}),
				}
				p.encounter = dnd.NewEncounter("Test", time.Now(), groups)
			},
			// ShortHelp shows: up, down, next turn (prev disabled at start), damage, heal, help
			// Note: prev turn is disabled at round 1 turn 0, and end key is only in FullHelp
			contains: []string{"→", "next turn", "?", "help"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := skeleton.NewSkeleton()
			page := newEncounterPage(s, nil, nil)
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

func TestEncounterPage_HelpKeysStateSpecific(t *testing.T) {
	s := skeleton.NewSkeleton()
	page := newEncounterPage(s, nil, nil)
	page.width = 80
	page.height = 24

	// Test empty state help keys
	shortHelp := page.ShortHelp()
	if len(shortHelp) == 0 {
		t.Error("Expected help keys in empty state")
	}

	// Test encounter form help keys
	page.encounterForm = newEncounterForm(nil, nil, 80, 24)
	formShortHelp := page.ShortHelp()
	if len(formShortHelp) == 0 {
		t.Error("Expected form help keys when in form state")
	}

	// Verify form keys are different from page keys
	if len(shortHelp) == len(formShortHelp) {
		same := true
		for i, key := range shortHelp {
			if i < len(formShortHelp) && key.Help().Key != formShortHelp[i].Help().Key {
				same = false
				break
			}
		}
		if same {
			t.Error("Expected different help keys when in form vs empty state")
		}
	}
}
