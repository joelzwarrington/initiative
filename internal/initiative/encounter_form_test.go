package initiative

import (
	"testing"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/joelzwarrington/initiative/dnd"
)

func TestGetCharacterOptions(t *testing.T) {
	tests := []struct {
		name       string
		characters map[string]*dnd.Character
		expected   []string // expected character names in sorted order
	}{
		{
			name: "sorts characters by name case-insensitively",
			characters: map[string]*dnd.Character{
				"uuid-1": dnd.NewCharacter("Zara"),
				"uuid-2": dnd.NewCharacter("alice"),
				"uuid-3": dnd.NewCharacter("Bob"),
				"uuid-4": dnd.NewCharacter("charlie"),
			},
			expected: []string{"alice", "Bob", "charlie", "Zara"},
		},
		{
			name: "handles single character",
			characters: map[string]*dnd.Character{
				"solo": dnd.NewCharacter("Hero"),
			},
			expected: []string{"Hero"},
		},
		{
			name:       "handles nil characters map",
			characters: nil,
			expected:   []string{},
		},
		{
			name:       "handles empty characters map",
			characters: map[string]*dnd.Character{},
			expected:   []string{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			options := getCharacterOptions(tt.characters)

			if len(options) != len(tt.expected) {
				t.Fatalf("expected %d options, got %d", len(tt.expected), len(options))
			}

			for i, option := range options {
				if option.Key != tt.expected[i] {
					t.Errorf("option %d: expected name %q, got %q", i, tt.expected[i], option.Key)
				}
			}
		})
	}
}

func TestEncounterFormFilteringKeyBind(t *testing.T) {
	characters := map[string]*dnd.Character{
		"1": dnd.NewCharacter("Alice"),
		"2": dnd.NewCharacter("Bob"),
	}

	tests := []struct {
		name          string
		setupKeys     []tea.KeyMsg
		expectEscHelp bool
	}{
		{
			name:          "initial state shows esc help",
			setupKeys:     []tea.KeyMsg{},
			expectEscHelp: true,
		},
		{
			name: "navigate to characters field and start filtering",
			setupKeys: []tea.KeyMsg{
				{Type: tea.KeyRunes, Runes: []rune("/")}, // Start filtering on characters field
			},
			expectEscHelp: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			form := newEncounterForm(characters, nil, 40, 10)

			// Navigate to the characters field using NextGroup()
			currentForm := form.getCurrentForm()
			if tt.name == "navigate to characters field and start filtering" {
				currentForm.NextGroup()
			}

			// Apply setup keys
			for _, keyMsg := range tt.setupKeys {
				updatedModel, _ := form.Update(keyMsg)
				form = updatedModel.(*encounterForm)
			}

			// Check if esc help is shown
			helpKeys := form.getHelpKeys()
			hasEscHelp := false
			for _, binding := range helpKeys {
				if binding.Help().Key == "esc" && binding.Help().Desc == "cancel" {
					hasEscHelp = true
					break
				}
			}

			if hasEscHelp != tt.expectEscHelp {
				t.Errorf("expected esc help: %v, got: %v", tt.expectEscHelp, hasEscHelp)
			}
		})
	}
}
