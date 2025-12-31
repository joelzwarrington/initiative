package initiative

import (
	"testing"

	"github.com/joelzwarrington/initiative/dnd"
)

func TestGetCharacterOptions(t *testing.T) {
	tests := []struct {
		name       string
		characters *map[string]dnd.Character
		expected   []string // expected character names in sorted order
	}{
		{
			name: "sorts characters by name case-insensitively",
			characters: &map[string]dnd.Character{
				"uuid-1": dnd.NewCharacter("Zara"),
				"uuid-2": dnd.NewCharacter("alice"),
				"uuid-3": dnd.NewCharacter("Bob"),
				"uuid-4": dnd.NewCharacter("charlie"),
			},
			expected: []string{"alice", "Bob", "charlie", "Zara"},
		},
		{
			name: "handles single character",
			characters: &map[string]dnd.Character{
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
			characters: &map[string]dnd.Character{},
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