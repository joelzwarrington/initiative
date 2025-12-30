package initiative

import (
	"testing"

	tea "github.com/charmbracelet/bubbletea"
)

func TestHitPointForm(t *testing.T) {
	tests := []struct {
		name          string
		groupIndex    int
		creatureIndex int
		creatureName  string
	}{
		{"Valid indices", 0, 0, "Goblin"},
		{"Different creature", 1, 2, "Orc"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			form := newHitPointForm(tt.groupIndex, tt.creatureIndex, tt.creatureName, 80, 24, true)

			if form.groupIndex != tt.groupIndex {
				t.Errorf("Expected groupIndex %d, got %d", tt.groupIndex, form.groupIndex)
			}

			if form.creatureIndex != tt.creatureIndex {
				t.Errorf("Expected creatureIndex %d, got %d", tt.creatureIndex, form.creatureIndex)
			}

			if form.creatureName != tt.creatureName {
				t.Errorf("Expected creatureName %s, got %s", tt.creatureName, form.creatureName)
			}
		})
	}
}

func TestHitPointAdjustmentType_String(t *testing.T) {
	tests := []struct {
		name     string
		adjType  HitPointAdjustmentType
		expected string
	}{
		{"Damage type", HitPointDamage, "damage"},
		{"Heal type", HitPointHeal, "heal"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.adjType.String(); got != tt.expected {
				t.Errorf("Expected %s, got %s", tt.expected, got)
			}
		})
	}
}

func TestHitPointFormMessages(t *testing.T) {
	form := newHitPointForm(0, 0, "Test Creature", 80, 24, true)

	t.Run("Cancel message", func(t *testing.T) {
		// Simulate ESC key press to abort form
		model, cmd := form.Update(tea.KeyMsg{Type: tea.KeyEsc})

		if cmd == nil {
			t.Error("Expected command to be returned")
		}

		msg := cmd()
		if _, ok := msg.(hitPointFormCancelledMsg); !ok {
			t.Error("Expected hitPointFormCancelledMsg")
		}

		_ = model
	})
}
