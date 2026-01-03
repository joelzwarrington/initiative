package dnd

import "testing"

func TestNewCharacter(t *testing.T) {
	name := "Frodo"
	char := NewCharacter(name)

	if char.Name() != name {
		t.Errorf("NewCharacter().Name() = %v, want %v", char.Name(), name)
	}
}

func TestCharacter_WithName(t *testing.T) {
	char := NewCharacter("Frodo")
	newChar := char.WithName("Gandalf")

	if newChar.Name() != "Gandalf" {
		t.Errorf("WithName().Name() = %v, want %v", newChar.Name(), "Gandalf")
	}
	if char.Name() != "Frodo" {
		t.Errorf("WithName() should not modify original, got %v", char.Name())
	}
}

func TestCharacter_Health(t *testing.T) {
	tests := []struct {
		name    string
		char    *Character
		wantHP  int
		wantMax int
	}{
		{
			name:    "new character with no health",
			char:    NewCharacter("Test"),
			wantHP:  0,
			wantMax: 0,
		},
		{
			name:    "new character with health",
			char:    NewCharacter("Test").WithHealth(20, 20),
			wantHP:  20,
			wantMax: 20,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.char.HP(); got != tt.wantHP {
				t.Errorf("Character.HP() = %v, want %v", got, tt.wantHP)
			}
			if got := tt.char.MaxHP(); got != tt.wantMax {
				t.Errorf("Character.MaxHP() = %v, want %v", got, tt.wantMax)
			}
		})
	}
}

func TestCharacter_SetHP(t *testing.T) {
	tests := []struct {
		name     string
		char     *Character
		newHP    int
		wantHP   int
		wantName string
	}{
		{
			name:     "set normal hit points",
			char:     NewCharacter("Test").WithHealth(20, 20),
			newHP:    15,
			wantHP:   15,
			wantName: "Test",
		},
		{
			name:     "set hit points to zero",
			char:     NewCharacter("Test").WithHealth(20, 20),
			newHP:    0,
			wantHP:   0,
			wantName: "Test",
		},
		{
			name:     "set negative hit points (should clamp to 0)",
			char:     NewCharacter("Test").WithHealth(20, 20),
			newHP:    -5,
			wantHP:   0,
			wantName: "Test",
		},
		{
			name:     "set hit points above maximum (should clamp)",
			char:     NewCharacter("Test").WithHealth(20, 20),
			newHP:    25,
			wantHP:   20,
			wantName: "Test",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			char := tt.char
			char.SetHP(tt.newHP)

			if got := char.HP(); got != tt.wantHP {
				t.Errorf("Character.HP() after SetHP(%v) = %v, want %v", tt.newHP, got, tt.wantHP)
			}
			if got := char.Name(); got != tt.wantName {
				t.Errorf("Character.Name() after SetHP() = %v, want %v", got, tt.wantName)
			}
		})
	}
}

func TestCharacter_AdjustHP(t *testing.T) {
	tests := []struct {
		name       string
		char       *Character
		adjustment int
		wantHP     int
		wantName   string
	}{
		{
			name:       "healing (positive adjustment)",
			char:       NewCharacter("Test").WithHealth(10, 20),
			adjustment: 5,
			wantHP:     15,
			wantName:   "Test",
		},
		{
			name:       "damage (negative adjustment)",
			char:       NewCharacter("Test").WithHealth(20, 20),
			adjustment: -8,
			wantHP:     12,
			wantName:   "Test",
		},
		{
			name:       "healing beyond maximum (should clamp)",
			char:       NewCharacter("Test").WithHealth(18, 20),
			adjustment: 5,
			wantHP:     20,
			wantName:   "Test",
		},
		{
			name:       "damage below zero (should clamp to 0)",
			char:       NewCharacter("Test").WithHealth(5, 20),
			adjustment: -10,
			wantHP:     0,
			wantName:   "Test",
		},
		{
			name:       "zero adjustment (no change)",
			char:       NewCharacter("Test").WithHealth(12, 20),
			adjustment: 0,
			wantHP:     12,
			wantName:   "Test",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			char := tt.char
			actualAdjustment := char.AdjustHP(tt.adjustment)

			// Test that we get the actual adjustment amount
			_ = actualAdjustment // We could test this if desired

			if got := char.HP(); got != tt.wantHP {
				t.Errorf("Character.HP() after AdjustHP(%v) = %v, want %v", tt.adjustment, got, tt.wantHP)
			}
			if got := char.Name(); got != tt.wantName {
				t.Errorf("Character.Name() after AdjustHP() = %v, want %v", got, tt.wantName)
			}
		})
	}
}
