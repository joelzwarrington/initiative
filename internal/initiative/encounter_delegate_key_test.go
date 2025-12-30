package initiative

import (
	"strings"
	"testing"

	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/joelzwarrington/initiative/dnd"
)

func TestCreatureItemKeyMap(t *testing.T) {
	keys := newCreatureItemKeyMap()
	
	// Test that damage key is 'd'
	if !key.Matches(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'d'}}, keys.dealDamage) {
		t.Error("Expected 'd' key to match dealDamage binding")
	}
	
	// Test that heal key is 'h'
	if !key.Matches(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'h'}}, keys.heal) {
		t.Error("Expected 'h' key to match heal binding")
	}
}

func TestAdjustHitPointsMessage(t *testing.T) {
	tests := []struct {
		name      string
		key       rune
		isDamage  bool
	}{
		{"Damage key", 'd', true},
		{"Heal key", 'h', false},
	}
	
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create test setup
			monster := dnd.NewMonster("Test Monster", dnd.StatBlock{})
			encounter := &dnd.Encounter{
				InitiativeGroups: []dnd.InitiativeGroup{
					{
						Initiative: 15,
						Creatures:  []dnd.Creature{monster},
					},
				},
			}
			
			delegate := newEncounterDelegate(80, 24)
			var buf strings.Builder
			delegate.Render(&buf, encounter) // Render to populate list
			
			// Get the first item and simulate selection
			items := delegate.list.Items()
			if len(items) == 0 {
				t.Skip("No items to test")
			}
			delegate.list.Select(0)
			
			// Create a mock creatureItemDelegate to test key handling
			itemDelegate := &creatureItemDelegate{
				keys: newCreatureItemKeyMap(),
			}
			
			// Simulate key press
			keyMsg := tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{tt.key}}
			cmd := itemDelegate.Update(keyMsg, &delegate.list)
			
			if cmd == nil {
				t.Error("Expected command to be returned")
				return
			}
			
			msg := cmd()
			if adjustMsg, ok := msg.(adjustHitPointsMsg); ok {
				if adjustMsg.isDamage != tt.isDamage {
					t.Errorf("Expected isDamage %v, got %v", tt.isDamage, adjustMsg.isDamage)
				}
				if adjustMsg.groupIndex != 0 {
					t.Errorf("Expected groupIndex 0, got %d", adjustMsg.groupIndex)
				}
				if adjustMsg.creatureIndex != 0 {
					t.Errorf("Expected creatureIndex 0, got %d", adjustMsg.creatureIndex)
				}
			} else {
				t.Errorf("Expected adjustHitPointsMsg, got %T", msg)
			}
		})
	}
}